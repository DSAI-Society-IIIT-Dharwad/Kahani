from milvus_service import milvus_client
from embedding_service import embedding_service
from llm_service import llm_service
from config import get_settings
import logging
from typing import List, Dict

logger = logging.getLogger(__name__)
settings = get_settings()


class RAGService:
    """RAG (Retrieval-Augmented Generation) Service"""
    
    def __init__(self):
        self.max_context_lines = settings.max_context_lines
    
    def retrieve_context(self, query: str, top_k: int = None, content_type: str = None) -> List[Dict]:
        """
        Retrieve relevant context from vector DB based on query
        
        Args:
            query: User's query or prompt
            top_k: Number of results to retrieve
            content_type: Filter by content type (story_line, lore, etc.)
        
        Returns:
            List of relevant context items
        """
        if top_k is None:
            top_k = self.max_context_lines
        
        # Check if Milvus is available
        if not milvus_client.is_connected():
            logger.warning("Milvus not connected - returning empty context")
            return []
        
        try:
            # Generate embedding for query
            query_embedding = embedding_service.generate_embedding(query)
            
            # Search in Milvus
            results = milvus_client.search_similar(
                query_embedding=query_embedding,
                top_k=top_k,
                content_type=content_type
            )
            
            context_items = []
            if results and len(results) > 0:
                for hits in results:
                    for hit in hits:
                        # Handle both dict (memory store) and entity (Milvus) formats
                        if hasattr(hit, 'entity') and hasattr(hit.entity, 'get'):
                            # Traditional Milvus format
                            context_items.append({
                                "text": hit.entity.get("text"),
                                "content_type": hit.entity.get("content_type"),
                                "embedding_id": hit.entity.get("embedding_id"),
                                "distance": hit.distance,
                                "score": 1 / (1 + hit.distance)
                            })
                        elif isinstance(hit.entity, dict):
                            # Memory store format
                            context_items.append({
                                "text": hit.entity.get("text"),
                                "content_type": hit.entity.get("content_type"),
                                "embedding_id": hit.entity.get("embedding_id"),
                                "distance": hit.distance,
                                "score": 1 / (1 + hit.distance)
                            })
            
            logger.info(f"Retrieved {len(context_items)} context items for query: {query[:50]}...")
            return context_items
            
        except Exception as e:
            logger.error(f"Context retrieval failed: {e}")
            return []
    
    def generate_with_context(self, user_prompt: str) -> Dict:
        """
        Complete RAG pipeline: retrieve context and generate response
        
        Args:
            user_prompt: User's story request
        
        Returns:
            Dictionary with suggestion and context used
        """
        try:
            # 1. Retrieve relevant context
            context_items = self.retrieve_context(user_prompt, content_type="story_line")
            
            # 2. Extract text from context
            context_texts = [item["text"] for item in context_items]
            
            # 3. Generate story suggestion using LLM
            suggestion = llm_service.generate_story_suggestion(
                user_prompt=user_prompt,
                context=context_texts
            )
            
            return {
                "suggestion": suggestion,
                "context_used": context_items,
                "context_count": len(context_items)
            }
            
        except Exception as e:
            logger.error(f"RAG generation failed: {e}")
            raise


# Global instance
rag_service = RAGService()
