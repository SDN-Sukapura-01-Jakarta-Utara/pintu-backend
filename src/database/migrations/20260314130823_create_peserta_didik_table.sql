-- Migration: create_peserta_didik_table
-- Created: 2026-03-14 13:08:23
-- Description: Create peserta_didik table for student data management

BEGIN;

CREATE TABLE peserta_didik (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(255) NOT NULL,
    nis VARCHAR(15) NOT NULL,
    jenis_kelamin VARCHAR(1) NOT NULL,
    nisn VARCHAR(25) NOT NULL,
    tempat_lahir VARCHAR(100),
    tanggal_lahir DATE,
    nik VARCHAR(25),
    agama VARCHAR(20),
    alamat TEXT,
    rt VARCHAR(5),
    rw VARCHAR(5),
    kelurahan VARCHAR(100),
    kecamatan VARCHAR(100),
    kode_pos VARCHAR(10),
    nama_ayah VARCHAR(255),
    nama_ibu VARCHAR(255),
    rombel_id INTEGER,
    tahun_pelajaran_id INTEGER,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    username VARCHAR(100),
    password VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_peserta_didik_rombel FOREIGN KEY (rombel_id) REFERENCES rombel(id) ON DELETE SET NULL,
    CONSTRAINT fk_peserta_didik_tahun_pelajaran FOREIGN KEY (tahun_pelajaran_id) REFERENCES tahun_pelajaran(id) ON DELETE SET NULL
);

-- Create indexes for better query performance
CREATE INDEX idx_peserta_didik_nama ON peserta_didik(nama);
CREATE INDEX idx_peserta_didik_nis ON peserta_didik(nis);
CREATE INDEX idx_peserta_didik_nisn ON peserta_didik(nisn);
CREATE INDEX idx_peserta_didik_rombel_id ON peserta_didik(rombel_id);
CREATE INDEX idx_peserta_didik_tahun_pelajaran_id ON peserta_didik(tahun_pelajaran_id);
CREATE INDEX idx_peserta_didik_status ON peserta_didik(status);

COMMIT;
