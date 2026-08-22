package services

import (
	"bytes"
	"fmt"
	"pintu-backend/src/dtos"
	"sort"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// DownloadPDFAbsensiSiswa generates PDF file for student attendance
func (s *AbsensiEkskulService) DownloadPDFAbsensiSiswa(req *dtos.AbsensiEkskulGetRequest) ([]byte, string, error) {
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

	// Calculate number of pertemuan
	totalPertemuan := len(dates)
	if totalPertemuan < 5 {
		totalPertemuan = 5 // Minimum 5 columns
	}

	// Build student attendance map
	studentMap := make(map[uint]map[string]string)
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

	// Collect and sort students
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

	sort.Slice(students, func(i, j int) bool {
		return students[i].Nama < students[j].Nama
	})

	// Create PDF - Portrait for better compatibility
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.SetAutoPageBreak(true, 10)
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 14)
	title := fmt.Sprintf("DAFTAR HADIR SISWA EKSTRAKURIKULER %s", strings.ToUpper(response.NamaEkstrakurikuler))
	pdf.CellFormat(190, 7, title, "", 1, "C", false, 0, "")

	// Subtitle
	if tahunText != "" {
		pdf.SetFont("Arial", "B", 12)
		subtitle := fmt.Sprintf("BULAN %s TAHUN %s", bulanText, tahunText)
		pdf.CellFormat(190, 6, subtitle, "", 1, "C", false, 0, "")
	}

	pdf.Ln(3)

	// Calculate column widths
	noWidth := 10.0
	namaWidth := 60.0
	nisWidth := 25.0
	rombelWidth := 20.0
	pertemuanWidth := 15.0

	// Header setup
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(211, 211, 211)

	// Row 1: Column headers with merge for No, Nama, NIS, Rombel
	pdf.CellFormat(noWidth, 10, "No", "1", 0, "C", true, 0, "")
	pdf.CellFormat(namaWidth, 10, "Nama", "1", 0, "C", true, 0, "")
	pdf.CellFormat(nisWidth, 10, "NIS", "1", 0, "C", true, 0, "")
	pdf.CellFormat(rombelWidth, 10, "Rombel", "1", 0, "C", true, 0, "")

	// P1, P2, ... headers (first row)
	for i := 0; i < totalPertemuan; i++ {
		pdf.CellFormat(pertemuanWidth, 5, fmt.Sprintf("P%d", i+1), "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Row 2: Dates under P columns only
	// Skip No, Nama, NIS, Rombel columns (they are merged)
	pdf.SetX(10 + noWidth + namaWidth + nisWidth + rombelWidth)
	pdf.SetFont("Arial", "B", 8)

	for i, date := range dates {
		if i < totalPertemuan {
			pdf.CellFormat(pertemuanWidth, 5, date.Format("02"), "1", 0, "C", true, 0, "")
		}
	}

	// Fill remaining P columns with "-"
	for i := len(dates); i < totalPertemuan; i++ {
		pdf.CellFormat(pertemuanWidth, 5, "-", "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Data rows
	pdf.SetFont("Arial", "", 8)
	pdf.SetFillColor(255, 255, 255)

	for idx, student := range students {
		pdf.CellFormat(noWidth, 6, fmt.Sprintf("%d", idx+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(namaWidth, 6, student.Nama, "1", 0, "L", false, 0, "")
		pdf.CellFormat(nisWidth, 6, student.NIS, "1", 0, "C", false, 0, "")
		pdf.CellFormat(rombelWidth, 6, student.NamaRombel, "1", 0, "C", false, 0, "")

		// Attendance marks
		for i, date := range dates {
			if i < totalPertemuan {
				dateStr := date.Format("2006-01-02")
				status := student.Attendances[dateStr]

				mark := "-"
				switch status {
				case "hadir":
					mark = "v"
				case "alpa":
					mark = "A"
				case "sakit":
					mark = "S"
				case "izin":
					mark = "I"
				}

				pdf.CellFormat(pertemuanWidth, 6, mark, "1", 0, "C", false, 0, "")
			}
		}

		// Fill remaining columns
		for i := len(dates); i < totalPertemuan; i++ {
			pdf.CellFormat(pertemuanWidth, 6, "-", "1", 0, "C", false, 0, "")
		}

		pdf.Ln(-1)
	}

	// Generate filename
	filename := fmt.Sprintf("Absensi_Siswa_%s", response.NamaEkstrakurikuler)
	if req.Bulan != nil && req.Tahun != nil {
		filename += fmt.Sprintf("_%s_%d", bulanText, *req.Tahun)
	}
	filename += ".pdf"

	// Get PDF bytes
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, "", fmt.Errorf("failed to generate PDF: %v", err)
	}

	return buf.Bytes(), filename, nil
}
