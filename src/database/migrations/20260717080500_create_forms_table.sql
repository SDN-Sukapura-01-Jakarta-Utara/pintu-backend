-- Migration: create_forms_table
-- Created: 2026-07-17 08:05:00
-- Description: Create forms table for dynamic form management

BEGIN;

CREATE TABLE forms (
    id BIGSERIAL PRIMARY KEY,
    judul VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    deskripsi TEXT,
    created_by_user_id INTEGER NOT NULL,
    created_by_role VARCHAR(50),
    is_active BOOLEAN DEFAULT TRUE,
    max_responses INTEGER,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    access_type VARCHAR(50) DEFAULT 'public',
    target_user_types JSONB,
    rombel_ids JSONB,
    allow_multiple_responses BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_forms_created_by_user_id ON forms(created_by_user_id);
CREATE INDEX idx_forms_is_active ON forms(is_active);
CREATE INDEX idx_forms_start_date ON forms(start_date);
CREATE INDEX idx_forms_end_date ON forms(end_date);
CREATE INDEX idx_forms_slug ON forms(slug);
CREATE INDEX idx_forms_access_type ON forms(access_type);

COMMIT;
