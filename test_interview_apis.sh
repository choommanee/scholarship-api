#!/bin/bash

# Test Interview Booking APIs
echo "======================================"
echo "Testing Interview Booking APIs"
echo "======================================"
echo ""

# First, login to get a valid token (using admin credentials)
echo "1. Logging in as admin..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@mahidol.ac.th",
    "password": "admin123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "❌ Login failed. Response:"
  echo $LOGIN_RESPONSE
  exit 1
fi

echo "✅ Login successful"
echo "Token: ${TOKEN:0:50}..."
echo ""

# Test 1: Get Statistics
echo "2. Testing GET /interview/statistics"
curl -s -X GET "http://localhost:8080/api/v1/interview/statistics" \
  -H "Authorization: Bearer $TOKEN" | head -20
echo -e "\n"

# Test 2: Get All Bookings
echo "3. Testing GET /interview/bookings"
curl -s -X GET "http://localhost:8080/api/v1/interview/bookings?page=1&limit=5" \
  -H "Authorization: Bearer $TOKEN" | head -30
echo -e "\n"

# Test 3: Get Booking by ID (if exists)
echo "4. Testing GET /interview/bookings/1"
curl -s -X GET "http://localhost:8080/api/v1/interview/bookings/1" \
  -H "Authorization: Bearer $TOKEN"
echo -e "\n"

# Test 4: Try to create a slot first (admin only)
echo "5. Testing POST /interview/slots (create slot)"
SLOT_RESPONSE=$(curl -s -X POST "http://localhost:8080/api/v1/interview/slots" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "scholarship_id": 1,
    "interviewer_id": "dfa4ee05-9b7b-4143-9155-58ba1c8494d2",
    "interview_date": "2025-11-15",
    "start_time": "09:00",
    "end_time": "10:00",
    "location": "Building A, Room 101",
    "building": "Building A",
    "floor": "1",
    "room": "101",
    "max_capacity": 5,
    "slot_type": "individual",
    "duration_minutes": 30
  }')

echo $SLOT_RESPONSE | head -30
echo -e "\n"

# Test 5: Get available slots
echo "6. Testing GET /interview/availability"
curl -s -X GET "http://localhost:8080/api/v1/interview/availability?scholarship_id=1" \
  -H "Authorization: Bearer $TOKEN" | head -40
echo -e "\n"

echo "======================================"
echo "✅ All tests completed!"
echo "======================================"
