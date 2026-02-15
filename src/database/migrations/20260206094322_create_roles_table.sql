-- Migration: create_roles_table
-- Created: 2026-02-06 09:43:22

BEGIN;

-- Create roles table
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    system_id INTEGER,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_roles_name ON roles(name);
CREATE INDEX idx_roles_system_id ON roles(system_id);

ALTER TABLE roles
ADD CONSTRAINT fk_roles_system_id
FOREIGN KEY (system_id) REFERENCES systems(id) ON DELETE SET NULL;

COMMIT;
