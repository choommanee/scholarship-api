#!/bin/bash

# Test the API fixes

API_URL="http://localhost:8080/api/v1"

# Login
echo "=== Logging in ==="
LOGIN_RESP=$(curl -s -X POST "$API_URL/auth/login" -H "Content-Type: application/json" -d '{"email":"teststudent@example.com","password":"password123"}')
TOKEN=$(echo $LOGIN_RESP | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "Login failed!"
  exit 1
fi

echo "Login successful!"
echo ""

# Test Problem 1: GetMyApplications
echo "=== Test 1: GET /applications/my (Problem 1) ==="
RESP=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X GET "$API_URL/applications/my" -H "Authorization: Bearer $TOKEN")
HTTP_CODE=$(echo "$RESP" | grep "HTTP_CODE:" | cut -d':' -f2)
BODY=$(echo "$RESP" | sed '/HTTP_CODE:/d')

echo "HTTP Status: $HTTP_CODE"
if [ "$HTTP_CODE" = "200" ]; then
  echo "✓ SUCCESS - GetMyApplications is working!"
else
  echo "✗ FAILED"
fi
echo "Response: $BODY"
echo ""

# Create a test scholarship first
echo "=== Creating test scholarship ==="
PGPASSWORD=postgres psql -h localhost -p 5434 -U postgres -d scholarship_db -c "
INSERT INTO scholarships (scholarship_name, scholarship_type, amount, academic_year, application_start_date, application_end_date, total_quota, available_quota, status)
VALUES ('Test Scholarship', 'merit', 10000, '2025', '2025-01-01', '2025-12-31', 10, 10, 'open')
ON CONFLICT DO NOTHING;
" > /dev/null 2>&1

# Create an application
echo "=== Creating test application ==="
CREATE_RESP=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST "$API_URL/applications" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"scholarship_id":1,"family_income":20000,"monthly_expenses":5000,"siblings_count":2}')
HTTP_CODE=$(echo "$CREATE_RESP" | grep "HTTP_CODE:" | cut -d':' -f2)
BODY=$(echo "$CREATE_RESP" | sed '/HTTP_CODE:/d')

echo "HTTP Status: $HTTP_CODE"
APP_ID=$(echo "$BODY" | grep -o '"application_id":[0-9]*' | cut -d':' -f2)
if [ -z "$APP_ID" ]; then
  APP_ID=$(echo "$BODY" | grep -o '"ApplicationID":[0-9]*' | cut -d':' -f2)
fi

if [ -z "$APP_ID" ]; then
  # Try to get from existing applications
  echo "Getting existing application ID..."
  GET_RESP=$(curl -s "$API_URL/applications/my" -H "Authorization: Bearer $TOKEN")
  APP_ID=$(echo "$GET_RESP" | grep -o '"id":"[0-9]*"' | head -1 | cut -d'"' -f4)
fi

if [ -z "$APP_ID" ]; then
  echo "Could not create or find application. Using ID=1"
  APP_ID=1
else
  echo "Application ID: $APP_ID"
fi
echo ""

# Test Problem 2: Application Details endpoints
echo "=== Test 2: POST /applications/$APP_ID/personal-info (Problem 2) ==="
RESP=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST "$API_URL/applications/$APP_ID/personal-info" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"first_name":"John","last_name":"Doe","citizen_id":"1234567890123","date_of_birth":"2000-01-01","gender":"male","nationality":"Thai","phone_number":"0812345678","email":"teststudent@example.com"}')
HTTP_CODE=$(echo "$RESP" | grep "HTTP_CODE:" | cut -d':' -f2)
BODY=$(echo "$RESP" | sed '/HTTP_CODE:/d')

echo "HTTP Status: $HTTP_CODE"
if [ "$HTTP_CODE" = "200" ]; then
  echo "✓ SUCCESS - SavePersonalInfo is working!"
elif [ "$HTTP_CODE" = "403" ]; then
  echo "✗ FAILED - Still getting 403 Forbidden (Problem 2 NOT fixed)"
else
  echo "Status: $HTTP_CODE"
fi
echo "Response: $BODY"
echo ""

# Test another endpoint
echo "=== Test 3: POST /applications/$APP_ID/addresses ==="
RESP=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST "$API_URL/applications/$APP_ID/addresses" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '[{"address_type":"current","house_number":"123","village":"Village 1","subdistrict":"Subdistrict","district":"District","province":"Bangkok","postal_code":"10100"}]')
HTTP_CODE=$(echo "$RESP" | grep "HTTP_CODE:" | cut -d':' -f2)
BODY=$(echo "$RESP" | sed '/HTTP_CODE:/d')

echo "HTTP Status: $HTTP_CODE"
if [ "$HTTP_CODE" = "200" ]; then
  echo "✓ SUCCESS - SaveAddresses is working!"
elif [ "$HTTP_CODE" = "403" ]; then
  echo "✗ FAILED - Still getting 403 Forbidden"
else
  echo "Status: $HTTP_CODE"
fi
echo "Response: $BODY"
echo ""

echo "=== SUMMARY ==="
echo "Check the results above to see if both problems are fixed."
