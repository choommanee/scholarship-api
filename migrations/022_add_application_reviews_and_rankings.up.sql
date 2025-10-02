-- Migration 022: Add Application Reviews and Rankings
-- For reviewing and ranking scholarship applications

-- Application Reviews Table
CREATE TABLE IF NOT EXISTS application_reviews (
    review_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id INTEGER NOT NULL REFERENCES scholarship_applications(application_id) ON DELETE CASCADE,
    reviewer_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    review_stage VARCHAR(50) NOT NULL CHECK (
        review_stage IN ('initial_screening', 'document_verification', 'eligibility_check', 'final_review')
    ),

    -- Review details
    review_status VARCHAR(50) DEFAULT 'pending' CHECK (
        review_status IN ('pending', 'approved', 'rejected', 'needs_revision', 'on_hold')
    ),
    review_score DECIMAL(5,2),         -- Overall score (0-100)
    max_score DECIMAL(5,2) DEFAULT 100,

    -- Detailed scoring
    academic_score DECIMAL(5,2),       -- คะแนนผลการเรียน
    financial_need_score DECIMAL(5,2), -- คะแนนความต้องการทางการเงิน
    extracurricular_score DECIMAL(5,2),-- คะแนนกิจกรรม
    essay_score DECIMAL(5,2),          -- คะแนนเรียงความ
    interview_score DECIMAL(5,2),      -- คะแนนสัมภาษณ์

    -- Comments
    strengths TEXT,                     -- จุดเด่น
    weaknesses TEXT,                    -- จุดที่ต้องปรับปรุง
    comments TEXT,                      -- ความเห็นเพิ่มเติม
    internal_notes TEXT,                -- บันทึกภายใน

    -- Recommendation
    recommendation VARCHAR(50) CHECK (
        recommendation IN ('strongly_recommend', 'recommend', 'neutral', 'not_recommend', 'strongly_not_recommend')
    ),

    -- Flags
    requires_interview BOOLEAN DEFAULT FALSE,
    priority_flag BOOLEAN DEFAULT FALSE,
    red_flag BOOLEAN DEFAULT FALSE,
    red_flag_reason TEXT,

    -- Metadata
    reviewed_at TIMESTAMP,
    time_spent_minutes INTEGER,        -- เวลาที่ใช้ในการรีวิว
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT unique_review_per_stage UNIQUE (application_id, reviewer_id, review_stage)
);

-- Application Rankings Table
CREATE TABLE IF NOT EXISTS application_rankings (
    ranking_id SERIAL PRIMARY KEY,
    round_id INTEGER REFERENCES scholarship_rounds(round_id) ON DELETE CASCADE,
    scholarship_id INTEGER REFERENCES scholarships(scholarship_id) ON DELETE CASCADE,
    application_id INTEGER NOT NULL REFERENCES scholarship_applications(application_id) ON DELETE CASCADE,

    -- Ranking details
    rank_position INTEGER NOT NULL,
    total_applicants INTEGER,

    -- Scores
    total_score DECIMAL(7,2) NOT NULL,
    academic_score DECIMAL(5,2),
    need_score DECIMAL(5,2),
    merit_score DECIMAL(5,2),
    interview_score DECIMAL(5,2),
    committee_score DECIMAL(5,2),

    -- Weighted score
    weighted_score DECIMAL(7,2),
    weight_academic DECIMAL(4,2) DEFAULT 0.30,
    weight_need DECIMAL(4,2) DEFAULT 0.20,
    weight_merit DECIMAL(4,2) DEFAULT 0.20,
    weight_interview DECIMAL(4,2) DEFAULT 0.20,
    weight_committee DECIMAL(4,2) DEFAULT 0.10,

    -- Status
    ranking_status VARCHAR(50) DEFAULT 'provisional' CHECK (
        ranking_status IN ('provisional', 'final', 'approved', 'rejected')
    ),

    -- Award status
    is_awarded BOOLEAN DEFAULT FALSE,
    award_amount DECIMAL(12,2),
    award_type VARCHAR(50),             -- 'full', 'partial', 'conditional'

    -- Waitlist
    is_waitlist BOOLEAN DEFAULT FALSE,
    waitlist_position INTEGER,

    -- Metadata
    ranked_by UUID REFERENCES users(user_id),
    ranked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    approved_by UUID REFERENCES users(user_id),
    approved_at TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT valid_rank_position CHECK (rank_position > 0),
    CONSTRAINT valid_total_score CHECK (total_score >= 0),
    CONSTRAINT valid_weights CHECK (
        weight_academic + weight_need + weight_merit + weight_interview + weight_committee = 1.00
    ),
    CONSTRAINT unique_application_ranking UNIQUE (round_id, scholarship_id, application_id)
);

-- Review Criteria Table (for configurable review criteria)
CREATE TABLE IF NOT EXISTS review_criteria (
    criteria_id SERIAL PRIMARY KEY,
    scholarship_id INTEGER REFERENCES scholarships(scholarship_id) ON DELETE CASCADE,
    criteria_name VARCHAR(100) NOT NULL,
    criteria_description TEXT,
    criteria_type VARCHAR(50) NOT NULL CHECK (
        criteria_type IN ('academic', 'financial', 'merit', 'essay', 'interview', 'other')
    ),
    max_score DECIMAL(5,2) NOT NULL,
    weight DECIMAL(4,2) NOT NULL,      -- Weight in final score (0.00-1.00)
    is_required BOOLEAN DEFAULT TRUE,
    is_active BOOLEAN DEFAULT TRUE,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for application_reviews
CREATE INDEX idx_application_reviews_application_id ON application_reviews(application_id);
CREATE INDEX idx_application_reviews_reviewer_id ON application_reviews(reviewer_id);
CREATE INDEX idx_application_reviews_review_stage ON application_reviews(review_stage);
CREATE INDEX idx_application_reviews_review_status ON application_reviews(review_status);
CREATE INDEX idx_application_reviews_reviewed_at ON application_reviews(reviewed_at);

-- Indexes for application_rankings
CREATE INDEX idx_application_rankings_round_id ON application_rankings(round_id);
CREATE INDEX idx_application_rankings_scholarship_id ON application_rankings(scholarship_id);
CREATE INDEX idx_application_rankings_application_id ON application_rankings(application_id);
CREATE INDEX idx_application_rankings_rank_position ON application_rankings(rank_position);
CREATE INDEX idx_application_rankings_total_score ON application_rankings(total_score DESC);
CREATE INDEX idx_application_rankings_status ON application_rankings(ranking_status);
CREATE INDEX idx_application_rankings_is_awarded ON application_rankings(is_awarded);
CREATE INDEX idx_application_rankings_is_waitlist ON application_rankings(is_waitlist);

-- Indexes for review_criteria
CREATE INDEX idx_review_criteria_scholarship_id ON review_criteria(scholarship_id);
CREATE INDEX idx_review_criteria_is_active ON review_criteria(is_active);
CREATE INDEX idx_review_criteria_display_order ON review_criteria(display_order);

-- Comments
COMMENT ON TABLE application_reviews IS 'Reviews of scholarship applications by reviewers';
COMMENT ON TABLE application_rankings IS 'Rankings of applications based on scores';
COMMENT ON TABLE review_criteria IS 'Configurable criteria for reviewing applications';
COMMENT ON COLUMN application_reviews.review_stage IS 'Stage: initial_screening, document_verification, eligibility_check, final_review';
COMMENT ON COLUMN application_rankings.weighted_score IS 'Final score after applying weights';
