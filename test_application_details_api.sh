#!/bin/bash

# Application Details API Testing Script
# ทดสอบ 9 endpoints สำหรับ Application Form

BASE_URL="http://localhost:8080/api/v1"
# Get fresh token via login
echo "🔑 Getting authentication token..."
LOGIN_RESPONSE=$(curl -s http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}')
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
    echo "❌ Failed to get authentication token"
    exit 1
fi

echo "✅ Token obtained successfully"
echo ""

echo "=================================================="
echo "🧪 Testing Application Details API Endpoints"
echo "=================================================="
echo ""

# Get student's applications first
echo "📋 Step 0: Get Student Applications"
echo "--------------------------------------------------"
APPLICATIONS=$(curl -s -X GET "${BASE_URL}/applications/my" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json")
echo "$APPLICATIONS" | jq '.'

# Extract first application ID
APP_ID=$(echo "$APPLICATIONS" | jq -r '.data[0].application_id // empty')

if [ -z "$APP_ID" ]; then
    echo "⚠️  No applications found. Creating a new application first..."

    # Get available scholarships
    SCHOLARSHIPS=$(curl -s -X GET "${BASE_URL}/scholarships" \
      -H "Authorization: Bearer ${TOKEN}")
    SCHOLARSHIP_ID=$(echo "$SCHOLARSHIPS" | jq -r '.data[0].scholarship_id // 1')

    # Create new application
    CREATE_RESPONSE=$(curl -s -X POST "${BASE_URL}/applications" \
      -H "Authorization: Bearer ${TOKEN}" \
      -H "Content-Type: application/json" \
      -d "{
        \"scholarship_id\": ${SCHOLARSHIP_ID},
        \"application_reason\": \"ทดสอบระบบ Application Details API\"
      }")

    APP_ID=$(echo "$CREATE_RESPONSE" | jq -r '.data.application_id // empty')
    echo "✅ Created new application: ID = ${APP_ID}"
fi

echo ""
echo "🎯 Using Application ID: ${APP_ID}"
echo ""
sleep 1

# Test 1: Save Personal Info
echo "=================================================="
echo "Test 1: POST /applications/${APP_ID}/personal-info"
echo "=================================================="
RESPONSE=$(curl -s -X POST "${BASE_URL}/applications/${APP_ID}/personal-info" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "prefix_th": "นาย",
    "first_name_th": "สมชาย",
    "last_name_th": "ใจดี",
    "prefix_en": "Mr.",
    "first_name_en": "Somchai",
    "last_name_en": "Jaidee",
    "email": "somchai.j@student.mahidol.ac.th",
    "phone": "0812345678",
    "student_id": "6388001",
    "faculty": "คณะเศรษฐศาสตร์",
    "department": "เศรษฐศาสตร์",
    "program": "เศรษฐศาสตรบัณฑิต",
    "year_level": 2,
    "admission_type": "TCAS รอบ 1"
  }')
echo "$RESPONSE" | jq '.'
echo ""
sleep 1

# Test 2: Save Addresses
echo "=================================================="
echo "Test 2: POST /applications/${APP_ID}/addresses"
echo "=================================================="
RESPONSE=$(curl -s -X POST "${BASE_URL}/applications/${APP_ID}/addresses" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "addresses": [
      {
        "address_type": "registered",
        "address_line1": "123 ถนนพญาไท",
        "address_line2": "แขวงพญาไท",
        "sub_district": "พญาไท",
        "district": "พญาไท",
        "province": "กรุงเทพมหานคร",
        "postal_code": "10400",
        "country": "ไทย"
      },
      {
        "address_type": "current",
        "address_line1": "456 ถนนรังสิต-นครนายก",
        "sub_district": "คลองหนึ่ง",
        "district": "คลองหลวง",
        "province": "ปทุมธานี",
        "postal_code": "12120",
        "country": "ไทย"
      }
    ]
  }')
echo "$RESPONSE" | jq '.'
echo ""
sleep 1

# Test 3: Save Education
echo "=================================================="
echo "Test 3: POST /applications/${APP_ID}/education"
echo "=================================================="
RESPONSE=$(curl -s -X POST "${BASE_URL}/applications/${APP_ID}/education" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "education_history": [
      {
        "education_level": "มัธยมศึกษาตอนปลาย",
        "institution_name": "โรงเรียนเตรียมอุดมศึกษา",
        "major": "วิทย์-คณิต",
        "gpa": 3.85,
        "graduation_year": 2022,
        "province": "กรุงเทพมหานคร"
      }
    ]
  }')
echo "$RESPONSE" | jq '.'
echo ""
sleep 1

# Test 4: Save Family
echo "=================================================="
echo "Test 4: POST /applications/${APP_ID}/family"
echo "=================================================="
RESPONSE=$(curl -s -X POST "${BASE_URL}/applications/${APP_ID}/family" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "family_members": [
      {
        "relationship": "บิดา",
        "prefix": "นาย",
        "first_name": "สมศักดิ์",
        "last_name": "ใจดี",
        "age": 55,
        "occupation": "ข้าราชการ",
        "monthly_income": 35000,
        "status": "มีชีวิต",
        "phone": "0891234567"
      },
      {
        "relationship": "มารดา",
        "prefix": "นาง",
        "first_name": "สมหญิง",
        "last_name": "ใจดี",
        "age": 52,
        "occupation": "แม่บ้าน",
        "monthly_income": 0,
        "status": "มีชีวิต",
        "phone": "0891234568"
      }
    ],
    "guardians": [
      {
        "prefix": "นาย",
        "first_name": "สมศักดิ์",
        "last_name": "ใจดี",
        "relationship": "บิดา",
        "occupation": "ข้าราชการ",
        "phone": "0891234567"
      }
    ],
    "siblings": [
      {
        "prefix": "นางสาว",
        "first_name": "สมหญิง",
        "last_name": "ใจดี",
        "age": 18,
        "education_level": "มัธยมศึกษาตอนปลาย",
        "occupation": "นักเรียน",
        "status": "กำลังศึกษา"
      }
    ],
    "living_situation": {
      "house_type": "บ้านเดี่ยว",
      "ownership": "เป็นเจ้าของ",
      "people_count": 4,
      "has_electricity": true,
      "has_water": true,
      "area_type": "ชุมชนเมือง"
    }
  }')
echo "$RESPONSE" | jq '.'
echo ""
sleep 1

# Test 5: Save Financial
echo "=================================================="
echo "Test 5: POST /applications/${APP_ID}/financial"
echo "=================================================="
RESPONSE=$(curl -s -X POST "${BASE_URL}/applications/${APP_ID}/financial" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "financial_info": {
      "total_family_income": 35000,
      "monthly_expenses": 25000,
      "debt_amount": 500000,
      "num_dependents": 3,
      "has_other_income": false
    },
    "assets": [
      {
        "asset_type": "บ้าน",
        "description": "บ้านเดี่ยว 2 ชั้น",
        "estimated_value": 3000000
      },
      {
        "asset_type": "รถยนต์",
        "description": "รถยนต์ 1 คัน",
        "estimated_value": 500000
      }
    ],
    "scholarship_history": [
      {
        "scholarship_name": "ทุนเรียนดี ปีการศึกษา 2565",
        "amount": 20000,
        "year_received": 2022,
        "source": "มหาวิทยาลัย"
      }
    ],
    "health_info": {
      "has_chronic_disease": false,
      "has_disability": false
    },
    "funding_needs": {
      "tuition_fee": 45000,
      "living_expenses": 5000,
      "books_supplies": 3000,
      "other_expenses": 2000,
      "total_needed": 55000,
      "justification": "ต้องการทุนเพื่อลดภาระค่าใช้จ่ายของครอบครัว"
    }
  }')
echo "$RESPONSE" | jq '.'
echo ""
sleep 1

# Test 6: Save Activities
echo "=================================================="
echo "Test 6: POST /applications/${APP_ID}/activities"
echo "=================================================="
RESPONSE=$(curl -s -X POST "${BASE_URL}/applications/${APP_ID}/activities" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "activities": [
      {
        "activity_type": "กิจกรรมชุมชน",
        "activity_name": "จิตอาสาพัฒนาชุมชน",
        "role": "หัวหน้าทีม",
        "year": 2023,
        "hours": 40,
        "description": "จัดกิจกรรมพัฒนาชุมชนบริเวณใกล้มหาวิทยาลัย"
      },
      {
        "activity_type": "กิจกรรมวิชาการ",
        "activity_name": "แข่งขันตอบปัญหาเศรษฐศาสตร์",
        "role": "ผู้เข้าร่วม",
        "year": 2023,
        "achievement": "รางวัลชนะเลิศ ระดับมหาวิทยาลัย"
      }
    ],
    "references": [
      {
        "prefix": "อาจารย์ ดร.",
        "first_name": "สมปอง",
        "last_name": "วิชาการ",
        "position": "อาจารย์",
        "organization": "คณะเศรษฐศาสตร์ มหาวิทยาลัยมหิดล",
        "phone": "021234567",
        "email": "sompong.w@mahidol.ac.th",
        "relationship": "อาจารย์ที่ปรึกษา"
      }
    ]
  }')
echo "$RESPONSE" | jq '.'
echo ""
sleep 1

# Test 7: Get Complete Form
echo "=================================================="
echo "Test 7: GET /applications/${APP_ID}/complete-form"
echo "=================================================="
RESPONSE=$(curl -s -X GET "${BASE_URL}/applications/${APP_ID}/complete-form" \
  -H "Authorization: Bearer ${TOKEN}")
echo "$RESPONSE" | jq '.'
echo ""
sleep 1

# Test 8: Save Complete Form (all at once)
echo "=================================================="
echo "Test 8: POST /applications/${APP_ID}/complete-form"
echo "=================================================="
RESPONSE=$(curl -s -X POST "${BASE_URL}/applications/${APP_ID}/complete-form" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "personal_info": {
      "prefix_th": "นาย",
      "first_name_th": "สมชาย",
      "last_name_th": "ใจดี (อัพเดท)",
      "email": "somchai.updated@student.mahidol.ac.th",
      "phone": "0812345678",
      "student_id": "6388001",
      "faculty": "คณะเศรษฐศาสตร์",
      "year_level": 2
    },
    "addresses": [
      {
        "address_type": "registered",
        "province": "กรุงเทพมหานคร",
        "district": "พญาไท"
      }
    ],
    "family_members": [
      {
        "relationship": "บิดา",
        "first_name": "สมศักดิ์",
        "last_name": "ใจดี",
        "occupation": "ข้าราชการ",
        "monthly_income": 35000
      }
    ]
  }')
echo "$RESPONSE" | jq '.'
echo ""
sleep 1

# Test 9: Submit Application
echo "=================================================="
echo "Test 9: PUT /applications/${APP_ID}/submit"
echo "=================================================="
RESPONSE=$(curl -s -X PUT "${BASE_URL}/applications/${APP_ID}/submit" \
  -H "Authorization: Bearer ${TOKEN}")
echo "$RESPONSE" | jq '.'
echo ""

echo "=================================================="
echo "✅ All Tests Completed!"
echo "=================================================="
echo ""
echo "📊 Summary:"
echo "  - Total endpoints tested: 9"
echo "  - Application ID used: ${APP_ID}"
echo "  - Base URL: ${BASE_URL}"
echo ""
echo "💡 Next steps:"
echo "  1. Review responses above"
echo "  2. Check database for saved data"
echo "  3. Verify all fields are correctly stored"
echo ""
