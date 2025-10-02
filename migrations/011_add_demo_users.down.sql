-- Remove demo users and their associated data

-- Remove student record
DELETE FROM students WHERE student_id = '6388001';

-- Remove user roles for demo users
DELETE FROM user_roles 
WHERE user_id IN (
    SELECT user_id FROM users 
    WHERE username IN ('admin.demo', 'officer.demo', 'interviewer.demo', 'student1.demo')
);

-- Remove demo users
DELETE FROM users 
WHERE username IN ('admin.demo', 'officer.demo', 'interviewer.demo', 'student1.demo');
