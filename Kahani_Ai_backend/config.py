from pydantic_settings import BaseSettings
from functools import lru_cache


class Settings(BaseSettings):
    # Groq API
    groq_api_key: str
    
    # Milvus
    milvus_host: str = "localhost"
    milvus_port: int = 19530
    milvus_collection_name: str = "story_embeddings"
    
    # Database
    database_url: str = "sqlite:///./kahani.db"
    
    # App
    app_name: str = "Kahani AI"
    app_version: str = "1.0.0"
    debug: bool = True
    
    # LLM
    llm_model: str = "llama-3.1-70b-versatile"
    embedding_model: str = "all-MiniLM-L6-v2"
    embedding_dim: int = 384
    max_context_lines: int = 10
    
    # Ngrok
    ngrok_auth_token: str = ""
    
    class Config:
        env_file = ".env"
        case_sensitive = False


@lru_cache()
def get_settings():
    return Settings()
