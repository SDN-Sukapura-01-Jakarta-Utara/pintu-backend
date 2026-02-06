-- Migration: create_roles_table
-- Created: 2026-02-06 09:43:22

BEGIN;

-- Create roles table
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    system VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_roles_name ON roles(name);
CREATE INDEX idx_roles_system ON roles(system);

COMMIT;
