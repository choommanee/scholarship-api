#!/bin/bash

# Admin API Testing Script
BASE_URL="http://localhost:8080"
ADMIN_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiODE4NjkxNjItZWYzMC00M2NjLTgyMmItNzkzYTljNTRlY2ZiIiwiZW1haWwiOiJhZG1pbkB1bml2ZXJzaXR5LmFjLnRoIiwidXNlcm5hbWUiOiJhZG1pbjEiLCJyb2xlcyI6WyJhZG1pbiJdLCJzdWIiOiI4MTg2OTE2Mi1lZjMwLTQzY2MtODIyYi03OTNhOWM1NGVjZmIiLCJleHAiOjE3NTk5NDIwOTEsImlhdCI6MTc1OTMzNzI5MX0.CfWj-gjFY9IpxqiTkriXSnQhMobWIC9TMDhZD8dIbe4"

echo "=== Testing Admin APIs ==="
echo ""

# Test Payment Methods
echo "1. Testing Get Payment Methods..."
curl -s $BASE_URL/api/v1/payments/methods \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq '.'
echo ""

# Test Analytics Dashboard
echo "2. Testing Analytics Dashboard..."
curl -s $BASE_URL/api/v1/analytics/dashboard \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq '.'
echo ""

# Test Scholarship Statistics
echo "3. Testing Scholarship Statistics..."
curl -s "$BASE_URL/api/v1/analytics/statistics?year=2023&round=1" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq '.'
echo ""

# Test All Statistics
echo "4. Testing Get All Statistics..."
curl -s $BASE_URL/api/v1/analytics/statistics/all \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq '.' | head -50
echo ""

# Test Processing Time
echo "5. Testing Average Processing Time..."
curl -s $BASE_URL/api/v1/analytics/processing-time \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq '.'
echo ""

# Test Bottlenecks
echo "6. Testing Bottlenecks Analysis..."
curl -s $BASE_URL/api/v1/analytics/bottlenecks \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq '.'
echo ""

echo "=== All Admin API tests completed ==="
