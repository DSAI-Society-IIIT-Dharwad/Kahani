#!/usr/bin/env python3
"""
Kahani AI Launcher with Integrated Ngrok
Automatically starts FastAPI app and creates ngrok tunnel
"""

import os
import sys
import time
import logging
import threading
import uvicorn
from contextlib import asynccontextmanager

# Add current directory to path
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

from ngrok_service import ngrok_service
from config import get_settings

# Setup logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

settings = get_settings()


def start_ngrok_tunnel():
    """Start ngrok tunnel in a separate thread"""
    def run_tunnel():
        # Wait a bit for FastAPI to start
        time.sleep(5)
        
        logger.info("🚀 Starting ngrok tunnel...")
        public_url = ngrok_service.start_tunnel(port=8000)
        
        if public_url:
            print("\n" + "="*60)
            print("🎉 KAHANI AI IS NOW LIVE ONLINE!")
            print("="*60)
            print(f"🌐 Public URL: {public_url}")
            print(f"🏠 Homepage: {public_url}/")
            print(f"✍️  Story UI: {public_url}/ui")
            print(f"📚 API Docs: {public_url}/docs")
            print(f"💚 Health Check: {public_url}/health")
            print("="*60)
            print("📱 Share this URL with anyone to access your AI storyteller!")
            print("🔄 Press Ctrl+C to stop both services")
            print("="*60 + "\n")
        else:
            logger.error("❌ Failed to start ngrok tunnel")
    
    # Start tunnel in background thread
    tunnel_thread = threading.Thread(target=run_tunnel, daemon=True)
    tunnel_thread.start()


def cleanup_on_exit():
    """Cleanup function"""
    logger.info("🔄 Shutting down Kahani AI...")
    ngrok_service.cleanup()
    logger.info("✅ Cleanup completed")


if __name__ == "__main__":
    try:
        print("🎭 Starting Kahani AI with Ngrok Integration...")
        
        # Start ngrok tunnel in background
        start_ngrok_tunnel()
        
        # Start FastAPI server
        logger.info("🚀 Starting FastAPI server on port 8000...")
        uvicorn.run(
            "main:app",
            host="0.0.0.0",
            port=8000,
            log_level="info",
            access_log=True
        )
        
    except KeyboardInterrupt:
        logger.info("👋 Received shutdown signal")
    except Exception as e:
        logger.error(f"❌ Error: {str(e)}")
    finally:
        cleanup_on_exit()
