-- Migration: create_systems_table
-- Created: 2026-02-05 09:43:22

BEGIN;

CREATE TABLE systems (
  id SERIAL PRIMARY KEY,
  nama VARCHAR(255) NOT NULL UNIQUE,
  description TEXT,
  status VARCHAR(50) DEFAULT 'active',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP,
  created_by_id INTEGER,
  updated_by_id INTEGER
);

COMMIT;
