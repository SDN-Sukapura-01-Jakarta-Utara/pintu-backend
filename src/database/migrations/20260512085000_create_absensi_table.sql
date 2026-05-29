-- Migration: create_absensi_table
-- Created: 2026-05-12 08:50:00
-- Description: Create absensi table for student attendance system with barcode scanning support

BEGIN;

CREATE TABLE absensi (
    id SERIAL PRIMARY KEY,
    peserta_didik_id INTEGER NOT NULL REFERENCES peserta_didik(id) ON DELETE CASCADE,
    rombel_id INTEGER NOT NULL REFERENCES rombel(id),
    tahun_pelajaran_id INTEGER NOT NULL REFERENCES tahun_pelajaran(id),
    semester INTEGER NOT NULL,  -- 1 atau 2
    tanggal DATE NOT NULL,
    bidang_studi_id INTEGER REFERENCES bidang_studi(id),  -- NULL = guru kelas, NOT NULL = guru mapel
    pertemuan_ke INTEGER,  -- NULL = guru kelas, NOT NULL = guru mapel (pertemuan ke-1, ke-2, dst)
    status VARCHAR(20) NOT NULL,  -- 'hadir', 'sakit', 'izin', 'alpa'
    waktu_absen TIMESTAMP,
    metode_input VARCHAR(20) NOT NULL,  -- 'scan' atau 'manual'
    keterangan TEXT,  -- Alasan sakit/izin/alpa
    file_surat TEXT,  -- Path file surat keterangan (untuk sakit/izin)
    dicatat_oleh_id INTEGER,  -- ID user (guru/admin) yang input
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,  -- Soft delete timestamp
    CONSTRAINT unique_absensi_per_mapel UNIQUE(peserta_didik_id, tanggal, bidang_studi_id)
);

CREATE INDEX idx_absensi_peserta_didik_id ON absensi(peserta_didik_id);
CREATE INDEX idx_absensi_rombel_id ON absensi(rombel_id);
CREATE INDEX idx_absensi_tanggal ON absensi(tanggal);
CREATE INDEX idx_absensi_status ON absensi(status);
CREATE INDEX idx_absensi_tahun_pelajaran_id ON absensi(tahun_pelajaran_id);
CREATE INDEX idx_absensi_semester ON absensi(semester);
CREATE INDEX idx_absensi_bidang_studi_id ON absensi(bidang_studi_id);
CREATE INDEX idx_absensi_pertemuan_ke ON absensi(pertemuan_ke);
CREATE INDEX idx_absensi_peserta_tanggal ON absensi(peserta_didik_id, tanggal);
CREATE INDEX idx_absensi_tanggal_status ON absensi(tanggal, status);
CREATE INDEX idx_absensi_tahun_semester ON absensi(tahun_pelajaran_id, semester);
CREATE INDEX idx_absensi_rombel_tanggal ON absensi(rombel_id, tanggal);
CREATE INDEX idx_absensi_bidang_tanggal ON absensi(bidang_studi_id, tanggal);
CREATE INDEX idx_absensi_bidang_pertemuan ON absensi(bidang_studi_id, pertemuan_ke);
CREATE INDEX idx_absensi_deleted_at ON absensi(deleted_at);

COMMIT;
