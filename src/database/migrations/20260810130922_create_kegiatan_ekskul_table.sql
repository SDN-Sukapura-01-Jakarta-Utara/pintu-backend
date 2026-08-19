-- Migration: create_kegiatan_ekskul_table
-- Created: 2026-08-10 13:09:22
-- Description: Create kegiatan_ekskul table for ekstrakurikuler activity sessions

BEGIN;

CREATE TABLE kegiatan_ekskul (
    id SERIAL PRIMARY KEY,
    ekstrakurikuler_id INTEGER NOT NULL,
    tahun_pelajaran_id INTEGER NOT NULL,
    tanggal_kegiatan DATE NOT NULL,
    waktu_mulai TIME,
    waktu_selesai TIME,
    materi_kegiatan TEXT NOT NULL,
    foto_kegiatan JSON,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    FOREIGN KEY (ekstrakurikuler_id) REFERENCES ekstrakurikuler(id) ON DELETE CASCADE,
    FOREIGN KEY (tahun_pelajaran_id) REFERENCES tahun_pelajaran(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Create indexes for better query performance
CREATE INDEX idx_kegiatan_ekskul_ekstrakurikuler ON kegiatan_ekskul(ekstrakurikuler_id);
CREATE INDEX idx_kegiatan_ekskul_tahun_pelajaran ON kegiatan_ekskul(tahun_pelajaran_id);
CREATE INDEX idx_kegiatan_ekskul_tanggal ON kegiatan_ekskul(tanggal_kegiatan);
CREATE INDEX idx_kegiatan_ekskul_deleted ON kegiatan_ekskul(deleted_at);

COMMIT;
