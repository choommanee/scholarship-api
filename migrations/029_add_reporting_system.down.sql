-- Migration 029 Down

DROP INDEX IF EXISTS idx_dashboard_widgets_is_active;
DROP INDEX IF EXISTS idx_dashboard_widgets_display_order;
DROP INDEX IF EXISTS idx_dashboard_widgets_widget_type;
DROP INDEX IF EXISTS idx_report_access_logs_accessed_at;
DROP INDEX IF EXISTS idx_report_access_logs_action;
DROP INDEX IF EXISTS idx_report_access_logs_user_id;
DROP INDEX IF EXISTS idx_report_access_logs_report_id;
DROP INDEX IF EXISTS idx_report_schedules_is_active;
DROP INDEX IF EXISTS idx_report_schedules_next_run_date;
DROP INDEX IF EXISTS idx_report_schedules_template_id;
DROP INDEX IF EXISTS idx_generated_reports_is_expired;
DROP INDEX IF EXISTS idx_generated_reports_end_date;
DROP INDEX IF EXISTS idx_generated_reports_start_date;
DROP INDEX IF EXISTS idx_generated_reports_status;
DROP INDEX IF EXISTS idx_generated_reports_generated_by;
DROP INDEX IF EXISTS idx_generated_reports_template_id;
DROP INDEX IF EXISTS idx_report_templates_created_by;
DROP INDEX IF EXISTS idx_report_templates_is_active;
DROP INDEX IF EXISTS idx_report_templates_report_type;

DROP TABLE IF EXISTS dashboard_widgets;
DROP TABLE IF EXISTS report_access_logs;
DROP TABLE IF EXISTS report_schedules;
DROP TABLE IF EXISTS generated_reports;
DROP TABLE IF EXISTS report_templates;
