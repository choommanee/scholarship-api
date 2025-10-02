-- Revert interview slot and booking enhancements

ALTER TABLE interview_bookings
    DROP COLUMN IF EXISTS cancellation_reason,
    DROP COLUMN IF EXISTS rescheduled_to_slot_id,
    DROP COLUMN IF EXISTS rescheduled_from_slot_id,
    DROP COLUMN IF EXISTS actual_duration_minutes,
    DROP COLUMN IF EXISTS check_out_time,
    DROP COLUMN IF EXISTS check_in_time,
    DROP COLUMN IF EXISTS reminder_sent_at,
    DROP COLUMN IF EXISTS officer_notes,
    DROP COLUMN IF EXISTS student_notes;

ALTER TABLE interview_slots
    DROP COLUMN IF EXISTS notes,
    DROP COLUMN IF EXISTS preparation_time,
    DROP COLUMN IF EXISTS duration_minutes,
    DROP COLUMN IF EXISTS slot_type,
    DROP COLUMN IF EXISTS room,
    DROP COLUMN IF EXISTS floor,
    DROP COLUMN IF EXISTS building;
