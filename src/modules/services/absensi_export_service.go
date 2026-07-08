package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"

	"github.com/xuri/excelize/v2"
)

// ExportAbsensiExcel exports absensi data to Excel file
func (s *AbsensiServiceImpl) ExportAbsensiExcel(req *dtos.ExportAbsensiExcelRequest) (*excelize.File, error) {
	// Validasi request
	if req.TipePeriode == "bulan" {
		if req.Bulan == nil || req.Tahun == nil {
			return nil, errors.New("bulan dan tahun wajib diisi untuk tipe_periode 'bulan'")
		}
	} else if req.TipePeriode == "semester" {
		if req.Semester == nil {
			return nil, errors.New("semester wajib diisi untuk tipe_periode 'semester'")
		}
		// Semester hanya untuk guru bidang studi
		if req.BidangStudiID == nil {
			return nil, errors.New("semester hanya tersedia untuk guru bidang studi")
		}
	}

	// Get rombel info
	var rombel models.Rombel
	if err := s.db.First(&rombel, req.RombelID).Error; err != nil {
		return nil, errors.New("data rombel tidak ditemukan")
	}

	// Get tahun pelajaran info
	var tahunPelajaran models.TahunPelajaran
	if err := s.db.First(&tahunPelajaran, req.TahunPelajaranID).Error; err != nil {
		return nil, errors.New("data tahun pelajaran tidak ditemukan")
	}

	// Get bidang studi info if applicable
	var bidangStudiNama string
	if req.BidangStudiID != nil {
		var bidangStudi models.BidangStudi
		if err := s.db.First(&bidangStudi, *req.BidangStudiID).Error; err != nil {
			return nil, errors.New("data bidang studi tidak ditemukan")
		}
		bidangStudiNama = bidangStudi.Name
	}

	// Determine if this is for guru kelas or guru bidang studi
	if req.BidangStudiID == nil {
		// Guru Kelas - Per Bulan only
		return s.exportGuruKelas(req, rombel.Name, tahunPelajaran.TahunPelajaran)
	} else {
		// Guru Bidang Studi - Per Bulan or Per Semester
		if req.TipePeriode == "bulan" {
			return s.exportGuruBidangStudiBulan(req, rombel.Name, tahunPelajaran.TahunPelajaran, bidangStudiNama)
		} else {
			return s.exportGuruBidangStudiSemester(req, rombel.Name, tahunPelajaran.TahunPelajaran, bidangStudiNama)
		}
	}
}

// exportGuruKelas exports absensi for guru kelas (per bulan, all dates)
func (s *AbsensiServiceImpl) exportGuruKelas(req *dtos.ExportAbsensiExcelRequest, rombelNama, tahunPelajaranNama string) (*excelize.File, error) {
	f := excelize.NewFile()
	sheetName := "Daftar Kehadiran"
	f.SetSheetName("Sheet1", sheetName)

	// Get start and end date for the month
	startDate := time.Date(*req.Tahun, time.Month(*req.Bulan), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1) // Last day of month
	daysInMonth := endDate.Day()

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

	// Get all absensi data for this month
	var absensiList []models.RekapitulasiAbsensi
	if err := s.db.Where("rombel_id = ? AND tahun_pelajaran_id = ? AND bidang_studi_id IS NULL", req.RombelID, req.TahunPelajaranID).
		Where("EXTRACT(MONTH FROM tanggal) = ? AND EXTRACT(YEAR FROM tanggal) = ?", *req.Bulan, *req.Tahun).
		Find(&absensiList).Error; err != nil {
		return nil, errors.New("gagal mengambil data absensi")
	}

	// Build absensi map: pesertaDidikRombelID -> tanggal (DD) -> status
	absensiMap := make(map[uint]map[int]string)
	for _, absensi := range absensiList {
		if _, exists := absensiMap[absensi.PesertaDidikRombelID]; !exists {
			absensiMap[absensi.PesertaDidikRombelID] = make(map[int]string)
		}
		day := absensi.Tanggal.Day()
		absensiMap[absensi.PesertaDidikRombelID][day] = absensi.Status
	}

	// Title row (row 1)
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 14},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	
	// Calculate last column for merge (NO + NAMA + P/L + Days + S + I + A + JUMLAH)
	// Total columns = 3 (NO, NAMA, P/L) + daysInMonth + 3 (S, I, A) + 1 (JUMLAH)
	lastCol := 3 + daysInMonth + 3 + 1
	lastColCell, _ := excelize.CoordinatesToCellName(lastCol, 1)
	
	title := fmt.Sprintf("DAFTAR KEHADIRAN KELAS %s TAHUN PELAJARAN %s", rombelNama, tahunPelajaranNama)
	f.SetCellValue(sheetName, "A1", title)
	f.MergeCell(sheetName, "A1", lastColCell)
	f.SetCellStyle(sheetName, "A1", lastColCell, titleStyle)

	// Subtitle row (row 2)
	bulanNama := []string{"", "JANUARI", "FEBRUARI", "MARET", "APRIL", "MEI", "JUNI",
		"JULI", "AGUSTUS", "SEPTEMBER", "OKTOBER", "NOVEMBER", "DESEMBER"}
	subtitle := fmt.Sprintf("BULAN %s TAHUN %d", bulanNama[*req.Bulan], *req.Tahun)
	lastColCell2, _ := excelize.CoordinatesToCellName(lastCol, 2)
	f.SetCellValue(sheetName, "A2", subtitle)
	f.MergeCell(sheetName, "A2", lastColCell2)
	f.SetCellStyle(sheetName, "A2", lastColCell2, titleStyle)

	// Header style (gray background, white text, bold, centered)
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#808080"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
		},
	})

	// Data style (centered with border)
	dataStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
		},
	})

	// Header Row 1 (row 4): NO, NAMA SISWA, P/L, JULI (merged), JUMLAH ABSEN (merged), JUMLAH
	headerRow1 := 4
	f.SetCellValue(sheetName, "A4", "NO")
	f.SetCellValue(sheetName, "B4", "NAMA SISWA")
	f.SetCellValue(sheetName, "C4", "P/L")
	
	// Merge NO, NAMA SISWA, P/L (rows 4-5)
	f.MergeCell(sheetName, "A4", "A5")
	f.MergeCell(sheetName, "B4", "B5")
	f.MergeCell(sheetName, "C4", "C5")
	
	// Month name column (merged across all date columns)
	monthStartCol := 4 // Column D (index 4)
	monthEndCol := monthStartCol + daysInMonth - 1
	monthStartCell, _ := excelize.CoordinatesToCellName(monthStartCol, headerRow1)
	monthEndCell, _ := excelize.CoordinatesToCellName(monthEndCol, headerRow1)
	f.SetCellValue(sheetName, monthStartCell, bulanNama[*req.Bulan])
	f.MergeCell(sheetName, monthStartCell, monthEndCell)
	
	// JUMLAH ABSEN column (merged across 3 columns: S, I, A)
	jumlahAbsenStartCol := monthEndCol + 1
	jumlahAbsenEndCol := jumlahAbsenStartCol + 2
	jumlahAbsenStartCell, _ := excelize.CoordinatesToCellName(jumlahAbsenStartCol, headerRow1)
	jumlahAbsenEndCell, _ := excelize.CoordinatesToCellName(jumlahAbsenEndCol, headerRow1)
	f.SetCellValue(sheetName, jumlahAbsenStartCell, "JUMLAH ABSEN")
	f.MergeCell(sheetName, jumlahAbsenStartCell, jumlahAbsenEndCell)
	
	// JUMLAH column (merged rows 4-5)
	jumlahCol := jumlahAbsenEndCol + 1
	jumlahCellR1, _ := excelize.CoordinatesToCellName(jumlahCol, headerRow1)
	jumlahCellR2, _ := excelize.CoordinatesToCellName(jumlahCol, headerRow1+1)
	f.SetCellValue(sheetName, jumlahCellR1, "JUMLAH")
	f.MergeCell(sheetName, jumlahCellR1, jumlahCellR2)

	// Header Row 2 (row 5): dates 1-31, S, I, A
	headerRow2 := 5
	for day := 1; day <= daysInMonth; day++ {
		col := monthStartCol + day - 1
		cell, _ := excelize.CoordinatesToCellName(col, headerRow2)
		f.SetCellValue(sheetName, cell, day)
	}
	
	// S, I, A columns
	sCell, _ := excelize.CoordinatesToCellName(jumlahAbsenStartCol, headerRow2)
	iCell, _ := excelize.CoordinatesToCellName(jumlahAbsenStartCol+1, headerRow2)
	aCell, _ := excelize.CoordinatesToCellName(jumlahAbsenStartCol+2, headerRow2)
	f.SetCellValue(sheetName, sCell, "S")
	f.SetCellValue(sheetName, iCell, "I")
	f.SetCellValue(sheetName, aCell, "A")

	// Apply header styles to all header cells
	for col := 1; col <= jumlahCol; col++ {
		cell1, _ := excelize.CoordinatesToCellName(col, headerRow1)
		cell2, _ := excelize.CoordinatesToCellName(col, headerRow2)
		f.SetCellStyle(sheetName, cell1, cell1, headerStyle)
		f.SetCellStyle(sheetName, cell2, cell2, headerStyle)
	}

	// Data rows starting from row 6
	dataStartRow := 6
	for idx, pdr := range pesertaDidikRombels {
		row := dataStartRow + idx
		siswa := pdr.PesertaDidik

		// NO
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), idx+1)
		
		// NAMA SISWA
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), siswa.Nama)
		
		// P/L (Jenis Kelamin)
		jenisKelamin := siswa.JenisKelamin
		if jenisKelamin == "L" {
			jenisKelamin = "L"
		} else if jenisKelamin == "P" {
			jenisKelamin = "P"
		}
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), jenisKelamin)

		// Attendance per day
		countS, countI, countA := 0, 0, 0
		for day := 1; day <= daysInMonth; day++ {
			col := monthStartCol + day - 1
			cell, _ := excelize.CoordinatesToCellName(col, row)
			
			if status, exists := absensiMap[pdr.ID][day]; exists {
				switch status {
				case "hadir":
					f.SetCellValue(sheetName, cell, "✓")
				case "sakit":
					f.SetCellValue(sheetName, cell, "S")
					countS++
				case "izin":
					f.SetCellValue(sheetName, cell, "I")
					countI++
				case "alpa":
					f.SetCellValue(sheetName, cell, "A")
					countA++
				}
			} else {
				f.SetCellValue(sheetName, cell, "-")
			}
			f.SetCellStyle(sheetName, cell, cell, dataStyle)
		}

		// Jumlah Absen: S, I, A
		sCell, _ := excelize.CoordinatesToCellName(jumlahAbsenStartCol, row)
		iCell, _ := excelize.CoordinatesToCellName(jumlahAbsenStartCol+1, row)
		aCell, _ := excelize.CoordinatesToCellName(jumlahAbsenStartCol+2, row)
		f.SetCellValue(sheetName, sCell, countS)
		f.SetCellValue(sheetName, iCell, countI)
		f.SetCellValue(sheetName, aCell, countA)
		f.SetCellStyle(sheetName, sCell, sCell, dataStyle)
		f.SetCellStyle(sheetName, iCell, iCell, dataStyle)
		f.SetCellStyle(sheetName, aCell, aCell, dataStyle)

		// JUMLAH (total S + I + A)
		jumlahCell, _ := excelize.CoordinatesToCellName(jumlahCol, row)
		f.SetCellValue(sheetName, jumlahCell, countS+countI+countA)
		f.SetCellStyle(sheetName, jumlahCell, jumlahCell, dataStyle)

		// Apply data style to NO, NAMA, P/L
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), dataStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), dataStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), dataStyle)
	}

	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 5)  // NO
	f.SetColWidth(sheetName, "B", "B", 30) // NAMA SISWA
	f.SetColWidth(sheetName, "C", "C", 5)  // P/L
	for col := monthStartCol; col <= monthEndCol; col++ {
		colName, _ := excelize.ColumnNumberToName(col)
		f.SetColWidth(sheetName, colName, colName, 4) // Date columns
	}
	for col := jumlahAbsenStartCol; col <= jumlahAbsenEndCol; col++ {
		colName, _ := excelize.ColumnNumberToName(col)
		f.SetColWidth(sheetName, colName, colName, 4) // S, I, A
	}
	jumlahColName, _ := excelize.ColumnNumberToName(jumlahCol)
	f.SetColWidth(sheetName, jumlahColName, jumlahColName, 8) // JUMLAH

	return f, nil
}

// exportGuruBidangStudiBulan exports absensi for guru bidang studi per bulan
func (s *AbsensiServiceImpl) exportGuruBidangStudiBulan(req *dtos.ExportAbsensiExcelRequest, rombelNama, tahunPelajaranNama, bidangStudiNama string) (*excelize.File, error) {
	f := excelize.NewFile()
	sheetName := "Daftar Kehadiran"
	f.SetSheetName("Sheet1", sheetName)

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

	// Get all absensi data for this month and bidang studi
	var absensiList []models.RekapitulasiAbsensi
	if err := s.db.Where("rombel_id = ? AND tahun_pelajaran_id = ? AND bidang_studi_id = ?", req.RombelID, req.TahunPelajaranID, *req.BidangStudiID).
		Where("EXTRACT(MONTH FROM tanggal) = ? AND EXTRACT(YEAR FROM tanggal) = ?", *req.Bulan, *req.Tahun).
		Order("pertemuan_ke ASC, tanggal ASC").
		Find(&absensiList).Error; err != nil {
		return nil, errors.New("gagal mengambil data absensi")
	}

	// Determine max pertemuan_ke in this month (default minimum 5)
	maxPertemuan := 5
	pertemuanTanggalMap := make(map[int]string) // pertemuan_ke -> tanggal
	for _, absensi := range absensiList {
		if absensi.PertemuanKe != nil && *absensi.PertemuanKe > maxPertemuan {
			maxPertemuan = *absensi.PertemuanKe
		}
		// Store tanggal for each pertemuan (take first occurrence)
		if absensi.PertemuanKe != nil {
			if _, exists := pertemuanTanggalMap[*absensi.PertemuanKe]; !exists {
				pertemuanTanggalMap[*absensi.PertemuanKe] = absensi.Tanggal.Format("02")
			}
		}
	}

	// Build absensi map: pesertaDidikRombelID -> pertemuan_ke -> status
	absensiMap := make(map[uint]map[int]string)
	for _, absensi := range absensiList {
		if _, exists := absensiMap[absensi.PesertaDidikRombelID]; !exists {
			absensiMap[absensi.PesertaDidikRombelID] = make(map[int]string)
		}
		if absensi.PertemuanKe != nil {
			absensiMap[absensi.PesertaDidikRombelID][*absensi.PertemuanKe] = absensi.Status
		}
	}

	// Title row (row 1)
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 14},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	
	// Calculate last column for merge (NO + NAMA + P/L + Pertemuan + S + I + A + JUMLAH)
	// Total columns = 3 (NO, NAMA, P/L) + maxPertemuan + 3 (S, I, A) + 1 (JUMLAH)
	lastCol := 3 + maxPertemuan + 3 + 1
	lastColCell, _ := excelize.CoordinatesToCellName(lastCol, 1)
	
	title := fmt.Sprintf("DAFTAR KEHADIRAN KELAS %s TAHUN PELAJARAN %s", rombelNama, tahunPelajaranNama)
	f.SetCellValue(sheetName, "A1", title)
	f.MergeCell(sheetName, "A1", lastColCell)
	f.SetCellStyle(sheetName, "A1", lastColCell, titleStyle)

	// Subtitle row (row 2)
	bulanNama := []string{"", "JANUARI", "FEBRUARI", "MARET", "APRIL", "MEI", "JUNI",
		"JULI", "AGUSTUS", "SEPTEMBER", "OKTOBER", "NOVEMBER", "DESEMBER"}
	subtitle := fmt.Sprintf("%s - BULAN %s TAHUN %d", strings.ToUpper(bidangStudiNama), bulanNama[*req.Bulan], *req.Tahun)
	lastColCell2, _ := excelize.CoordinatesToCellName(lastCol, 2)
	f.SetCellValue(sheetName, "A2", subtitle)
	f.MergeCell(sheetName, "A2", lastColCell2)
	f.SetCellStyle(sheetName, "A2", lastColCell2, titleStyle)

	// Header style (gray background, white text, bold, centered)
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#808080"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
		},
	})

	// Data style (centered with border)
	dataStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
		},
	})

	// Header Row 1 (row 4): NO, NAMA SISWA, P/L, BULAN (merged), JUMLAH ABSEN (merged), JUMLAH
	headerRow1 := 4
	f.SetCellValue(sheetName, "A4", "NO")
	f.SetCellValue(sheetName, "B4", "NAMA SISWA")
	f.SetCellValue(sheetName, "C4", "P/L")
	
	// Merge NO, NAMA SISWA, P/L (rows 4-6)
	f.MergeCell(sheetName, "A4", "A6")
	f.MergeCell(sheetName, "B4", "B6")
	f.MergeCell(sheetName, "C4", "C6")
	
	// Month name column (merged across all pertemuan columns)
	monthStartCol := 4 // Column D (index 4)
	monthEndCol := monthStartCol + maxPertemuan - 1
	monthStartCell, _ := excelize.CoordinatesToCellName(monthStartCol, headerRow1)
	monthEndCell, _ := excelize.CoordinatesToCellName(monthEndCol, headerRow1)
	f.SetCellValue(sheetName, monthStartCell, bulanNama[*req.Bulan])
	f.MergeCell(sheetName, monthStartCell, monthEndCell)
	
	// JUMLAH ABSEN column (merged across 3 columns: S, I, A) and 3 rows
	jumlahAbsenStartCol := monthEndCol + 1
	jumlahAbsenEndCol := jumlahAbsenStartCol + 2
	jumlahAbsenStartCell, _ := excelize.CoordinatesToCellName(jumlahAbsenStartCol, headerRow1)
	jumlahAbsenEndCell, _ := excelize.CoordinatesToCellName(jumlahAbsenEndCol, headerRow1)
	f.SetCellValue(sheetName, jumlahAbsenStartCell, "JUMLAH ABSEN")
	f.MergeCell(sheetName, jumlahAbsenStartCell, jumlahAbsenEndCell)
	
	// Merge JUMLAH ABSEN vertically to row 5
	jumlahAbsenR2, _ := excelize.CoordinatesToCellName(jumlahAbsenStartCol, headerRow1+1)
	f.MergeCell(sheetName, jumlahAbsenStartCell, jumlahAbsenR2)
	
	// JUMLAH column (merged rows 4-6)
	jumlahCol := jumlahAbsenEndCol + 1
	jumlahCellR1, _ := excelize.CoordinatesToCellName(jumlahCol, headerRow1)
	jumlahCellR3, _ := excelize.CoordinatesToCellName(jumlahCol, headerRow1+2)
	f.SetCellValue(sheetName, jumlahCellR1, "JUMLAH")
	f.MergeCell(sheetName, jumlahCellR1, jumlahCellR3)

	// Header Row 2 (row 5): P1, P2, P3, ... Px
	headerRow2 := 5
	for p := 1; p <= maxPertemuan; p++ {
		col := monthStartCol + p - 1
		cell, _ := excelize.CoordinatesToCellName(col, headerRow2)
		f.SetCellValue(sheetName, cell, fmt.Sprintf("P%d", p))
	}

	// Header Row 3 (row 6): dates for each pertemuan, S, I, A
	headerRow3 := 6
	for p := 1; p <= maxPertemuan; p++ {
		col := monthStartCol + p - 1
		cell, _ := excelize.CoordinatesToCellName(col, headerRow3)
		if tanggal, exists := pertemuanTanggalMap[p]; exists {
			f.SetCellValue(sheetName, cell, tanggal)
		} else {
			f.SetCellValue(sheetName, cell, "-")
		}
	}
	
	// S, I, A columns (row 6)
	sCell, _ := excelize.CoordinatesToCellName(jumlahAbsenStartCol, headerRow3)
	iCell, _ := excelize.CoordinatesToCellName(jumlahAbsenStartCol+1, headerRow3)
	aCell, _ := excelize.CoordinatesToCellName(jumlahAbsenStartCol+2, headerRow3)
	f.SetCellValue(sheetName, sCell, "S")
	f.SetCellValue(sheetName, iCell, "I")
	f.SetCellValue(sheetName, aCell, "A")

	// Apply header styles to all header cells
	for col := 1; col <= jumlahCol; col++ {
		cell1, _ := excelize.CoordinatesToCellName(col, headerRow1)
		cell2, _ := excelize.CoordinatesToCellName(col, headerRow2)
		cell3, _ := excelize.CoordinatesToCellName(col, headerRow3)
		f.SetCellStyle(sheetName, cell1, cell1, headerStyle)
		f.SetCellStyle(sheetName, cell2, cell2, headerStyle)
		f.SetCellStyle(sheetName, cell3, cell3, headerStyle)
	}

	// Data rows starting from row 7
	dataStartRow := 7
	for idx, pdr := range pesertaDidikRombels {
		row := dataStartRow + idx
		siswa := pdr.PesertaDidik

		// NO
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), idx+1)
		
		// NAMA SISWA
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), siswa.Nama)
		
		// P/L
		jenisKelamin := siswa.JenisKelamin
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), jenisKelamin)

		// Attendance per pertemuan
		countS, countI, countA := 0, 0, 0
		for p := 1; p <= maxPertemuan; p++ {
			col := monthStartCol + p - 1
			cell, _ := excelize.CoordinatesToCellName(col, row)
			
			if status, exists := absensiMap[pdr.ID][p]; exists {
				switch status {
				case "hadir":
					f.SetCellValue(sheetName, cell, "✓")
				case "sakit":
					f.SetCellValue(sheetName, cell, "S")
					countS++
				case "izin":
					f.SetCellValue(sheetName, cell, "I")
					countI++
				case "alpa":
					f.SetCellValue(sheetName, cell, "A")
					countA++
				}
			} else {
				f.SetCellValue(sheetName, cell, "-")
			}
			f.SetCellStyle(sheetName, cell, cell, dataStyle)
		}

		// Jumlah Absen: S, I, A
		sCell, _ := excelize.CoordinatesToCellName(jumlahAbsenStartCol, row)
		iCell, _ := excelize.CoordinatesToCellName(jumlahAbsenStartCol+1, row)
		aCell, _ := excelize.CoordinatesToCellName(jumlahAbsenStartCol+2, row)
		f.SetCellValue(sheetName, sCell, countS)
		f.SetCellValue(sheetName, iCell, countI)
		f.SetCellValue(sheetName, aCell, countA)
		f.SetCellStyle(sheetName, sCell, sCell, dataStyle)
		f.SetCellStyle(sheetName, iCell, iCell, dataStyle)
		f.SetCellStyle(sheetName, aCell, aCell, dataStyle)

		// JUMLAH (total S + I + A)
		jumlahCell, _ := excelize.CoordinatesToCellName(jumlahCol, row)
		f.SetCellValue(sheetName, jumlahCell, countS+countI+countA)
		f.SetCellStyle(sheetName, jumlahCell, jumlahCell, dataStyle)

		// Apply data style
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), dataStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), dataStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), dataStyle)
	}

	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 5)  // NO
	f.SetColWidth(sheetName, "B", "B", 30) // NAMA SISWA
	f.SetColWidth(sheetName, "C", "C", 5)  // P/L
	for col := monthStartCol; col <= monthEndCol; col++ {
		colName, _ := excelize.ColumnNumberToName(col)
		f.SetColWidth(sheetName, colName, colName, 5) // Pertemuan columns
	}
	for col := jumlahAbsenStartCol; col <= jumlahAbsenEndCol; col++ {
		colName, _ := excelize.ColumnNumberToName(col)
		f.SetColWidth(sheetName, colName, colName, 4) // S, I, A
	}
	jumlahColName, _ := excelize.ColumnNumberToName(jumlahCol)
	f.SetColWidth(sheetName, jumlahColName, jumlahColName, 8) // JUMLAH

	return f, nil
}
