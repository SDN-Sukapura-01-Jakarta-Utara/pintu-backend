package services

import (
	"errors"
	"fmt"
	"strings"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"

	"github.com/xuri/excelize/v2"
)

// exportGuruBidangStudiSemester exports absensi for guru bidang studi per semester
func (s *AbsensiServiceImpl) exportGuruBidangStudiSemester(req *dtos.ExportAbsensiExcelRequest, rombelNama, tahunPelajaranNama, bidangStudiNama string) (*excelize.File, error) {
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

	// Get all absensi data for this semester and bidang studi
	var absensiList []models.RekapitulasiAbsensi
	if err := s.db.Where("rombel_id = ? AND tahun_pelajaran_id = ? AND bidang_studi_id = ? AND semester = ?", 
		req.RombelID, req.TahunPelajaranID, *req.BidangStudiID, *req.Semester).
		Order("tanggal ASC, pertemuan_ke ASC").
		Find(&absensiList).Error; err != nil {
		return nil, errors.New("gagal mengambil data absensi")
	}

	// Group absensi by month
	type MonthData struct {
		Bulan            int
		MaxPertemuan     int
		PertemuanTanggal map[int]string // pertemuan_ke -> tanggal (DD)
	}

	monthsMap := make(map[int]*MonthData) // bulan -> MonthData
	
	// Initialize all months for the semester (even if no data yet)
	var semesterMonths []int
	if *req.Semester == 1 {
		// Semester 1: Juli (7) sampai Desember (12)
		semesterMonths = []int{7, 8, 9, 10, 11, 12}
	} else {
		// Semester 2: Januari (1) sampai Juni (6)
		semesterMonths = []int{1, 2, 3, 4, 5, 6}
	}
	
	// Initialize all months with default values
	for _, bulan := range semesterMonths {
		monthsMap[bulan] = &MonthData{
			Bulan:            bulan,
			MaxPertemuan:     5, // Default minimum 5 pertemuan
			PertemuanTanggal: make(map[int]string),
		}
	}
	
	// Update with actual data from absensiList
	for _, absensi := range absensiList {
		bulan := int(absensi.Tanggal.Month())
		
		// Only process if bulan is in the semester
		if _, exists := monthsMap[bulan]; !exists {
			continue // Skip if month not in this semester
		}
		
		if absensi.PertemuanKe != nil {
			// Update max pertemuan for this month
			if *absensi.PertemuanKe > monthsMap[bulan].MaxPertemuan {
				monthsMap[bulan].MaxPertemuan = *absensi.PertemuanKe
			}
			
			// Store tanggal (take first occurrence)
			if _, exists := monthsMap[bulan].PertemuanTanggal[*absensi.PertemuanKe]; !exists {
				monthsMap[bulan].PertemuanTanggal[*absensi.PertemuanKe] = absensi.Tanggal.Format("02")
			}
		}
	}

	// Use semesterMonths as the sorted list (already in order)
	months := semesterMonths
	// Build absensi map: pesertaDidikRombelID -> bulan -> pertemuan_ke -> status
	absensiMap := make(map[uint]map[int]map[int]string)
	for _, absensi := range absensiList {
		bulan := int(absensi.Tanggal.Month())
		
		if _, exists := absensiMap[absensi.PesertaDidikRombelID]; !exists {
			absensiMap[absensi.PesertaDidikRombelID] = make(map[int]map[int]string)
		}
		if _, exists := absensiMap[absensi.PesertaDidikRombelID][bulan]; !exists {
			absensiMap[absensi.PesertaDidikRombelID][bulan] = make(map[int]string)
		}
		
		if absensi.PertemuanKe != nil {
			absensiMap[absensi.PesertaDidikRombelID][bulan][*absensi.PertemuanKe] = absensi.Status
		}
	}

	// Title row (row 1)
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 14},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	
	// Calculate total pertemuan across all months
	totalPertemuan := 0
	for _, monthData := range monthsMap {
		totalPertemuan += monthData.MaxPertemuan
	}
	
	// Calculate last column for merge (NO + NAMA + P/L + TotalPertemuan + S + I + A + JUMLAH)
	// Total columns = 3 (NO, NAMA, P/L) + totalPertemuan + 3 (S, I, A) + 1 (JUMLAH)
	lastCol := 3 + totalPertemuan + 3 + 1
	lastColCell, _ := excelize.CoordinatesToCellName(lastCol, 1)
	
	title := fmt.Sprintf("DAFTAR KEHADIRAN KELAS %s TAHUN PELAJARAN %s", rombelNama, tahunPelajaranNama)
	f.SetCellValue(sheetName, "A1", title)
	f.MergeCell(sheetName, "A1", lastColCell)
	f.SetCellStyle(sheetName, "A1", lastColCell, titleStyle)

	// Subtitle row (row 2)
	subtitle := fmt.Sprintf("%s - SEMESTER %d", strings.ToUpper(bidangStudiNama), *req.Semester)
	lastColCell2, _ := excelize.CoordinatesToCellName(lastCol, 2)
	f.SetCellValue(sheetName, "A2", subtitle)
	f.MergeCell(sheetName, "A2", lastColCell2)
	f.SetCellStyle(sheetName, "A2", lastColCell2, titleStyle)

	// Header style
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

	dataStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
		},
	})

	bulanNama := []string{"", "JANUARI", "FEBRUARI", "MARET", "APRIL", "MEI", "JUNI",
		"JULI", "AGUSTUS", "SEPTEMBER", "OKTOBER", "NOVEMBER", "DESEMBER"}

	// Header Row 1 (row 4): NO, NAMA SISWA, P/L, BULAN1, BULAN2, ..., JUMLAH ABSEN, JUMLAH
	headerRow1 := 4
	f.SetCellValue(sheetName, "A4", "NO")
	f.SetCellValue(sheetName, "B4", "NAMA SISWA")
	f.SetCellValue(sheetName, "C4", "P/L")
	
	// Merge NO, NAMA SISWA, P/L (rows 4-6)
	f.MergeCell(sheetName, "A4", "A6")
	f.MergeCell(sheetName, "B4", "B6")
	f.MergeCell(sheetName, "C4", "C6")

	// Calculate column positions for each month
	currentCol := 4 // Start after P/L
	monthColMap := make(map[int]int) // bulan -> start column
	
	for _, bulan := range months {
		monthData := monthsMap[bulan]
		monthStartCol := currentCol
		monthEndCol := currentCol + monthData.MaxPertemuan - 1
		monthColMap[bulan] = monthStartCol
		
		// Set month name and merge
		monthStartCell, _ := excelize.CoordinatesToCellName(monthStartCol, headerRow1)
		monthEndCell, _ := excelize.CoordinatesToCellName(monthEndCol, headerRow1)
		f.SetCellValue(sheetName, monthStartCell, bulanNama[bulan])
		f.MergeCell(sheetName, monthStartCell, monthEndCell)
		
		currentCol = monthEndCol + 1
	}

	// JUMLAH ABSEN column
	jumlahAbsenStartCol := currentCol
	jumlahAbsenEndCol := currentCol + 2
	jumlahAbsenStartCell, _ := excelize.CoordinatesToCellName(jumlahAbsenStartCol, headerRow1)
	jumlahAbsenEndCell, _ := excelize.CoordinatesToCellName(jumlahAbsenEndCol, headerRow1)
	f.SetCellValue(sheetName, jumlahAbsenStartCell, "JUMLAH ABSEN")
	f.MergeCell(sheetName, jumlahAbsenStartCell, jumlahAbsenEndCell)
	jumlahAbsenR2, _ := excelize.CoordinatesToCellName(jumlahAbsenStartCol, headerRow1+1)
	f.MergeCell(sheetName, jumlahAbsenStartCell, jumlahAbsenR2)
	
	// JUMLAH column
	jumlahCol := jumlahAbsenEndCol + 1
	jumlahCellR1, _ := excelize.CoordinatesToCellName(jumlahCol, headerRow1)
	jumlahCellR3, _ := excelize.CoordinatesToCellName(jumlahCol, headerRow1+2)
	f.SetCellValue(sheetName, jumlahCellR1, "JUMLAH")
	f.MergeCell(sheetName, jumlahCellR1, jumlahCellR3)

	// Header Row 2 (row 5): P1, P2, P3, ... for each month
	headerRow2 := 5
	for _, bulan := range months {
		monthData := monthsMap[bulan]
		monthStartCol := monthColMap[bulan]
		
		for p := 1; p <= monthData.MaxPertemuan; p++ {
			col := monthStartCol + p - 1
			cell, _ := excelize.CoordinatesToCellName(col, headerRow2)
			f.SetCellValue(sheetName, cell, fmt.Sprintf("P%d", p))
		}
	}

	// Header Row 3 (row 6): dates for each pertemuan in each month, S, I, A
	headerRow3 := 6
	for _, bulan := range months {
		monthData := monthsMap[bulan]
		monthStartCol := monthColMap[bulan]
		
		for p := 1; p <= monthData.MaxPertemuan; p++ {
			col := monthStartCol + p - 1
			cell, _ := excelize.CoordinatesToCellName(col, headerRow3)
			
			if tanggal, exists := monthData.PertemuanTanggal[p]; exists {
				f.SetCellValue(sheetName, cell, tanggal)
			} else {
				f.SetCellValue(sheetName, cell, "-")
			}
		}
	}
	
	// S, I, A columns
	sCell, _ := excelize.CoordinatesToCellName(jumlahAbsenStartCol, headerRow3)
	iCell, _ := excelize.CoordinatesToCellName(jumlahAbsenStartCol+1, headerRow3)
	aCell, _ := excelize.CoordinatesToCellName(jumlahAbsenStartCol+2, headerRow3)
	f.SetCellValue(sheetName, sCell, "S")
	f.SetCellValue(sheetName, iCell, "I")
	f.SetCellValue(sheetName, aCell, "A")

	// Apply header styles
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

		// NO, NAMA, P/L
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), idx+1)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), siswa.Nama)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), siswa.JenisKelamin)

		// Attendance per month and pertemuan
		totalS, totalI, totalA := 0, 0, 0
		
		for _, bulan := range months {
			monthData := monthsMap[bulan]
			monthStartCol := monthColMap[bulan]
			
			for p := 1; p <= monthData.MaxPertemuan; p++ {
				col := monthStartCol + p - 1
				cell, _ := excelize.CoordinatesToCellName(col, row)
				
				if status, exists := absensiMap[pdr.ID][bulan][p]; exists {
					switch status {
					case "hadir":
						f.SetCellValue(sheetName, cell, "✓")
					case "sakit":
						f.SetCellValue(sheetName, cell, "S")
						totalS++
					case "izin":
						f.SetCellValue(sheetName, cell, "I")
						totalI++
					case "alpa":
						f.SetCellValue(sheetName, cell, "A")
						totalA++
					}
				} else {
					f.SetCellValue(sheetName, cell, "-")
				}
				f.SetCellStyle(sheetName, cell, cell, dataStyle)
			}
		}

		// Jumlah Absen: S, I, A
		sCell, _ := excelize.CoordinatesToCellName(jumlahAbsenStartCol, row)
		iCell, _ := excelize.CoordinatesToCellName(jumlahAbsenStartCol+1, row)
		aCell, _ := excelize.CoordinatesToCellName(jumlahAbsenStartCol+2, row)
		f.SetCellValue(sheetName, sCell, totalS)
		f.SetCellValue(sheetName, iCell, totalI)
		f.SetCellValue(sheetName, aCell, totalA)
		f.SetCellStyle(sheetName, sCell, sCell, dataStyle)
		f.SetCellStyle(sheetName, iCell, iCell, dataStyle)
		f.SetCellStyle(sheetName, aCell, aCell, dataStyle)

		// JUMLAH
		jumlahCell, _ := excelize.CoordinatesToCellName(jumlahCol, row)
		f.SetCellValue(sheetName, jumlahCell, totalS+totalI+totalA)
		f.SetCellStyle(sheetName, jumlahCell, jumlahCell, dataStyle)

		// Apply data style to NO, NAMA, P/L
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), dataStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), dataStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), dataStyle)
	}

	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 5)
	f.SetColWidth(sheetName, "B", "B", 30)
	f.SetColWidth(sheetName, "C", "C", 5)
	
	// Set width for all pertemuan columns
	for _, bulan := range months {
		monthData := monthsMap[bulan]
		monthStartCol := monthColMap[bulan]
		for p := 0; p < monthData.MaxPertemuan; p++ {
			colName, _ := excelize.ColumnNumberToName(monthStartCol + p)
			f.SetColWidth(sheetName, colName, colName, 5)
		}
	}
	
	// Set width for S, I, A, JUMLAH
	for col := jumlahAbsenStartCol; col <= jumlahAbsenEndCol; col++ {
		colName, _ := excelize.ColumnNumberToName(col)
		f.SetColWidth(sheetName, colName, colName, 4)
	}
	jumlahColName, _ := excelize.ColumnNumberToName(jumlahCol)
	f.SetColWidth(sheetName, jumlahColName, jumlahColName, 8)

	return f, nil
}
