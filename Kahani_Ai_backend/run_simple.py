#!/usr/bin/env python3
"""
Simple Kahani AI Launcher
Starts FastAPI app and provides instructions for ngrok
"""

import os
import sys
import time
import logging
import subprocess
import signal
import uvicorn

# Setup logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

def start_ngrok_process():
    """Start ngrok as a separate process"""
    try:
        logger.info("🚀 Starting ngrok tunnel...")
        ngrok_process = subprocess.Popen(
            ["ngrok", "http", "8000"],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE
        )
        
        # Give ngrok time to start
        time.sleep(3)
        
        # Check if process is still running
        if ngrok_process.poll() is None:
            logger.info("✅ Ngrok tunnel started successfully!")
            logger.info("🌐 Check ngrok web interface at: http://localhost:4040")
            return ngrok_process
        else:
            stdout, stderr = ngrok_process.communicate()
            logger.error(f"❌ Ngrok failed to start: {stderr.decode()}")
            return None
            
    except Exception as e:
        logger.error(f"❌ Failed to start ngrok: {str(e)}")
        return None

def cleanup(ngrok_process=None):
    """Cleanup function"""
    logger.info("🔄 Shutting down Kahani AI...")
    
    if ngrok_process and ngrok_process.poll() is None:
        logger.info("Stopping ngrok...")
        ngrok_process.terminate()
        ngrok_process.wait()
    
    logger.info("✅ Cleanup completed")

if __name__ == "__main__":
    ngrok_process = None
    
    try:
        print("🎭 Starting Kahani AI with Ngrok Integration...")
        print("=" * 50)
        
        # Start ngrok in background
        ngrok_process = start_ngrok_process()
        
        if ngrok_process:
            print("\n🎉 SERVICES STARTING...")
            print("🌐 Ngrok tunnel: Starting...")
            print("🚀 FastAPI server: Starting on port 8000...")
            print("\n📱 Once both services are ready:")
            print("   • Visit http://localhost:4040 to see your public URL")
            print("   • Your API will be available at that public URL")
            print("   • Press Ctrl+C to stop both services")
            print("=" * 50 + "\n")
        
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
        cleanup(ngrok_process)
