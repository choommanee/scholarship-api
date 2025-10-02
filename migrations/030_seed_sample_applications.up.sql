-- Migration 030: Seed Sample Application Data
-- This migration creates complete sample application data for testing

-- Insert sample personal info for application 1
INSERT INTO application_personal_info (
    application_id, prefix_th, prefix_en, first_name_th, last_name_th,
    first_name_en, last_name_en, email, phone, line_id,
    citizen_id, student_id, faculty, department, major, year_level,
    admission_type
) VALUES (
    1, 'นาย', 'Mr.', 'สมชาย', 'ใจดี',
    'Somchai', 'Jaidee', 'somchai.j@student.mahidol.ac.th', '081-234-5678', 'somchai_j',
    '1234567890123', '6512345678', 'คณะเศรษฐศาสตร์', 'ภาควิชาเศรษฐศาสตร์', 'เศรษฐศาสตร์', 3,
    'Portfolio'
) ON CONFLICT (application_id) DO UPDATE SET
    prefix_th = EXCLUDED.prefix_th,
    first_name_th = EXCLUDED.first_name_th,
    last_name_th = EXCLUDED.last_name_th,
    email = EXCLUDED.email,
    phone = EXCLUDED.phone;

-- Insert addresses for application 1
INSERT INTO application_addresses (
    application_id, address_type, house_number, village_number, alley, road,
    subdistrict, district, province, postal_code, address_line1
) VALUES
(1, 'registered', '123', '5', 'สุขสันต์ 3', 'พญาไท', 'ทุ่งพญาไท', 'ราชเทวี', 'กรุงเทพมหานคร', '10400',
 '123 หมู่ 5 ซอยสุขสันต์ 3 ถนนพญาไท แขวงทุ่งพญาไท เขตราชเทวี กรุงเทพมหานคร 10400'),
(1, 'current', '123', '5', 'สุขสันต์ 3', 'พญาไท', 'ทุ่งพญาไท', 'ราชเทวี', 'กรุงเทพมหานคร', '10400',
 '123 หมู่ 5 ซอยสุขสันต์ 3 ถนนพญาไท แขวงทุ่งพญาไท เขตราชเทวี กรุงเทพมหานคร 10400');

-- Insert education history for application 1
INSERT INTO application_education_history (
    application_id, education_level, school_name, school_province, gpa, graduation_year
) VALUES
(1, 'มัธยมปลาย', 'โรงเรียนสตรีวิทยา', 'กรุงเทพมหานคร', 3.85, '2565'),
(1, 'มัธยมต้น', 'โรงเรียนสตรีวิทยา', 'กรุงเทพมหานคร', 3.90, '2562');

-- Insert family members for application 1
INSERT INTO application_family_members (
    application_id, relationship, title, first_name, last_name, age,
    living_status, occupation, position, workplace, workplace_province, monthly_income, phone
) VALUES
(1, 'father', 'นาย', 'สมศักดิ์', 'ใจดี', 52, 'alive',
 'พนักงานบริษัท', 'พนักงานทั่วไป', 'บริษัท ABC จำกัด', 'กรุงเทพมหานคร', 18000.00, '081-111-1111'),
(1, 'mother', 'นาง', 'สมหญิง', 'ใจดี', 50, 'alive',
 'ค้าขาย', 'เจ้าของร้าน', 'ร้านขายของชำ', 'กรุงเทพมหานคร', 12000.00, '081-222-2222');

-- Insert financial info for application 1
INSERT INTO application_financial_info (
    application_id, monthly_allowance, daily_travel_cost, monthly_dorm_cost,
    other_monthly_costs, has_income, income_source, monthly_income
) VALUES (
    1, 15000.00, 50.00, 3000.00,
    2000.00, true, 'ทำงานพาร์ทไทม์', 3000.00
) ON CONFLICT (application_id) DO UPDATE SET
    monthly_allowance = EXCLUDED.monthly_allowance,
    daily_travel_cost = EXCLUDED.daily_travel_cost,
    monthly_dorm_cost = EXCLUDED.monthly_dorm_cost;

-- Insert scholarship history for application 1
INSERT INTO application_scholarship_history (
    application_id, scholarship_name, scholarship_type, amount, academic_year
) VALUES
(1, 'ทุนพัฒนาทักษะภาษาอังกฤษ', 'merit', 10000.00, '2566'),
(1, 'ทุนสนับสนุนนักศึกษาจากมูลนิธิ ABC', 'need', 15000.00, '2565');

-- Insert activities for application 1
INSERT INTO application_activities (
    application_id, activity_type, activity_name, description, achievement, award_level, year
) VALUES
(1, 'volunteer', 'ชมรมนักศึกษาอาสา',
 'จัดกิจกรรมบำเพ็ญประโยชน์ช่วยเหลือสังคม พัฒนาชุมชน และจัดอบรมให้ความรู้แก่เยาวชน ดำรงตำแหน่งประธานชมรม',
 'รางวัลชมรมดีเด่น ประจำปี 2566', 'university', '2566'),
(1, 'community_service', 'โครงการพัฒนาชุมชน',
 'พัฒนาคุณภาพชีวิตชุมชนในพื้นที่ห่างไกล สอนภาษาอังกฤษและคอมพิวเตอร์แก่เด็กนักเรียน ดำรงตำแหน่งหัวหน้าทีม',
 'โครงการสำเร็จลุล่วงด้วยดี', 'university', '2566'),
(1, 'academic_competition', 'การแข่งขันเศรษฐศาสตร์ระดับชาติ',
 'แข่งขันนำเสนอแผนธุรกิจด้านเศรษฐกิจสร้างสรรค์ ได้รับรางวัลรองชนะเลิศอันดับ 1',
 'รางวัลรองชนะเลิศอันดับ 1 ระดับประเทศ', 'national', '2565');

-- Insert sample documents for application 1
INSERT INTO application_documents (
    application_id, document_type, document_name, file_path, file_size, mime_type,
    upload_status, verification_status, original_filename
) VALUES
(1, 'transcript', 'ใบแสดงผลการศึกษา', '/uploads/applications/1/transcript.pdf',
 256000, 'application/pdf', 'uploaded', 'approved', 'transcript_6512345678.pdf'),
(1, 'income_certificate', 'หนังสือรับรองรายได้', '/uploads/applications/1/income.pdf',
 128000, 'application/pdf', 'uploaded', 'approved', 'income_cert.pdf'),
(1, 'id_card', 'สำเนาบัตรประชาชน', '/uploads/applications/1/id_card.pdf',
 64000, 'application/pdf', 'uploaded', 'approved', 'id_card.pdf'),
(1, 'photo', 'รูปถ่าย', '/uploads/applications/1/photo.jpg',
 512000, 'image/jpeg', 'uploaded', 'approved', 'student_photo.jpg'),
(1, 'recommendation', 'จดหมายรับรอง', '/uploads/applications/1/recommendation.pdf',
 96000, 'application/pdf', 'uploaded', 'pending', 'recommendation_letter.pdf');

-- Insert assets for application 1
INSERT INTO application_assets (
    application_id, asset_type, category, description, value, monthly_cost
) VALUES
(1, 'own_house', 'asset', 'บ้านพักอาศัย ชั้นเดียว พื้นที่ 50 ตร.ว.', 1500000.00, NULL),
(1, 'rent_land', 'liability', 'ที่ดินเช่าทำกิน 2 ไร่', NULL, 3000.00);

-- Insert living situation for application 1
INSERT INTO application_living_situation (
    application_id, living_with, living_details
) VALUES (
    1, 'parents', 'อาศัยอยู่กับบิดา มารดา และน้องชาย 1 คน ในบ้านของครอบครัวเอง'
) ON CONFLICT (application_id) DO UPDATE SET
    living_with = EXCLUDED.living_with,
    living_details = EXCLUDED.living_details;

-- Insert health info for application 1
INSERT INTO application_health_info (
    application_id, has_health_issues, health_condition, affects_study, monthly_medical_cost
) VALUES (
    1, false, 'สุขภาพแข็งแรงดี', false, 0.00
) ON CONFLICT (application_id) DO UPDATE SET
    has_health_issues = EXCLUDED.has_health_issues,
    health_condition = EXCLUDED.health_condition;

-- Insert funding needs for application 1
INSERT INTO application_funding_needs (
    application_id, tuition_support, monthly_support, book_support,
    dorm_support, other_support, other_details, total_requested, necessity_reason
) VALUES (
    1, 15000.00, 10000.00, 3000.00,
    0.00, 2000.00, 'ค่าใช้จ่ายในการทำกิจกรรมพัฒนาทักษะ', 30000.00,
    'ครอบครัวมีรายได้น้อย มีพี่น้องต้องเลี้ยงดู ต้องการทุนเพื่อช่วยเหลือครอบครัวและสามารถศึกษาต่อได้อย่างมีคุณภาพ'
) ON CONFLICT (application_id) DO UPDATE SET
    total_requested = EXCLUDED.total_requested,
    tuition_support = EXCLUDED.tuition_support;

-- Insert references for application 1
INSERT INTO application_references (
    application_id, title, first_name, last_name, relationship, phone, email
) VALUES
(1, 'อาจารย์', 'สมศรี', 'รักเรียน', 'อาจารย์ที่ปรึกษา', '02-123-4567', 'somsri.r@mahidol.ac.th');

-- Update application status and scores
UPDATE scholarship_applications
SET
    application_status = 'submitted',
    priority_score = 85.00,
    automated_score = 85.00,
    submitted_at = CURRENT_TIMESTAMP - INTERVAL '5 days',
    updated_at = CURRENT_TIMESTAMP
WHERE application_id = 1;

-- Create application 2 with different data
INSERT INTO application_personal_info (
    application_id, prefix_th, prefix_en, first_name_th, last_name_th,
    first_name_en, last_name_en, email, phone, line_id,
    citizen_id, student_id, faculty, department, major, year_level,
    admission_type
) VALUES (
    2, 'นางสาว', 'Ms.', 'สมหญิง', 'ขยันเรียน',
    'Somying', 'Khayanlean', 'somying.k@student.mahidol.ac.th', '089-876-5432', 'somying_k',
    '9876543210987', '6512345679', 'คณะเศรษฐศาสตร์', 'ภาควิชาเศรษฐศาสตร์', 'เศรษฐศาสตร์', 2,
    'Admission'
) ON CONFLICT (application_id) DO UPDATE SET
    prefix_th = EXCLUDED.prefix_th,
    first_name_th = EXCLUDED.first_name_th,
    last_name_th = EXCLUDED.last_name_th;

INSERT INTO application_financial_info (
    application_id, monthly_allowance, daily_travel_cost, monthly_dorm_cost,
    other_monthly_costs, has_income, income_source, monthly_income
) VALUES (
    2, 8000.00, 30.00, 0.00,
    1500.00, false, NULL, 0.00
) ON CONFLICT (application_id) DO UPDATE SET
    monthly_allowance = EXCLUDED.monthly_allowance,
    daily_travel_cost = EXCLUDED.daily_travel_cost;

INSERT INTO application_activities (
    application_id, activity_type, activity_name, description, achievement, award_level, year
) VALUES
(2, 'music', 'ชุมนุมดนตรี',
 'เล่นดนตรีในวงดนตรีของคณะ แสดงในงานต่างๆ ของมหาวิทยาลัย เป็นสมาชิกวงดนตรี',
 'แสดงในงานปีใหม่ของคณะ', 'faculty', '2566');

INSERT INTO application_documents (
    application_id, document_type, document_name, file_path, file_size, mime_type,
    upload_status, verification_status, original_filename
) VALUES
(2, 'transcript', 'ใบแสดงผลการศึกษา', '/uploads/applications/2/transcript.pdf',
 256000, 'application/pdf', 'uploaded', 'approved', 'transcript_6512345679.pdf'),
(2, 'income_certificate', 'หนังสือรับรองรายได้', '/uploads/applications/2/income.pdf',
 128000, 'application/pdf', 'uploaded', 'approved', 'income_cert.pdf'),
(2, 'id_card', 'สำเนาบัตรประชาชน', '/uploads/applications/2/id_card.pdf',
 64000, 'application/pdf', 'uploaded', 'approved', 'id_card.pdf'),
(2, 'photo', 'รูปถ่าย', '/uploads/applications/2/photo.jpg',
 512000, 'image/jpeg', 'uploaded', 'approved', 'student_photo.jpg'),
(2, 'recommendation', 'จดหมายรับรอง', '/uploads/applications/2/recommendation.pdf',
 96000, 'application/pdf', 'uploaded', 'pending', 'recommendation_letter.pdf');

UPDATE scholarship_applications
SET
    application_status = 'submitted',
    priority_score = 78.00,
    automated_score = 78.00,
    submitted_at = CURRENT_TIMESTAMP - INTERVAL '3 days',
    updated_at = CURRENT_TIMESTAMP
WHERE application_id = 2;
