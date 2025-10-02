#!/bin/bash

echo "==========================================================="
echo "🧪 Application Details API - Complete Test"
echo "==========================================================="
echo ""

# Use existing student user
STUDENT_EMAIL="student1@university.ac.th"
STUDENT_PASS="password123"
STUDENT_USER_ID="f63b84a0-051c-4651-b9e7-621e28c2ef05"
STUDENT_ID="6388001"
SCHOLARSHIP_ID=4

echo "🔑 Step 1: Login as student"
LOGIN=$(curl -s http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"${STUDENT_EMAIL}\",\"password\":\"${STUDENT_PASS}\"}")

TOKEN=$(echo "$LOGIN" | jq -r '.token')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
    echo "❌ Login failed"
    echo "$LOGIN" | jq '.'
    exit 1
fi

echo "✅ Login successful"
echo "Token: ${TOKEN:0:50}..."
echo ""

# Ensure student record exists
echo "📝 Step 2: Ensure student record exists"
psql scholarship_db -c "INSERT INTO students (student_id, user_id, faculty_code, year_level) VALUES ('${STUDENT_ID}', '${STUDENT_USER_ID}', 'ECON', 2) ON CONFLICT (student_id) DO UPDATE SET user_id = EXCLUDED.user_id;" > /dev/null 2>&1
echo "✅ Student record ready"
echo ""

# Create application or use existing
echo "📄 Step 3: Get/Create application"
MY_APPS=$(curl -s http://localhost:8080/api/v1/applications/my \
  -H "Authorization: Bearer ${TOKEN}")
APP_ID=$(echo "$MY_APPS" | jq -r '.applications[0].application_id // empty')

if [ -z "$APP_ID" ] || [ "$APP_ID" = "null" ]; then
    APP_ID=$(psql scholarship_db -t -c "INSERT INTO scholarship_applications (scholarship_id, student_id, application_status) VALUES (${SCHOLARSHIP_ID}, '${STUDENT_ID}', 'draft') RETURNING application_id;" 2>&1 | tr -d ' \n')
fi

echo "✅ Application ID: ${APP_ID}"
echo ""

# Test 1: Save Personal Info
echo "========================================================="
echo "Test 1: POST /applications/${APP_ID}/personal-info"
echo "========================================================="
RESULT=$(curl -s -X POST "http://localhost:8080/api/v1/applications/${APP_ID}/personal-info" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "prefix_th": "นาย",
    "first_name_th": "สมชาย",
    "last_name_th": "ใจดี",
    "email": "somchai@test.com",
    "phone": "0812345678",
    "student_id": "6388001",
    "faculty": "คณะเศรษฐศาสตร์",
    "year_level": 2
  }')
echo "$RESULT" | jq '.'
SUCCESS=$(echo "$RESULT" | jq -r '.success // false')
[ "$SUCCESS" = "true" ] && echo "✅ PASS" || echo "❌ FAIL"
echo ""

# Test 2: Save Addresses
echo "========================================================="
echo "Test 2: POST /applications/${APP_ID}/addresses"
echo "========================================================="
RESULT=$(curl -s -X POST "http://localhost:8080/api/v1/applications/${APP_ID}/addresses" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "addresses": [
      {
        "address_type": "registered",
        "province": "กรุงเทพมหานคร",
        "district": "พญาไท",
        "postal_code": "10400"
      }
    ]
  }')
echo "$RESULT" | jq '.'
SUCCESS=$(echo "$RESULT" | jq -r '.success // false')
[ "$SUCCESS" = "true" ] && echo "✅ PASS" || echo "❌ FAIL"
echo ""

# Test 3: Get Complete Form
echo "========================================================="
echo "Test 3: GET /applications/${APP_ID}/complete-form"
echo "========================================================="
RESULT=$(curl -s -X GET "http://localhost:8080/api/v1/applications/${APP_ID}/complete-form" \
  -H "Authorization: Bearer ${TOKEN}")
echo "$RESULT" | jq '.'
SUCCESS=$(echo "$RESULT" | jq -r '.success // false')
[ "$SUCCESS" = "true" ] && echo "✅ PASS" || echo "❌ FAIL"
echo ""

echo "==========================================================="
echo "✅ Testing Complete!"
echo "==========================================================="
