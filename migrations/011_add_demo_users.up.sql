-- Add demo users for testing that match frontend expectations
-- Password hash for 'password123': $2a$10$xkFxzHAF7rB9PfQhHuWqPeQgvRQQ.CP1sLXoVo1UhjmR5GZW0CQCi

-- Insert demo users
INSERT INTO users (user_id, username, email, password_hash, first_name, last_name, phone, is_active) VALUES
(uuid_generate_v4(), 'admin.demo', 'admin@university.ac.th', '$2a$10$xkFxzHAF7rB9PfQhHuWqPeQgvRQQ.CP1sLXoVo1UhjmR5GZW0CQCi', 'ผู้ดูแลระบบ', 'Demo', '02-000-0001', true),
(uuid_generate_v4(), 'officer.demo', 'officer@university.ac.th', '$2a$10$xkFxzHAF7rB9PfQhHuWqPeQgvRQQ.CP1sLXoVo1UhjmR5GZW0CQCi', 'เจ้าหน้าที่ทุน', 'Demo', '02-000-0002', true),
(uuid_generate_v4(), 'interviewer.demo', 'interviewer@university.ac.th', '$2a$10$xkFxzHAF7rB9PfQhHuWqPeQgvRQQ.CP1sLXoVo1UhjmR5GZW0CQCi', 'ผู้สัมภาษณ์', 'Demo', '02-000-0003', true),
(uuid_generate_v4(), 'student1.demo', 'student1@university.ac.th', '$2a$10$xkFxzHAF7rB9PfQhHuWqPeQgvRQQ.CP1sLXoVo1UhjmR5GZW0CQCi', 'นักศึกษา', 'Demo', '081-000-0001', true)
ON CONFLICT (email) DO NOTHING;

-- Assign roles to demo users
INSERT INTO user_roles (user_id, role_id)
SELECT u.user_id, r.role_id
FROM users u, roles r
WHERE u.username = 'admin.demo' AND r.role_name = 'admin'
UNION ALL
SELECT u.user_id, r.role_id
FROM users u, roles r
WHERE u.username = 'officer.demo' AND r.role_name = 'scholarship_officer'
UNION ALL
SELECT u.user_id, r.role_id
FROM users u, roles r
WHERE u.username = 'interviewer.demo' AND r.role_name = 'interviewer'
UNION ALL
SELECT u.user_id, r.role_id
FROM users u, roles r
WHERE u.username = 'student1.demo' AND r.role_name = 'student'
ON CONFLICT (user_id, role_id) DO NOTHING;

-- Add student record for demo student
INSERT INTO students (student_id, user_id, faculty_code, department_code, year_level, gpa, admission_year)
SELECT '6388001', u.user_id, 'DEMO', 'DEMO', 3, 3.50, 2023
FROM users u WHERE u.username = 'student1.demo'
ON CONFLICT (student_id) DO NOTHING;

-- Update last login for demo users to make them appear active
UPDATE users 
SET last_login = NOW() 
WHERE username IN ('admin.demo', 'officer.demo', 'interviewer.demo', 'student1.demo');
