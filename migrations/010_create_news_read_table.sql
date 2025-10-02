-- Migration: Create news_read table for tracking which news articles users have read
-- Version: 010
-- Description: Add news_read table to track user reading status

-- Create news_read table
CREATE TABLE IF NOT EXISTS news_read (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    news_id UUID NOT NULL REFERENCES news(id) ON DELETE CASCADE,
    read_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure one read record per user per news article
    UNIQUE(user_id, news_id)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_news_read_user_id ON news_read(user_id);
CREATE INDEX IF NOT EXISTS idx_news_read_news_id ON news_read(news_id);
CREATE INDEX IF NOT EXISTS idx_news_read_read_at ON news_read(read_at);

-- Add trigger for updated_at
CREATE OR REPLACE FUNCTION update_news_read_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_news_read_updated_at
    BEFORE UPDATE ON news_read
    FOR EACH ROW
    EXECUTE FUNCTION update_news_read_updated_at();
