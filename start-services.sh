#!/bin/bash

# NCV Services Startup Script
# This script starts all necessary services for the NCV dashboard

set -e

echo "🚀 Starting NCV Services..."

# Navigate to project root
cd "$(dirname "$0")"

# Check if virtual environment exists
if [ ! -d "venv" ]; then
    echo "❌ Virtual environment not found. Please run: python3 -m venv venv"
    exit 1
fi

# Check if dependencies are installed
if ! ./venv/bin/python -c "import requests, substrateinterface" 2>/dev/null; then
    echo "📦 Installing Python dependencies..."
    ./venv/bin/pip install requests substrate-interface
fi

# Start log server in background
echo "📡 Starting log server on port 8081..."
./venv/bin/python log_server.py &
LOG_SERVER_PID=$!
echo "   Log server PID: $LOG_SERVER_PID"

# Wait for log server to start
sleep 2

# Test log server
if curl -s http://localhost:8081/logs > /dev/null; then
    echo "✅ Log server is running"
else
    echo "❌ Log server failed to start"
    kill $LOG_SERVER_PID 2>/dev/null || true
    exit 1
fi

# Optional: Run demo to generate initial logs
read -p "🎯 Run demo to generate verification logs? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🔍 Running verification demo..."
    ./venv/bin/python demo.py
    echo "✅ Demo completed"
fi

echo ""
echo "✨ All services started successfully!"
echo ""
echo "📊 To view the dashboard:"
echo "   cd frontend"
echo "   npm start"
echo ""
echo "🛑 To stop the log server:"
echo "   kill $LOG_SERVER_PID"
echo ""
echo "💡 Log server is running at: http://localhost:8081/logs"
