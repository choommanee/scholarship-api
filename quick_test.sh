#!/bin/bash

# Use existing valid token from previous tests
# This is a student token, we'll test what we can
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZGZhNGVlMDUtOWI3Yi00MTQzLTkxNTUtNThiYTFjODQ5NGQyIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwidXNlcm5hbWUiOiIiLCJyb2xlcyI6WyJzdHVkZW50Il0sInN1YiI6ImRmYTRlZTA1LTliN2ItNDE0My05MTU1LTU4YmExYzg0OTRkMiIsImV4cCI6MTc1OTkzODMzNywiaWF0IjoxNzU5MzMzNTM3fQ.vlZGL3kZNUo45v_-H2g9pHaWrtTlEFN4Yc0HzG5cbYM"

echo "Testing Interview Booking APIs"
echo "================================"

# Test GET /interview/bookings (should fail - no permission)
echo -e "\n1. GET /interview/bookings (expect 403 - student role)"
curl -s "http://localhost:8080/api/v1/interview/bookings" -H "Authorization: Bearer $TOKEN"

# Test GET /interview/bookings/:id (should work if JWT valid)
echo -e "\n\n2. GET /interview/bookings/1 (expect success or 404)"
curl -s "http://localhost:8080/api/v1/interview/bookings/1" -H "Authorization: Bearer $TOKEN"

# Test GET /interview/statistics (should fail - no permission)
echo -e "\n\n3. GET /interview/statistics (expect 403 - student role)"
curl -s "http://localhost:8080/api/v1/interview/statistics" -H "Authorization: Bearer $TOKEN"

# Test GET /interview/availability (should work)
echo -e "\n\n4. GET /interview/availability?scholarship_id=1 (expect success)"
curl -s "http://localhost:8080/api/v1/interview/availability?scholarship_id=1" -H "Authorization: Bearer $TOKEN"

echo -e "\n\n================================"
echo "✅ Tests completed!"
