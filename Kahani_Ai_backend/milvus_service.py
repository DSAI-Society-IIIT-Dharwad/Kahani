from pymilvus import (
    connections,
    Collection,
    CollectionSchema,
    FieldSchema,
    DataType,
    utility
)
from config import get_settings
import logging
import os
import numpy as np
from typing import List, Dict

logger = logging.getLogger(__name__)
settings = get_settings()


class MilvusClient:
    """Milvus Vector Database Client with in-memory fallback"""
    
    def __init__(self):
        self.collection_name = settings.milvus_collection_name
        self.collection = None
        self.connected = False
        # In-memory fallback storage
        self.memory_store = []
        
    def connect(self):
        """Connect to Milvus with in-memory fallback"""
        try:
            # Try traditional Milvus server first
            connections.connect(
                alias="default",
                host=settings.milvus_host,
                port=settings.milvus_port,
                timeout=5  # Short timeout for quick fallback
            )
            logger.info(f"Connected to Milvus server at {settings.milvus_host}:{settings.milvus_port}")
            self._init_collection()
            self.connected = True
            return True
                
        except Exception as e:
            logger.warning(f"Traditional Milvus failed: {e}")
            logger.info("Using in-memory vector store as fallback")
            
            # Use in-memory storage as fallback
            self.connected = True  # Mark as connected for in-memory mode
            logger.info("✅ In-memory vector store initialized - RAG will work with limited persistence")
            return True
    
    def _cosine_similarity(self, vec1: List[float], vec2: List[float]) -> float:
        """Calculate cosine similarity between two vectors"""
        try:
            vec1 = np.array(vec1)
            vec2 = np.array(vec2)
            return np.dot(vec1, vec2) / (np.linalg.norm(vec1) * np.linalg.norm(vec2))
        except:
            return 0.0
    
    def _init_collection(self):
        """Initialize or load collection (traditional Milvus)"""
        if utility.has_collection(self.collection_name):
            self.collection = Collection(self.collection_name)
            logger.info(f"Loaded existing collection: {self.collection_name}")
        else:
            self._create_collection()
    
    def _create_collection(self):
        """Create new collection with schema"""
        fields = [
            FieldSchema(name="id", dtype=DataType.INT64, is_primary=True, auto_id=True),
            FieldSchema(name="embedding_id", dtype=DataType.VARCHAR, max_length=100),
            FieldSchema(name="text", dtype=DataType.VARCHAR, max_length=5000),
            FieldSchema(name="content_type", dtype=DataType.VARCHAR, max_length=50),  # story_line, lore, etc.
            FieldSchema(name="embedding", dtype=DataType.FLOAT_VECTOR, dim=settings.embedding_dim)
        ]
        
        schema = CollectionSchema(
            fields=fields,
            description="Story embeddings for RAG retrieval"
        )
        
        self.collection = Collection(
            name=self.collection_name,
            schema=schema
        )
        
        # Create index for vector search
        index_params = {
            "metric_type": "L2",
            "index_type": "IVF_FLAT",
            "params": {"nlist": 128}
        }
        
        self.collection.create_index(
            field_name="embedding",
            index_params=index_params
        )
        
        logger.info(f"Created new collection: {self.collection_name}")
    
    def insert_embedding(self, embedding_id: str, text: str, embedding: list, content_type: str = "story_line"):
        """Insert an embedding into Milvus or memory store"""
        try:
            if not self.connected:
                logger.warning("No vector store available")
                return False
                
            if self.collection:
                # Use traditional Milvus
                data = [
                    [embedding_id],
                    [text],
                    [content_type],
                    [embedding]
                ]
                
                self.collection.insert(data)
                self.collection.flush()
                logger.info(f"Inserted embedding (Milvus): {embedding_id}")
            else:
                # Use in-memory storage
                self.memory_store.append({
                    "embedding_id": embedding_id,
                    "text": text,
                    "content_type": content_type,
                    "embedding": embedding
                })
                logger.info(f"Inserted embedding (Memory): {embedding_id}")
            
            return True
        except Exception as e:
            logger.error(f"Failed to insert embedding: {e}")
            return False
    
    def search_similar(self, query_embedding: list, top_k: int = 5, content_type: str = None):
        """Search for similar embeddings"""
        try:
            if not self.connected:
                logger.warning("No vector store available")
                return []
                
            if self.collection:
                # Use traditional Milvus
                self.collection.load()
                
                search_params = {
                    "metric_type": "L2",
                    "params": {"nprobe": 10}
                }
                
                # Add filter if content_type specified
                expr = f'content_type == "{content_type}"' if content_type else None
                
                results = self.collection.search(
                    data=[query_embedding],
                    anns_field="embedding",
                    param=search_params,
                    limit=top_k,
                    expr=expr,
                    output_fields=["embedding_id", "text", "content_type"]
                )
                return results
            else:
                # Use in-memory search
                candidates = self.memory_store
                
                # Filter by content type if specified
                if content_type:
                    candidates = [item for item in candidates if item["content_type"] == content_type]
                
                # Calculate similarities
                similarities = []
                for item in candidates:
                    similarity = self._cosine_similarity(query_embedding, item["embedding"])
                    similarities.append((similarity, item))
                
                # Sort by similarity (descending) and take top_k
                similarities.sort(key=lambda x: x[0], reverse=True)
                top_items = similarities[:top_k]
                
                # Format results to match Milvus format
                class MockResult:
                    def __init__(self, items):
                        self.items = items
                    
                    def __iter__(self):
                        for similarity, item in self.items:
                            yield MockHit(item, 1.0 - similarity)  # Convert similarity to distance
                
                class MockHit:
                    def __init__(self, item, distance):
                        self.entity = item
                        self.distance = distance
                
                return [MockResult(top_items)] if top_items else []
                
        except Exception as e:
            logger.error(f"Search failed: {e}")
            return []
    
    def disconnect(self):
        """Disconnect from Milvus"""
        try:
            if self.collection:
                connections.disconnect("default")
                logger.info("Disconnected from Milvus")
            else:
                logger.info("Cleared in-memory vector store")
                self.memory_store.clear()
            self.connected = False
        except Exception as e:
            logger.warning(f"Error during disconnect: {e}")
    
    def is_connected(self):
        """Check if vector store is available"""
        return self.connected


# Global instance
milvus_client = MilvusClient()
