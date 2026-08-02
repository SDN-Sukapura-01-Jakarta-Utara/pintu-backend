-- Migration: create_form_responses_table
-- Created: 2026-07-17 08:07:00
-- Description: Create form_responses table for storing form submissions

BEGIN;

CREATE TABLE form_responses (
    id BIGSERIAL PRIMARY KEY,
    form_id BIGINT NOT NULL,
    submitted_by_user_id INTEGER,
    submitted_as_role VARCHAR(50), -- "pendidik", "tendik", "murid"
    ip_address VARCHAR(45),
    user_agent TEXT,
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_form_responses_form FOREIGN KEY (form_id) REFERENCES forms(id) ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX idx_form_responses_form_id ON form_responses(form_id);
CREATE INDEX idx_form_responses_submitted_by_user_id ON form_responses(submitted_by_user_id);
CREATE INDEX idx_form_responses_ip_address ON form_responses(ip_address);

-- Index to prevent multiple submissions from same user (checked in application logic)
CREATE INDEX idx_form_responses_form_user ON form_responses(form_id, submitted_by_user_id);

COMMIT;
