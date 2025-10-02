#!/bin/bash

echo "==========================================================="
echo "🧪 Application Details API - Working Test"
echo "==========================================================="
echo ""

STUDENT_EMAIL="student1@university.ac.th"
STUDENT_PASS="password123"
STUDENT_ID="6411111111"  # Existing student ID
SCHOLARSHIP_ID=4

echo "🔑 Login as student"
LOGIN=$(curl -s http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"${STUDENT_EMAIL}\",\"password\":\"${STUDENT_PASS}\"}")
TOKEN=$(echo "$LOGIN" | jq -r '.token')
echo "✅ Token: ${TOKEN:0:50}..."
echo ""

echo "📄 Get/Create application"
APP_ID=$(psql scholarship_db -t -c "INSERT INTO scholarship_applications (scholarship_id, student_id, application_status) VALUES (${SCHOLARSHIP_ID}, '${STUDENT_ID}', 'draft') ON CONFLICT DO NOTHING RETURNING application_id;" 2>&1 | tr -d ' \n')

if [ -z "$APP_ID" ]; then
    APP_ID=$(psql scholarship_db -t -c "SELECT application_id FROM scholarship_applications WHERE student_id = '${STUDENT_ID}' ORDER BY created_at DESC LIMIT 1;" | tr -d ' \n')
fi

echo "✅ Application ID: ${APP_ID}"
echo ""

echo "========================================="
echo "Test 1: Save Personal Info"
echo "========================================="
curl -s -X POST "http://localhost:8080/api/v1/applications/${APP_ID}/personal-info" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"prefix_th":"นาย","first_name_th":"สมชาย","last_name_th":"ใจดี","email":"somchai@test.com","phone":"0812345678","faculty":"คณะเศรษฐศาสตร์","year_level":2}' | jq '.'
echo ""

echo "========================================="
echo "Test 2: Save Addresses"
echo "========================================="
curl -s -X POST "http://localhost:8080/api/v1/applications/${APP_ID}/addresses" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"addresses":[{"address_type":"registered","province":"กรุงเทพมหานคร","district":"พญาไท"}]}' | jq '.'
echo ""

echo "========================================="
echo "Test 3: Get Complete Form"
echo "========================================="
curl -s -X GET "http://localhost:8080/api/v1/applications/${APP_ID}/complete-form" \
  -H "Authorization: Bearer ${TOKEN}" | jq '.success, .data.personal_info, .data.addresses'
echo ""

echo "✅ All tests completed!"
