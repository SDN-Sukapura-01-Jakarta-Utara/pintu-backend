-- Migration: create_monitoring_inklusi_semester_table
-- Created: 2026-09-13 13:25:37
-- Description: Create monitoring_inklusi_semester table for semester-based evaluation of inclusive students

BEGIN;

CREATE TABLE monitoring_inklusi_semester (
    id SERIAL PRIMARY KEY,
    anak_inklusi_rombel_id INTEGER NOT NULL,
    semester SMALLINT NOT NULL,
    tahun_pelajaran_id INTEGER NOT NULL,
    deskripsi TEXT NOT NULL,
    predikat VARCHAR(50),
    rekomendasi TEXT,
    rekomendasi_guru_selanjutnya TEXT,
    guru_pengisi_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by_id INTEGER,
    updated_by_id INTEGER,
    deleted_at TIMESTAMP,
    
    CONSTRAINT fk_monitoring_inklusi_semester_anak_inklusi_rombel 
        FOREIGN KEY (anak_inklusi_rombel_id) 
        REFERENCES anak_inklusi_rombel(id) 
        ON DELETE CASCADE,
        
    CONSTRAINT fk_monitoring_inklusi_semester_tahun_pelajaran 
        FOREIGN KEY (tahun_pelajaran_id) 
        REFERENCES tahun_pelajaran(id) 
        ON DELETE CASCADE,
        
    CONSTRAINT fk_monitoring_inklusi_semester_guru_pengisi 
        FOREIGN KEY (guru_pengisi_id) 
        REFERENCES kepegawaian(id) 
        ON DELETE SET NULL
);

-- Create indexes for better query performance
CREATE INDEX idx_monitoring_inklusi_semester_anak_inklusi_rombel_id ON monitoring_inklusi_semester(anak_inklusi_rombel_id);
CREATE INDEX idx_monitoring_inklusi_semester_tahun_pelajaran_id ON monitoring_inklusi_semester(tahun_pelajaran_id);
CREATE INDEX idx_monitoring_inklusi_semester_semester ON monitoring_inklusi_semester(semester);
CREATE INDEX idx_monitoring_inklusi_semester_predikat ON monitoring_inklusi_semester(predikat);
CREATE INDEX idx_monitoring_inklusi_semester_guru_pengisi_id ON monitoring_inklusi_semester(guru_pengisi_id);

-- Create unique constraint to prevent duplicate entries per semester
CREATE UNIQUE INDEX idx_monitoring_inklusi_semester_unique ON monitoring_inklusi_semester(anak_inklusi_rombel_id, semester, tahun_pelajaran_id) WHERE deleted_at IS NULL;

COMMIT;
