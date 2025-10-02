-- Migration: Add Import System Tables
-- Created: 2025-10-01
-- Description: Adds import details and data mapping configuration for Excel import

-- 1. Import Details Table (extends import_logs)
CREATE TABLE IF NOT EXISTS import_details (
    detail_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    import_id INTEGER REFERENCES import_logs(import_id) ON DELETE CASCADE,
    row_number INTEGER NOT NULL,
    raw_data JSONB NOT NULL,
    processed_data JSONB,
    status VARCHAR(30) NOT NULL,
    error_message TEXT,
    warnings TEXT[],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP
);

-- 2. Data Mapping Config Table
CREATE TABLE IF NOT EXISTS data_mapping_config (
    config_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source_field VARCHAR(100) NOT NULL,
    target_field VARCHAR(100) NOT NULL,
    data_type VARCHAR(50) NOT NULL,
    transformation_rule TEXT,
    validation_rule TEXT,
    is_required BOOLEAN DEFAULT false,
    default_value TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_import_details_import ON import_details(import_id);
CREATE INDEX IF NOT EXISTS idx_import_details_status ON import_details(status);
CREATE INDEX IF NOT EXISTS idx_data_mapping_config_source ON data_mapping_config(source_field);
CREATE INDEX IF NOT EXISTS idx_data_mapping_config_target ON data_mapping_config(target_field);

-- Insert default data mapping configurations for scholarship import
INSERT INTO data_mapping_config (source_field, target_field, data_type, is_required, default_value) VALUES
('รหัสนักศึกษา', 'student_id', 'string', true, NULL),
('ชื่อ', 'first_name', 'string', true, NULL),
('นามสกุล', 'last_name', 'string', true, NULL),
('อีเมล', 'email', 'email', true, NULL),
('เบอร์โทร', 'phone_number', 'string', false, NULL),
('คณะ', 'faculty', 'string', true, NULL),
('สาขา', 'department', 'string', false, NULL),
('ชั้นปี', 'year', 'integer', true, NULL),
('เกรดเฉลี่ย', 'gpa', 'decimal', true, NULL),
('รายได้ครอบครัว', 'family_income', 'decimal', false, '0'),
('จำนวนพี่น้อง', 'siblings_count', 'integer', false, '0')
ON CONFLICT DO NOTHING;

COMMENT ON TABLE import_details IS 'รายละเอียดการนำเข้าข้อมูลแต่ละแถว';
COMMENT ON TABLE data_mapping_config IS 'การตั้งค่าการแมปข้อมูลจากไฟล์ Excel';
