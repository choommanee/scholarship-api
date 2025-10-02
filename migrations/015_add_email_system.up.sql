-- Migration: Add Email System Tables
-- Created: 2025-10-01
-- Description: Adds email queue and email templates for notification system

-- 1. Email Templates Table
CREATE TABLE IF NOT EXISTS email_templates (
    template_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    template_name VARCHAR(100) NOT NULL,
    subject VARCHAR(500) NOT NULL,
    body TEXT NOT NULL,
    variables JSONB,
    template_type VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Email Queue Table
CREATE TABLE IF NOT EXISTS email_queue (
    queue_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    recipient_email VARCHAR(255) NOT NULL,
    recipient_name VARCHAR(255),
    sender_email VARCHAR(255) NOT NULL,
    subject VARCHAR(500) NOT NULL,
    body TEXT NOT NULL,
    template_id UUID REFERENCES email_templates(template_id),
    priority INTEGER DEFAULT 5,
    status VARCHAR(30) DEFAULT 'pending',
    sent_at TIMESTAMP,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_email_queue_status ON email_queue(status, priority DESC);
CREATE INDEX IF NOT EXISTS idx_email_queue_created ON email_queue(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_email_templates_type ON email_templates(template_type, is_active);

-- Insert default email templates
INSERT INTO email_templates (template_name, subject, body, template_type, variables) VALUES
(
    'application_submitted',
    'ยืนยันการส่งใบสมัครทุนการศึกษา - {{scholarship_name}}',
    'เรียน {{student_name}}<br><br>ระบบได้รับใบสมัครทุน {{scholarship_name}} ของคุณเรียบร้อยแล้ว<br>หมายเลขใบสมัคร: {{application_id}}<br><br>ขอบคุณครับ',
    'application',
    '{"student_name": "string", "scholarship_name": "string", "application_id": "string"}'::jsonb
),
(
    'interview_scheduled',
    'แจ้งตารางสัมภาษณ์ทุนการศึกษา',
    'เรียน {{student_name}}<br><br>คุณได้รับการนัดสัมภาษณ์ทุนการศึกษา<br>วันที่: {{interview_date}}<br>เวลา: {{interview_time}}<br>สถานที่: {{location}}<br><br>กรุณามาตรงเวลา',
    'interview',
    '{"student_name": "string", "interview_date": "string", "interview_time": "string", "location": "string"}'::jsonb
),
(
    'scholarship_approved',
    'ยินดีด้วย! คุณได้รับทุนการศึกษา',
    'เรียน {{student_name}}<br><br>ยินดีด้วย! คุณได้รับทุน {{scholarship_name}}<br>จำนวนเงิน: {{amount}} บาท<br><br>กรุณาติดต่อเจ้าหน้าที่เพื่อดำเนินการต่อไป',
    'allocation',
    '{"student_name": "string", "scholarship_name": "string", "amount": "string"}'::jsonb
),
(
    'payment_disbursed',
    'แจ้งการจ่ายเงินทุนการศึกษา',
    'เรียน {{student_name}}<br><br>เงินทุนการศึกษาจำนวน {{amount}} บาท ได้โอนเข้าบัญชีของคุณเรียบร้อยแล้ว<br>เลขที่อ้างอิง: {{reference}}<br><br>ขอบคุณครับ',
    'payment',
    '{"student_name": "string", "amount": "string", "reference": "string"}'::jsonb
)
ON CONFLICT DO NOTHING;

COMMENT ON TABLE email_templates IS 'แม่แบบอีเมลสำหรับส่งการแจ้งเตือน';
COMMENT ON TABLE email_queue IS 'คิวการส่งอีเมล';
