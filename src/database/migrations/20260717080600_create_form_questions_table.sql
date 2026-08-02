-- Migration: create_form_questions_table
-- Created: 2026-07-17 08:06:00
-- Description: Create form_questions table for storing form questions

BEGIN;

CREATE TYPE form_question_type AS ENUM (
    'text',
    'textarea', 
    'number',
    'email',
    'phone',
    'radio',
    'checkbox',
    'select',
    'file',
    'date',
    'time',
    'datetime'
);

CREATE TABLE form_questions (
    id BIGSERIAL PRIMARY KEY,
    form_id BIGINT NOT NULL,
    urutan INTEGER NOT NULL,
    label TEXT NOT NULL,
    placeholder VARCHAR(255),
    tipe form_question_type NOT NULL,
    is_required BOOLEAN DEFAULT FALSE,
    options JSONB,
    validation_rules JSONB,
    file_config JSONB,
    dokumen VARCHAR(255),
    link TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_form_questions_form FOREIGN KEY (form_id) REFERENCES forms(id) ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX idx_form_questions_form_id ON form_questions(form_id);
CREATE INDEX idx_form_questions_urutan ON form_questions(form_id, urutan);
CREATE INDEX idx_form_questions_dokumen ON form_questions(dokumen);

COMMIT;
