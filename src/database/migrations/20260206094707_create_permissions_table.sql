-- Migration: create_permissions_table
-- Created: 2026-02-06 09:47:07

BEGIN;

-- Create permissions table
CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    system VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_permissions_name ON permissions(name);
CREATE INDEX idx_permissions_group_name ON permissions(group_name);
CREATE INDEX idx_permissions_system ON permissions(system);

COMMIT;
