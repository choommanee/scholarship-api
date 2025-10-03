# Postman Collection Guide - Scholarship Application System

## Quick Start

This guide provides ready-to-use Postman/HTTPie/cURL requests for testing all new endpoints.

---

## Environment Variables

Set these in Postman environment:

```
BASE_URL = http://localhost:8080/api/v1
TOKEN = (your JWT token after login)
APPLICATION_ID = (created application ID)
DOCUMENT_ID = (uploaded document ID)
SCHOLARSHIP_ID = 1
```

---

## 1. Authentication (Prerequisite)

### Login
```
POST {{BASE_URL}}/auth/login
Content-Type: application/json

{
  "email": "student@example.com",
  "password": "password123"
}
```

**Save the token from response!**

---

## 2. Check Eligibility

### Check Scholarship Eligibility
```
POST {{BASE_URL}}/scholarships/{{SCHOLARSHIP_ID}}/check-eligibility
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

{
  "student_data": {
    "gpa": 3.5,
    "year_level": 3,
    "faculty": "Engineering",
    "department": "Computer Engineering",
    "family_income": 25000
  }
}
```

**Tests:**
- Status code is 200
- Response has `success: true`
- Response has `is_eligible` boolean
- Response has `eligibility_score` number

---

## 3. Draft Management

### Create Draft Application
```
POST {{BASE_URL}}/applications/draft
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

{
  "scholarship_id": {{SCHOLARSHIP_ID}}
}
```

**Tests:**
- Status code is 201
- Response has `application_id`
- Application status is "draft"

**Save application_id!**

---

### Get Draft Application
```
GET {{BASE_URL}}/applications/draft?scholarship_id={{SCHOLARSHIP_ID}}
Authorization: Bearer {{TOKEN}}
```

**Tests:**
- Status code is 200
- Response has `current_step`
- Response has `draft_data`

---

## 4. Section Saving

### Save Personal Info
```
POST {{BASE_URL}}/applications/{{APPLICATION_ID}}/sections/personal_info
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

{
  "prefix_th": "นาย",
  "prefix_en": "Mr.",
  "first_name_th": "สมชาย",
  "last_name_th": "ใจดี",
  "first_name_en": "Somchai",
  "last_name_en": "Jaidee",
  "email": "somchai@example.com",
  "phone": "0812345678",
  "line_id": "somchai123",
  "citizen_id": "1234567890123",
  "student_id": "6412345678",
  "faculty": "Engineering",
  "department": "Computer Engineering",
  "major": "Software Engineering",
  "year_level": 3,
  "admission_type": "Direct Admission"
}
```

---

### Save Address Info
```
POST {{BASE_URL}}/applications/{{APPLICATION_ID}}/sections/address_info
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

[
  {
    "address_type": "current",
    "house_number": "123",
    "village_number": "5",
    "alley": "ซอย 10",
    "road": "ถนนพระราม 4",
    "subdistrict": "คลองเตย",
    "district": "คลองเตย",
    "province": "กรุงเทพมหานคร",
    "postal_code": "10110"
  },
  {
    "address_type": "permanent",
    "house_number": "456",
    "district": "เมือง",
    "province": "เชียงใหม่",
    "postal_code": "50000"
  }
]
```

---

### Save Education History
```
POST {{BASE_URL}}/applications/{{APPLICATION_ID}}/sections/education_history
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

[
  {
    "education_level": "high_school",
    "school_name": "โรงเรียนสามเสนวิทยาลัย",
    "school_province": "กรุงเทพมหานคร",
    "gpa": 3.85,
    "graduation_year": "2564"
  },
  {
    "education_level": "undergraduate",
    "school_name": "มหาวิทยาลัยเทคโนโลยีพระจอมเกล้าธนบุรี",
    "school_province": "กรุงเทพมหานคร",
    "gpa": 3.45
  }
]
```

---

### Save Family Info
```
POST {{BASE_URL}}/applications/{{APPLICATION_ID}}/sections/family_info
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

{
  "members": [
    {
      "relationship": "father",
      "title": "นาย",
      "first_name": "สมศักดิ์",
      "last_name": "ใจดี",
      "age": 55,
      "living_status": "alive",
      "occupation": "พนักงานบริษัทเอกชน",
      "position": "หัวหน้าแผนก",
      "workplace": "บริษัท ABC จำกัด",
      "workplace_province": "กรุงเทพมหานคร",
      "monthly_income": 35000,
      "phone": "0898765432"
    },
    {
      "relationship": "mother",
      "title": "นาง",
      "first_name": "สมหญิง",
      "last_name": "ใจดี",
      "age": 52,
      "living_status": "alive",
      "occupation": "แม่บ้าน",
      "monthly_income": 0
    }
  ],
  "guardians": [],
  "siblings": [
    {
      "sibling_order": 1,
      "gender": "male",
      "school_or_workplace": "โรงเรียนมัธยม",
      "education_level": "high_school",
      "is_working": false,
      "monthly_income": 0
    }
  ],
  "living_situation": {
    "living_with": "parents",
    "living_details": "อาศัยอยู่กับบิดามารดา"
  }
}
```

---

### Save Financial Info
```
POST {{BASE_URL}}/applications/{{APPLICATION_ID}}/sections/financial_info
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

{
  "financial_info": {
    "monthly_allowance": 5000,
    "daily_travel_cost": 50,
    "monthly_dorm_cost": 0,
    "other_monthly_costs": 2000,
    "has_income": false,
    "income_source": null,
    "monthly_income": 0,
    "financial_notes": "ได้รับค่าใช้จ่ายจากบิดามารดา"
  },
  "assets": [
    {
      "asset_type": "property",
      "category": "house",
      "description": "บ้านที่อยู่อาศัย",
      "value": 2000000
    },
    {
      "asset_type": "vehicle",
      "category": "motorcycle",
      "description": "รถมอเตอร์ไซค์",
      "value": 30000
    }
  ],
  "scholarship_history": [],
  "health_info": {
    "has_health_issues": false,
    "affects_study": false
  },
  "funding_needs": {
    "tuition_support": 20000,
    "monthly_support": 5000,
    "book_support": 3000,
    "dorm_support": 0,
    "other_support": 2000,
    "total_requested": 30000,
    "necessity_reason": "ครอบครัวมีรายได้จำกัด ต้องการทุนเพื่อลดภาระค่าใช้จ่าย"
  }
}
```

---

### Save Activities & Skills
```
POST {{BASE_URL}}/applications/{{APPLICATION_ID}}/sections/activities_skills
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

{
  "activities": [
    {
      "activity_type": "sport",
      "activity_name": "ฟุตบอล",
      "description": "นักกีฬาฟุตบอลทีมมหาวิทยาลัย",
      "achievement": "รางวัลชนะเลิศการแข่งขันระดับภาค",
      "award_level": "regional",
      "year": "2567"
    },
    {
      "activity_type": "academic",
      "activity_name": "โครงงานวิทยาศาสตร์",
      "description": "โครงงานพัฒนาแอพพลิเคชั่น",
      "achievement": "รางวัลรองชนะเลิศอันดับ 1",
      "award_level": "national",
      "year": "2567"
    }
  ],
  "references": [
    {
      "title": "ผศ.",
      "first_name": "สมพร",
      "last_name": "ใจดี",
      "relationship": "อาจารย์ที่ปรึกษา",
      "address": "มหาวิทยาลัยเทคโนโลยีพระจอมเกล้าธนบุรี",
      "phone": "021234567",
      "email": "somporn@kmutt.ac.th"
    }
  ]
}
```

---

## 5. Document Management

### Upload ID Card
```
POST {{BASE_URL}}/documents/applications/{{APPLICATION_ID}}/upload-enhanced
Authorization: Bearer {{TOKEN}}
Content-Type: multipart/form-data

form-data:
  document_type: id_card
  file: [select PDF file]
```

---

### Upload Transcript
```
POST {{BASE_URL}}/documents/applications/{{APPLICATION_ID}}/upload-enhanced
Authorization: Bearer {{TOKEN}}
Content-Type: multipart/form-data

form-data:
  document_type: transcript
  file: [select PDF file]
```

**Save document_id from response!**

---

### Download Document
```
GET {{BASE_URL}}/documents/{{DOCUMENT_ID}}/download-enhanced
Authorization: Bearer {{TOKEN}}
```

**Tests:**
- Status code is 200
- Content-Type matches uploaded file
- File downloads successfully

---

### Delete Document
```
DELETE {{BASE_URL}}/documents/{{DOCUMENT_ID}}/delete-enhanced
Authorization: Bearer {{TOKEN}}
```

**Tests:**
- Status code is 200
- Response has `success: true`
- Subsequent download returns 404

---

## 6. Submit Application

### Submit Application
```
POST {{BASE_URL}}/applications/{{APPLICATION_ID}}/submit-enhanced
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

{
  "terms_accepted": true,
  "declaration_accepted": true
}
```

**Tests:**
- Status code is 200
- Response has `reference_number`
- Application status is "submitted"
- Subsequent edits return error

---

## 7. Error Test Cases

### Test 1: Submit without required sections
```
POST {{BASE_URL}}/applications/{{APPLICATION_ID}}/submit-enhanced
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

{
  "terms_accepted": true,
  "declaration_accepted": true
}
```

**Expected:** 400 Bad Request with validation errors

---

### Test 2: Submit without terms
```
POST {{BASE_URL}}/applications/{{APPLICATION_ID}}/submit-enhanced
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

{
  "terms_accepted": false,
  "declaration_accepted": false
}
```

**Expected:** 400 Bad Request

---

### Test 3: Delete document after submission
```
DELETE {{BASE_URL}}/documents/{{DOCUMENT_ID}}/delete-enhanced
Authorization: Bearer {{TOKEN}}
```

**Expected:** 400 Bad Request (not in draft status)

---

### Test 4: Upload oversized file
Upload file > 10MB (or > 5MB for id_card)

**Expected:** 400 Bad Request (file size limit)

---

### Test 5: Upload invalid file type
Upload .exe or .zip file

**Expected:** 400 Bad Request (invalid file type)

---

### Test 6: Access another student's application
Use different student token, try to access first student's application

**Expected:** 403 Forbidden

---

## 8. Complete Workflow Test Sequence

Execute in this order:

1. ✅ Login
2. ✅ Check Eligibility
3. ✅ Create Draft
4. ✅ Save Personal Info
5. ✅ Save Address Info
6. ✅ Save Education History
7. ✅ Save Family Info
8. ✅ Save Financial Info
9. ✅ Save Activities & Skills
10. ✅ Upload ID Card
11. ✅ Upload Transcript
12. ✅ Get Draft (verify all data)
13. ✅ Submit Application
14. ✅ Verify submission (status, reference number)

---

## 9. Postman Collection JSON

Save this as a Postman collection:

```json
{
  "info": {
    "name": "Scholarship Application System - New Endpoints",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "variable": [
    {
      "key": "BASE_URL",
      "value": "http://localhost:8080/api/v1"
    },
    {
      "key": "TOKEN",
      "value": ""
    },
    {
      "key": "APPLICATION_ID",
      "value": ""
    },
    {
      "key": "SCHOLARSHIP_ID",
      "value": "1"
    },
    {
      "key": "DOCUMENT_ID",
      "value": ""
    }
  ],
  "item": [
    {
      "name": "1. Authentication",
      "item": [
        {
          "name": "Login",
          "request": {
            "method": "POST",
            "header": [
              {
                "key": "Content-Type",
                "value": "application/json"
              }
            ],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"email\": \"student@example.com\",\n  \"password\": \"password123\"\n}"
            },
            "url": {
              "raw": "{{BASE_URL}}/auth/login",
              "host": ["{{BASE_URL}}"],
              "path": ["auth", "login"]
            }
          }
        }
      ]
    },
    {
      "name": "2. Eligibility",
      "item": [
        {
          "name": "Check Eligibility",
          "request": {
            "method": "POST",
            "header": [
              {
                "key": "Authorization",
                "value": "Bearer {{TOKEN}}"
              },
              {
                "key": "Content-Type",
                "value": "application/json"
              }
            ],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"student_data\": {\n    \"gpa\": 3.5,\n    \"year_level\": 3,\n    \"faculty\": \"Engineering\",\n    \"family_income\": 25000\n  }\n}"
            },
            "url": {
              "raw": "{{BASE_URL}}/scholarships/{{SCHOLARSHIP_ID}}/check-eligibility",
              "host": ["{{BASE_URL}}"],
              "path": ["scholarships", "{{SCHOLARSHIP_ID}}", "check-eligibility"]
            }
          }
        }
      ]
    }
  ]
}
```

---

## 10. HTTPie Examples

If you prefer HTTPie:

```bash
# Login
http POST localhost:8080/api/v1/auth/login \
  email=student@example.com password=password123

# Check Eligibility
http POST localhost:8080/api/v1/scholarships/1/check-eligibility \
  Authorization:"Bearer $TOKEN" \
  student_data:='{"gpa": 3.5, "year_level": 3, "faculty": "Engineering"}'

# Create Draft
http POST localhost:8080/api/v1/applications/draft \
  Authorization:"Bearer $TOKEN" \
  scholarship_id:=1

# Upload Document
http --form POST localhost:8080/api/v1/documents/applications/1/upload-enhanced \
  Authorization:"Bearer $TOKEN" \
  document_type=id_card \
  file@id_card.pdf
```

---

## Support

For issues or questions, refer to:
- API_ENDPOINTS_GUIDE.md - Complete API documentation
- IMPLEMENTATION_REPORT.md - Technical details
