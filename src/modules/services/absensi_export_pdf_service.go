package services

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"

	"github.com/jung-kurt/gofpdf"
)

// ExportAbsensiPDF exports absensi to PDF
func (s *AbsensiServiceImpl) ExportAbsensiPDF(req *dtos.ExportAbsensiExcelRequest) ([]byte, error) {
	// Validate request
	if req.TipePeriode != "bulan" && req.TipePeriode != "semester" {
		return nil, errors.New("tipe_periode harus 'bulan' atau 'semester'")
	}

	// Guru kelas cannot use semester
	if req.BidangStudiID == nil && req.TipePeriode == "semester" {
		return nil, errors.New("guru kelas hanya bisa export per bulan, tidak bisa per semester")
	}

	// Validate required fields based on tipe_periode
	if req.TipePeriode == "bulan" {
		if req.Bulan == nil || req.Tahun == nil {
			return nil, errors.New("bulan dan tahun wajib diisi untuk tipe_periode bulan")
		}
	} else if req.TipePeriode == "semester" {
		if req.Semester == nil {
			return nil, errors.New("semester wajib diisi untuk tipe_periode semester")
		}
	}

	// Get rombel
	var rombel models.Rombel
	if err := s.db.First(&rombel, req.RombelID).Error; err != nil {
		return nil, errors.New("rombel tidak ditemukan")
	}

	// Get tahun pelajaran
	var tahunPelajaran models.TahunPelajaran
	if err := s.db.First(&tahunPelajaran, req.TahunPelajaranID).Error; err != nil {
		return nil, errors.New("tahun pelajaran tidak ditemukan")
	}

	// Get bidang studi if provided
	bidangStudiNama := ""
	if req.BidangStudiID != nil {
		var bidangStudi models.BidangStudi
		if err := s.db.First(&bidangStudi, *req.BidangStudiID).Error; err != nil {
			return nil, errors.New("bidang studi tidak ditemukan")
		}
		bidangStudiNama = bidangStudi.Name
	}

	// Determine which export to use
	if req.BidangStudiID == nil {
		// Guru Kelas - Per Bulan only
		return s.exportPDFGuruKelas(req, rombel.Name, tahunPelajaran.TahunPelajaran)
	} else {
		// Guru Bidang Studi - Per Bulan or Per Semester
		if req.TipePeriode == "bulan" {
			return s.exportPDFGuruBidangStudiBulan(req, rombel.Name, tahunPelajaran.TahunPelajaran, bidangStudiNama)
		} else {
			return s.exportPDFGuruBidangStudiSemester(req, rombel.Name, tahunPelajaran.TahunPelajaran, bidangStudiNama)
		}
	}
}

// exportPDFGuruKelas exports PDF for guru kelas (per bulan)
func (s *AbsensiServiceImpl) exportPDFGuruKelas(req *dtos.ExportAbsensiExcelRequest, rombelNama, tahunPelajaranNama string) ([]byte, error) {
	// Get all students in rombel (only active students)
	var pesertaDidikRombels []models.PesertaDidikRombel
	if err := s.db.Preload("PesertaDidik", "status = ?", "active").
		Where("rombel_id = ? AND tahun_pelajaran_id = ? AND status = ?", req.RombelID, req.TahunPelajaranID, "active").
		Order("peserta_didik_id ASC").
		Find(&pesertaDidikRombels).Error; err != nil {
		return nil, errors.New("gagal mengambil data siswa")
	}

	// Filter out students where PesertaDidik is nil or inactive
	filteredPesertaDidikRombels := []models.PesertaDidikRombel{}
	for _, pdr := range pesertaDidikRombels {
		if pdr.PesertaDidik != nil && pdr.PesertaDidik.Status == "active" {
			filteredPesertaDidikRombels = append(filteredPesertaDidikRombels, pdr)
		}
	}
	pesertaDidikRombels = filteredPesertaDidikRombels

	if len(pesertaDidikRombels) == 0 {
		return nil, errors.New("tidak ada siswa aktif di rombel ini")
	}

	// Get start and end date for the month
	startDate := time.Date(*req.Tahun, time.Month(*req.Bulan), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)
	daysInMonth := endDate.Day()

	// Get all absensi data for this month
	var absensiList []models.RekapitulasiAbsensi
	if err := s.db.Where("rombel_id = ? AND tahun_pelajaran_id = ? AND bidang_studi_id IS NULL", req.RombelID, req.TahunPelajaranID).
		Where("EXTRACT(MONTH FROM tanggal) = ? AND EXTRACT(YEAR FROM tanggal) = ?", *req.Bulan, *req.Tahun).
		Find(&absensiList).Error; err != nil {
		return nil, errors.New("gagal mengambil data absensi")
	}

	// Build absensi map
	absensiMap := make(map[uint]map[int]string)
	for _, absensi := range absensiList {
		if _, exists := absensiMap[absensi.PesertaDidikRombelID]; !exists {
			absensiMap[absensi.PesertaDidikRombelID] = make(map[int]string)
		}
		day := absensi.Tanggal.Day()
		absensiMap[absensi.PesertaDidikRombelID][day] = absensi.Status
	}

	// Create PDF - Landscape A4
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(5, 10, 5)
	pdf.SetAutoPageBreak(true, 10)
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 12)
	title := fmt.Sprintf("DAFTAR KEHADIRAN KELAS %s TAHUN PELAJARAN %s", rombelNama, tahunPelajaranNama)
	pdf.CellFormat(287, 7, title, "", 1, "C", false, 0, "")

	// Subtitle
	bulanNama := []string{"", "JANUARI", "FEBRUARI", "MARET", "APRIL", "MEI", "JUNI",
		"JULI", "AGUSTUS", "SEPTEMBER", "OKTOBER", "NOVEMBER", "DESEMBER"}
	subtitle := fmt.Sprintf("BULAN %s TAHUN %d", bulanNama[*req.Bulan], *req.Tahun)
	pdf.CellFormat(287, 7, subtitle, "", 1, "C", false, 0, "")
	pdf.Ln(2)

	// Calculate column widths - increased sizes
	// Total width: 287mm (A4 Landscape with 5mm margins)
	noWidth := 10.0
	namaWidth := 50.0
	plWidth := 10.0
	dateWidth := 6.5 // Per date column - increased
	siaWidth := 8.0  // S, I, A columns - increased
	jumlahWidth := 15.0 // Increased

	totalDateWidth := float64(daysInMonth) * dateWidth
	totalWidth := noWidth + namaWidth + plWidth + totalDateWidth + (3 * siaWidth) + jumlahWidth

	// Adjust if exceeds page width
	if totalWidth > 287 {
		dateWidth = (287 - noWidth - namaWidth - plWidth - (3 * siaWidth) - jumlahWidth) / float64(daysInMonth)
		totalWidth = 287
	}

	// Calculate X offset to center table
	xOffset := (297 - totalWidth) / 2 // A4 landscape = 297mm width

	// Header
	pdf.SetFillColor(128, 128, 128)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 9) // Increased from 7

	// Save starting Y position
	startY := pdf.GetY()
	
	// Set X position to center table
	pdf.SetX(xOffset)
	
	// Row 1 headers - with merged cells
	pdf.CellFormat(noWidth, 12, "NO", "1", 0, "C", true, 0, "") // Merged 2 rows
	pdf.CellFormat(namaWidth, 12, "NAMA SISWA", "1", 0, "C", true, 0, "") // Merged 2 rows
	pdf.CellFormat(plWidth, 12, "P/L", "1", 0, "C", true, 0, "") // Merged 2 rows
	
	// Bulan header (spans row 1 only)
	pdf.CellFormat(float64(daysInMonth)*dateWidth, 6, bulanNama[*req.Bulan], "1", 0, "C", true, 0, "")
	
	// JUMLAH ABSEN header (with manual line break display)
	pdf.SetFont("Arial", "B", 8) // Increased from 6
	pdf.CellFormat(3*siaWidth, 6, "JUMLAH ABSEN", "1", 0, "C", true, 0, "")
	pdf.SetFont("Arial", "B", 9)
	
	// JUMLAH header (merged 2 rows)
	pdf.CellFormat(jumlahWidth, 12, "JUMLAH", "1", 0, "C", true, 0, "")
	
	// Move to row 2
	pdf.Ln(6)
	pdf.SetXY(xOffset+noWidth+namaWidth+plWidth, startY+6) // Position at bulan column, row 2

	// Row 2 headers (dates and S, I, A)
	for day := 1; day <= daysInMonth; day++ {
		pdf.CellFormat(dateWidth, 6, fmt.Sprintf("%d", day), "1", 0, "C", true, 0, "")
	}

	pdf.CellFormat(siaWidth, 6, "S", "1", 0, "C", true, 0, "")
	pdf.CellFormat(siaWidth, 6, "I", "1", 0, "C", true, 0, "")
	pdf.CellFormat(siaWidth, 6, "A", "1", 0, "C", true, 0, "")
	pdf.Ln(-1)

	// Data rows
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", 8) // Increased from 6

	for idx, pdr := range pesertaDidikRombels {
		if pdr.PesertaDidik == nil {
			continue
		}

		// Calculate summary
		totalSakit, totalIzin, totalAlpa := 0, 0, 0
		for day := 1; day <= daysInMonth; day++ {
			status, _ := absensiMap[pdr.ID][day]
			switch status {
			case "sakit":
				totalSakit++
			case "izin":
				totalIzin++
			case "alpa":
				totalAlpa++
			}
		}
		totalJumlah := totalSakit + totalIzin + totalAlpa

		// Row data - set X offset for centering
		pdf.SetX(xOffset)
		pdf.CellFormat(noWidth, 6, fmt.Sprintf("%d", idx+1), "1", 0, "C", false, 0, "") // Increased height from 5 to 6
		pdf.CellFormat(namaWidth, 6, pdr.PesertaDidik.Nama, "1", 0, "L", false, 0, "")
		
		jenisKelamin := ""
		if pdr.PesertaDidik.JenisKelamin == "L" {
			jenisKelamin = "L"
		} else {
			jenisKelamin = "P"
		}
		pdf.CellFormat(plWidth, 6, jenisKelamin, "1", 0, "C", false, 0, "")

		// Date columns
		for day := 1; day <= daysInMonth; day++ {
			status, exists := absensiMap[pdr.ID][day]
			mark := "-"
			if exists {
				switch status {
				case "hadir":
					mark = "v" // Using simple 'v' for checkmark
				case "sakit":
					mark = "S"
				case "izin":
					mark = "I"
				case "alpa":
					mark = "A"
				}
			}
			pdf.CellFormat(dateWidth, 6, mark, "1", 0, "C", false, 0, "")
		}

		// Summary columns
		pdf.CellFormat(siaWidth, 6, fmt.Sprintf("%d", totalSakit), "1", 0, "C", false, 0, "")
		pdf.CellFormat(siaWidth, 6, fmt.Sprintf("%d", totalIzin), "1", 0, "C", false, 0, "")
		pdf.CellFormat(siaWidth, 6, fmt.Sprintf("%d", totalAlpa), "1", 0, "C", false, 0, "")
		pdf.CellFormat(jumlahWidth, 6, fmt.Sprintf("%d", totalJumlah), "1", 1, "C", false, 0, "")
	}

	// Add signature section at the bottom right
	pdf.Ln(10) // Space before signature
	
	// Get Kepala Sekolah info from konfigurasi_absensi
	var konfigAbsensi models.KonfigurasiAbsensi
	namaKepsek := "___________________"
	nipKepsek := "___________________"
	if err := s.db.First(&konfigAbsensi).Error; err == nil {
		if konfigAbsensi.NamaKepsek != nil && *konfigAbsensi.NamaKepsek != "" {
			namaKepsek = *konfigAbsensi.NamaKepsek
		}
		if konfigAbsensi.NIPKepsek != nil && *konfigAbsensi.NIPKepsek != "" {
			nipKepsek = *konfigAbsensi.NIPKepsek
		}
	}

	// Position signature on the right side (approximately 200mm from left)
	signatureX := 200.0
	pdf.SetX(signatureX)
	
	// Reset font and color for signature
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", 10)
	
	// Mengetahui,
	pdf.CellFormat(60, 6, "Mengetahui,", "", 1, "C", false, 0, "")
	pdf.SetX(signatureX)
	
	// Kepala SDN Sukapura 01
	pdf.CellFormat(60, 6, "Kepala SDN Sukapura 01", "", 1, "C", false, 0, "")
	pdf.SetX(signatureX)
	
	// Empty space for signature (3 lines)
	pdf.Ln(18)
	pdf.SetX(signatureX)
	
	// Nama Kepsek
	pdf.CellFormat(60, 6, namaKepsek, "", 1, "C", false, 0, "")
	pdf.SetX(signatureX)
	
	// NIP
	pdf.CellFormat(60, 6, fmt.Sprintf("NIP. %s", nipKepsek), "", 1, "C", false, 0, "")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// exportPDFGuruBidangStudiBulan exports PDF for guru bidang studi per bulan
func (s *AbsensiServiceImpl) exportPDFGuruBidangStudiBulan(req *dtos.ExportAbsensiExcelRequest, rombelNama, tahunPelajaranNama, bidangStudiNama string) ([]byte, error) {
	// Get all students
	var pesertaDidikRombels []models.PesertaDidikRombel
	if err := s.db.Preload("PesertaDidik", "status = ?", "active").
		Where("rombel_id = ? AND tahun_pelajaran_id = ? AND status = ?", req.RombelID, req.TahunPelajaranID, "active").
		Order("peserta_didik_id ASC").
		Find(&pesertaDidikRombels).Error; err != nil {
		return nil, errors.New("gagal mengambil data siswa")
	}

	filteredPesertaDidikRombels := []models.PesertaDidikRombel{}
	for _, pdr := range pesertaDidikRombels {
		if pdr.PesertaDidik != nil && pdr.PesertaDidik.Status == "active" {
			filteredPesertaDidikRombels = append(filteredPesertaDidikRombels, pdr)
		}
	}
	pesertaDidikRombels = filteredPesertaDidikRombels

	if len(pesertaDidikRombels) == 0 {
		return nil, errors.New("tidak ada siswa aktif di rombel ini")
	}

	// Get absensi data
	var absensiList []models.RekapitulasiAbsensi
	if err := s.db.Where("rombel_id = ? AND tahun_pelajaran_id = ? AND bidang_studi_id = ?", req.RombelID, req.TahunPelajaranID, *req.BidangStudiID).
		Where("EXTRACT(MONTH FROM tanggal) = ? AND EXTRACT(YEAR FROM tanggal) = ?", *req.Bulan, *req.Tahun).
		Order("pertemuan_ke ASC, tanggal ASC").
		Find(&absensiList).Error; err != nil {
		return nil, errors.New("gagal mengambil data absensi")
	}

	// Determine max pertemuan
	maxPertemuan := 5
	for _, absensi := range absensiList {
		if absensi.PertemuanKe != nil && *absensi.PertemuanKe > maxPertemuan {
			maxPertemuan = *absensi.PertemuanKe
		}
	}

	// Build absensi map
	absensiMap := make(map[uint]map[int]string)
	for _, absensi := range absensiList {
		if _, exists := absensiMap[absensi.PesertaDidikRombelID]; !exists {
			absensiMap[absensi.PesertaDidikRombelID] = make(map[int]string)
		}
		if absensi.PertemuanKe != nil {
			absensiMap[absensi.PesertaDidikRombelID][*absensi.PertemuanKe] = absensi.Status
		}
	}

	// Create PDF
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(5, 10, 5)
	pdf.SetAutoPageBreak(true, 10)
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 12)
	title := fmt.Sprintf("DAFTAR KEHADIRAN KELAS %s TAHUN PELAJARAN %s", rombelNama, tahunPelajaranNama)
	pdf.CellFormat(287, 7, title, "", 1, "C", false, 0, "")

	// Subtitle
	bulanNama := []string{"", "JANUARI", "FEBRUARI", "MARET", "APRIL", "MEI", "JUNI",
		"JULI", "AGUSTUS", "SEPTEMBER", "OKTOBER", "NOVEMBER", "DESEMBER"}
	subtitle := fmt.Sprintf("%s - BULAN %s TAHUN %d", strings.ToUpper(bidangStudiNama), bulanNama[*req.Bulan], *req.Tahun)
	pdf.CellFormat(287, 7, subtitle, "", 1, "C", false, 0, "")
	pdf.Ln(2)

	// Calculate column widths - increased sizes
	noWidth := 10.0
	namaWidth := 60.0
	plWidth := 10.0
	pertemuanWidth := 10.0 // Increased
	siaWidth := 8.0 // Increased
	jumlahWidth := 15.0 // Increased

	totalWidth := noWidth + namaWidth + plWidth + (float64(maxPertemuan) * pertemuanWidth) + (3 * siaWidth) + jumlahWidth
	
	// Calculate X offset to center table
	xOffset := (297 - totalWidth) / 2
	if xOffset < 5 {
		xOffset = 5
	}

	// Header
	pdf.SetFillColor(128, 128, 128)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 9) // Increased from 7

	// Save starting Y position
	startY := pdf.GetY()
	
	// Set X position to center table
	pdf.SetX(xOffset)
	
	// Row 1 headers - with merged cells
	pdf.CellFormat(noWidth, 12, "NO", "1", 0, "C", true, 0, "") // Merged 2 rows
	pdf.CellFormat(namaWidth, 12, "NAMA SISWA", "1", 0, "C", true, 0, "") // Merged 2 rows
	pdf.CellFormat(plWidth, 12, "P/L", "1", 0, "C", true, 0, "") // Merged 2 rows
	
	// PERTEMUAN header (spans row 1 only)
	pdf.CellFormat(float64(maxPertemuan)*pertemuanWidth, 6, "PERTEMUAN", "1", 0, "C", true, 0, "")
	
	// JUMLAH ABSEN header
	pdf.SetFont("Arial", "B", 8) // Increased from 6
	pdf.CellFormat(3*siaWidth, 6, "JUMLAH ABSEN", "1", 0, "C", true, 0, "")
	pdf.SetFont("Arial", "B", 9)
	
	// JUMLAH header (merged 2 rows)
	pdf.CellFormat(jumlahWidth, 12, "JUMLAH", "1", 0, "C", true, 0, "")
	
	// Move to row 2
	pdf.Ln(6)
	pdf.SetXY(xOffset+noWidth+namaWidth+plWidth, startY+6) // Position at pertemuan column, row 2

	// Row 2 - pertemuan numbers
	for p := 1; p <= maxPertemuan; p++ {
		pdf.CellFormat(pertemuanWidth, 6, fmt.Sprintf("%d", p), "1", 0, "C", true, 0, "")
	}

	pdf.CellFormat(siaWidth, 6, "S", "1", 0, "C", true, 0, "")
	pdf.CellFormat(siaWidth, 6, "I", "1", 0, "C", true, 0, "")
	pdf.CellFormat(siaWidth, 6, "A", "1", 0, "C", true, 0, "")
	pdf.Ln(-1)

	// Data rows
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", 8) // Increased from 6

	for idx, pdr := range pesertaDidikRombels {
		if pdr.PesertaDidik == nil {
			continue
		}

		totalSakit, totalIzin, totalAlpa := 0, 0, 0
		for p := 1; p <= maxPertemuan; p++ {
			status, _ := absensiMap[pdr.ID][p]
			switch status {
			case "sakit":
				totalSakit++
			case "izin":
				totalIzin++
			case "alpa":
				totalAlpa++
			}
		}
		totalJumlah := totalSakit + totalIzin + totalAlpa

		// Set X offset for centering
		pdf.SetX(xOffset)
		pdf.CellFormat(noWidth, 6, fmt.Sprintf("%d", idx+1), "1", 0, "C", false, 0, "") // Increased height from 5 to 6
		pdf.CellFormat(namaWidth, 6, pdr.PesertaDidik.Nama, "1", 0, "L", false, 0, "")
		
		jenisKelamin := "L"
		if pdr.PesertaDidik.JenisKelamin == "P" {
			jenisKelamin = "P"
		}
		pdf.CellFormat(plWidth, 6, jenisKelamin, "1", 0, "C", false, 0, "")

		for p := 1; p <= maxPertemuan; p++ {
			status, exists := absensiMap[pdr.ID][p]
			mark := "-"
			if exists {
				switch status {
				case "hadir":
					mark = "v" // Using simple 'v' for checkmark
				case "sakit":
					mark = "S"
				case "izin":
					mark = "I"
				case "alpa":
					mark = "A"
				}
			}
			pdf.CellFormat(pertemuanWidth, 6, mark, "1", 0, "C", false, 0, "")
		}

		pdf.CellFormat(siaWidth, 6, fmt.Sprintf("%d", totalSakit), "1", 0, "C", false, 0, "")
		pdf.CellFormat(siaWidth, 6, fmt.Sprintf("%d", totalIzin), "1", 0, "C", false, 0, "")
		pdf.CellFormat(siaWidth, 6, fmt.Sprintf("%d", totalAlpa), "1", 0, "C", false, 0, "")
		pdf.CellFormat(jumlahWidth, 6, fmt.Sprintf("%d", totalJumlah), "1", 1, "C", false, 0, "")
	}

	// Add signature section at the bottom right
	pdf.Ln(10) // Space before signature
	
	// Get Kepala Sekolah info from konfigurasi_absensi
	var konfigAbsensi models.KonfigurasiAbsensi
	namaKepsek := "___________________"
	nipKepsek := "___________________"
	if err := s.db.First(&konfigAbsensi).Error; err == nil {
		if konfigAbsensi.NamaKepsek != nil && *konfigAbsensi.NamaKepsek != "" {
			namaKepsek = *konfigAbsensi.NamaKepsek
		}
		if konfigAbsensi.NIPKepsek != nil && *konfigAbsensi.NIPKepsek != "" {
			nipKepsek = *konfigAbsensi.NIPKepsek
		}
	}

	// Position signature on the right side (approximately 200mm from left)
	signatureX := 200.0
	pdf.SetX(signatureX)
	
	// Reset font and color for signature
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", 10)
	
	// Mengetahui,
	pdf.CellFormat(60, 6, "Mengetahui,", "", 1, "C", false, 0, "")
	pdf.SetX(signatureX)
	
	// Kepala SDN Sukapura 01
	pdf.CellFormat(60, 6, "Kepala SDN Sukapura 01", "", 1, "C", false, 0, "")
	pdf.SetX(signatureX)
	
	// Empty space for signature (3 lines)
	pdf.Ln(18)
	pdf.SetX(signatureX)
	
	// Nama Kepsek
	pdf.CellFormat(60, 6, namaKepsek, "", 1, "C", false, 0, "")
	pdf.SetX(signatureX)
	
	// NIP
	pdf.CellFormat(60, 6, fmt.Sprintf("NIP. %s", nipKepsek), "", 1, "C", false, 0, "")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// exportPDFGuruBidangStudiSemester exports PDF for guru bidang studi per semester
func (s *AbsensiServiceImpl) exportPDFGuruBidangStudiSemester(req *dtos.ExportAbsensiExcelRequest, rombelNama, tahunPelajaranNama, bidangStudiNama string) ([]byte, error) {
	// Get all students
	var pesertaDidikRombels []models.PesertaDidikRombel
	if err := s.db.Preload("PesertaDidik", "status = ?", "active").
		Where("rombel_id = ? AND tahun_pelajaran_id = ? AND status = ?", req.RombelID, req.TahunPelajaranID, "active").
		Order("peserta_didik_id ASC").
		Find(&pesertaDidikRombels).Error; err != nil {
		return nil, errors.New("gagal mengambil data siswa")
	}

	filteredPesertaDidikRombels := []models.PesertaDidikRombel{}
	for _, pdr := range pesertaDidikRombels {
		if pdr.PesertaDidik != nil && pdr.PesertaDidik.Status == "active" {
			filteredPesertaDidikRombels = append(filteredPesertaDidikRombels, pdr)
		}
	}
	pesertaDidikRombels = filteredPesertaDidikRombels

	if len(pesertaDidikRombels) == 0 {
		return nil, errors.New("tidak ada siswa aktif di rombel ini")
	}

	// Determine semester months
	startMonth, endMonth := 7, 12
	if *req.Semester == 2 {
		startMonth, endMonth = 1, 6
	}

	// Get absensi data for the semester
	var absensiList []models.RekapitulasiAbsensi
	if err := s.db.Where("rombel_id = ? AND tahun_pelajaran_id = ? AND bidang_studi_id = ? AND semester = ?",
		req.RombelID, req.TahunPelajaranID, *req.BidangStudiID, *req.Semester).
		Order("tanggal ASC, pertemuan_ke ASC").
		Find(&absensiList).Error; err != nil {
		return nil, errors.New("gagal mengambil data absensi")
	}

	// Build month-pertemuan map
	monthPertemuanMap := make(map[int]map[uint]map[int]string) // [month][student_id][pertemuan] = status
	for month := startMonth; month <= endMonth; month++ {
		monthPertemuanMap[month] = make(map[uint]map[int]string)
	}

	for _, absensi := range absensiList {
		month := int(absensi.Tanggal.Month())
		if _, exists := monthPertemuanMap[month][absensi.PesertaDidikRombelID]; !exists {
			monthPertemuanMap[month][absensi.PesertaDidikRombelID] = make(map[int]string)
		}
		if absensi.PertemuanKe != nil {
			monthPertemuanMap[month][absensi.PesertaDidikRombelID][*absensi.PertemuanKe] = absensi.Status
		}
	}

	// Create PDF
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(5, 10, 5)
	pdf.SetAutoPageBreak(true, 10)
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 12)
	title := fmt.Sprintf("DAFTAR KEHADIRAN KELAS %s TAHUN PELAJARAN %s", rombelNama, tahunPelajaranNama)
	pdf.CellFormat(287, 7, title, "", 1, "C", false, 0, "")

	// Subtitle
	bulanNama := []string{"", "JANUARI", "FEBRUARI", "MARET", "APRIL", "MEI", "JUNI",
		"JULI", "AGUSTUS", "SEPTEMBER", "OKTOBER", "NOVEMBER", "DESEMBER"}
	subtitle := fmt.Sprintf("%s - SEMESTER %d", strings.ToUpper(bidangStudiNama), *req.Semester)
	pdf.CellFormat(287, 7, subtitle, "", 1, "C", false, 0, "")
	pdf.Ln(2)

	// Calculate column widths - increased for better readability
	noWidth := 8.0
	namaWidth := 45.0
	plWidth := 8.0
	pertemuanWidth := 6.0 // Increased from 5.0
	pertemuanPerMonth := 5
	siaWidth := 7.0 // Increased from 5.0
	jumlahWidth := 12.0 // Increased from 10.0

	totalWidth := noWidth + namaWidth + plWidth + (float64((endMonth-startMonth+1)*pertemuanPerMonth) * pertemuanWidth) + (3 * siaWidth) + jumlahWidth
	
	// Calculate X offset to center table
	xOffset := (297 - totalWidth) / 2
	if xOffset < 5 {
		xOffset = 5
	}

	// Header
	pdf.SetFillColor(128, 128, 128)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 7) // Increased from 6

	// Save starting Y position
	startY := pdf.GetY()

	// Set X position to center table
	pdf.SetX(xOffset)
	
	// Row 1 headers - with merged cells
	pdf.CellFormat(noWidth, 12, "NO", "1", 0, "C", true, 0, "")
	pdf.CellFormat(namaWidth, 12, "NAMA SISWA", "1", 0, "C", true, 0, "")
	pdf.CellFormat(plWidth, 12, "P/L", "1", 0, "C", true, 0, "")

	// Month headers
	for month := startMonth; month <= endMonth; month++ {
		pdf.CellFormat(float64(pertemuanPerMonth)*pertemuanWidth, 6, bulanNama[month], "1", 0, "C", true, 0, "")
	}

	// JUMLAH ABSEN header
	pdf.SetFont("Arial", "B", 6)
	pdf.CellFormat(3*siaWidth, 6, "JML ABSEN", "1", 0, "C", true, 0, "")
	pdf.SetFont("Arial", "B", 7)

	// JUMLAH header
	pdf.CellFormat(jumlahWidth, 12, "JUMLAH", "1", 0, "C", true, 0, "")

	// Move to row 2
	pdf.Ln(6)
	pdf.SetXY(xOffset+noWidth+namaWidth+plWidth, startY+6)

	// Row 2 - pertemuan numbers for each month
	for month := startMonth; month <= endMonth; month++ {
		for p := 1; p <= pertemuanPerMonth; p++ {
			pdf.CellFormat(pertemuanWidth, 6, fmt.Sprintf("%d", p), "1", 0, "C", true, 0, "")
		}
	}

	pdf.CellFormat(siaWidth, 6, "S", "1", 0, "C", true, 0, "")
	pdf.CellFormat(siaWidth, 6, "I", "1", 0, "C", true, 0, "")
	pdf.CellFormat(siaWidth, 6, "A", "1", 0, "C", true, 0, "")
	pdf.Ln(-1)

	// Data rows
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", 6) // Increased from 5

	for idx, pdr := range pesertaDidikRombels {
		if pdr.PesertaDidik == nil {
			continue
		}

		// Calculate totals across all months
		totalSakit, totalIzin, totalAlpa := 0, 0, 0
		for month := startMonth; month <= endMonth; month++ {
			for p := 1; p <= pertemuanPerMonth; p++ {
				if studentData, exists := monthPertemuanMap[month][pdr.ID]; exists {
					status := studentData[p]
					switch status {
					case "sakit":
						totalSakit++
					case "izin":
						totalIzin++
					case "alpa":
						totalAlpa++
					}
				}
			}
		}
		totalJumlah := totalSakit + totalIzin + totalAlpa

		// Set X offset for centering
		pdf.SetX(xOffset)
		pdf.CellFormat(noWidth, 6, fmt.Sprintf("%d", idx+1), "1", 0, "C", false, 0, "") // Increased height from 5 to 6
		pdf.CellFormat(namaWidth, 6, pdr.PesertaDidik.Nama, "1", 0, "L", false, 0, "")

		jenisKelamin := "L"
		if pdr.PesertaDidik.JenisKelamin == "P" {
			jenisKelamin = "P"
		}
		pdf.CellFormat(plWidth, 6, jenisKelamin, "1", 0, "C", false, 0, "")

		// Data for each month and pertemuan
		for month := startMonth; month <= endMonth; month++ {
			for p := 1; p <= pertemuanPerMonth; p++ {
				mark := "-"
				if studentData, exists := monthPertemuanMap[month][pdr.ID]; exists {
					status := studentData[p]
					switch status {
					case "hadir":
						mark = "v"
					case "sakit":
						mark = "S"
					case "izin":
						mark = "I"
					case "alpa":
						mark = "A"
					}
				}
				pdf.CellFormat(pertemuanWidth, 6, mark, "1", 0, "C", false, 0, "")
			}
		}

		pdf.CellFormat(siaWidth, 6, fmt.Sprintf("%d", totalSakit), "1", 0, "C", false, 0, "")
		pdf.CellFormat(siaWidth, 6, fmt.Sprintf("%d", totalIzin), "1", 0, "C", false, 0, "")
		pdf.CellFormat(siaWidth, 6, fmt.Sprintf("%d", totalAlpa), "1", 0, "C", false, 0, "")
		pdf.CellFormat(jumlahWidth, 6, fmt.Sprintf("%d", totalJumlah), "1", 1, "C", false, 0, "")
	}

	// Add signature section at the bottom right
	pdf.Ln(10) // Space before signature
	
	// Get Kepala Sekolah info from konfigurasi_absensi
	var konfigAbsensi models.KonfigurasiAbsensi
	namaKepsek := "___________________"
	nipKepsek := "___________________"
	if err := s.db.First(&konfigAbsensi).Error; err == nil {
		if konfigAbsensi.NamaKepsek != nil && *konfigAbsensi.NamaKepsek != "" {
			namaKepsek = *konfigAbsensi.NamaKepsek
		}
		if konfigAbsensi.NIPKepsek != nil && *konfigAbsensi.NIPKepsek != "" {
			nipKepsek = *konfigAbsensi.NIPKepsek
		}
	}

	// Position signature on the right side (approximately 200mm from left)
	signatureX := 200.0
	pdf.SetX(signatureX)
	
	// Reset font and color for signature
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", 10)
	
	// Mengetahui,
	pdf.CellFormat(60, 6, "Mengetahui,", "", 1, "C", false, 0, "")
	pdf.SetX(signatureX)
	
	// Kepala SDN Sukapura 01
	pdf.CellFormat(60, 6, "Kepala SDN Sukapura 01", "", 1, "C", false, 0, "")
	pdf.SetX(signatureX)
	
	// Empty space for signature (3 lines)
	pdf.Ln(18)
	pdf.SetX(signatureX)
	
	// Nama Kepsek
	pdf.CellFormat(60, 6, namaKepsek, "", 1, "C", false, 0, "")
	pdf.SetX(signatureX)
	
	// NIP
	pdf.CellFormat(60, 6, fmt.Sprintf("NIP. %s", nipKepsek), "", 1, "C", false, 0, "")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
