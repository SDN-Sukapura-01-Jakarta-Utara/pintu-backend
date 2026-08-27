-- Migration: drop_kegiatan_ekskul_foreign_keys
-- Created: 2026-08-28 06:57:31
-- Description: Drop foreign key constraints on created_by_id and updated_by_id from kegiatan_ekskul table to allow both users and pelatih as creators

BEGIN;

-- Drop foreign key constraints
ALTER TABLE kegiatan_ekskul 
    DROP CONSTRAINT IF EXISTS kegiatan_ekskul_created_by_id_fkey;

ALTER TABLE kegiatan_ekskul 
    DROP CONSTRAINT IF EXISTS kegiatan_ekskul_updated_by_id_fkey;

COMMIT;
