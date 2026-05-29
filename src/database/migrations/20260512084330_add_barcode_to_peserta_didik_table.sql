-- Migration: add_barcode_to_peserta_didik_table
-- Created: 2026-05-12 08:43:30
-- Description: Add barcode and barcode_generated_at columns to peserta_didik table for attendance system

BEGIN;

ALTER TABLE peserta_didik 
ADD COLUMN IF NOT EXISTS barcode VARCHAR(100),
ADD COLUMN IF NOT EXISTS barcode_generated_at TIMESTAMP;

-- Create index for faster barcode lookup during scan (without UNIQUE constraint)
CREATE INDEX IF NOT EXISTS idx_peserta_didik_barcode ON peserta_didik(barcode) WHERE barcode IS NOT NULL;

COMMIT;
