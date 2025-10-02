-- Migration: Add Background Jobs System
-- Created: 2025-10-01
-- Description: Adds job queue and background tasks for async processing

-- 1. Job Queue Table
CREATE TABLE IF NOT EXISTS job_queue (
    job_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    job_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    priority INTEGER DEFAULT 5,
    status VARCHAR(30) DEFAULT 'pending',
    attempts INTEGER DEFAULT 0,
    max_attempts INTEGER DEFAULT 3,
    scheduled_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    error_message TEXT
);

-- 2. Background Tasks Table
CREATE TABLE IF NOT EXISTS background_tasks (
    task_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_name VARCHAR(100) NOT NULL,
    task_type VARCHAR(50) NOT NULL,
    schedule VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    last_run TIMESTAMP,
    next_run TIMESTAMP,
    run_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_job_queue_status ON job_queue(status, priority DESC, scheduled_at);
CREATE INDEX IF NOT EXISTS idx_job_queue_type ON job_queue(job_type);
CREATE INDEX IF NOT EXISTS idx_background_tasks_schedule ON background_tasks(is_active, next_run);

-- Add trigger for updated_at
CREATE TRIGGER update_background_tasks_updated_at
    BEFORE UPDATE ON background_tasks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert default background tasks
INSERT INTO background_tasks (task_name, task_type, schedule, is_active, next_run) VALUES
('Send Email Queue', 'email', '*/5 * * * *', true, CURRENT_TIMESTAMP),
('Process Payment Reminders', 'payment', '0 9 * * *', true, CURRENT_TIMESTAMP + INTERVAL '1 day'),
('Generate Daily Statistics', 'analytics', '0 23 * * *', true, CURRENT_TIMESTAMP),
('Cleanup Old Files', 'maintenance', '0 2 * * 0', true, CURRENT_TIMESTAMP + INTERVAL '7 days'),
('Application Deadline Reminders', 'notification', '0 8 * * *', true, CURRENT_TIMESTAMP + INTERVAL '1 day')
ON CONFLICT DO NOTHING;

COMMENT ON TABLE job_queue IS 'คิวงานพื้นหลังที่ต้องประมวลผลแบบ async';
COMMENT ON TABLE background_tasks IS 'งานที่รันตามตารางเวลา (cron jobs)';
