-- Migration: create_jumbotron_table
-- Created: 2026-02-10 07:57:12
-- Description: Create jumbotron table for banner/image management

BEGIN;

CREATE TABLE jumbotron (
    id SERIAL PRIMARY KEY,
    file VARCHAR(512) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP
);

COMMIT;
