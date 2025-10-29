from pydantic import BaseModel, Field
from typing import List, Optional, Dict
from datetime import datetime


class StoryLineCreate(BaseModel):
    """Request to create a new story line"""
    user_prompt: str = Field(..., description="User's request for next story line")
    user_id: str = Field(default="default_user", description="User identifier")
    
    
class StoryLineEdit(BaseModel):
    """User's edited version of suggested line"""
    llm_proposed: str = Field(..., description="Original LLM suggestion")
    final_text: str = Field(..., description="User's final edited text")
    user_id: str = Field(default="default_user", description="User identifier")


class StoryLineResponse(BaseModel):
    """Response with story line suggestion"""
    id: Optional[int] = None
    suggestion: str
    context_used: List[Dict]
    context_count: int
    verified: bool = False
    embedding_id: Optional[str] = None
    

class StoryLineVerifyRequest(BaseModel):
    """Request to verify and store a story line"""
    line_id: int
    signature: str = Field(default="user_signed", description="User signature/hash")


class LoreExtractRequest(BaseModel):
    """Request to extract lore from story lines"""
    line_ids: List[int] = Field(..., description="Story line IDs to extract lore from")


class LoreResponse(BaseModel):
    """Response with extracted lore"""
    characters: List[Dict]
    locations: List[Dict]
    events: List[Dict]
    items: List[Dict]
    total_entries: int


class CanonicalizeRequest(BaseModel):
    """Request to canonicalize story"""
    line_ids: Optional[List[int]] = Field(None, description="Specific line IDs, or all if None")
    title: Optional[str] = Field(None, description="Story title")


class CanonicalStoryResponse(BaseModel):
    """Response with canonicalized story"""
    id: int
    title: Optional[str]
    full_text: str
    original_lines_count: int
    created_at: datetime


class ContextRetrievalRequest(BaseModel):
    """Request to retrieve context"""
    query: str
    top_k: int = Field(default=5, description="Number of results")
    content_type: Optional[str] = Field(None, description="Filter by content type")


class HealthResponse(BaseModel):
    """Health check response"""
    status: str
    milvus_connected: bool
    database_ok: bool
