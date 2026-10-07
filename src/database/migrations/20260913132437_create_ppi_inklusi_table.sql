-- Migration: create_ppi_inklusi_table
-- Created: 2026-09-13 13:24:37
-- Description: Create ppi_inklusi table for Program Pembelajaran Individual (Individual Learning Program)

BEGIN;

CREATE TABLE ppi_inklusi (
    id SERIAL PRIMARY KEY,
    anak_inklusi_rombel_id INTEGER NOT NULL,
    tahun_pelajaran_id INTEGER NOT NULL,
    aspek_pembelajaran VARCHAR(100) NOT NULL,
    nama_program VARCHAR(255) NOT NULL,
    tujuan_pembelajaran TEXT NOT NULL,
    strategi_pembelajaran TEXT NOT NULL,
    target_waktu VARCHAR(100),
    capaian_saat_ini TEXT,
    status_capaian VARCHAR(100) DEFAULT 'Belum Dimulai',
    guru_pembuat_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP,
    
    CONSTRAINT fk_ppi_inklusi_anak_inklusi_rombel 
        FOREIGN KEY (anak_inklusi_rombel_id) 
        REFERENCES anak_inklusi_rombel(id) 
        ON DELETE CASCADE,
        
    CONSTRAINT fk_ppi_inklusi_tahun_pelajaran 
        FOREIGN KEY (tahun_pelajaran_id) 
        REFERENCES tahun_pelajaran(id) 
        ON DELETE CASCADE,
        
    CONSTRAINT fk_ppi_inklusi_guru_pembuat 
        FOREIGN KEY (guru_pembuat_id) 
        REFERENCES kepegawaian(id) 
        ON DELETE SET NULL
);

-- Create indexes for better query performance
CREATE INDEX idx_ppi_inklusi_anak_inklusi_rombel_id ON ppi_inklusi(anak_inklusi_rombel_id);
CREATE INDEX idx_ppi_inklusi_aspek_pembelajaran ON ppi_inklusi(aspek_pembelajaran);
CREATE INDEX idx_ppi_inklusi_tahun_pelajaran_id ON ppi_inklusi(tahun_pelajaran_id);
CREATE INDEX idx_ppi_inklusi_status_capaian ON ppi_inklusi(status_capaian);
CREATE INDEX idx_ppi_inklusi_guru_pembuat_id ON ppi_inklusi(guru_pembuat_id);

COMMIT;
