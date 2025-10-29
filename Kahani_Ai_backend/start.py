#!/usr/bin/env python3
"""
Kahani AI Quick Start
Simple script to start the application and show ngrok commands
"""

import os
import sys
import uvicorn
import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

def show_deployment_info():
    """Show deployment instructions"""
    print("\n" + "=" * 60)
    print("🎭 KAHANI AI - DEPLOYMENT INSTRUCTIONS")
    print("=" * 60)
    print("📱 TO DEPLOY ONLINE WITH NGROK:")
    print()
    print("1. Open a NEW terminal window")
    print("2. Run this command:")
    print("   ngrok http 8000")
    print()
    print("3. Copy the public URL from ngrok (e.g., https://abc123.ngrok-free.app)")
    print("4. Share that URL to access your AI storyteller!")
    print()
    print("🌐 YOUR ENDPOINTS WILL BE:")
    print("   • Homepage: https://your-ngrok-url/")
    print("   • Story UI: https://your-ngrok-url/ui")
    print("   • API Docs: https://your-ngrok-url/docs")
    print("   • Health: https://your-ngrok-url/health")
    print()
    print("💡 LOCAL ACCESS:")
    print("   • http://localhost:8000 (Homepage)")
    print("   • http://localhost:8000/ui (Story UI)")
    print("   • http://localhost:8000/docs (API Docs)")
    print("=" * 60 + "\n")

if __name__ == "__main__":
    try:
        print("🎭 Starting Kahani AI...")
        
        # Show deployment instructions
        show_deployment_info()
        
        # Start FastAPI server
        logger.info("🚀 Starting FastAPI server on port 8000...")
        print("🚀 Server starting... Ready for connections!")
        print("📝 To deploy online, follow the ngrok instructions above")
        print("🔄 Press Ctrl+C to stop the server\n")
        
        uvicorn.run(
            "main:app",
            host="0.0.0.0",
            port=8000,
            log_level="info",
            access_log=True
        )
        
    except KeyboardInterrupt:
        logger.info("👋 Shutting down Kahani AI...")
    except Exception as e:
        logger.error(f"❌ Error: {str(e)}")
        sys.exit(1)
