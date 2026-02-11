-- Migration: create_sarana_prasaranas_table
-- Created: 2026-02-11 12:58:04
-- Description: Create sarana_prasaranas table for school facilities management

BEGIN;

CREATE TABLE sarana_prasaranas (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    foto VARCHAR(512) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP
);

COMMIT;
