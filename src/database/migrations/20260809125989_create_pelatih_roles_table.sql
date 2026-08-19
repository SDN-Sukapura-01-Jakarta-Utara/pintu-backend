-- Migration: create_pelatih_roles_table
-- Created: 2026-08-17 13:00:00
-- Description: Create pelatih_roles pivot table for role assignments

BEGIN;

-- Create pelatih_roles pivot table
CREATE TABLE pelatih_roles (
    id SERIAL PRIMARY KEY,
    pelatih_id INTEGER NOT NULL REFERENCES pelatih(id) ON DELETE CASCADE,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_pelatih_roles_pelatih_id ON pelatih_roles(pelatih_id);
CREATE INDEX idx_pelatih_roles_role_id ON pelatih_roles(role_id);

-- Unique constraint to prevent duplicate role assignments
CREATE UNIQUE INDEX idx_unique_pelatih_role ON pelatih_roles(pelatih_id, role_id);

COMMIT;
