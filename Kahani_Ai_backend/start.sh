#!/bin/bash

# Kahani AI - Quick Start Script

echo "🎭 Kahani AI - RAG Story Generation System"
echo "=========================================="
echo ""

# Check Python version
echo "Checking Python version..."
python3 --version

# Create virtual environment
echo ""
echo "Creating virtual environment..."
python3 -m venv venv

# Activate virtual environment
echo "Activating virtual environment..."
source venv/bin/activate

# Install dependencies
echo ""
echo "Installing dependencies..."
pip install --upgrade pip
pip install -r requirements.txt

# Check if Milvus is running
echo ""
echo "Checking Milvus connection..."
if curl -s http://localhost:19530 > /dev/null 2>&1; then
    echo "✅ Milvus is running"
else
    echo "❌ Milvus is not running!"
    echo ""
    echo "Please start Milvus with Docker:"
    echo "  docker-compose up -d"
    echo ""
    echo "Or install Milvus standalone:"
    echo "  See: https://milvus.io/docs/install_standalone-docker.md"
    echo ""
    read -p "Do you want to start Milvus now? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        docker-compose up -d
        echo "Waiting for Milvus to start..."
        sleep 10
    fi
fi

# Check .env file
echo ""
if [ -f .env ]; then
    echo "✅ .env file exists"
else
    echo "⚠️  .env file not found. Creating from example..."
    cp .env.example .env
    echo "Please edit .env and add your GROQ_API_KEY"
fi

echo ""
echo "=========================================="
echo "Setup complete! 🎉"
echo ""
echo "To start the application:"
echo "  python main.py"
echo ""
echo "Or use uvicorn:"
echo "  uvicorn main:app --reload --host 0.0.0.0 --port 8000"
echo ""
echo "Then open: http://localhost:8000/ui"
echo "=========================================="
