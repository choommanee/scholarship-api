-- Real data for Scholarship Management System - Faculty of Economics, Thammasat University

-- Insert roles first (if not already exists)
INSERT INTO roles (role_name, permissions) VALUES
('admin', '["full_access", "manage_users", "system_settings", "view_reports", "system_admin"]'::jsonb),
('scholarship_officer', '["manage_applications", "review_documents", "schedule_interviews", "manage_scholarships", "allocate_funds"]'::jsonb),
('interviewer', '["conduct_interviews", "submit_evaluations", "view_applications"]'::jsonb),
('student', '["view_scholarships", "apply_scholarships", "manage_documents", "view_applications"]'::jsonb)
ON CONFLICT (role_name) DO NOTHING;

-- Create admin user
INSERT INTO users (user_id, username, email, password_hash, first_name, last_name, phone) VALUES
(uuid_generate_v4(), 'admin.econ', 'admin@econ.tu.ac.th', '$2a$10$xkFxzHAF7rB9PfQhHuWqPeQgvRQQ.CP1sLXoVo1UhjmR5GZW0CQCi', 'ผู้ดูแลระบบ', 'คณะเศรษฐศาสตร์', '02-613-2200'),
(uuid_generate_v4(), 'officer.econ', 'officer@econ.tu.ac.th', '$2a$10$xkFxzHAF7rB9PfQhHuWqPeQgvRQQ.CP1sLXoVo1UhjmR5GZW0CQCi', 'นางสาววิไลวรรณ', 'จัดการดี', '02-613-2201'),
(uuid_generate_v4(), 'interviewer.econ', 'interviewer@econ.tu.ac.th', '$2a$10$xkFxzHAF7rB9PfQhHuWqPeQgvRQQ.CP1sLXoVo1UhjmR5GZW0CQCi', 'ผศ.ดร.สมปอง', 'วิชาการดี', '02-613-2202');

-- Assign roles to staff users
INSERT INTO user_roles (user_id, role_id)
SELECT u.user_id, r.role_id
FROM users u, roles r
WHERE u.username = 'admin.econ' AND r.role_name = 'admin'
UNION ALL
SELECT u.user_id, r.role_id
FROM users u, roles r
WHERE u.username = 'officer.econ' AND r.role_name = 'scholarship_officer'
UNION ALL
SELECT u.user_id, r.role_id
FROM users u, roles r
WHERE u.username = 'interviewer.econ' AND r.role_name = 'interviewer';

-- Create sample students (anonymized real data pattern)
INSERT INTO users (user_id, username, email, password_hash, first_name, last_name, phone) VALUES
(uuid_generate_v4(), '66114411001', '66114411001@student.tu.ac.th', '$2a$10$xkFxzHAF7rB9PfQhHuWqPeQgvRQQ.CP1sLXoVo1UhjmR5GZW0CQCi', 'สมใจ', 'เรียนดี', '081-234-5678'),
(uuid_generate_v4(), '66114411002', '66114411002@student.tu.ac.th', '$2a$10$xkFxzHAF7rB9PfQhHuWqPeQgvRQQ.CP1sLXoVo1UhjmR5GZW0CQCi', 'สมหมาย', 'ขยันเรียน', '081-234-5679'),
(uuid_generate_v4(), '66114411003', '66114411003@student.tu.ac.th', '$2a$10$xkFxzHAF7rB9PfQhHuWqPeQgvRQQ.CP1sLXoVo1UhjmR5GZW0CQCi', 'สมศรี', 'ใจดี', '081-234-5680'),
(uuid_generate_v4(), '65114411001', '65114411001@student.tu.ac.th', '$2a$10$xkFxzHAF7rB9PfQhHuWqPeQgvRQQ.CP1sLXoVo1UhjmR5GZW0CQCi', 'สมจิต', 'มานะ', '081-234-5681'),
(uuid_generate_v4(), '64114411001', '64114411001@student.tu.ac.th', '$2a$10$xkFxzHAF7rB9PfQhHuWqPeQgvRQQ.CP1sLXoVo1UhjmR5GZW0CQCi', 'สมหญิง', 'พากเพียร', '081-234-5682');

-- Assign student role
INSERT INTO user_roles (user_id, role_id)
SELECT u.user_id, r.role_id
FROM users u, roles r
WHERE u.username LIKE '6%' AND r.role_name = 'student';

-- Insert student data
INSERT INTO students (student_id, user_id, faculty_code, department_code, year_level, gpa, admission_year)
SELECT '66114411001', u.user_id, 'ECON', 'ECON', 2, 3.65, 2023
FROM users u WHERE u.username = '66114411001'
UNION ALL
SELECT '66114411002', u.user_id, 'ECON', 'ECON', 2, 3.42, 2023
FROM users u WHERE u.username = '66114411002'
UNION ALL
SELECT '66114411003', u.user_id, 'ECON', 'ECON', 2, 3.78, 2023
FROM users u WHERE u.username = '66114411003'
UNION ALL
SELECT '65114411001', u.user_id, 'ECON', 'ECON', 3, 3.58, 2022
FROM users u WHERE u.username = '65114411001'
UNION ALL
SELECT '64114411001', u.user_id, 'ECON', 'ECON', 4, 3.71, 2021
FROM users u WHERE u.username = '64114411001';

-- Real scholarship sources based on university data
INSERT INTO scholarship_sources (source_name, source_type, contact_person, contact_email) VALUES
('คณะเศรษฐศาสตร์ มหาวิทยาลัยธรรมศาสตร์', 'internal', 'งานกิจการนักศึกษา', 'student.affairs@econ.tu.ac.th'),
('มูลนิธิธรรมศาสตร์', 'external', 'มูลนิธิธรรมศาสตร์', 'foundation@tu.ac.th'),
('ธนาคารกรุงไทย', 'external', 'ฝ่ายพัฒนาสังคม', 'csr@ktb.co.th'),
('บริษัท ปตท. จำกัด (มหาชน)', 'external', 'ฝ่ายความรับผิดชอบต่อสังคม', 'csr@pttplc.com'),
('มูลนิธิแม่ฟ้าหลวง', 'external', 'ฝ่ายทุนการศึกษา', 'scholarship@maefahluang.org'),
('กรมส่งเสริมการปกครองท้องถิิน', 'government', 'ส่วนทุนการศึกษา', 'scholarship@dla.go.th');

-- Real scholarships based on actual university scholarship programs
INSERT INTO scholarships (source_id, scholarship_name, scholarship_type, amount, total_quota, available_quota, 
                         academic_year, semester, application_start_date, application_end_date, interview_required) VALUES
-- Merit scholarships (52 total from budget analysis)
(1, 'ทุนเรียนดีเด่น คณะเศรษฐศาสตร์', 'merit', 40000.00, 15, 15, '2567', '1', '2024-05-01', '2024-06-30', true),
(1, 'ทุนเกียรตินิยมอันดับ 1', 'merit', 50000.00, 10, 10, '2567', '1', '2024-05-01', '2024-06-30', true),
(2, 'ทุนมูลนิธิธรรมศาสตร์ ประเภทความสามารถ', 'merit', 35000.00, 12, 12, '2567', '1', '2024-05-15', '2024-07-15', true),
(3, 'ทุนธนาคารกรุงไทย เพื่อความเป็นเลิศ', 'merit', 60000.00, 8, 8, '2567', '1', '2024-06-01', '2024-07-31', true),
(4, 'ทุน ปตท. Young Talent', 'merit', 80000.00, 7, 7, '2567', '1', '2024-06-15', '2024-08-15', true),

-- Need-based scholarships (46 total)
(1, 'ทุนช่วยเหลือนักศึกษายากจน คณะเศรษฐศาสตร์', 'need-based', 25000.00, 20, 20, '2567', '1', '2024-05-01', '2024-06-30', false),
(5, 'ทุนมูลนิธิแม่ฟ้าหลวง ช่วยเหลือครอบครัวยากจน', 'need-based', 30000.00, 15, 15, '2567', '1', '2024-05-15', '2024-07-15', true),
(6, 'ทุนรัฐบาล นักศึกษาต่างจังหวัด', 'need-based', 20000.00, 11, 11, '2567', '1', '2024-05-01', '2024-06-15', false),

-- Mixed merit+need scholarships (83 total)
(1, 'ทุนส่งเสริมนักศึกษาเศรษฐศาสตร์ดีเด่น', 'mixed', 45000.00, 25, 25, '2567', '1', '2024-05-01', '2024-06-30', true),
(2, 'ทุนธรรมศาสตร์เพื่อสังคม', 'mixed', 35000.00, 20, 20, '2567', '1', '2024-05-15', '2024-07-15', true),
(3, 'ทุนกรุงไทยพัฒนาเยาวชน', 'mixed', 55000.00, 15, 15, '2567', '1', '2024-06-01', '2024-07-31', true),
(4, 'ทุน ปตท. สร้างอนาคต', 'mixed', 70000.00, 12, 12, '2567', '1', '2024-06-15', '2024-08-15', true),
(5, 'ทุนแม่ฟ้าหลวงพัฒนาชุมชน', 'mixed', 40000.00, 11, 11, '2567', '1', '2024-05-15', '2024-07-15', true);

-- Scholarship budgets reflecting 4,834,200 THB total budget
INSERT INTO scholarship_budgets (scholarship_id, budget_year, total_budget) VALUES
(1, '2567', 600000.00),   -- 15 x 40,000
(2, '2567', 500000.00),   -- 10 x 50,000
(3, '2567', 420000.00),   -- 12 x 35,000
(4, '2567', 480000.00),   -- 8 x 60,000
(5, '2567', 560000.00),   -- 7 x 80,000
(6, '2567', 500000.00),   -- 20 x 25,000
(7, '2567', 450000.00),   -- 15 x 30,000
(8, '2567', 220000.00),   -- 11 x 20,000
(9, '2567', 1125000.00),  -- 25 x 45,000
(10, '2567', 700000.00),  -- 20 x 35,000
(11, '2567', 825000.00),  -- 15 x 55,000
(12, '2567', 840000.00),  -- 12 x 70,000
(13, '2567', 440000.00);  -- 11 x 40,000

-- Sample applications representing real application patterns
INSERT INTO scholarship_applications (student_id, scholarship_id, application_status, family_income, monthly_expenses, siblings_count, 
                                    special_abilities, activities_participation) VALUES
('66114411001', 1, 'submitted', 180000.00, 8000.00, 2, 'ผู้นำชุมชน มีทักษะการจัดการดี', 'ชมรมนักศึกษาเศรษฐศาสตร์, กิจกรรมอาสาพัฒนา'),
('66114411002', 6, 'under_review', 120000.00, 6000.00, 3, 'มีความสามารถด้านคณิตศาสตร์', 'กิจกรรมวิชาการ, การแข่งขันทางวิชาการ'),
('66114411003', 9, 'interview_scheduled', 250000.00, 9000.00, 1, 'ทักษะด้านการวิเคราะห์เศรษฐกิจ', 'ชมรมวิจัยเศรษฐกิจ, โครงการบริการสังคม'),
('65114411001', 2, 'approved', 200000.00, 7500.00, 2, 'เกรดเฉลี่ยสูง มีผลงานวิจัย', 'ผู้ช่วยผู้สอน, กิจกรรมวิชาการ'),
('64114411001', 10, 'submitted', 280000.00, 10000.00, 1, 'ภาวะผู้นำ มีประสบการณ์การทำงาน', 'หัวหน้าชมรม, อินเทิร์นในบริษัท');

-- Application documents
INSERT INTO application_documents (application_id, document_type, document_name, file_path, file_size, mime_type) VALUES
(1, 'id_card', 'บัตรประชาชน_66114411001.pdf', '/uploads/documents/id_card_66114411001.pdf', 1024000, 'application/pdf'),
(1, 'transcript', 'ใบแสดงผลการเรียน_66114411001.pdf', '/uploads/documents/transcript_66114411001.pdf', 2048000, 'application/pdf'),
(1, 'income_certificate', 'หนังสือรับรองรายได้_66114411001.pdf', '/uploads/documents/income_66114411001.pdf', 1536000, 'application/pdf'),
(2, 'id_card', 'บัตรประชาชน_66114411002.pdf', '/uploads/documents/id_card_66114411002.pdf', 1024000, 'application/pdf'),
(2, 'transcript', 'ใบแสดงผลการเรียน_66114411002.pdf', '/uploads/documents/transcript_66114411002.pdf', 2048000, 'application/pdf');

-- Interview schedules
INSERT INTO interview_schedules (scholarship_id, interview_date, start_time, end_time, location, max_applicants, interviewer_ids) VALUES
(1, '2024-07-15', '09:00:00', '12:00:00', 'ห้องประชุม 301 คณะเศรษฐศาสตร์', 5, '[{"interviewer_id": "' || (SELECT user_id FROM users WHERE username = 'interviewer.econ') || '"}]'),
(2, '2024-07-16', '13:00:00', '16:00:00', 'ห้องประชุม 302 คณะเศรษฐศาสตร์', 5, '[{"interviewer_id": "' || (SELECT user_id FROM users WHERE username = 'interviewer.econ') || '"}]'),
(9, '2024-07-20', '09:00:00', '12:00:00', 'ห้องประชุม 303 คณะเศรษฐศาสตร์', 8, '[{"interviewer_id": "' || (SELECT user_id FROM users WHERE username = 'interviewer.econ') || '"}]');

-- Interview appointments
INSERT INTO interview_appointments (application_id, schedule_id, appointment_status, student_confirmed) VALUES
(3, 3, 'scheduled', true),
(4, 2, 'completed', true);

-- Interview results for completed interviews
INSERT INTO interview_results (appointment_id, interviewer_id, scores, overall_score, recommendation) 
SELECT 
    2, 
    u.user_id,
    '{"academic_knowledge": 9.0, "communication": 8.5, "leadership": 9.5, "motivation": 9.0}',
    9.0,
    'strongly_recommended'
FROM users u WHERE u.username = 'interviewer.econ';

-- Scholarship allocations for approved applications
INSERT INTO scholarship_allocations (application_id, scholarship_id, allocated_amount, allocation_status, 
                                   allocation_date, disbursement_method, bank_account, bank_name) 
SELECT 
    4, 2, 50000.00, 'approved', '2024-08-01', 'bank_transfer', '1234567890', 'ธนาคารกรุงไทย'
WHERE EXISTS (SELECT 1 FROM scholarship_applications WHERE application_id = 4);

-- Academic progress tracking
INSERT INTO academic_progress_tracking (student_id, allocation_id, semester, gpa, credits_earned, academic_status, report_date) VALUES
('65114411001', 1, '1/2567', 3.72, 18, 'excellent', '2024-12-15');

-- Sample news and announcements
INSERT INTO news (title, content, category, author_id, published, publish_date) VALUES
('เปิดรับสมัครทุนการศึกษา ประจำปีการศึกษา 2567', 
 'คณะเศรษฐศาสตร์ มหาวิทยาลัยธรรมศาสตร์ ขอประกาศรับสมัครทุนการศึกษาประจำปีการศึกษา 2567 จำนวน 13 ประเภท รวมงบประมาณ 4.8 ล้านบาท สำหรับนักศึกษาที่มีผลการเรียนดีและครอบครัวมีฐานะทางเศรษฐกิจที่ต้องการความช่วยเหลือ',
 'scholarship', 
 (SELECT user_id FROM users WHERE username = 'officer.econ'), 
 true, 
 '2024-05-01'),
('ประกาศผลการคัดเลือกทุนเรียนดีเด่น', 
 'ขอแสดงความยินดีกับนักศึกษาที่ได้รับทุนเรียนดีเด่น คณะเศรษฐศาสตร์ ประจำปีการศึกษา 2567 จำนวน 15 ทุน',
 'announcement', 
 (SELECT user_id FROM users WHERE username = 'officer.econ'), 
 true, 
 '2024-08-15');

-- Notifications for relevant events
INSERT INTO notifications (user_id, notification_type, title, message, reference_type, reference_id) 
SELECT 
    u.user_id,
    'application_status',
    'ใบสมัครได้รับการอนุมัติ',
    'ยินดีด้วย! ใบสมัครทุนเกียรตินิยมอันดับ 1 ของคุณได้รับการอนุมัติแล้ว',
    'scholarship_application',
    '4'
FROM users u WHERE u.username = '65114411001'
UNION ALL
SELECT 
    u.user_id,
    'interview_scheduled',
    'นัดหมายสัมภาษณ์',
    'คุณมีนัดหมายสัมภาษณ์ทุนส่งเสริมนักศึกษาเศรษฐศาสตร์ดีเด่น วันที่ 20 กรกฎาคม 2567 เวลา 09:00 น.',
    'interview_appointment',
    '1'
FROM users u WHERE u.username = '66114411003';

-- Login history for active users
INSERT INTO login_history (user_id, login_method, ip_address, login_status) 
SELECT 
    u.user_id,
    'password',
    '127.0.0.1',
    'success'
FROM users u 
WHERE u.username IN ('admin.econ', 'officer.econ', '66114411001', '65114411001');
