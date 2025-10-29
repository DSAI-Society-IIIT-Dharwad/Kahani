from fastapi import FastAPI, Depends, HTTPException, BackgroundTasks
from fastapi.staticfiles import StaticFiles
from fastapi.responses import HTMLResponse, FileResponse
from sqlalchemy.orm import Session
from contextlib import asynccontextmanager
import logging
import hashlib
from datetime import datetime
from typing import List
import os

from config import get_settings
from database import get_db, init_db
from models import StoryLine, LoreEntry, CanonicalStory
from schemas import (
    StoryLineCreate,
    StoryLineEdit,
    StoryLineResponse,
    StoryLineVerifyRequest,
    LoreExtractRequest,
    LoreResponse,
    CanonicalizeRequest,
    CanonicalStoryResponse,
    ContextRetrievalRequest,
    HealthResponse
)
from milvus_service import milvus_client
from embedding_service import embedding_service
from rag_service import rag_service
from llm_service import llm_service

# Setup logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

settings = get_settings()


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Startup and shutdown events"""
    # Startup
    logger.info("Starting Kahani AI...")
    
    # Initialize database
    try:
        init_db()
        logger.info("Database initialized successfully")
    except Exception as e:
        logger.error(f"Database initialization failed: {e}")
        raise
    
    # Connect to Milvus (non-blocking)
    try:
        milvus_connected = milvus_client.connect()
        if milvus_connected:
            logger.info("Vector database connected successfully")
        else:
            logger.warning("Vector database not available - continuing without RAG features")
    except Exception as e:
        logger.warning(f"Vector database connection failed: {e} - continuing without RAG features")
    
    logger.info("Application started successfully")
    
    yield
    
    # Shutdown
    logger.info("Shutting down...")
    try:
        milvus_client.disconnect()
    except Exception as e:
        logger.warning(f"Error during Milvus disconnect: {e}")
    logger.info("Shutdown complete")


app = FastAPI(
    title=settings.app_name,
    version=settings.app_version,
    lifespan=lifespan
)

# Mount static files
from fastapi.staticfiles import StaticFiles
import os

static_path = os.path.join(os.path.dirname(__file__), "static")
if os.path.exists(static_path):
    app.mount("/static", StaticFiles(directory=static_path), name="static")


# ============= Health Check =============

@app.get("/", response_class=HTMLResponse)
async def root():
    """Root endpoint - redirect to docs"""
    return """
    <html>
        <head><title>Kahani AI - Story Generation System</title></head>
        <body style="font-family: Arial, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px;">
            <h1>🎭 Kahani AI - Story Generation System</h1>
            <p>Welcome to the RAG-based storytelling system!</p>
            <h2>Quick Links:</h2>
            <ul>
                <li><a href="/docs">📚 API Documentation (Swagger UI)</a></li>
                <li><a href="/redoc">📖 API Documentation (ReDoc)</a></li>
                <li><a href="/health">💚 Health Check</a></li>
                <li><a href="/ui">✍️ Story Writer UI</a></li>
            </ul>
            <h2>System Architecture:</h2>
            <p>User asks → RAG pulls context → LLM-1 suggests → User edits/signs → 
            Backend stores + embeds → Periodic extractors update vector DB → 
            LLM-2 canonicalizes final story</p>
        </body>
    </html>
    """


@app.get("/ui", response_class=FileResponse)
async def story_ui():
    """Serve the story writer UI"""
    static_path = os.path.join(os.path.dirname(__file__), "static", "index.html")
    return FileResponse(static_path)


@app.get("/health", response_model=HealthResponse)
async def health_check(db: Session = Depends(get_db)):
    """Health check endpoint"""
    try:
        # Check database
        db.execute("SELECT 1")
        db_ok = True
    except:
        db_ok = False
    
    # Check Milvus
    milvus_ok = milvus_client.is_connected()
    
    return HealthResponse(
        status="healthy" if (db_ok and milvus_ok) else "degraded",
        milvus_connected=milvus_ok,
        database_ok=db_ok
    )


# ============= Story Line Endpoints =============

@app.post("/api/story/suggest", response_model=StoryLineResponse)
async def suggest_story_line(
    request: StoryLineCreate,
    db: Session = Depends(get_db)
):
    """
    Step 1: User asks → RAG retrieves context → LLM-1 proposes story line(s)
    """
    try:
        # Use RAG to generate suggestion with context
        result = rag_service.generate_with_context(request.user_prompt)
        
        # Create a story line record (unverified)
        story_line = StoryLine(
            user_id=request.user_id,
            line_text=result["suggestion"],
            line_number=db.query(StoryLine).count() + 1,
            context_used=str(result["context_used"]),
            llm_proposed=result["suggestion"],
            user_edited=False,
            verified=False
        )
        
        db.add(story_line)
        db.commit()
        db.refresh(story_line)
        
        return StoryLineResponse(
            id=story_line.id,
            suggestion=result["suggestion"],
            context_used=result["context_used"],
            context_count=result["context_count"],
            verified=False
        )
        
    except Exception as e:
        logger.error(f"Story suggestion failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/story/edit", response_model=StoryLineResponse)
async def edit_and_sign_story_line(
    request: StoryLineEdit,
    background_tasks: BackgroundTasks,
    db: Session = Depends(get_db)
):
    """
    Step 2: User picks/edits the suggested line and signs it
    """
    try:
        # Create story line with user's final text
        story_line = StoryLine(
            user_id=request.user_id,
            line_text=request.final_text,
            line_number=db.query(StoryLine).count() + 1,
            llm_proposed=request.llm_proposed,
            user_edited=(request.final_text != request.llm_proposed),
            verified=True,  # Auto-verify when user signs
            signature=hashlib.sha256(request.final_text.encode()).hexdigest()[:16]
        )
        
        db.add(story_line)
        db.commit()
        db.refresh(story_line)
        
        # Background: Create embedding and store in Milvus
        background_tasks.add_task(
            store_embedding,
            story_line.id,
            request.final_text,
            db
        )
        
        return StoryLineResponse(
            id=story_line.id,
            suggestion=request.final_text,
            context_used=[],
            context_count=0,
            verified=True,
            embedding_id=f"story_line_{story_line.id}"
        )
        
    except Exception as e:
        logger.error(f"Story line edit failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/story/verify/{line_id}")
async def verify_story_line(
    line_id: int,
    request: StoryLineVerifyRequest,
    background_tasks: BackgroundTasks,
    db: Session = Depends(get_db)
):
    """
    Step 3: Verify and store a story line, create embedding
    """
    story_line = db.query(StoryLine).filter(StoryLine.id == line_id).first()
    
    if not story_line:
        raise HTTPException(status_code=404, detail="Story line not found")
    
    # Update verification
    story_line.verified = True
    story_line.signature = request.signature
    db.commit()
    
    # Background: Create and store embedding
    background_tasks.add_task(
        store_embedding,
        story_line.id,
        story_line.line_text,
        db
    )
    
    return {"status": "verified", "line_id": line_id}


def store_embedding(line_id: int, text: str, db: Session):
    """Background task to create and store embedding"""
    try:
        # Generate embedding
        embedding = embedding_service.generate_embedding(text)
        
        # Store in Milvus
        embedding_id = f"story_line_{line_id}"
        milvus_client.insert_embedding(
            embedding_id=embedding_id,
            text=text,
            embedding=embedding,
            content_type="story_line"
        )
        
        # Update story line with embedding ID
        story_line = db.query(StoryLine).filter(StoryLine.id == line_id).first()
        if story_line:
            story_line.embedding_id = embedding_id
            db.commit()
        
        logger.info(f"Stored embedding for story line {line_id}")
        
    except Exception as e:
        logger.error(f"Failed to store embedding: {e}")


@app.get("/api/story/lines", response_model=List[dict])
async def get_all_story_lines(
    verified_only: bool = False,
    db: Session = Depends(get_db)
):
    """Get all story lines"""
    query = db.query(StoryLine)
    
    if verified_only:
        query = query.filter(StoryLine.verified == True)
    
    lines = query.order_by(StoryLine.line_number).all()
    
    return [
        {
            "id": line.id,
            "line_number": line.line_number,
            "text": line.line_text,
            "verified": line.verified,
            "user_edited": line.user_edited,
            "created_at": line.created_at
        }
        for line in lines
    ]


# ============= Lore Extraction Endpoints =============

@app.post("/api/lore/extract", response_model=LoreResponse)
async def extract_lore(
    request: LoreExtractRequest,
    background_tasks: BackgroundTasks,
    db: Session = Depends(get_db)
):
    """
    Step 4: Extract lore/entities from story lines
    """
    # Get story lines
    story_lines = db.query(StoryLine).filter(
        StoryLine.id.in_(request.line_ids)
    ).all()
    
    if not story_lines:
        raise HTTPException(status_code=404, detail="No story lines found")
    
    # Extract text
    texts = [line.line_text for line in story_lines]
    
    # Use LLM to extract lore
    lore_data = llm_service.extract_lore(texts)
    
    # Store lore entries and create embeddings in background
    background_tasks.add_task(
        store_lore_entries,
        lore_data,
        request.line_ids,
        db
    )
    
    total = sum([
        len(lore_data.get("characters", [])),
        len(lore_data.get("locations", [])),
        len(lore_data.get("events", [])),
        len(lore_data.get("items", []))
    ])
    
    return LoreResponse(
        characters=lore_data.get("characters", []),
        locations=lore_data.get("locations", []),
        events=lore_data.get("events", []),
        items=lore_data.get("items", []),
        total_entries=total
    )


def store_lore_entries(lore_data: dict, source_line_ids: List[int], db: Session):
    """Background task to store lore entries with embeddings"""
    try:
        entity_types = ["characters", "locations", "events", "items"]
        
        for entity_type in entity_types:
            entities = lore_data.get(entity_type, [])
            
            for entity in entities:
                name = entity.get("name", "")
                description = entity.get("description", "")
                
                # Create lore entry
                lore_entry = LoreEntry(
                    entity_type=entity_type[:-1],  # Remove 's'
                    entity_name=name,
                    description=description,
                    source_lines=str(source_line_ids),
                    confidence=0.8
                )
                
                db.add(lore_entry)
                db.commit()
                db.refresh(lore_entry)
                
                # Create embedding for lore
                text_for_embedding = f"{name}: {description}"
                embedding = embedding_service.generate_embedding(text_for_embedding)
                
                embedding_id = f"lore_{entity_type}_{lore_entry.id}"
                milvus_client.insert_embedding(
                    embedding_id=embedding_id,
                    text=text_for_embedding,
                    embedding=embedding,
                    content_type=f"lore_{entity_type}"
                )
                
                lore_entry.embedding_id = embedding_id
                db.commit()
        
        logger.info(f"Stored lore entries from lines {source_line_ids}")
        
    except Exception as e:
        logger.error(f"Failed to store lore entries: {e}")


@app.get("/api/lore/all")
async def get_all_lore(db: Session = Depends(get_db)):
    """Get all lore entries"""
    lore_entries = db.query(LoreEntry).all()
    
    # Group by type
    grouped = {
        "characters": [],
        "locations": [],
        "events": [],
        "items": []
    }
    
    for entry in lore_entries:
        entity_type = f"{entry.entity_type}s"
        if entity_type in grouped:
            grouped[entity_type].append({
                "id": entry.id,
                "name": entry.entity_name,
                "description": entry.description,
                "confidence": entry.confidence
            })
    
    return grouped


# ============= Canonicalization Endpoints =============

@app.post("/api/story/canonicalize", response_model=CanonicalStoryResponse)
async def canonicalize_story(
    request: CanonicalizeRequest,
    db: Session = Depends(get_db)
):
    """
    Final step: LLM-2 creates canonical version of the story
    """
    # Get story lines
    query = db.query(StoryLine).filter(StoryLine.verified == True)
    
    if request.line_ids:
        query = query.filter(StoryLine.id.in_(request.line_ids))
    
    story_lines = query.order_by(StoryLine.line_number).all()
    
    if not story_lines:
        raise HTTPException(status_code=404, detail="No verified story lines found")
    
    # Extract text
    texts = [line.line_text for line in story_lines]
    
    # Use LLM-2 to canonicalize
    canonical_text = llm_service.canonicalize_story(texts)
    
    # Store canonical version
    canonical_story = CanonicalStory(
        title=request.title or "Untitled Story",
        full_text=canonical_text,
        canonicalized_by=f"LLM-2 ({settings.llm_model})",
        original_lines_count=len(story_lines),
        finalized_at=datetime.utcnow()
    )
    
    db.add(canonical_story)
    db.commit()
    db.refresh(canonical_story)
    
    return CanonicalStoryResponse(
        id=canonical_story.id,
        title=canonical_story.title,
        full_text=canonical_story.full_text,
        original_lines_count=canonical_story.original_lines_count,
        created_at=canonical_story.created_at
    )


@app.get("/api/story/canonical/{story_id}")
async def get_canonical_story(story_id: int, db: Session = Depends(get_db)):
    """Get a canonical story by ID"""
    story = db.query(CanonicalStory).filter(CanonicalStory.id == story_id).first()
    
    if not story:
        raise HTTPException(status_code=404, detail="Canonical story not found")
    
    return {
        "id": story.id,
        "title": story.title,
        "full_text": story.full_text,
        "original_lines_count": story.original_lines_count,
        "created_at": story.created_at,
        "finalized_at": story.finalized_at
    }


# ============= Context Retrieval Endpoint =============

@app.post("/api/context/retrieve")
async def retrieve_context(request: ContextRetrievalRequest):
    """Retrieve relevant context from vector DB"""
    context_items = rag_service.retrieve_context(
        query=request.query,
        top_k=request.top_k,
        content_type=request.content_type
    )
    
    return {
        "query": request.query,
        "results": context_items,
        "count": len(context_items)
    }


if __name__ == "__main__":
    import uvicorn
    import threading
    import time
    
    # Optional ngrok integration
    ENABLE_NGROK = os.getenv('ENABLE_NGROK', 'false').lower() == 'true'
    
    if ENABLE_NGROK:
        try:
            from ngrok_service import ngrok_service
            
            def start_ngrok():
                # Wait for FastAPI to start
                time.sleep(5)
                public_url = ngrok_service.start_tunnel(port=8000)
                if public_url:
                    print(f"\n🌐 Kahani AI is live at: {public_url}")
                    print(f"📚 API Docs: {public_url}/docs")
                    print(f"✍️ Story UI: {public_url}/ui\n")
            
            # Start ngrok in background
            ngrok_thread = threading.Thread(target=start_ngrok, daemon=True)
            ngrok_thread.start()
            
        except ImportError:
            logger.warning("Ngrok service not available. Install pyngrok to enable.")
    
    # Start FastAPI server
    uvicorn.run(app, host="0.0.0.0", port=8000)
