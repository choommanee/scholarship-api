# Scholarship Application System - API Endpoints Guide

## Overview
This document describes all the newly implemented backend API endpoints for the scholarship application system.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
All protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

---

## 1. Draft Application Management

### 1.1 Create Draft Application

**Endpoint:** `POST /api/v1/applications/draft`
**Auth Required:** Yes (Student role)
**Description:** Create a new draft application or return existing one

**Request Body:**
```json
{
  "scholarship_id": 1
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Draft created successfully",
  "data": {
    "application_id": 123,
    "scholarship_id": 1,
    "student_id": "6412345678",
    "application_status": "draft",
    "created_at": "2025-01-02T10:00:00Z"
  }
}
```

**Example cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/applications/draft \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"scholarship_id": 1}'
```

---

### 1.2 Get Draft Application

**Endpoint:** `GET /api/v1/applications/draft?scholarship_id=1`
**Auth Required:** Yes (Student role)
**Description:** Retrieve draft application with all saved data

**Query Parameters:**
- `scholarship_id` (required): The scholarship ID

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "application_id": 123,
    "scholarship_id": 1,
    "current_step": 3,
    "draft_data": {
      "personal_info": {...},
      "addresses": [...],
      "education_history": [...]
    },
    "updated_at": "2025-01-02T11:30:00Z"
  }
}
```

**Example cURL:**
```bash
curl -X GET "http://localhost:8080/api/v1/applications/draft?scholarship_id=1" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 2. Section Save API

### 2.1 Save Application Section

**Endpoint:** `POST /api/v1/applications/{id}/sections/{section_name}`
**Auth Required:** Yes (Student role)
**Description:** Save a specific section of the application

**Supported Sections:**
- `personal_info`
- `address_info`
- `education_history`
- `family_info`
- `financial_info`
- `activities_skills`

**Request Body (example for personal_info):**
```json
{
  "prefix_th": "นาย",
  "first_name_th": "สมชาย",
  "last_name_th": "ใจดี",
  "email": "somchai@example.com",
  "phone": "0812345678",
  "citizen_id": "1234567890123",
  "student_id": "6412345678"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Section saved successfully",
  "data": {
    "section": "personal_info",
    "saved_at": "2025-01-02T11:35:00Z"
  }
}
```

**Example cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/applications/123/sections/personal_info \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "prefix_th": "นาย",
    "first_name_th": "สมชาย",
    "last_name_th": "ใจดี",
    "email": "somchai@example.com",
    "phone": "0812345678"
  }'
```

---

## 3. Submit Application

### 3.1 Submit Draft for Review

**Endpoint:** `POST /api/v1/applications/{id}/submit-enhanced`
**Auth Required:** Yes (Student role)
**Description:** Submit the draft application for review

**Request Body:**
```json
{
  "terms_accepted": true,
  "declaration_accepted": true
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Application submitted successfully",
  "data": {
    "application_id": 123,
    "application_status": "submitted",
    "submitted_at": "2025-01-02T12:00:00Z",
    "reference_number": "SCH-2025-000123"
  }
}
```

**Example cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/applications/123/submit-enhanced \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "terms_accepted": true,
    "declaration_accepted": true
  }'
```

---

## 4. Eligibility Check

### 4.1 Check Scholarship Eligibility

**Endpoint:** `POST /api/v1/scholarships/{id}/check-eligibility`
**Auth Required:** Yes
**Description:** Check if student meets scholarship eligibility criteria

**Request Body:**
```json
{
  "student_data": {
    "gpa": 3.25,
    "year_level": 3,
    "faculty": "Engineering",
    "department": "Computer Engineering",
    "family_income": 25000
  }
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "is_eligible": true,
    "eligibility_score": 85,
    "criteria_results": [
      {
        "criteria": "min_gpa",
        "required": 2.5,
        "actual": 3.25,
        "passed": true
      },
      {
        "criteria": "max_family_income",
        "required": 30000,
        "actual": 25000,
        "passed": true
      }
    ],
    "missing_requirements": []
  }
}
```

**Example cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/scholarships/1/check-eligibility \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "student_data": {
      "gpa": 3.25,
      "year_level": 3,
      "faculty": "Engineering",
      "family_income": 25000
    }
  }'
```

---

## 5. Document Management

### 5.1 Upload Document

**Endpoint:** `POST /api/v1/documents/applications/{id}/upload-enhanced`
**Auth Required:** Yes (Student role)
**Description:** Upload a document for an application

**Form Data:**
- `document_type`: string (id_card, transcript, etc.)
- `file`: file (PDF, JPEG, PNG - max 10MB)

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Document uploaded successfully",
  "data": {
    "document_id": 456,
    "document_type": "id_card",
    "file_name": "id_card_somchai.pdf",
    "file_path": "/uploads/applications/123/id_card_somchai.pdf",
    "file_size": 524288,
    "mime_type": "application/pdf",
    "verification_status": "pending",
    "uploaded_at": "2025-01-02T11:40:00Z"
  }
}
```

**Example cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/documents/applications/123/upload-enhanced \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "document_type=id_card" \
  -F "file=@/path/to/document.pdf"
```

---

### 5.2 Delete Document

**Endpoint:** `DELETE /api/v1/documents/{id}/delete-enhanced`
**Auth Required:** Yes (Student role)
**Description:** Delete an uploaded document

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Document deleted successfully"
}
```

**Example cURL:**
```bash
curl -X DELETE http://localhost:8080/api/v1/documents/456/delete-enhanced \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

### 5.3 Download Document

**Endpoint:** `GET /api/v1/documents/{id}/download-enhanced`
**Auth Required:** Yes
**Description:** Download a document file

**Response:** File download (binary)

**Example cURL:**
```bash
curl -X GET http://localhost:8080/api/v1/documents/456/download-enhanced \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -o downloaded_document.pdf
```

---

## 6. Public Scholarship Endpoints (Already Implemented)

### 6.1 Get All Scholarships

**Endpoint:** `GET /api/v1/scholarships`
**Auth Required:** No
**Description:** Get list of scholarships with pagination

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 10)
- `type`: Filter by type (optional)
- `search`: Search keyword (optional)

**Example cURL:**
```bash
curl -X GET "http://localhost:8080/api/v1/scholarships?page=1&limit=10"
```

---

### 6.2 Get Scholarship Details

**Endpoint:** `GET /api/v1/scholarships/{id}`
**Auth Required:** No
**Description:** Get detailed information about a specific scholarship

**Example cURL:**
```bash
curl -X GET http://localhost:8080/api/v1/scholarships/1
```

---

## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request
```json
{
  "error": "Invalid request body"
}
```

### 401 Unauthorized
```json
{
  "error": "User ID not found in context"
}
```

### 403 Forbidden
```json
{
  "error": "You can only modify your own applications"
}
```

### 404 Not Found
```json
{
  "error": "Application not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Failed to save application"
}
```

---

## Testing Workflow

### Complete Application Workflow Test

1. **Login as Student**
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email": "student@example.com", "password": "password"}'
   ```

2. **Check Eligibility**
   ```bash
   curl -X POST http://localhost:8080/api/v1/scholarships/1/check-eligibility \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"student_data": {"gpa": 3.5, "year_level": 3, "faculty": "Engineering", "family_income": 20000}}'
   ```

3. **Create Draft Application**
   ```bash
   curl -X POST http://localhost:8080/api/v1/applications/draft \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"scholarship_id": 1}'
   ```

4. **Save Personal Info**
   ```bash
   curl -X POST http://localhost:8080/api/v1/applications/123/sections/personal_info \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "prefix_th": "นาย",
       "first_name_th": "สมชาย",
       "last_name_th": "ใจดี",
       "email": "somchai@example.com",
       "phone": "0812345678"
     }'
   ```

5. **Upload Document**
   ```bash
   curl -X POST http://localhost:8080/api/v1/documents/applications/123/upload-enhanced \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -F "document_type=id_card" \
     -F "file=@id_card.pdf"
   ```

6. **Submit Application**
   ```bash
   curl -X POST http://localhost:8080/api/v1/applications/123/submit-enhanced \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"terms_accepted": true, "declaration_accepted": true}'
   ```

---

## File Structure

### Handler Files Created
- `/api/internal/handlers/application_draft_handler.go` - Draft management
- `/api/internal/handlers/application_section_handler.go` - Section save
- `/api/internal/handlers/application_submit_handler.go` - Submit application
- `/api/internal/handlers/document_enhanced_handler.go` - Document management

### Router Updates
- `/api/internal/router/router.go` - Added all new routes

### Upload Directory
- `/api/uploads/applications/` - Application documents
- `/api/uploads/documents/` - General documents

---

## Notes

1. **File Upload Limits:**
   - Default: 10MB
   - ID Card/Transcript: 5MB
   - House Photos: 20MB

2. **Allowed File Types:**
   - PDF (.pdf)
   - JPEG (.jpg, .jpeg)
   - PNG (.png)

3. **Required Documents:**
   - ID Card (id_card)
   - Transcript (transcript)

4. **Application Status Flow:**
   - `draft` → `submitted` → `under_review` → `approved/rejected`

5. **Section Names:**
   - personal_info
   - address_info
   - education_history
   - family_info
   - financial_info
   - activities_skills

---

## Database Tables Used

- `scholarship_applications` - Main application records
- `application_personal_info` - Personal information
- `application_addresses` - Address information
- `application_education_history` - Education records
- `application_family_members` - Family information
- `application_documents` - Document metadata
- `application_workflow` - Workflow tracking
- `scholarships` - Scholarship information
- `students` - Student records
- `users` - User accounts

---

## Security Considerations

1. All endpoints require JWT authentication (except public scholarship endpoints)
2. Student role required for application operations
3. Ownership verification prevents access to other students' applications
4. File type and size validation on uploads
5. Draft status check before modifications
6. Terms and declaration acceptance required for submission

---

## Next Steps

1. Implement email notifications for submission
2. Add webhook for document verification
3. Implement application progress tracking
4. Add real-time validation feedback
5. Implement auto-save functionality
6. Add application analytics

---

## Support

For issues or questions, please contact the development team.
