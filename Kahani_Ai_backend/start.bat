@echo off
REM Kahani AI - Quick Start Script for Windows

echo ========================================
echo Kahani AI - RAG Story Generation System
echo ========================================
echo.

REM Check Python version
echo Checking Python version...
python --version
echo.

REM Create virtual environment
echo Creating virtual environment...
python -m venv venv

REM Activate virtual environment
echo Activating virtual environment...
call venv\Scripts\activate.bat

REM Install dependencies
echo.
echo Installing dependencies...
pip install --upgrade pip
pip install -r requirements.txt

REM Check .env file
echo.
if exist .env (
    echo .env file exists
) else (
    echo .env file not found. Creating from example...
    copy .env.example .env
    echo Please edit .env and add your GROQ_API_KEY
)

echo.
echo ========================================
echo Setup complete!
echo.
echo To start the application:
echo   python main.py
echo.
echo Or use uvicorn:
echo   uvicorn main:app --reload --host 0.0.0.0 --port 8000
echo.
echo Then open: http://localhost:8000/ui
echo ========================================
echo.
echo Note: Make sure Milvus is running:
echo   docker-compose up -d
echo ========================================

pause
