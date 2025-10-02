# 🚀 Quick Start Guide

## วิธีรัน Development Server

### วิธีที่ 1: ใช้ Make (แนะนำ)
```bash
make dev
```

### วิธีที่ 2: ใช้ Script โดยตรง
```bash
./dev.sh
```

### วิธีที่ 3: ใช้ Go โดยตรง
```bash
go run main.go
```

## 🔗 Links สำคัญ

- **API Documentation**: http://localhost:8080/swagger/
- **Health Check**: http://localhost:8080/health
- **Database**: PostgreSQL on port 5434

## 📋 Commands พื้นฐาน

```bash
# รัน server
make run

# สร้าง documentation
make docs

# Build binary
make build

# รัน tests
make test

# รัน database migration
make migrate

# ดู commands ทั้งหมด
make help
```

## 🔐 Test การเข้าสู่ระบบ

### 1. Register User ใหม่
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User",
    "student_id": "6012345678"
  }'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 3. ใช้ Token ที่ได้
```bash
curl -X GET http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 🎯 API Endpoints หลัก

### Authentication
- `POST /api/v1/auth/login` - เข้าสู่ระบบ
- `POST /api/v1/auth/register` - สมัครสมาชิก

### Scholarships
- `GET /api/v1/scholarships` - ดูทุนการศึกษา
- `POST /api/v1/scholarships` - สร้างทุน (Admin/Officer)

### Applications
- `GET /api/v1/applications/my` - ใบสมัครของฉัน (Student)
- `POST /api/v1/applications` - สมัครทุน (Student)

### Documents
- `POST /api/v1/documents/applications/:id/upload` - อัพโหลดเอกสาร
- `GET /api/v1/documents/types` - ประเภทเอกสาร

## 🛠️ Troubleshooting

### Database Connection Error
```bash
# ตรวจสอบว่า PostgreSQL รันอยู่
pg_isready -h localhost -p 5434

# รัน migration ใหม่
make migrate
```

### Port Already in Use
```bash
# หา process ที่ใช้ port 8080
lsof -ti:8080

# ฆ่า process
kill -9 $(lsof -ti:8080)
```

### Missing Dependencies
```bash
# ติดตั้ง dependencies ใหม่
make install
```

## 📱 Development Tools

### Hot Reload (Optional)
```bash
# ติดตั้ง Air สำหรับ hot reload
make install-air

# รันด้วย hot reload
air
```

### Database Tools
```bash
# เชื่อมต่อ database
psql "postgres://postgres:postgres@localhost:5434/scholarship_db"

# ดู tables
\dt

# ดูข้อมูล users
SELECT * FROM users LIMIT 5;
```

---

🎉 **พร้อมใช้งาน!** เข้าไปดู Swagger Documentation ที่ http://localhost:8080/swagger/ เพื่อทดสอบ API