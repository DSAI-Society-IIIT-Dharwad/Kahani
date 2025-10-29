from sqlalchemy import Column, Integer, String, Text, DateTime, Boolean, Float
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.sql import func
from datetime import datetime

Base = declarative_base()


class StoryLine(Base):
    """Stores user-signed story lines"""
    __tablename__ = "story_lines"
    
    id = Column(Integer, primary_key=True, index=True)
    user_id = Column(String, index=True)  # Could be user identifier
    line_text = Column(Text, nullable=False)
    line_number = Column(Integer, nullable=False)  # Sequential line number
    
    # RAG Context used
    context_used = Column(Text)  # JSON or text of context pulled
    
    # LLM proposed vs user final
    llm_proposed = Column(Text)  # What LLM-1 suggested
    user_edited = Column(Boolean, default=False)  # Did user edit?
    
    # Verification
    verified = Column(Boolean, default=False)
    signature = Column(String)  # Hash or signature for verification
    
    # Timestamps
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), onupdate=func.now())
    
    # Embedding reference
    embedding_id = Column(String, index=True)  # Reference to Milvus vector


class LoreEntry(Base):
    """Stores extracted lore/entities from story"""
    __tablename__ = "lore_entries"
    
    id = Column(Integer, primary_key=True, index=True)
    entity_type = Column(String, index=True)  # character, location, event, etc.
    entity_name = Column(String, index=True)
    description = Column(Text)
    
    # Context
    source_lines = Column(Text)  # Which story lines this came from
    confidence = Column(Float, default=0.0)
    
    # Timestamps
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), onupdate=func.now())
    
    # Embedding reference
    embedding_id = Column(String, index=True)


class CanonicalStory(Base):
    """Stores finalized canonical versions"""
    __tablename__ = "canonical_stories"
    
    id = Column(Integer, primary_key=True, index=True)
    title = Column(String)
    full_text = Column(Text, nullable=False)
    
    # Canonicalization metadata
    canonicalized_by = Column(String)  # LLM-2
    original_lines_count = Column(Integer)
    
    # Timestamps
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    finalized_at = Column(DateTime(timezone=True))
    
    version = Column(Integer, default=1)
