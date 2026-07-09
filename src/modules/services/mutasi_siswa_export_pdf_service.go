package services

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/repositories"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// ExportFormulirPDF exports formulir pendaftaran mutasi siswa to PDF
func (s *MutasiSiswaServiceImpl) ExportFormulirPDF(id uint) ([]byte, error) {
	// Get mutasi siswa data
	mutasi, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("data mutasi siswa tidak ditemukan")
	}

	// Get konfigurasi - access db directly from repository implementation
	repoImpl, ok := s.repository.(*repositories.MutasiSiswaRepositoryImpl)
	if !ok {
		return nil, fmt.Errorf("invalid repository implementation")
	}
	
	// Get db field using reflection-free approach - create temp konfigurasi repo
	konfigurasiRepo := repositories.NewKonfigurasiMutasiSiswaRepository(repoImpl.GetDB())
	konfigurasi, err := konfigurasiRepo.GetByID(1)
	if err != nil {
		return nil, fmt.Errorf("konfigurasi mutasi siswa belum diatur")
	}

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(5, 5, 5)
	pdf.SetAutoPageBreak(true, 10)

	// Add first page
	pdf.AddPage()

	// Download and add kop - positioned higher
	kopURL := "https://pintu-storage.sdnsukapura01.sch.id/kop.png"
	resp, err := http.Get(kopURL)
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		kopBytes, _ := io.ReadAll(resp.Body)
		
		// Register kop image
		opt := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}
		pdf.RegisterImageOptionsReader("kop", opt, bytes.NewReader(kopBytes))
		
		// Add kop image - full width, positioned higher
		pdf.ImageOptions("kop", 5, 3, 200, 0, false, opt, 0, "")
		pdf.Ln(45) // Further increased space to push title down more (was 35)
	} else {
		pdf.Ln(5)
	}

	// Title - larger font and more space
	pdf.SetFont("Times", "BU", 14)
	pdf.CellFormat(200, 7, "FORMULIR PENDAFTARAN CALON MURID BARU", "", 1, "C", false, 0, "")
	pdf.Ln(4) // Space after title

	// No Pendaftaran with border - centered and wrapped
	pdf.SetFont("Times", "", 14) // Same size as title
	currentYear := time.Now().Year()
	noPendaftaran := fmt.Sprintf("No. Pendaftaran: %s/CMB/%d", mutasi.NomorPendaftaran, currentYear)
	
	// Calculate text width for wrapping border
	textWidth := pdf.GetStringWidth(noPendaftaran) + 10 // Add padding
	
	// Center position
	pageWidth := 210.0 // A4 width
	xPos := (pageWidth - textWidth) / 2
	
	// Draw border box for no pendaftaran - wrapped around text
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetLineWidth(0.5)
	_, y := pdf.GetXY()
	pdf.Rect(xPos, y, textWidth, 10, "D")
	pdf.SetXY(xPos+5, y+3)
	pdf.CellFormat(textWidth-10, 4, noPendaftaran, "", 1, "C", false, 0, "")
	pdf.Ln(8) // Increased space before SISWA table (was 5)

	// SISWA Table
	s.addSiswaTable(pdf, mutasi)
	pdf.Ln(5)

	// ORANG TUA Table
	s.addOrangTuaTable(pdf, mutasi)
	pdf.Ln(5)

	// Add second page for WALI, ASAL MULA ANAK, and signatures
	pdf.AddPage()
	pdf.Ln(10) // Add spacing from top margin

	// WALI Table
	s.addWaliTable(pdf, mutasi)
	pdf.Ln(5)

	// ASAL MULA ANAK Table
	s.addAsalMulaAnakTable(pdf, mutasi)
	pdf.Ln(8)

	// Signatures
	s.addSignatures(pdf, konfigurasi)

	// Add third page for Catatan
	pdf.AddPage()
	s.addCatatanPage(pdf)

	// Output PDF
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// addSiswaTable adds SISWA section table
func (s *MutasiSiswaServiceImpl) addSiswaTable(pdf *gofpdf.Fpdf, mutasi *models.MutasiSiswa) {
	// Table width calculations - narrower table
	tableWidth := 170.0
	leftColWidth := 70.0
	rightColWidth := 100.0
	leftMargin := (210.0 - tableWidth) / 2 // Center the table
	
	pdf.SetFont("Times", "B", 12) // Larger font for header
	pdf.SetFillColor(245, 245, 220) // Cream color
	
	// Header
	pdf.SetX(leftMargin)
	pdf.CellFormat(tableWidth, 8, "SISWA", "1", 1, "L", true, 0, "")

	pdf.SetFont("Times", "", 11) // Larger font for content
	pdf.SetFillColor(255, 255, 255) // White
	rowHeight := 7.0 // Increased row height

	// 1. Nama
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "1. Nama", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, "", "1", 1, "L", true, 0, "")
	
	namaPanggilan := ""
	if mutasi.NamaPanggilan != nil {
		namaPanggilan = *mutasi.NamaPanggilan
	}
	
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "    a. Nama Lengkap", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+mutasi.NamaLengkap, "1", 1, "L", true, 0, "")
	
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "    b. Nama Panggilan", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+namaPanggilan, "1", 1, "L", true, 0, "")

	// 2. NISN
	nisn := ""
	if mutasi.NISN != nil {
		nisn = *mutasi.NISN
	}
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "2. NISN", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+nisn, "1", 1, "L", true, 0, "")

	// 3. Tempat, Tanggal Lahir
	bulanIndo := []string{"", "JANUARI", "FEBRUARI", "MARET", "APRIL", "MEI", "JUNI",
		"JULI", "AGUSTUS", "SEPTEMBER", "OKTOBER", "NOVEMBER", "DESEMBER"}
	tanggalLahir := fmt.Sprintf("%d %s %d", 
		mutasi.TanggalLahir.Day(), 
		bulanIndo[mutasi.TanggalLahir.Month()], 
		mutasi.TanggalLahir.Year())
	
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "3. Tempat, Tanggal Lahir", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+mutasi.TempatLahir+", "+tanggalLahir, "1", 1, "L", true, 0, "")

	// 4. Jenis Kelamin
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "4. Jenis Kelamin", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+mutasi.JenisKelamin, "1", 1, "L", true, 0, "")

	// 5. Agama
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "5. Agama", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+mutasi.Agama, "1", 1, "L", true, 0, "")

	// 6. Golongan Darah
	goldar := ""
	if mutasi.GolonganDarah != nil {
		goldar = *mutasi.GolonganDarah
	}
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "6. Golongan Darah", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+goldar, "1", 1, "L", true, 0, "")

	// 7. Anak Ke
	anakKe := ""
	if mutasi.AnakKe != nil {
		anakKe = strconv.Itoa(*mutasi.AnakKe)
	}
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "7. Anak Ke", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+anakKe, "1", 1, "L", true, 0, "")

	// 8. Jumlah Saudara
	jumlahSaudara := ""
	if mutasi.JumlahSaudara != nil {
		jumlahSaudara = strconv.Itoa(*mutasi.JumlahSaudara)
	}
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "8. Jumlah Saudara", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+jumlahSaudara, "1", 1, "L", true, 0, "")

	// 9. Status Anak
	statusAnak := ""
	if mutasi.StatusAnak != nil {
		statusAnak = *mutasi.StatusAnak
	}
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "9. Status Anak dalam Keluarga", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+statusAnak, "1", 1, "L", true, 0, "")

	// 10. Alamat Siswa - with proper height synchronization
	alamatLengkap := mutasi.Alamat
	if mutasi.RT != nil && *mutasi.RT != "" {
		alamatLengkap += ", RT. " + *mutasi.RT
	}
	if mutasi.RW != nil && *mutasi.RW != "" {
		alamatLengkap += ", RW. " + *mutasi.RW
	}
	if mutasi.Kelurahan != nil && *mutasi.Kelurahan != "" {
		alamatLengkap += ", " + *mutasi.Kelurahan
	}
	if mutasi.Kecamatan != nil && *mutasi.Kecamatan != "" {
		alamatLengkap += ", " + *mutasi.Kecamatan
	}
	if mutasi.Kota != nil && *mutasi.Kota != "" {
		alamatLengkap += ", " + *mutasi.Kota
	}
	if mutasi.Provinsi != nil && *mutasi.Provinsi != "" {
		alamatLengkap += ", " + *mutasi.Provinsi
	}
	
	// Calculate height needed for alamat
	lines := pdf.SplitLines([]byte(": "+alamatLengkap), rightColWidth-2)
	numLines := len(lines)
	addressHeight := float64(numLines) * rowHeight
	
	// Draw left cell with calculated height
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, addressHeight, "10. Alamat Siswa", "1", 0, "L", true, 0, "")
	
	// Draw right cell with MultiCell
	_, y := pdf.GetXY()
	pdf.MultiCell(rightColWidth, rowHeight, ": "+alamatLengkap, "1", "L", true)
	
	// Move to end of address section
	pdf.SetXY(leftMargin, y+addressHeight)
}

// addOrangTuaTable adds ORANG TUA section table
func (s *MutasiSiswaServiceImpl) addOrangTuaTable(pdf *gofpdf.Fpdf, mutasi *models.MutasiSiswa) {
	// Table width calculations - same as SISWA table
	tableWidth := 170.0
	leftColWidth := 70.0
	rightColWidth := 100.0
	leftMargin := (210.0 - tableWidth) / 2
	rowHeight := 7.0
	
	pdf.SetFont("Times", "B", 12)
	pdf.SetFillColor(245, 245, 220) // Cream
	
	pdf.SetX(leftMargin)
	pdf.CellFormat(tableWidth, 8, "ORANG TUA", "1", 1, "L", true, 0, "")

	pdf.SetFont("Times", "", 11)
	pdf.SetFillColor(255, 255, 255) // White

	// 1. Nama Orang Tua
	namaAyah := ""
	if mutasi.NamaAyah != nil {
		namaAyah = *mutasi.NamaAyah
	}
	namaIbu := ""
	if mutasi.NamaIbu != nil {
		namaIbu = *mutasi.NamaIbu
	}
	
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "1. Nama Orang Tua", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, "", "1", 1, "L", true, 0, "")
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "    a. Ayah", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+namaAyah, "1", 1, "L", true, 0, "")
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "    b. Ibu Kandung", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+namaIbu, "1", 1, "L", true, 0, "")

	// 2. Pendidikan Tertinggi
	pendAyah := ""
	if mutasi.PendidikanAyah != nil {
		pendAyah = *mutasi.PendidikanAyah
	}
	pendIbu := ""
	if mutasi.PendidikanIbu != nil {
		pendIbu = *mutasi.PendidikanIbu
	}
	
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "2. Pendidikan Tertinggi", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, "", "1", 1, "L", true, 0, "")
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "    a. Ayah", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+pendAyah, "1", 1, "L", true, 0, "")
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "    b. Ibu", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+pendIbu, "1", 1, "L", true, 0, "")

	// 3. Pekerjaan
	pekerjaanAyah := ""
	if mutasi.PekerjaanAyah != nil {
		pekerjaanAyah = *mutasi.PekerjaanAyah
	}
	pekerjaanIbu := ""
	if mutasi.PekerjaanIbu != nil {
		pekerjaanIbu = *mutasi.PekerjaanIbu
	}
	
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "3. Pekerjaan", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, "", "1", 1, "L", true, 0, "")
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "    a. Ayah", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+pekerjaanAyah, "1", 1, "L", true, 0, "")
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "    b. Ibu", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+pekerjaanIbu, "1", 1, "L", true, 0, "")

	// 4. Penghasilan Perbulan
	penghasilanAyah := ""
	if mutasi.PenghasilanAyah != nil {
		penghasilanAyah = fmt.Sprintf("Rp %.0f", *mutasi.PenghasilanAyah)
	}
	penghasilanIbu := ""
	if mutasi.PenghasilanIbu != nil {
		penghasilanIbu = fmt.Sprintf("Rp %.0f", *mutasi.PenghasilanIbu)
	}
	
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "4. Penghasilan Perbulan", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, "", "1", 1, "L", true, 0, "")
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "    a. Ayah", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+penghasilanAyah, "1", 1, "L", true, 0, "")
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "    b. Ibu", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+penghasilanIbu, "1", 1, "L", true, 0, "")

	// 5. Nomor Hp/Whatsapp
	nomorHP := ""
	if mutasi.NomorHPOrtu != nil {
		nomorHP = *mutasi.NomorHPOrtu
	}
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "5. Nomor Hp/Whatsapp", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+nomorHP, "1", 1, "L", true, 0, "")
}

// addWaliTable adds WALI section table
func (s *MutasiSiswaServiceImpl) addWaliTable(pdf *gofpdf.Fpdf, mutasi *models.MutasiSiswa) {
	// Table width calculations - same as other tables
	tableWidth := 170.0
	leftColWidth := 70.0
	rightColWidth := 100.0
	leftMargin := (210.0 - tableWidth) / 2
	rowHeight := 7.0
	
	pdf.SetFont("Times", "B", 12)
	pdf.SetFillColor(245, 245, 220) // Cream
	
	pdf.SetX(leftMargin)
	pdf.CellFormat(tableWidth, 8, "WALI", "1", 1, "L", true, 0, "")

	pdf.SetFont("Times", "", 11)
	pdf.SetFillColor(255, 255, 255) // White

	namaWali := ""
	if mutasi.NamaWali != nil {
		namaWali = *mutasi.NamaWali
	}
	pendidikanWali := ""
	if mutasi.PendidikanWali != nil {
		pendidikanWali = *mutasi.PendidikanWali
	}
	hubunganWali := ""
	if mutasi.HubunganWali != nil {
		hubunganWali = *mutasi.HubunganWali
	}
	pekerjaanWali := ""
	if mutasi.PekerjaanWali != nil {
		pekerjaanWali = *mutasi.PekerjaanWali
	}
	nomorHPWali := ""
	if mutasi.NomorHPWali != nil {
		nomorHPWali = *mutasi.NomorHPWali
	}

	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "1. Nama Wali", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+namaWali, "1", 1, "L", true, 0, "")
	
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "2. Pendidikan Tertinggi Wali", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+pendidikanWali, "1", 1, "L", true, 0, "")
	
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "3. Hubungan Wali Terhadap Anak", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+hubunganWali, "1", 1, "L", true, 0, "")
	
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "4. Pekerjaan Wali", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+pekerjaanWali, "1", 1, "L", true, 0, "")
	
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "5. Nomor Hp/Whatsapp", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+nomorHPWali, "1", 1, "L", true, 0, "")
}

// addAsalMulaAnakTable adds ASAL MULA ANAK section table
func (s *MutasiSiswaServiceImpl) addAsalMulaAnakTable(pdf *gofpdf.Fpdf, mutasi *models.MutasiSiswa) {
	// Table width calculations - same as other tables
	tableWidth := 170.0
	leftColWidth := 70.0
	rightColWidth := 100.0
	leftMargin := (210.0 - tableWidth) / 2
	rowHeight := 7.0
	
	pdf.SetFont("Times", "B", 12)
	pdf.SetFillColor(245, 245, 220) // Cream
	
	pdf.SetX(leftMargin)
	pdf.CellFormat(tableWidth, 8, "ASAL MULA ANAK", "1", 1, "L", true, 0, "")

	pdf.SetFont("Times", "", 11)
	pdf.SetFillColor(255, 255, 255) // White

	pindahanKelas := ""
	if mutasi.PindahanKelas != nil {
		pindahanKelas = "Pindahan kelas " + strconv.Itoa(*mutasi.PindahanKelas)
	}
	asalSekolah := ""
	if mutasi.AsalSekolah != nil {
		asalSekolah = *mutasi.AsalSekolah
	}
	namaAsalSekolah := ""
	if mutasi.NamaAsalSekolah != nil {
		namaAsalSekolah = *mutasi.NamaAsalSekolah
	}

	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "1. Masuk sekolah ini sebagai", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+pindahanKelas, "1", 1, "L", true, 0, "")
	
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "2. Asal Anak", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+asalSekolah, "1", 1, "L", true, 0, "")
	
	pdf.SetX(leftMargin)
	pdf.CellFormat(leftColWidth, rowHeight, "3. Nama Asal Sekolah", "1", 0, "L", true, 0, "")
	pdf.CellFormat(rightColWidth, rowHeight, ": "+namaAsalSekolah, "1", 1, "L", true, 0, "")
}

// addSignatures adds signature section
func (s *MutasiSiswaServiceImpl) addSignatures(pdf *gofpdf.Fpdf, konfigurasi *models.KonfigurasiMutasiSiswa) {
	pdf.SetFont("Times", "", 11) // Larger font

	// Left: Orang Tua/Wali
	_, y := pdf.GetXY()
	pdf.SetXY(30, y)
	pdf.CellFormat(70, 6, "Orang Tua/Wali", "", 1, "L", false, 0, "")
	pdf.SetXY(30, y+6)
	pdf.Ln(25) // Increased space for signature
	pdf.SetXY(30, y+31)
	pdf.CellFormat(70, 6, ".................................", "", 1, "L", false, 0, "")

	// Right: Ketua Panitia
	bulanIndo := []string{"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	now := time.Now()
	tanggalStr := fmt.Sprintf("Jakarta, %d %s %d", now.Day(), bulanIndo[now.Month()], now.Year())
	
	pdf.SetXY(120, y)
	pdf.CellFormat(70, 6, tanggalStr, "", 1, "L", false, 0, "")
	pdf.SetXY(120, y+6)
	pdf.CellFormat(70, 6, "Ketua Panitia PMB", "", 1, "L", false, 0, "")
	pdf.SetXY(120, y+12)
	pdf.Ln(19) // Increased space for signature
	pdf.SetXY(120, y+31)
	pdf.CellFormat(70, 6, konfigurasi.NamaKetuaPanitia, "", 1, "L", false, 0, "")
	pdf.SetXY(120, y+37)
	pdf.CellFormat(70, 6, "NIP. "+konfigurasi.NIPKetuaPanitia, "", 1, "L", false, 0, "")

	// Center: Kepala Sekolah
	pdf.Ln(15) // Extra space before kepala sekolah
	pdf.SetFont("Times", "", 11)
	pdf.CellFormat(200, 6, "Mengetahui,", "", 1, "C", false, 0, "")
	pdf.CellFormat(200, 6, "Kepala SDN Sukapura 01", "", 1, "C", false, 0, "")
	pdf.Ln(25) // Increased space for signature
	pdf.CellFormat(200, 6, konfigurasi.NamaKepalaSekolah, "", 1, "C", false, 0, "")
	pdf.CellFormat(200, 6, "NIP. "+konfigurasi.NIPKepalaSekolah, "", 1, "C", false, 0, "")
}

// addCatatanPage adds second page with notes
func (s *MutasiSiswaServiceImpl) addCatatanPage(pdf *gofpdf.Fpdf) {
	pdf.SetFont("Times", "B", 12)
	pdf.CellFormat(200, 7, "Catatan:", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	pdf.SetFont("Times", "", 11)
	pdf.CellFormat(200, 6, "1. Fotokopi Rapor yang telah dilegalisir", "", 1, "L", false, 0, "")
	pdf.CellFormat(200, 6, "2. Fotokopi Akte Kelahiran Anak", "", 1, "L", false, 0, "")
	pdf.CellFormat(200, 6, "3. Fotokopi Kartu Keluarga", "", 1, "L", false, 0, "")
	pdf.CellFormat(200, 6, "4. Map Plastik", "", 1, "L", false, 0, "")
	pdf.CellFormat(200, 6, "    - Untuk Calon Peserta Didik Laki-Laki : Merah", "", 1, "L", false, 0, "")
	pdf.CellFormat(200, 6, "    - Untuk Calon Peserta Didik Perempuan : Putih", "", 1, "L", false, 0, "")
}
