#!/usr/bin/env python3
"""
Milvus Lite Server Starter
Starts Milvus Lite server and keeps it running
"""

import time
import logging
import sys
import signal
from milvus import default_server

# Setup logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

def signal_handler(signum, frame):
    """Handle shutdown signals"""
    logger.info("Received shutdown signal, stopping Milvus server...")
    try:
        default_server.stop()
        logger.info("Milvus server stopped successfully")
    except Exception as e:
        logger.error(f"Error stopping Milvus server: {e}")
    sys.exit(0)

def start_milvus_server():
    """Start Milvus Lite server"""
    try:
        logger.info("🚀 Starting Milvus Lite server...")
        
        # Register signal handlers
        signal.signal(signal.SIGINT, signal_handler)
        signal.signal(signal.SIGTERM, signal_handler)
        
        # Start the server
        default_server.start()
        logger.info("✅ Milvus Lite server started successfully on localhost:19530")
        
        # Keep the server running
        logger.info("Server is running. Press Ctrl+C to stop.")
        try:
            while True:
                time.sleep(1)
        except KeyboardInterrupt:
            signal_handler(signal.SIGINT, None)
            
    except Exception as e:
        logger.error(f"❌ Failed to start Milvus server: {e}")
        sys.exit(1)

if __name__ == "__main__":
    start_milvus_server()
