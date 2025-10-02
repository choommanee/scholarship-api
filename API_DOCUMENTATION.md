# API Documentation - Scholarship Management System

**Base URL:** `http://localhost:8080`
**Version:** 1.0.0
**Last Updated:** 1 ตุลาคม 2025

---

## 📑 สารบัญ

1. [Authentication](#authentication)
2. [Scholarships](#scholarships)
3. [Applications](#applications)
4. [Payments](#payments)
5. [Analytics](#analytics)
6. [User Management](#user-management)
7. [News](#news)
8. [Error Codes](#error-codes)

---

## 🔐 Authentication

### Register User
สร้าง user account ใหม่

**Endpoint:** `POST /api/v1/auth/register`

**Request Body:**
```json
{
  "email": "student@university.ac.th",
  "password": "password123",
  "full_name": "John Doe",
  "role": "student"
}
```

**Response (201 Created):**
```json
{
  "user_id": "uuid",
  "email": "student@university.ac.th",
  "role": "student",
  "created_at": "2025-10-01T12:00:00Z"
}
```

**Errors:**
- `409` - Email already exists
- `400` - Invalid request body

---

### Login
เข้าสู่ระบบและรับ JWT token

**Endpoint:** `POST /api/v1/auth/login`

**Request Body:**
```json
{
  "email": "student@university.ac.th",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2025-10-08T12:00:00Z",
  "user": {
    "user_id": "uuid",
    "email": "student@university.ac.th",
    "full_name": "John Doe",
    "role": "student",
    "roles": ["student"]
  }
}
```

**Errors:**
- `401` - Invalid credentials
- `400` - Invalid request body

---

## 🎓 Scholarships

### Get Scholarships List
ดึงรายการทุนการศึกษาทั้งหมด (Public API)

**Endpoint:** `GET /api/v1/scholarships`

**Query Parameters:**
- `limit` (int) - จำนวนรายการต่อหน้า (default: 10)
- `offset` (int) - เริ่มจากรายการที่ (default: 0)
- `search` (string) - ค้นหาชื่อทุน
- `type` (string) - ประเภททุน (merit, need-based, etc.)
- `academic_year` (string) - ปีการศึกษา
- `active_only` (bool) - แสดงเฉพาะทุนที่เปิดรับสมัคร (default: true)

**Example Request:**
```bash
GET /api/v1/scholarships?limit=10&offset=0&type=merit&academic_year=2023
```

**Response (200 OK):**
```json
{
  "scholarships": [
    {
      "scholarship_id": 1,
      "source_id": 1,
      "scholarship_name": "Merit Scholarship 2023",
      "scholarship_type": "merit",
      "amount": 50000,
      "total_quota": 10,
      "available_quota": 8,
      "academic_year": "2023",
      "semester": "1",
      "eligibility_criteria": "GPA >= 3.50",
      "required_documents": "Transcript, ID Card",
      "application_start_date": "2023-05-01T00:00:00Z",
      "application_end_date": "2023-06-30T00:00:00Z",
      "interview_required": true,
      "is_active": true,
      "source": {
        "source_id": 1,
        "source_name": "University Foundation",
        "source_type": "internal"
      }
    }
  ],
  "total": 6,
  "limit": 10,
  "offset": 0
}
```

---

### Get Scholarship Details
ดึงข้อมูลทุนการศึกษารายการเดียว

**Endpoint:** `GET /api/v1/scholarships/:id`

**Response (200 OK):**
```json
{
  "scholarship_id": 1,
  "source_id": 1,
  "scholarship_name": "Merit Scholarship 2023",
  "scholarship_type": "merit",
  "amount": 50000,
  "total_quota": 10,
  "available_quota": 8,
  "academic_year": "2023",
  "semester": "1",
  "eligibility_criteria": "GPA >= 3.50",
  "required_documents": "Transcript, ID Card, Recommendation Letter",
  "application_start_date": "2023-05-01T00:00:00Z",
  "application_end_date": "2023-06-30T00:00:00Z",
  "interview_required": true,
  "is_active": true,
  "created_by": "uuid",
  "created_at": "2023-04-01T10:00:00Z",
  "updated_at": "2023-04-01T10:00:00Z",
  "source": {
    "source_id": 1,
    "source_name": "University Foundation",
    "source_type": "internal",
    "contact_person": "John Doe",
    "contact_email": "foundation@university.ac.th"
  }
}
```

**Errors:**
- `404` - Scholarship not found

---

### Create Scholarship
สร้างทุนการศึกษาใหม่ (Admin/Officer only)

**Endpoint:** `POST /api/v1/scholarships`

**Authorization:** `Bearer <token>` (admin or scholarship_officer)

**Request Body:**
```json
{
  "source_id": 1,
  "scholarship_name": "Excellence Scholarship 2024",
  "scholarship_type": "merit",
  "amount": 60000,
  "total_quota": 15,
  "academic_year": "2024",
  "semester": "1",
  "eligibility_criteria": "GPA >= 3.75, Leadership activities",
  "required_documents": "Transcript, Portfolio, Recommendation",
  "application_start_date": "2024-05-01T00:00:00Z",
  "application_end_date": "2024-06-30T00:00:00Z",
  "interview_required": true
}
```

**Response (201 Created):**
```json
{
  "message": "Scholarship created successfully",
  "scholarship": {
    "scholarship_id": 11,
    "source_id": 1,
    "scholarship_name": "Excellence Scholarship 2024",
    "amount": 60000,
    "is_active": true,
    "created_at": "2025-10-01T12:00:00Z"
  }
}
```

**Errors:**
- `403` - Insufficient permissions
- `400` - Invalid request body

---

### Get Available Scholarships
ดึงรายการทุนที่เปิดรับสมัครในขณะนี้

**Endpoint:** `GET /api/v1/scholarships/available`

**Response (200 OK):**
```json
{
  "scholarships": [
    {
      "scholarship_id": 1,
      "scholarship_name": "Merit Scholarship 2023",
      "amount": 50000,
      "available_quota": 8,
      "application_end_date": "2023-06-30T00:00:00Z"
    }
  ],
  "count": 3
}
```

---

## 📝 Applications

### Submit Application
ยื่นใบสมัครทุนการศึกษา (Student only)

**Endpoint:** `POST /api/v1/applications`

**Authorization:** `Bearer <token>` (student)

**Request Body:**
```json
{
  "scholarship_id": 1,
  "gpa": 3.75,
  "income": 25000,
  "statement": "เหตุผลที่ต้องการรับทุน...",
  "documents": [
    {
      "document_type": "transcript",
      "file_path": "/uploads/transcript.pdf"
    }
  ]
}
```

**Response (201 Created):**
```json
{
  "message": "Application submitted successfully",
  "application": {
    "application_id": 123,
    "scholarship_id": 1,
    "student_id": "uuid",
    "status": "submitted",
    "submitted_at": "2025-10-01T12:00:00Z"
  }
}
```

---

### Get My Applications
ดึงรายการใบสมัครของตัวเอง (Student)

**Endpoint:** `GET /api/v1/applications/my-applications`

**Authorization:** `Bearer <token>` (student)

**Response (200 OK):**
```json
{
  "applications": [
    {
      "application_id": 123,
      "scholarship": {
        "scholarship_id": 1,
        "scholarship_name": "Merit Scholarship 2023",
        "amount": 50000
      },
      "status": "under_review",
      "submitted_at": "2025-10-01T12:00:00Z",
      "current_step": "document_verification"
    }
  ],
  "total": 2
}
```

---

### Get Application Details
ดึงข้อมูลใบสมัครรายการเดียว

**Endpoint:** `GET /api/v1/applications/:id`

**Authorization:** `Bearer <token>` (student, officer, admin)

**Response (200 OK):**
```json
{
  "application_id": 123,
  "scholarship_id": 1,
  "student_id": "uuid",
  "status": "approved",
  "gpa": 3.75,
  "income": 25000,
  "statement": "เหตุผลที่ต้องการรับทุน...",
  "submitted_at": "2025-10-01T12:00:00Z",
  "reviewed_at": "2025-10-05T14:30:00Z",
  "reviewed_by": "uuid",
  "review_notes": "คุณสมบัติครบถ้วน",
  "scholarship": {
    "scholarship_name": "Merit Scholarship 2023",
    "amount": 50000
  },
  "student": {
    "student_id": "1234567890",
    "first_name": "John",
    "last_name": "Doe"
  },
  "documents": [
    {
      "document_id": "uuid",
      "document_type": "transcript",
      "file_path": "/uploads/transcript.pdf",
      "uploaded_at": "2025-10-01T12:05:00Z"
    }
  ]
}
```

---

## 💰 Payments

### Get Payment Methods
ดึงรายการวิธีการจ่ายเงิน (Admin/Officer only)

**Endpoint:** `GET /api/v1/payments/methods`

**Authorization:** `Bearer <token>` (admin or scholarship_officer)

**Response (200 OK):**
```json
[
  {
    "method_id": "uuid",
    "method_name": "Bank Transfer",
    "method_code": "bank_transfer",
    "description": "โอนเงินผ่านธนาคาร",
    "is_active": true,
    "configuration": {},
    "created_at": "2025-10-01T10:00:00Z"
  },
  {
    "method_id": "uuid",
    "method_name": "Cheque",
    "method_code": "cheque",
    "description": "เช็คธนาคาร",
    "is_active": true
  },
  {
    "method_id": "uuid",
    "method_name": "Mobile Banking",
    "method_code": "mobile_banking",
    "description": "Mobile Banking",
    "is_active": true
  }
]
```

**Errors:**
- `403` - Insufficient permissions (requires admin or scholarship_officer role)

---

### Create Payment Transaction
สร้างธุรกรรมการจ่ายเงิน (Admin/Officer only)

**Endpoint:** `POST /api/v1/payments/transactions`

**Authorization:** `Bearer <token>` (admin or scholarship_officer)

**Request Body:**
```json
{
  "allocation_id": 456,
  "amount": 50000,
  "payment_method": "bank_transfer",
  "bank_code": "BBL",
  "account_number": "1234567890",
  "payment_date": "2025-10-15T10:00:00Z",
  "reference_number": "REF123456",
  "notes": "การจ่ายงวดที่ 1"
}
```

**Response (201 Created):**
```json
{
  "transaction_id": "uuid",
  "allocation_id": 456,
  "amount": 50000,
  "payment_method": "bank_transfer",
  "payment_status": "pending",
  "payment_date": "2025-10-15T10:00:00Z",
  "reference_number": "REF123456",
  "created_at": "2025-10-01T12:00:00Z"
}
```

---

### Get Payment Transaction
ดึงข้อมูลธุรกรรมการจ่ายเงิน

**Endpoint:** `GET /api/v1/payments/transactions/:id`

**Authorization:** `Bearer <token>` (admin or scholarship_officer)

**Response (200 OK):**
```json
{
  "transaction_id": "uuid",
  "allocation_id": 456,
  "amount": 50000,
  "payment_method": "bank_transfer",
  "bank_code": "BBL",
  "account_number": "1234567890",
  "payment_date": "2025-10-15T10:00:00Z",
  "payment_status": "completed",
  "reference_number": "REF123456",
  "notes": "การจ่ายงวดที่ 1",
  "created_at": "2025-10-01T12:00:00Z",
  "updated_at": "2025-10-15T14:30:00Z"
}
```

---

### Update Transaction Status
อัพเดตสถานะการจ่ายเงิน (Admin/Officer only)

**Endpoint:** `PUT /api/v1/payments/transactions/:id/status`

**Authorization:** `Bearer <token>` (admin or scholarship_officer)

**Request Body:**
```json
{
  "status": "completed"
}
```

**Valid Status Values:**
- `pending` - รอดำเนินการ
- `processing` - กำลังดำเนินการ
- `completed` - สำเร็จ
- `failed` - ไม่สำเร็จ
- `cancelled` - ยกเลิก

**Response (200 OK):**
```json
{
  "message": "Status updated successfully"
}
```

---

### Create Disbursement Schedule
สร้างตารางการจ่ายเงินแบบงวด (Admin/Officer only)

**Endpoint:** `POST /api/v1/payments/disbursements`

**Authorization:** `Bearer <token>` (admin or scholarship_officer)

**Request Body:**
```json
{
  "allocation_id": 456,
  "installment_number": 1,
  "due_date": "2025-11-15T00:00:00Z",
  "amount": 25000
}
```

**Response (201 Created):**
```json
{
  "schedule_id": "uuid",
  "allocation_id": 456,
  "installment_number": 1,
  "due_date": "2025-11-15T00:00:00Z",
  "amount": 25000,
  "status": "scheduled",
  "created_at": "2025-10-01T12:00:00Z"
}
```

---

### Get Pending Disbursements
ดึงรายการจ่ายเงินที่ครบกำหนด (Admin/Officer only)

**Endpoint:** `GET /api/v1/payments/disbursements/pending`

**Authorization:** `Bearer <token>` (admin or scholarship_officer)

**Response (200 OK):**
```json
{
  "disbursements": [
    {
      "schedule_id": "uuid",
      "allocation_id": 456,
      "installment_number": 1,
      "due_date": "2025-11-15T00:00:00Z",
      "amount": 25000,
      "status": "scheduled"
    }
  ],
  "total": 5
}
```

---

## 📊 Analytics

### Get Dashboard Summary
ดึงข้อมูลสรุปสำหรับ dashboard (Admin/Officer only)

**Endpoint:** `GET /api/v1/analytics/dashboard`

**Authorization:** `Bearer <token>` (admin or scholarship_officer)

**Response (200 OK):**
```json
{
  "total_periods": 3,
  "latest_statistics": {
    "academic_year": "2023",
    "scholarship_round": "1",
    "total_applications": 150,
    "approved_applications": 45,
    "rejected_applications": 30,
    "success_rate": 30.0
  },
  "average_processing_time": 15.5,
  "bottlenecks": {
    "document_verification": 25,
    "interview_scheduling": 18,
    "committee_review": 12
  }
}
```

---

### Get Scholarship Statistics
ดึงสถิติทุนการศึกษา (Admin/Officer only)

**Endpoint:** `GET /api/v1/analytics/statistics`

**Query Parameters:**
- `year` (string) - ปีการศึกษา (required)
- `round` (string) - รอบการสมัคร (required)

**Example:**
```bash
GET /api/v1/analytics/statistics?year=2023&round=1
```

**Response (200 OK):**
```json
{
  "stat_id": "uuid",
  "academic_year": "2023",
  "scholarship_round": "1",
  "total_applications": 150,
  "approved_applications": 45,
  "rejected_applications": 30,
  "pending_applications": 75,
  "total_budget": 5000000,
  "allocated_budget": 2250000,
  "remaining_budget": 2750000,
  "average_amount": 50000,
  "success_rate": 30.0,
  "processing_time_avg": 15.5,
  "total_faculties": 12,
  "most_popular_scholarship": "Merit Scholarship",
  "created_at": "2025-10-01T10:00:00Z"
}
```

**Errors:**
- `404` - Statistics not found

---

### Get All Statistics
ดึงสถิติทั้งหมด (Admin/Officer only)

**Endpoint:** `GET /api/v1/analytics/statistics/all`

**Authorization:** `Bearer <token>` (admin or scholarship_officer)

**Response (200 OK):**
```json
[
  {
    "academic_year": "2023",
    "scholarship_round": "1",
    "total_applications": 150,
    "success_rate": 30.0
  },
  {
    "academic_year": "2023",
    "scholarship_round": "2",
    "total_applications": 180,
    "success_rate": 35.0
  }
]
```

---

### Get Average Processing Time
คำนวณเวลาเฉลี่ยในการดำเนินการ (Admin/Officer only)

**Endpoint:** `GET /api/v1/analytics/processing-time`

**Authorization:** `Bearer <token>` (admin or scholarship_officer)

**Response (200 OK):**
```json
{
  "average_processing_time_days": 15.5
}
```

---

### Get Bottlenecks Analysis
วิเคราะห์จุดคอขวดในกระบวนการ (Admin/Officer only)

**Endpoint:** `GET /api/v1/analytics/bottlenecks`

**Authorization:** `Bearer <token>` (admin or scholarship_officer)

**Response (200 OK):**
```json
{
  "bottlenecks": {
    "document_verification": 25,
    "interview_scheduling": 18,
    "committee_review": 12,
    "final_approval": 8
  }
}
```

---

## 👤 User Management

### Get User Profile
ดึงข้อมูลโปรไฟล์ของตัวเอง

**Endpoint:** `GET /api/v1/user/profile`

**Authorization:** `Bearer <token>`

**Response (200 OK):**
```json
{
  "user_id": "uuid",
  "email": "student@university.ac.th",
  "username": "john_doe",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "0812345678",
  "is_active": true,
  "role": "student",
  "roles": ["student"],
  "student": {
    "student_id": "1234567890",
    "faculty": "Engineering",
    "department": "Computer Engineering",
    "year_of_study": 3,
    "gpa": 3.75
  },
  "created_at": "2023-01-01T10:00:00Z",
  "last_login": "2025-10-01T12:00:00Z"
}
```

---

### Update Profile
อัพเดตข้อมูลโปรไฟล์

**Endpoint:** `PUT /api/v1/user/profile`

**Authorization:** `Bearer <token>`

**Request Body:**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "phone": "0812345678"
}
```

**Response (200 OK):**
```json
{
  "message": "Profile updated successfully",
  "user": {
    "user_id": "uuid",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "0812345678"
  }
}
```

---

## 📰 News

### Get News List
ดึงรายการข่าวสาร (Public API)

**Endpoint:** `GET /api/v1/news`

**Query Parameters:**
- `limit` (int) - จำนวนรายการต่อหน้า (default: 10)
- `offset` (int) - เริ่มจากรายการที่ (default: 0)

**Response (200 OK):**
```json
{
  "news": [
    {
      "news_id": 1,
      "title": "เปิดรับสมัครทุนการศึกษา 2024",
      "content": "มหาวิทยาลัยเปิดรับสมัคร...",
      "category": "scholarship",
      "published_at": "2025-10-01T10:00:00Z",
      "author": "Admin"
    }
  ],
  "total": 10,
  "limit": 10,
  "offset": 0
}
```

---

### Get News Details
ดึงข้อมูลข่าวรายการเดียว (Public API)

**Endpoint:** `GET /api/v1/news/:id`

**Response (200 OK):**
```json
{
  "news_id": 1,
  "title": "เปิดรับสมัครทุนการศึกษา 2024",
  "content": "มหาวิทยาลัยเปิดรับสมัครทุนการศึกษา ประจำปี 2024...",
  "category": "scholarship",
  "published_at": "2025-10-01T10:00:00Z",
  "author": "Admin",
  "views": 1250,
  "created_at": "2025-09-30T15:00:00Z"
}
```

---

## ⚠️ Error Codes

### HTTP Status Codes

| Code | Status | Description |
|------|--------|-------------|
| 200 | OK | สำเร็จ |
| 201 | Created | สร้างข้อมูลสำเร็จ |
| 400 | Bad Request | ข้อมูล request ไม่ถูกต้อง |
| 401 | Unauthorized | ไม่ได้ login หรือ token หมดอายุ |
| 403 | Forbidden | ไม่มีสิทธิ์เข้าถึง |
| 404 | Not Found | ไม่พบข้อมูล |
| 409 | Conflict | ข้อมูลซ้ำ |
| 500 | Internal Server Error | เกิดข้อผิดพลาดที่ server |

### Error Response Format

```json
{
  "error": "Error message description"
}
```

**Example Errors:**

```json
// 401 Unauthorized
{
  "error": "Invalid token or token expired"
}

// 403 Forbidden
{
  "error": "Insufficient permissions"
}

// 404 Not Found
{
  "error": "Scholarship not found"
}

// 400 Bad Request
{
  "error": "Invalid request body"
}

// 409 Conflict
{
  "error": "Email already exists"
}
```

---

## 🔒 Authorization

### Headers

ทุก endpoint ที่ต้องการ authentication จะต้องใส่ header:

```
Authorization: Bearer <your_jwt_token>
```

### Roles และสิทธิ์

| Role | Permissions |
|------|-------------|
| `student` | - View scholarships<br>- Submit applications<br>- View own applications<br>- Update profile |
| `scholarship_officer` | - All student permissions<br>- View all applications<br>- Review applications<br>- Manage payments<br>- View analytics |
| `admin` | - All permissions<br>- Create/Edit scholarships<br>- Manage users<br>- System configuration |
| `interviewer` | - View assigned applications<br>- Submit interview scores |
| `advisor` | - View student applications<br>- Provide recommendations |

---

## 📌 Rate Limiting

- Default: 100 requests per minute per IP
- Authenticated: 1000 requests per minute per user

**Headers:**
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1696147200
```

---

## 🧪 Testing

### Test Credentials

**Student Account:**
- Email: `test@example.com`
- Password: `password123`
- Role: student

**Admin Account:**
- Email: `admin@university.ac.th`
- Password: `admin123`
- Role: admin

### cURL Examples

**Login:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@university.ac.th","password":"admin123"}'
```

**Get Scholarships:**
```bash
curl http://localhost:8080/api/v1/scholarships?limit=5
```

**Get Payment Methods (with auth):**
```bash
curl http://localhost:8080/api/v1/payments/methods \
  -H "Authorization: Bearer <your_token>"
```

---

## 📖 Additional Resources

- **Swagger UI:** `http://localhost:8080/swagger`
- **Health Check:** `http://localhost:8080/health`
- **Test Report:** `./TEST_REPORT.md`
- **Progress Report:** `./PROGRESS_REPORT_UPDATED.md`

---

**Last Updated:** 1 ตุลาคม 2025
**API Version:** 1.0.0
**Maintained by:** Scholarship System Development Team
