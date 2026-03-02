-- Migration: create_kepegawaian_table
-- Created: 2026-03-01 21:55:20
-- Description: Create kepegawaian table for staff/employee data management

BEGIN;

CREATE TABLE kepegawaian (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(255) NOT NULL,
    username VARCHAR(100),
    password VARCHAR(255),
    nip VARCHAR(30),
    nkki VARCHAR(15),
    foto VARCHAR(512),
    kategori VARCHAR(20),
    jabatan VARCHAR(100),
    rombel_guru_kelas_id INTEGER,
    rombel_bidang_studi JSONB DEFAULT '[]'::jsonb,
    kk TEXT,
    akta_lahir TEXT,
    ktp TEXT,
    ijazah_sd TEXT,
    ijazah_smp TEXT,
    ijazah_sma TEXT,
    ijazah_s1 TEXT,
    ijazah_s2 TEXT,
    ijazah_s3 TEXT,
    sertifikat_pendidik TEXT,
    sertifikat_lainnya JSONB DEFAULT '[]'::jsonb,
    sk TEXT,
    dokumen_lainnya JSONB DEFAULT '[]'::jsonb,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_kepegawaian_rombel FOREIGN KEY (rombel_guru_kelas_id) REFERENCES rombel(id) ON DELETE SET NULL
);

-- Create indexes for better query performance
CREATE INDEX idx_kepegawaian_nama ON kepegawaian(nama);
CREATE INDEX idx_kepegawaian_rombel_guru_kelas_id ON kepegawaian(rombel_guru_kelas_id);

COMMIT;
