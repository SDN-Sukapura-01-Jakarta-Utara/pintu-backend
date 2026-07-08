-- Migration: create_rekapitulasi_absensi_table
-- Created: 2026-07-06 21:35:31
-- Description: Create rekapitulasi_absensi table for teacher attendance recording per rombel and subject

BEGIN;

CREATE TABLE rekapitulasi_absensi (
    id SERIAL PRIMARY KEY,
    peserta_didik_rombel_id INTEGER NOT NULL REFERENCES peserta_didik_rombel(id) ON DELETE CASCADE,
    rombel_id INTEGER NOT NULL REFERENCES rombel(id),
    tahun_pelajaran_id INTEGER NOT NULL REFERENCES tahun_pelajaran(id),
    semester INTEGER NOT NULL,  -- 1 atau 2
    tanggal DATE NOT NULL,
    bidang_studi_id INTEGER REFERENCES bidang_studi(id),  -- NULL = guru kelas, NOT NULL = guru mapel
    pertemuan_ke INTEGER,  -- NULL = guru kelas, NOT NULL = guru mapel (pertemuan ke-1, ke-2, dst)
    status VARCHAR(20) NOT NULL,  -- 'hadir', 'sakit', 'izin', 'alpa'
    waktu_absen TIMESTAMP,
    metode_input VARCHAR(20) NOT NULL,  -- 'auto' (dari tabel absensi) atau 'manual' (input guru)
    keterangan TEXT,  -- Alasan sakit/izin/alpa
    file_surat TEXT,  -- Path file surat keterangan (untuk sakit/izin)
    dicatat_oleh_id INTEGER,  -- ID user (guru/admin) yang input
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,  -- Soft delete timestamp
    CONSTRAINT unique_rekap_per_mapel UNIQUE(peserta_didik_rombel_id, tanggal, bidang_studi_id)
);

CREATE INDEX idx_rekap_peserta_didik_rombel_id ON rekapitulasi_absensi(peserta_didik_rombel_id);
CREATE INDEX idx_rekap_rombel_id ON rekapitulasi_absensi(rombel_id);
CREATE INDEX idx_rekap_tanggal ON rekapitulasi_absensi(tanggal);
CREATE INDEX idx_rekap_status ON rekapitulasi_absensi(status);
CREATE INDEX idx_rekap_tahun_pelajaran_id ON rekapitulasi_absensi(tahun_pelajaran_id);
CREATE INDEX idx_rekap_semester ON rekapitulasi_absensi(semester);
CREATE INDEX idx_rekap_bidang_studi_id ON rekapitulasi_absensi(bidang_studi_id);
CREATE INDEX idx_rekap_pertemuan_ke ON rekapitulasi_absensi(pertemuan_ke);
CREATE INDEX idx_rekap_peserta_tanggal ON rekapitulasi_absensi(peserta_didik_rombel_id, tanggal);
CREATE INDEX idx_rekap_tanggal_status ON rekapitulasi_absensi(tanggal, status);
CREATE INDEX idx_rekap_tahun_semester ON rekapitulasi_absensi(tahun_pelajaran_id, semester);
CREATE INDEX idx_rekap_rombel_tanggal ON rekapitulasi_absensi(rombel_id, tanggal);
CREATE INDEX idx_rekap_bidang_tanggal ON rekapitulasi_absensi(bidang_studi_id, tanggal);
CREATE INDEX idx_rekap_bidang_pertemuan ON rekapitulasi_absensi(bidang_studi_id, pertemuan_ke);

COMMIT;
