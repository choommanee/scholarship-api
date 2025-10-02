-- Interview Slots Management (Simplified)
CREATE TABLE IF NOT EXISTS interview_slots (
    id SERIAL PRIMARY KEY,
    scholarship_id INT NOT NULL,
    interviewer_id UUID NOT NULL,
    interview_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    location VARCHAR(255),
    max_capacity INT DEFAULT 1,
    current_bookings INT DEFAULT 0,
    is_available BOOLEAN DEFAULT TRUE,
    created_by UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (scholarship_id) REFERENCES scholarships(scholarship_id) ON DELETE CASCADE,
    FOREIGN KEY (interviewer_id) REFERENCES users(user_id),
    FOREIGN KEY (created_by) REFERENCES users(user_id)
);

-- Interview Bookings (Simplified)
CREATE TABLE IF NOT EXISTS interview_bookings (
    id SERIAL PRIMARY KEY,
    slot_id INT NOT NULL,
    application_id INT NOT NULL,
    student_id UUID NOT NULL,
    booking_status VARCHAR(20) DEFAULT 'booked',
    booked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    confirmed_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (slot_id) REFERENCES interview_slots(id) ON DELETE CASCADE,
    FOREIGN KEY (application_id) REFERENCES scholarship_applications(application_id) ON DELETE CASCADE,
    FOREIGN KEY (student_id) REFERENCES users(user_id)
);

-- Interview Scores (Simplified) - Renamed to avoid conflict
CREATE TABLE IF NOT EXISTS interview_scores_detailed (
    id SERIAL PRIMARY KEY,
    booking_id INT NOT NULL,
    interviewer_id UUID NOT NULL,
    criteria_name VARCHAR(100) NOT NULL,
    score DECIMAL(5,2) NOT NULL,
    max_score DECIMAL(5,2) NOT NULL DEFAULT 100.00,
    comments TEXT,
    scored_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (booking_id) REFERENCES interview_bookings(id) ON DELETE CASCADE,
    FOREIGN KEY (interviewer_id) REFERENCES users(user_id)
);

-- Interview Results Summary (Different name to avoid conflict)
CREATE TABLE IF NOT EXISTS interview_results_summary (
    id SERIAL PRIMARY KEY,
    booking_id INT NOT NULL,
    total_score DECIMAL(7,2) NOT NULL,
    max_possible_score DECIMAL(7,2) NOT NULL,
    score_percentage DECIMAL(5,2) NOT NULL,
    recommendation VARCHAR(20) NOT NULL,
    interviewer_feedback TEXT,
    submitted_by UUID,
    submitted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (booking_id) REFERENCES interview_bookings(id) ON DELETE CASCADE,
    FOREIGN KEY (submitted_by) REFERENCES users(user_id)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_interview_slots_date ON interview_slots(interview_date, start_time);
CREATE INDEX IF NOT EXISTS idx_interview_slots_interviewer ON interview_slots(interviewer_id);
CREATE INDEX IF NOT EXISTS idx_interview_bookings_application ON interview_bookings(application_id);
CREATE INDEX IF NOT EXISTS idx_interview_bookings_student ON interview_bookings(student_id);
CREATE INDEX IF NOT EXISTS idx_interview_bookings_status ON interview_bookings(booking_status);
CREATE INDEX IF NOT EXISTS idx_interview_scores_detailed_booking ON interview_scores_detailed(booking_id);
CREATE INDEX IF NOT EXISTS idx_interview_results_summary_booking ON interview_results_summary(booking_id);

-- Insert sample data (optional)
INSERT INTO interview_slots (scholarship_id, interviewer_id, interview_date, start_time, end_time, location, created_by)
SELECT 
    1,
    (SELECT user_id FROM users WHERE email LIKE '%@%' LIMIT 1),
    CURRENT_DATE + INTERVAL '7 days',
    '09:00:00',
    '09:30:00',
    'Building A, Room 101',
    (SELECT user_id FROM users WHERE email LIKE '%@%' LIMIT 1)
WHERE EXISTS (SELECT 1 FROM scholarships WHERE scholarship_id = 1)
AND EXISTS (SELECT 1 FROM users WHERE email LIKE '%@%'); 