#!/bin/bash

echo "🧪 Testing Swagger Documentation..."

# Start server in background
go run main.go &
SERVER_PID=$!

# Wait for server to start
sleep 3

echo ""
echo "📖 Testing Swagger UI..."
curl -s -I http://localhost:8080/swagger/ | head -1

echo ""  
echo "📋 Testing API endpoints..."
echo "✅ Health check:"
curl -s http://localhost:8080/health | jq .

echo ""
echo "📄 API documentation available at:"
echo "   http://localhost:8080/swagger/"

echo ""
echo "🛑 Stopping server..."
kill $SERVER_PID 2>/dev/null || true

echo "✅ Test completed!"