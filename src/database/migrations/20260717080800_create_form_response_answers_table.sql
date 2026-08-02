-- Migration: create_form_response_answers_table
-- Created: 2026-07-17 08:08:00
-- Description: Create form_response_answers table for storing individual answers

BEGIN;

CREATE TABLE form_response_answers (
    id BIGSERIAL PRIMARY KEY,
    response_id BIGINT NOT NULL,
    question_id BIGINT NOT NULL,
    jawaban_text TEXT,
    jawaban_json JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_form_response_answers_response FOREIGN KEY (response_id) REFERENCES form_responses(id) ON DELETE CASCADE,
    CONSTRAINT fk_form_response_answers_question FOREIGN KEY (question_id) REFERENCES form_questions(id) ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX idx_form_response_answers_response_id ON form_response_answers(response_id);
CREATE INDEX idx_form_response_answers_question_id ON form_response_answers(question_id);

COMMIT;
