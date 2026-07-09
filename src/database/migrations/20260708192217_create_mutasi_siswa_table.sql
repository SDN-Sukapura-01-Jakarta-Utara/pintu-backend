-- Migration: create_mutasi_siswa_table
-- Created: 2026-07-08 19:22:17

BEGIN;

CREATE TABLE mutasi_siswa (
    id SERIAL PRIMARY KEY,
    nomor_pendaftaran VARCHAR(50) NOT NULL,
    tahun_pelajaran_id INTEGER NOT NULL,
    semester INTEGER NOT NULL,
    
    -- Data Pribadi Siswa
    nama_lengkap VARCHAR(255) NOT NULL,
    nama_panggilan VARCHAR(100),
    nisn VARCHAR(25),
    tempat_lahir VARCHAR(100) NOT NULL,
    tanggal_lahir DATE NOT NULL,
    jenis_kelamin VARCHAR(20) NOT NULL,
    agama VARCHAR(50) NOT NULL,
    golongan_darah VARCHAR(5),
    anak_ke INTEGER,
    jumlah_saudara INTEGER,
    status_anak VARCHAR(50),
    
    -- Alamat
    alamat TEXT NOT NULL,
    rt VARCHAR(10),
    rw VARCHAR(10),
    kelurahan VARCHAR(100),
    kecamatan VARCHAR(100),
    kota VARCHAR(100),
    provinsi VARCHAR(100),
    
    -- Data Orang Tua
    nama_ayah VARCHAR(255),
    nama_ibu VARCHAR(255),
    pendidikan_ayah VARCHAR(100),
    pendidikan_ibu VARCHAR(100),
    pekerjaan_ayah VARCHAR(255),
    pekerjaan_ibu VARCHAR(255),
    penghasilan_ayah DECIMAL(15,2),
    penghasilan_ibu DECIMAL(15,2),
    nomor_hp_ortu VARCHAR(20),
    
    -- Data Wali
    nama_wali VARCHAR(255),
    pendidikan_wali VARCHAR(100),
    hubungan_wali VARCHAR(100),
    pekerjaan_wali VARCHAR(255),
    nomor_hp_wali VARCHAR(20),
    
    -- Data Sekolah Asal
    pindahan_kelas INTEGER,
    asal_sekolah VARCHAR(50),
    nama_asal_sekolah VARCHAR(255),
    
    -- Dokumen Pendukung
    rapor VARCHAR(255),
    akte_kelahiran VARCHAR(255),
    kartu_keluarga VARCHAR(255),
    sptjm VARCHAR(255),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    CONSTRAINT fk_mutasi_siswa_tahun_pelajaran FOREIGN KEY (tahun_pelajaran_id) 
        REFERENCES tahun_pelajaran(id) ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX idx_mutasi_siswa_nomor_pendaftaran ON mutasi_siswa(nomor_pendaftaran);
CREATE INDEX idx_mutasi_siswa_tahun_pelajaran_id ON mutasi_siswa(tahun_pelajaran_id);
CREATE INDEX idx_mutasi_siswa_nama_lengkap ON mutasi_siswa(nama_lengkap);
CREATE INDEX idx_mutasi_siswa_tahun_semester ON mutasi_siswa(tahun_pelajaran_id, semester);

COMMIT;
