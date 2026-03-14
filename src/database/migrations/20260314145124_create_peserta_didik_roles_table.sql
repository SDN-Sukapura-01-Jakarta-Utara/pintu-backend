-- Migration: create_peserta_didik_roles_table
-- Created: 2026-03-14 14:51:24
-- Description: Create peserta_didik_roles pivot table for role assignments

BEGIN;

-- Create peserta_didik_roles pivot table
CREATE TABLE peserta_didik_roles (
    id SERIAL PRIMARY KEY,
    peserta_didik_id INTEGER NOT NULL REFERENCES peserta_didik(id) ON DELETE CASCADE,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_peserta_didik_roles_peserta_didik_id ON peserta_didik_roles(peserta_didik_id);
CREATE INDEX idx_peserta_didik_roles_role_id ON peserta_didik_roles(role_id);

COMMIT;
