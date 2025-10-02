#!/bin/bash

# API Testing Script
BASE_URL="http://localhost:8080"

echo "=== Testing Scholarship Management System APIs ==="
echo ""

# Get token
echo "1. Testing Login..."
LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')
echo "✓ Login successful, got token"
echo ""

# Test Scholarships
echo "2. Testing Get Scholarships..."
curl -s $BASE_URL/api/v1/scholarships \
  -H "Authorization: Bearer $TOKEN" | jq '.' | head -30
echo ""

# Test Payment Methods
echo "3. Testing Get Payment Methods (requires admin)..."
curl -s $BASE_URL/api/v1/payments/methods \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo ""

# Test Analytics Dashboard
echo "4. Testing Analytics Dashboard (requires admin)..."
curl -s $BASE_URL/api/v1/analytics/dashboard \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo ""

# Test User Profile
echo "5. Testing Get User Profile..."
curl -s $BASE_URL/api/v1/user/profile \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo ""

# Test News (public)
echo "6. Testing Get News (public)..."
curl -s $BASE_URL/api/v1/news | jq '.' | head -20
echo ""

echo "=== All tests completed ==="
