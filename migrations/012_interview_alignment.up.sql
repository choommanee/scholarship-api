-- Align interview_slots and interview_bookings with enhanced schema

ALTER TABLE interview_slots
    ADD COLUMN IF NOT EXISTS building VARCHAR(100),
    ADD COLUMN IF NOT EXISTS floor VARCHAR(50),
    ADD COLUMN IF NOT EXISTS room VARCHAR(100),
    ADD COLUMN IF NOT EXISTS slot_type VARCHAR(20) DEFAULT 'individual',
    ADD COLUMN IF NOT EXISTS duration_minutes INT DEFAULT 30,
    ADD COLUMN IF NOT EXISTS preparation_time INT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS notes TEXT;

-- Ensure existing rows have sensible defaults
UPDATE interview_slots
SET
    slot_type = COALESCE(slot_type, 'individual'),
    duration_minutes = COALESCE(duration_minutes, 30),
    preparation_time = COALESCE(preparation_time, 0)
WHERE TRUE;

ALTER TABLE interview_bookings
    ADD COLUMN IF NOT EXISTS student_notes TEXT,
    ADD COLUMN IF NOT EXISTS officer_notes TEXT,
    ADD COLUMN IF NOT EXISTS reminder_sent_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS check_in_time TIMESTAMP,
    ADD COLUMN IF NOT EXISTS check_out_time TIMESTAMP,
    ADD COLUMN IF NOT EXISTS actual_duration_minutes INT,
    ADD COLUMN IF NOT EXISTS rescheduled_from_slot_id INT REFERENCES interview_slots(id),
    ADD COLUMN IF NOT EXISTS rescheduled_to_slot_id INT REFERENCES interview_slots(id),
    ADD COLUMN IF NOT EXISTS cancellation_reason TEXT;

-- Preserve legacy notes into student_notes if present
UPDATE interview_bookings
SET student_notes = notes
WHERE student_notes IS NULL AND notes IS NOT NULL;
