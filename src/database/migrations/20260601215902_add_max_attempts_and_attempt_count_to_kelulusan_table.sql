-- Migration: add_max_attempts_and_attempt_count_to_kelulusan_table
-- Created: 2026-06-01 21:59:02

BEGIN;

ALTER TABLE kelulusan
ADD COLUMN IF NOT EXISTS max_attempts INTEGER NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS attempt_count INTEGER NOT NULL DEFAULT 0;

COMMIT;
