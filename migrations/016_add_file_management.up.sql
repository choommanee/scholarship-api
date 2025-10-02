-- Migration: Add File Management Tables
-- Created: 2025-10-01
-- Description: Adds file storage, document versions, and file access logs

-- 1. File Storage Table
CREATE TABLE IF NOT EXISTS file_storage (
    file_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    original_name VARCHAR(255) NOT NULL,
    stored_name VARCHAR(255) NOT NULL,
    stored_path VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    file_hash VARCHAR(64) NOT NULL,
    uploaded_by UUID REFERENCES users(user_id),
    related_table VARCHAR(50),
    related_id UUID,
    storage_type VARCHAR(30) DEFAULT 'local',
    access_level VARCHAR(30) DEFAULT 'private',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Document Versions Table
CREATE TABLE IF NOT EXISTS document_versions (
    version_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    file_id UUID REFERENCES file_storage(file_id) ON DELETE CASCADE,
    version_number INTEGER NOT NULL,
    change_description TEXT,
    uploaded_by UUID REFERENCES users(user_id),
    file_size BIGINT NOT NULL,
    is_current BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    replaced_at TIMESTAMP
);

-- 3. File Access Logs Table
CREATE TABLE IF NOT EXISTS file_access_logs (
    access_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    file_id UUID REFERENCES file_storage(file_id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(user_id),
    access_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    action VARCHAR(50) NOT NULL,
    ip_address INET,
    user_agent TEXT,
    success BOOLEAN DEFAULT true
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_file_storage_uploaded_by ON file_storage(uploaded_by);
CREATE INDEX IF NOT EXISTS idx_file_storage_related ON file_storage(related_table, related_id);
CREATE INDEX IF NOT EXISTS idx_file_storage_hash ON file_storage(file_hash);
CREATE INDEX IF NOT EXISTS idx_document_versions_file ON document_versions(file_id, version_number DESC);
CREATE INDEX IF NOT EXISTS idx_document_versions_current ON document_versions(file_id, is_current);
CREATE INDEX IF NOT EXISTS idx_file_access_logs_file ON file_access_logs(file_id, access_time DESC);
CREATE INDEX IF NOT EXISTS idx_file_access_logs_user ON file_access_logs(user_id, access_time DESC);

-- Add trigger for updated_at
CREATE TRIGGER update_file_storage_updated_at
    BEFORE UPDATE ON file_storage
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE file_storage IS 'การจัดเก็บไฟล์และเอกสาร';
COMMENT ON TABLE document_versions IS 'เวอร์ชันของเอกสารที่มีการแก้ไข';
COMMENT ON TABLE file_access_logs IS 'บันทึกการเข้าถึงไฟล์';
