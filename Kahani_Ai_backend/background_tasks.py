from apscheduler.schedulers.background import BackgroundScheduler
from sqlalchemy.orm import Session
from database import SessionLocal
from models import StoryLine
from llm_service import llm_service
from embedding_service import embedding_service
from milvus_service import milvus_client
import logging

logger = logging.getLogger(__name__)


class BackgroundTasksService:
    """Periodic background tasks for lore extraction and summarization"""
    
    def __init__(self):
        self.scheduler = BackgroundScheduler()
    
    def start(self):
        """Start the scheduler"""
        # Extract lore every 30 minutes
        self.scheduler.add_job(
            func=self.periodic_lore_extraction,
            trigger="interval",
            minutes=30,
            id="lore_extraction"
        )
        
        # Update summaries every hour
        self.scheduler.add_job(
            func=self.periodic_summary_update,
            trigger="interval",
            hours=1,
            id="summary_update"
        )
        
        self.scheduler.start()
        logger.info("Background tasks scheduler started")
    
    def stop(self):
        """Stop the scheduler"""
        self.scheduler.shutdown()
        logger.info("Background tasks scheduler stopped")
    
    def periodic_lore_extraction(self):
        """
        Periodic task: Extract lore from recent story lines
        """
        logger.info("Running periodic lore extraction...")
        db = SessionLocal()
        
        try:
            # Get recent verified story lines without lore extraction
            # (You could add a flag to track which lines have been processed)
            recent_lines = db.query(StoryLine).filter(
                StoryLine.verified == True
            ).order_by(StoryLine.created_at.desc()).limit(20).all()
            
            if not recent_lines:
                logger.info("No recent lines to process")
                return
            
            texts = [line.line_text for line in recent_lines]
            
            # Extract lore
            lore_data = llm_service.extract_lore(texts)
            
            # Store lore with embeddings
            from main import store_lore_entries
            line_ids = [line.id for line in recent_lines]
            store_lore_entries(lore_data, line_ids, db)
            
            logger.info(f"Extracted lore from {len(recent_lines)} story lines")
            
        except Exception as e:
            logger.error(f"Periodic lore extraction failed: {e}")
        finally:
            db.close()
    
    def periodic_summary_update(self):
        """
        Periodic task: Generate and store summaries of story sections
        """
        logger.info("Running periodic summary update...")
        db = SessionLocal()
        
        try:
            # Get all verified story lines
            all_lines = db.query(StoryLine).filter(
                StoryLine.verified == True
            ).order_by(StoryLine.line_number).all()
            
            if len(all_lines) < 5:
                logger.info("Not enough lines for summary")
                return
            
            # Create summaries for chunks of 10 lines
            chunk_size = 10
            for i in range(0, len(all_lines), chunk_size):
                chunk = all_lines[i:i+chunk_size]
                texts = [line.line_text for line in chunk]
                
                # Generate summary
                summary = llm_service.summarize_story(texts)
                
                # Create embedding for summary
                embedding = embedding_service.generate_embedding(summary)
                
                # Store summary in Milvus
                embedding_id = f"summary_{i}_{i+chunk_size}"
                milvus_client.insert_embedding(
                    embedding_id=embedding_id,
                    text=summary,
                    embedding=embedding,
                    content_type="summary"
                )
                
                logger.info(f"Created summary for lines {i}-{i+chunk_size}")
            
        except Exception as e:
            logger.error(f"Periodic summary update failed: {e}")
        finally:
            db.close()


# Global instance
background_tasks_service = BackgroundTasksService()
