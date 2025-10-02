-- Migration: Add Analytics and Reporting Tables
-- Created: 2025-10-01
-- Description: Adds scholarship statistics and application analytics tables

-- 1. Scholarship Statistics Table
CREATE TABLE IF NOT EXISTS scholarship_statistics (
    stat_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    academic_year VARCHAR(10) NOT NULL,
    scholarship_round VARCHAR(30) NOT NULL,
    total_applications INTEGER NOT NULL,
    approved_applications INTEGER NOT NULL,
    rejected_applications INTEGER NOT NULL,
    total_budget DECIMAL(15,2) NOT NULL,
    allocated_budget DECIMAL(15,2) NOT NULL,
    remaining_budget DECIMAL(15,2) NOT NULL,
    average_amount DECIMAL(12,2) NOT NULL,
    success_rate DECIMAL(5,2) NOT NULL,
    processing_time_avg INTEGER NOT NULL,
    total_faculties INTEGER NOT NULL,
    most_popular_scholarship VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Application Analytics Table
CREATE TABLE IF NOT EXISTS application_analytics (
    analytics_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    application_id INTEGER REFERENCES scholarship_applications(application_id),
    processing_time INTEGER NOT NULL,
    total_steps INTEGER NOT NULL,
    completed_steps INTEGER NOT NULL,
    bottleneck_step VARCHAR(50),
    time_in_each_step JSONB,
    document_upload_time INTEGER,
    review_time INTEGER,
    interview_score DECIMAL(5,2),
    final_score DECIMAL(5,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_scholarship_statistics_year_round ON scholarship_statistics(academic_year, scholarship_round);
CREATE INDEX IF NOT EXISTS idx_application_analytics_application ON application_analytics(application_id);

-- Create unique constraint for statistics
CREATE UNIQUE INDEX IF NOT EXISTS idx_scholarship_statistics_unique ON scholarship_statistics(academic_year, scholarship_round);

COMMENT ON TABLE scholarship_statistics IS 'สถิติทุนการศึกษาแยกตามปีการศึกษาและรอบ';
COMMENT ON TABLE application_analytics IS 'การวิเคราะห์ประสิทธิภาพการดำเนินงานใบสมัคร';
