-- Seed data for demo accounts
-- Password for all accounts: password123
-- Bcrypt hash: $2a$10$YourHashHere (will be generated)

-- Delete existing demo users if they exist
DELETE FROM user_roles WHERE user_id IN (
    SELECT user_id FROM users WHERE email IN (
        'admin@university.ac.th',
        'officer@university.ac.th',
        'interviewer@university.ac.th',
        'student1@university.ac.th'
    )
);

DELETE FROM students WHERE user_id IN (
    SELECT user_id FROM users WHERE email IN (
        'student1@university.ac.th'
    )
);

DELETE FROM users WHERE email IN (
    'admin@university.ac.th',
    'officer@university.ac.th',
    'interviewer@university.ac.th',
    'student1@university.ac.th'
);

-- Insert demo users
-- Password hash for 'password123'
INSERT INTO users (user_id, username, email, password_hash, first_name, last_name, is_active, created_at, updated_at)
VALUES
    ('a1111111-1111-1111-1111-111111111111', 'admin', 'admin@university.ac.th',
     '$2a$10$82AdwmcvcuTRlWuZ27UCb.A0vTyt0JwKpNdiWGI2sndWC9TT2ZMQ6', 'Admin', 'User', true, NOW(), NOW()),

    ('a2222222-2222-2222-2222-222222222222', 'officer', 'officer@university.ac.th',
     '$2a$10$82AdwmcvcuTRlWuZ27UCb.A0vTyt0JwKpNdiWGI2sndWC9TT2ZMQ6', 'Scholarship', 'Officer', true, NOW(), NOW()),

    ('a3333333-3333-3333-3333-333333333333', 'interviewer', 'interviewer@university.ac.th',
     '$2a$10$82AdwmcvcuTRlWuZ27UCb.A0vTyt0JwKpNdiWGI2sndWC9TT2ZMQ6', 'Interview', 'Staff', true, NOW(), NOW()),

    ('a4444444-4444-4444-4444-444444444444', 'student1', 'student1@university.ac.th',
     '$2a$10$82AdwmcvcuTRlWuZ27UCb.A0vTyt0JwKpNdiWGI2sndWC9TT2ZMQ6', 'Student', 'One', true, NOW(), NOW())
ON CONFLICT (email) DO NOTHING;

-- Ensure roles exist
INSERT INTO roles (role_id, role_name, role_description, permissions)
VALUES
    (1, 'admin', 'System Administrator', '["all"]'),
    (2, 'scholarship_officer', 'Scholarship Officer', '["manage_scholarships", "review_applications"]'),
    (3, 'interviewer', 'Interviewer', '["view_applications", "conduct_interviews"]'),
    (4, 'student', 'Student', '["apply_scholarship", "view_own_applications"]')
ON CONFLICT (role_id) DO NOTHING;

-- Assign roles to demo users
INSERT INTO user_roles (user_id, role_id, assigned_by, is_active, assigned_at)
VALUES
    ('a1111111-1111-1111-1111-111111111111', 1, 'a1111111-1111-1111-1111-111111111111', true, NOW()),
    ('a2222222-2222-2222-2222-222222222222', 2, 'a1111111-1111-1111-1111-111111111111', true, NOW()),
    ('a3333333-3333-3333-3333-333333333333', 3, 'a1111111-1111-1111-1111-111111111111', true, NOW()),
    ('a4444444-4444-4444-4444-444444444444', 4, 'a1111111-1111-1111-1111-111111111111', true, NOW())
ON CONFLICT (user_id, role_id) DO UPDATE SET is_active = true;

-- Create student profile for student user
INSERT INTO students (student_id, user_id, faculty_code, department_code, year_level, gpa, admission_year, student_status)
VALUES
    ('6412345678', 'a4444444-4444-4444-4444-444444444444', 'ECON', 'ECON', 2, 3.50, 2023, 'active')
ON CONFLICT (student_id) DO NOTHING;

-- Add some sample scholarships
INSERT INTO scholarships (scholarship_name, scholarship_type, amount, total_quota, available_quota,
    academic_year, semester, application_start_date, application_end_date, interview_required, created_by)
VALUES
    ('ทุนพัฒนาศักยภาพนักศึกษา', 'merit-based', 100000.00, 10, 10,
     '2567', '1', CURRENT_DATE, CURRENT_DATE + INTERVAL '30 days', true, 'a1111111-1111-1111-1111-111111111111'),

    ('ทุนช่วยเหลือนักศึกษายากจน', 'need-based', 150000.00, 15, 15,
     '2567', '1', CURRENT_DATE, CURRENT_DATE + INTERVAL '30 days', false, 'a1111111-1111-1111-1111-111111111111'),

    ('ทุนเรียนดี กีฬาเด่น', 'activity-based', 80000.00, 5, 5,
     '2567', '1', CURRENT_DATE, CURRENT_DATE + INTERVAL '30 days', true, 'a1111111-1111-1111-1111-111111111111')
ON CONFLICT DO NOTHING;

-- Log seed completion
DO $$
BEGIN
    RAISE NOTICE '✅ Demo data seeded successfully!';
    RAISE NOTICE '📧 Demo Accounts:';
    RAISE NOTICE '   Admin: admin@university.ac.th / password123';
    RAISE NOTICE '   Officer: officer@university.ac.th / password123';
    RAISE NOTICE '   Interviewer: interviewer@university.ac.th / password123';
    RAISE NOTICE '   Student: student1@university.ac.th / password123';
END $$;
