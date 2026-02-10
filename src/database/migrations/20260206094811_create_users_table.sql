-- Migration: create_users_table
-- Created: 2026-02-06 07:45:47

BEGIN;

-- Create status enum type
CREATE TYPE user_status AS ENUM ('active', 'inactive', 'suspended');

-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(255) NOT NULL,
    username VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role_id INTEGER,
    accessible_system JSONB DEFAULT '[]'::jsonb,
    status user_status DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_role_id ON users(role_id);

ALTER TABLE users
ADD CONSTRAINT fk_users_role_id
FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE SET NULL;

COMMIT;
