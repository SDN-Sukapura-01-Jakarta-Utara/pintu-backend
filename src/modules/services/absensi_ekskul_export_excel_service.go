package services

import (
	"fmt"
	"pintu-backend/src/dtos"
	"sort"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// DownloadExcelAbsensiSiswa generates Excel file for student attendance
func (s *AbsensiEkskulService) DownloadExcelAbsensiSiswa(req *dtos.AbsensiEkskulGetRequest) (*excelize.File, string, error) {
	// Get all attendance data without pagination
	response, err := s.GetAbsensiSiswa(req)
	if err != nil {
		return nil, "", err
	}

	// Get month and year for title
	bulanText := "SEMUA BULAN"
	tahunText := ""
	if req.Bulan != nil && req.Tahun != nil {
		bulanNames := []string{"", "JANUARI", "FEBRUARI", "MARET", "APRIL", "MEI", "JUNI", "JULI", "AGUSTUS", "SEPTEMBER", "OKTOBER", "NOVEMBER", "DESEMBER"}
		if *req.Bulan >= 1 && *req.Bulan <= 12 {
			bulanText = bulanNames[*req.Bulan]
		}
		tahunText = fmt.Sprintf("%d", *req.Tahun)
	} else if req.Tahun != nil {
		tahunText = fmt.Sprintf("%d", *req.Tahun)
	}

	// Create new Excel file
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheetName := "Absensi Siswa"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, "", err
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 5)  // No
	f.SetColWidth(sheetName, "B", "B", 30) // Nama
	f.SetColWidth(sheetName, "C", "C", 12) // NIS
	f.SetColWidth(sheetName, "D", "D", 10) // Rombel

	// Title style
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 14},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})

	// Subtitle style
	subtitleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 12},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	// Header style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#D3D3D3"}, Pattern: 1},
	})

	// Data style
	dataStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	// Center style for attendance marks
	centerStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	currentRow := 1

	// Title
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("DAFTAR HADIR SISWA EKSTRAKURIKULER %s", strings.ToUpper(response.NamaEkstrakurikuler)))
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("A%d", currentRow), titleStyle)
	currentRow++

	// Subtitle
	if tahunText != "" {
		subtitle := fmt.Sprintf("BULAN %s TAHUN %s", bulanText, tahunText)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow), subtitle)
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("A%d", currentRow), subtitleStyle)
		currentRow++
	}

	currentRow++ // Empty row

	// Collect unique dates from kegiatan
	dateMap := make(map[string]time.Time)
	for _, kegiatan := range response.Kegiatan {
		parsedDate, _ := time.Parse("2006-01-02", kegiatan.TanggalKegiatan)
		dateMap[kegiatan.TanggalKegiatan] = parsedDate
	}

	// Sort dates
	var dates []time.Time
	for _, date := range dateMap {
		dates = append(dates, date)
	}
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	// Calculate number of pertemuan (P1, P2, ...)
	totalPertemuan := len(dates)
	if totalPertemuan < 5 {
		totalPertemuan = 5 // Minimum 5 columns
	}

	// Merge cells for title and subtitle
	lastCol := string(rune('D' + totalPertemuan))
	f.MergeCell(sheetName, "A1", fmt.Sprintf("%s1", lastCol))
	f.MergeCell(sheetName, "A2", fmt.Sprintf("%s2", lastCol))

	// Set column width for pertemuan columns
	for i := 0; i < totalPertemuan; i++ {
		col := string(rune('E' + i))
		f.SetColWidth(sheetName, col, col, 8)
	}

	headerRow1 := currentRow
	headerRow2 := currentRow + 1

	// Header Row 1: No, Nama, NIS, Rombel, P1, P2, ...
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", headerRow1), "No")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", headerRow1), "Nama")
	f.SetCellValue(sheetName, fmt.Sprintf("C%d", headerRow1), "NIS")
	f.SetCellValue(sheetName, fmt.Sprintf("D%d", headerRow1), "Rombel")

	for i := 0; i < totalPertemuan; i++ {
		col := string(rune('E' + i))
		f.SetCellValue(sheetName, fmt.Sprintf("%s%d", col, headerRow1), fmt.Sprintf("P%d", i+1))
	}

	// Header Row 2: Dates (only for P columns)
	for i, date := range dates {
		col := string(rune('E' + i))
		f.SetCellValue(sheetName, fmt.Sprintf("%s%d", col, headerRow2), date.Format("02"))
	}

	// Fill remaining P columns with empty dates if less than 5
	for i := len(dates); i < totalPertemuan; i++ {
		col := string(rune('E' + i))
		f.SetCellValue(sheetName, fmt.Sprintf("%s%d", col, headerRow2), "-")
	}

	// Merge cells for No, Nama, NIS, Rombel (spanning 2 rows)
	f.MergeCell(sheetName, fmt.Sprintf("A%d", headerRow1), fmt.Sprintf("A%d", headerRow2))
	f.MergeCell(sheetName, fmt.Sprintf("B%d", headerRow1), fmt.Sprintf("B%d", headerRow2))
	f.MergeCell(sheetName, fmt.Sprintf("C%d", headerRow1), fmt.Sprintf("C%d", headerRow2))
	f.MergeCell(sheetName, fmt.Sprintf("D%d", headerRow1), fmt.Sprintf("D%d", headerRow2))

	// Apply header style to both header rows
	for col := 'A'; col <= rune('D'+totalPertemuan); col++ {
		f.SetCellStyle(sheetName, fmt.Sprintf("%c%d", col, headerRow1), fmt.Sprintf("%c%d", col, headerRow1), headerStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("%c%d", col, headerRow2), fmt.Sprintf("%c%d", col, headerRow2), headerStyle)
	}

	currentRow = headerRow2 + 1

	// Build student attendance map
	studentMap := make(map[uint]map[string]string) // peserta_didik_rombel_id -> date -> status
	studentInfo := make(map[uint]dtos.AbsensiSiswaDetail)

	for _, kegiatan := range response.Kegiatan {
		for _, absensi := range kegiatan.AbsensiSiswa {
			if _, exists := studentMap[absensi.PesertaDidikRombelID]; !exists {
				studentMap[absensi.PesertaDidikRombelID] = make(map[string]string)
			}
			studentMap[absensi.PesertaDidikRombelID][kegiatan.TanggalKegiatan] = absensi.Status
			studentInfo[absensi.PesertaDidikRombelID] = absensi
		}
	}

	// Collect unique students and sort
	type StudentRow struct {
		ID          uint
		Nama        string
		NIS         string
		NamaRombel  string
		Attendances map[string]string
	}

	var students []StudentRow
	for id, info := range studentInfo {
		students = append(students, StudentRow{
			ID:          id,
			Nama:        info.NamaSiswa,
			NIS:         info.NIS,
			NamaRombel:  info.NamaRombel,
			Attendances: studentMap[id],
		})
	}

	// Sort students by rombel first, then by name alphabetically
	sort.Slice(students, func(i, j int) bool {
		if students[i].NamaRombel != students[j].NamaRombel {
			return students[i].NamaRombel < students[j].NamaRombel
		}
		return students[i].Nama < students[j].Nama
	})

	// Write student data
	for idx, student := range students {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow), idx+1)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", currentRow), student.Nama)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", currentRow), student.NIS)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", currentRow), student.NamaRombel)

		// Apply data style
		for col := 'A'; col <= 'D'; col++ {
			f.SetCellStyle(sheetName, fmt.Sprintf("%c%d", col, currentRow), fmt.Sprintf("%c%d", col, currentRow), dataStyle)
		}

		// Write attendance marks
		for i, date := range dates {
			col := string(rune('E' + i))
			dateStr := date.Format("2006-01-02")
			status := student.Attendances[dateStr]

			mark := "-"
			switch status {
			case "hadir":
				mark = "✓"
			case "alpa":
				mark = "A"
			case "sakit":
				mark = "S"
			case "izin":
				mark = "I"
			}

			f.SetCellValue(sheetName, fmt.Sprintf("%s%d", col, currentRow), mark)
			f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", col, currentRow), fmt.Sprintf("%s%d", col, currentRow), centerStyle)
		}

		// Fill remaining columns with "-"
		for i := len(dates); i < totalPertemuan; i++ {
			col := string(rune('E' + i))
			f.SetCellValue(sheetName, fmt.Sprintf("%s%d", col, currentRow), "-")
			f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", col, currentRow), fmt.Sprintf("%s%d", col, currentRow), centerStyle)
		}

		currentRow++
	}

	// Generate filename
	filename := fmt.Sprintf("Absensi_Siswa_%s", response.NamaEkstrakurikuler)
	if req.Bulan != nil && req.Tahun != nil {
		filename += fmt.Sprintf("_%s_%d", bulanText, *req.Tahun)
	}
	filename += ".xlsx"

	return f, filename, nil
}
