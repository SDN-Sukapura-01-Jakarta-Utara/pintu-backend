-- Migration: create_pelatih_ekstrakurikuler_table
-- Created: 2026-08-12 12:51:38
-- Description: Create pelatih_ekstrakurikuler junction table for mapping coaches to ekstrakurikuler (permanent assignment)

BEGIN;

CREATE TABLE pelatih_ekstrakurikuler (
    id SERIAL PRIMARY KEY,
    pelatih_id INTEGER NOT NULL,
    ekstrakurikuler_id INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP,
    
    FOREIGN KEY (pelatih_id) REFERENCES pelatih(id) ON DELETE CASCADE,
    FOREIGN KEY (ekstrakurikuler_id) REFERENCES ekstrakurikuler(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Create indexes for better query performance
CREATE INDEX idx_pelatih_ekskul_pelatih ON pelatih_ekstrakurikuler(pelatih_id);
CREATE INDEX idx_pelatih_ekskul_ekstrakurikuler ON pelatih_ekstrakurikuler(ekstrakurikuler_id);
CREATE INDEX idx_pelatih_ekskul_deleted ON pelatih_ekstrakurikuler(deleted_at);

-- Unique constraint: one assignment per pelatih per ekstrakurikuler (when active)
CREATE UNIQUE INDEX idx_unique_pelatih_ekskul_active 
ON pelatih_ekstrakurikuler(pelatih_id, ekstrakurikuler_id) 
WHERE deleted_at IS NULL;

COMMIT;
