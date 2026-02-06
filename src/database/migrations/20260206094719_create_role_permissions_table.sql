-- Migration: create_role_permissions_table
-- Created: 2026-02-06 09:47:19

BEGIN;

-- Create role_permissions pivot table
CREATE TABLE role_permissions (
    id SERIAL PRIMARY KEY,
    role_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create unique constraint
CREATE UNIQUE INDEX idx_role_permissions_unique ON role_permissions(role_id, permission_id);

-- Create indexes
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);

-- Add foreign keys
ALTER TABLE role_permissions
ADD CONSTRAINT fk_role_permissions_role_id
FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE;

ALTER TABLE role_permissions
ADD CONSTRAINT fk_role_permissions_permission_id
FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE;

COMMIT;
