-- Migration: create_kepegawaian_roles_table
-- Created: 2026-03-02 12:00:00
-- Description: Create kepegawaian_roles pivot table for role assignments

BEGIN;

-- Create kepegawaian_roles pivot table
CREATE TABLE kepegawaian_roles (
    id SERIAL PRIMARY KEY,
    kepegawaian_id INTEGER NOT NULL REFERENCES kepegawaian(id) ON DELETE CASCADE,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_kepegawaian_roles_kepegawaian_id ON kepegawaian_roles(kepegawaian_id);
CREATE INDEX idx_kepegawaian_roles_role_id ON kepegawaian_roles(role_id);

COMMIT;
