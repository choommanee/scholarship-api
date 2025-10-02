# Scholarship Management System API

ระบบบริหารจัดการทุนการศึกษา สำหรับคณะเศรษฐศาสตร์ มหาวิทยาลัยธรรมศาสตร์

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Fiber](https://img.shields.io/badge/Fiber-v2.52.8-00ACD7?style=flat)](https://gofiber.io)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-336791?style=flat&logo=postgresql)](https://www.postgresql.org)

---

## 📋 สารบัญ

- [Features](#-features)
- [Tech Stack](#-tech-stack)
- [เริ่มต้นใช้งาน](#-เริ่มต้นใช้งาน)
- [การตั้งค่า](#-การตั้งค่า)
- [ฐานข้อมูล](#-ฐานข้อมูล)
- [API Documentation](#-api-documentation)
- [การทดสอบ](#-การทดสอบ)

---

## 🌟 Features

### ฟีเจอร์หลัก
- ✅ **User Management** - ระบบผู้ใช้ พร้อม Role-based Access Control
- ✅ **Scholarship Management** - จัดการทุนการศึกษา
- ✅ **Application System** - ยื่นและติดตามใบสมัคร
- ✅ **Payment System** - ระบบการจ่ายเงิน รองรับหลายวิธี
- ✅ **Analytics & Reporting** - สถิติและรายงาน
- ✅ **Email Notifications** - แจ้งเตือนทางอีเมล
- ✅ **Document Management** - จัดการเอกสาร พร้อม version control

### คุณสมบัติทางเทคนิค
- 🏗️ **Clean Architecture** - โครงสร้างโค้ดที่ชัดเจน
- 🔐 **JWT Authentication** - ระบบความปลอดภัย
- 📊 **41 ตารางในฐานข้อมูล** - ครอบคลุมทุกความต้องการ
- 🚀 **246 API Endpoints** - RESTful APIs
- ✅ **ทดสอบและแก้ไขแล้ว** - พร้อมใช้งาน

---

## 🛠 Tech Stack

- **Language:** Go 1.21+
- **Framework:** Fiber v2.52.8
- **Database:** PostgreSQL 14+
- **Authentication:** JWT

---

## 🚀 เริ่มต้นใช้งาน

### ความต้องการของระบบ

- Go 1.21+
- PostgreSQL 14+

### ติดตั้ง

1. **Clone repository**
\`\`\`bash
cd /Users/sakdachoommanee/Documents/fund\ system/fund/api
\`\`\`

2. **ติดตั้ง dependencies**
\`\`\`bash
go mod download
\`\`\`

3. **ตั้งค่า environment**
\`\`\`bash
# แก้ไขไฟล์ .env ด้วยข้อมูลฐานข้อมูลของคุณ
\`\`\`

4. **รัน server**
\`\`\`bash
go run main.go
\`\`\`

Server จะทำงานที่: \`http://localhost:8080\`

---

## ⚙️ การตั้งค่า

### Environment Variables (.env)

\`\`\`env
# Database
DATABASE_URL=postgres://postgres:postgres@localhost:5434/scholarship_db?sslmode=disable

# JWT
JWT_SECRET=your-secret-key-here

# Server
PORT=8080
ENVIRONMENT=development

# Upload
UPLOAD_PATH=./uploads
MAX_FILE_SIZE=10485760
\`\`\`

---

## 💾 ฐานข้อมูล

### โครงสร้าง

**41 ตาราง** แบ่งเป็นกลุ่ม:
- Core Tables (12 tables)
- Interview System (3 tables)  
- Payment System (5 tables)
- Analytics System (2 tables)
- Email System (2 tables)
- และอื่นๆ

### Migration

\`\`\`bash
# รัน migrations
cd migrations
./run_migrations.sh
\`\`\`

---

## 📚 API Documentation

### Base URL
\`http://localhost:8080\`

### Public Endpoints
- \`GET /health\` - Health check
- \`GET /api/v1/scholarships\` - รายการทุน
- \`POST /api/v1/auth/login\` - Login

### Protected Endpoints
- \`GET /api/v1/user/profile\` - โปรไฟล์
- \`POST /api/v1/applications\` - ยื่นใบสมัคร
- \`GET /api/v1/payments/methods\` - วิธีจ่ายเงิน (Admin)
- \`GET /api/v1/analytics/dashboard\` - Dashboard (Admin)

**เอกสารเพิ่มเติม:**
- [API_DOCUMENTATION.md](./API_DOCUMENTATION.md) - เอกสาร API แบบละเอียด
- Swagger UI: http://localhost:8080/swagger

---

## 🧪 การทดสอบ

### บัญชีทดสอบ

**นักศึกษา:**
- Email: \`test@example.com\`
- Password: \`password123\`

**Admin:**
- Email: \`admin@university.ac.th\`
- Password: \`admin123\`

### รัน Tests

\`\`\`bash
# ทดสอบ API พื้นฐาน
./test_api.sh

# ทดสอบ Admin APIs
./test_admin_apis.sh
\`\`\`

**ผลการทดสอบ:** ดูที่ [TEST_REPORT.md](./TEST_REPORT.md)

---

## 📊 สถานะโครงการ

### ✅ พร้อมใช้งาน 95%

#### เสร็จสมบูรณ์:
- ✅ Database Schema (41 ตาราง)
- ✅ Models, Repository, Handlers
- ✅ 246 API Endpoints
- ✅ Authentication & Authorization
- ✅ Testing & Bug Fixes
- ✅ Documentation

**รายละเอียด:** [PROGRESS_REPORT_UPDATED.md](./PROGRESS_REPORT_UPDATED.md)

---

## 📝 เอกสารเพิ่มเติม

- **API Documentation:** [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)
- **Test Report:** [TEST_REPORT.md](./TEST_REPORT.md)
- **Progress Report:** [PROGRESS_REPORT_UPDATED.md](./PROGRESS_REPORT_UPDATED.md)
- **Swagger UI:** http://localhost:8080/swagger

---

## 📞 ติดต่อ

- **Issues:** GitHub Issues
- **Email:** support@university.ac.th

---

**สร้างด้วย ❤️ โดย คณะเศรษฐศาสตร์ มหาวิทยาลัยธรรมศาสตร์**

**Test Reports:**
- [TEST_REPORT.md](./TEST_REPORT.md) - Basic API tests
- [COMPREHENSIVE_TEST_REPORT.md](./COMPREHENSIVE_TEST_REPORT.md) - **Full detailed tests (92.31% pass rate)** ✅

**Test Scripts:**
- `test_api.sh` - Basic API tests
- `test_admin_apis.sh` - Admin API tests
- `comprehensive_api_test.sh` - Comprehensive tests with all field validation

**อัพเดตล่าสุด:** 2 ตุลาคม 2025
