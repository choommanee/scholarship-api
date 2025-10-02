-- Migration 029: Reporting System Tables

-- 1. Report Templates (แม่แบบรายงาน)
CREATE TABLE IF NOT EXISTS report_templates (
    template_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- ข้อมูลแม่แบบ
    template_name VARCHAR(255) NOT NULL,
    template_code VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,

    -- ประเภทรายงาน
    report_type VARCHAR(50) NOT NULL CHECK (
        report_type IN (
            'scholarship_allocation',    -- รายงานการจัดสรรทุน
            'student_performance',        -- รายงานผลการเรียนนักศึกษา
            'budget_summary',             -- สรุปงบประมาณ
            'disbursement_summary',       -- สรุปการโอนเงิน
            'application_statistics',     -- สถิติการสมัคร
            'interview_summary',          -- สรุปการสัมภาษณ์
            'advisor_workload',           -- ภาระงานอาจารย์ที่ปรึกษา
            'custom'                      -- รายงานแบบกำหนดเอง
        )
    ),

    -- การตั้งค่ารายงาน
    template_config JSONB,              -- Configuration JSON
    query_template TEXT,                -- SQL Query Template

    -- รูปแบบผลลัพธ์
    output_format VARCHAR(50)[],        -- ['pdf', 'excel', 'csv']
    default_format VARCHAR(50) DEFAULT 'pdf',

    -- การแสดงผล
    columns JSONB,                      -- คอลัมน์ที่แสดง
    filters JSONB,                      -- ตัวกรองข้อมูล
    sorting JSONB,                      -- การเรียงลำดับ
    grouping JSONB,                     -- การจัดกลุ่ม

    -- สิทธิ์การใช้งาน
    accessible_roles VARCHAR(50)[],     -- บทบาทที่สามารถใช้รายงานนี้ได้

    -- สถานะ
    is_active BOOLEAN DEFAULT TRUE,
    is_system BOOLEAN DEFAULT FALSE,    -- รายงานระบบ (ห้ามลบ)

    -- ผู้สร้าง
    created_by UUID REFERENCES users(user_id),
    updated_by UUID REFERENCES users(user_id),

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Generated Reports (รายงานที่สร้างแล้ว)
CREATE TABLE IF NOT EXISTS generated_reports (
    report_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID REFERENCES report_templates(template_id) ON DELETE SET NULL,

    -- ข้อมูลรายงาน
    report_name VARCHAR(255) NOT NULL,
    report_description TEXT,

    -- ช่วงเวลา
    report_period VARCHAR(100),         -- "ปีการศึกษา 2567", "รอบที่ 1/2567"
    start_date DATE,
    end_date DATE,

    -- ตัวกรอง
    filter_params JSONB,                -- พารามิเตอร์ที่ใช้กรอง

    -- ไฟล์รายงาน
    file_path TEXT,
    file_name VARCHAR(255),
    file_size INTEGER,                  -- bytes
    file_format VARCHAR(50),            -- pdf, excel, csv
    mime_type VARCHAR(100),

    -- สถิติ
    total_records INTEGER,
    total_pages INTEGER,

    -- สถานะ
    status VARCHAR(50) DEFAULT 'generating' CHECK (
        status IN ('generating', 'completed', 'failed', 'expired')
    ),

    -- การหมดอายุ
    expires_at TIMESTAMP,               -- รายงานหมดอายุเมื่อไร
    is_expired BOOLEAN DEFAULT FALSE,

    -- ข้อผิดพลาด
    error_message TEXT,

    -- ผู้สร้าง
    generated_by UUID NOT NULL REFERENCES users(user_id),

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. Report Schedules (ตารางการสร้างรายงานอัตโนมัติ)
CREATE TABLE IF NOT EXISTS report_schedules (
    schedule_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID NOT NULL REFERENCES report_templates(template_id) ON DELETE CASCADE,

    -- การตั้งเวลา
    schedule_name VARCHAR(255) NOT NULL,
    description TEXT,

    -- ความถี่
    frequency VARCHAR(50) NOT NULL CHECK (
        frequency IN ('daily', 'weekly', 'monthly', 'quarterly', 'yearly', 'custom')
    ),

    -- Cron Expression (สำหรับ custom)
    cron_expression VARCHAR(100),

    -- วันที่ต่อไป
    next_run_date TIMESTAMP,
    last_run_date TIMESTAMP,

    -- ตัวกรอง
    default_filters JSONB,

    -- ผู้รับรายงาน
    recipients JSONB,                   -- Array ของ email addresses

    -- สถานะ
    is_active BOOLEAN DEFAULT TRUE,

    -- ผู้สร้าง
    created_by UUID REFERENCES users(user_id),

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. Report Access Logs (Log การเข้าถึงรายงาน)
CREATE TABLE IF NOT EXISTS report_access_logs (
    log_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_id UUID REFERENCES generated_reports(report_id) ON DELETE SET NULL,
    template_id UUID REFERENCES report_templates(template_id) ON DELETE SET NULL,

    -- ผู้เข้าถึง
    user_id UUID REFERENCES users(user_id),

    -- การกระทำ
    action VARCHAR(50) NOT NULL,        -- view, download, export, delete

    -- ข้อมูลเพิ่มเติม
    ip_address VARCHAR(45),
    user_agent TEXT,

    -- Metadata
    accessed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 5. Dashboard Widgets (วิดเจ็ตในหน้า Dashboard)
CREATE TABLE IF NOT EXISTS dashboard_widgets (
    widget_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- ข้อมูลวิดเจ็ต
    widget_name VARCHAR(255) NOT NULL,
    widget_type VARCHAR(50) NOT NULL,   -- chart, table, stat, progress
    description TEXT,

    -- Data Source
    data_source VARCHAR(50),            -- api, query, static
    query_template TEXT,                -- SQL Query
    api_endpoint VARCHAR(255),

    -- การตั้งค่า
    config JSONB,                       -- Configuration สำหรับ chart/table

    -- การแสดงผล
    display_order INTEGER,
    width INTEGER DEFAULT 6,            -- Grid width (1-12)
    height INTEGER DEFAULT 4,           -- Grid height

    -- สิทธิ์
    accessible_roles VARCHAR(50)[],

    -- Refresh
    refresh_interval INTEGER,           -- วินาที (null = ไม่ refresh)
    cache_duration INTEGER,             -- วินาที

    -- สถานะ
    is_active BOOLEAN DEFAULT TRUE,

    -- ผู้สร้าง
    created_by UUID REFERENCES users(user_id),

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for report_templates
CREATE INDEX idx_report_templates_report_type ON report_templates(report_type);
CREATE INDEX idx_report_templates_is_active ON report_templates(is_active);
CREATE INDEX idx_report_templates_created_by ON report_templates(created_by);

-- Indexes for generated_reports
CREATE INDEX idx_generated_reports_template_id ON generated_reports(template_id);
CREATE INDEX idx_generated_reports_generated_by ON generated_reports(generated_by);
CREATE INDEX idx_generated_reports_status ON generated_reports(status);
CREATE INDEX idx_generated_reports_start_date ON generated_reports(start_date);
CREATE INDEX idx_generated_reports_end_date ON generated_reports(end_date);
CREATE INDEX idx_generated_reports_is_expired ON generated_reports(is_expired);

-- Indexes for report_schedules
CREATE INDEX idx_report_schedules_template_id ON report_schedules(template_id);
CREATE INDEX idx_report_schedules_next_run_date ON report_schedules(next_run_date);
CREATE INDEX idx_report_schedules_is_active ON report_schedules(is_active);

-- Indexes for report_access_logs
CREATE INDEX idx_report_access_logs_report_id ON report_access_logs(report_id);
CREATE INDEX idx_report_access_logs_user_id ON report_access_logs(user_id);
CREATE INDEX idx_report_access_logs_action ON report_access_logs(action);
CREATE INDEX idx_report_access_logs_accessed_at ON report_access_logs(accessed_at);

-- Indexes for dashboard_widgets
CREATE INDEX idx_dashboard_widgets_widget_type ON dashboard_widgets(widget_type);
CREATE INDEX idx_dashboard_widgets_display_order ON dashboard_widgets(display_order);
CREATE INDEX idx_dashboard_widgets_is_active ON dashboard_widgets(is_active);

-- Comments
COMMENT ON TABLE report_templates IS 'แม่แบบรายงานต่างๆ ในระบบ';
COMMENT ON TABLE generated_reports IS 'รายงานที่ถูกสร้างแล้ว';
COMMENT ON TABLE report_schedules IS 'ตารางการสร้างรายงานอัตโนมัติ';
COMMENT ON TABLE report_access_logs IS 'บันทึกการเข้าถึงรายงาน';
COMMENT ON TABLE dashboard_widgets IS 'วิดเจ็ตสำหรับหน้า Dashboard';
