-- Migration: add_photo_to_peserta_didik_table
-- Created: 2026-07-05 12:00:00
-- Description: Add photo column to peserta_didik table for student card photo

BEGIN;

ALTER TABLE peserta_didik 
ADD COLUMN IF NOT EXISTS photo VARCHAR(255);

COMMIT;
