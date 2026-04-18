-- Migration: create_prestasi_table
-- Created: 2026-04-15 08:51:18
-- Description: Create prestasi table for achievement management

BEGIN;

CREATE TABLE prestasi (
    id SERIAL PRIMARY KEY,
    peserta_didik_id INTEGER,
    jenis VARCHAR(100) NOT NULL,
    nama_grup VARCHAR(255),
    nama_prestasi VARCHAR(255) NOT NULL,
    tingkat_prestasi VARCHAR(100),
    penyelenggara VARCHAR(255),
    tanggal_lomba DATE NOT NULL,
    juara VARCHAR(100) NOT NULL,
    keterangan TEXT,
    foto JSONB DEFAULT '[]',
    ekstrakurikuler_id INTEGER,
    tahun_pelajaran_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    CONSTRAINT fk_prestasi_ekstrakurikuler FOREIGN KEY (ekstrakurikuler_id) REFERENCES ekstrakurikuler(id) ON DELETE SET NULL,
    CONSTRAINT fk_prestasi_peserta_didik FOREIGN KEY (peserta_didik_id) REFERENCES peserta_didik(id) ON DELETE SET NULL,
    CONSTRAINT fk_prestasi_tahun_pelajaran FOREIGN KEY (tahun_pelajaran_id) REFERENCES tahun_pelajaran(id) ON DELETE SET NULL
);

-- Create indexes for better query performance
CREATE INDEX idx_prestasi_peserta_didik_id ON prestasi(peserta_didik_id);
CREATE INDEX idx_prestasi_ekstrakurikuler_id ON prestasi(ekstrakurikuler_id);
CREATE INDEX idx_prestasi_tahun_pelajaran_id ON prestasi(tahun_pelajaran_id);

COMMIT;