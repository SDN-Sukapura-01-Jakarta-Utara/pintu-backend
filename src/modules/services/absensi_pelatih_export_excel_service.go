package services

import (
	"fmt"
	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"sort"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// DownloadExcelAbsensiPelatih generates Excel file for pelatih attendance
func (s *AbsensiEkskulService) DownloadExcelAbsensiPelatih(req *dtos.AbsensiPelatihExportRequest) (*excelize.File, string, error) {
	// Define PelatihInfo struct
	type PelatihInfo struct {
		ID      uint
		Nama    string
		Telepon string
	}

	// Helper function to format date in Indonesian
	formatTanggalIndonesia := func(t time.Time) string {
		bulanNames := []string{"", "Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
		return fmt.Sprintf("%d %s %d", t.Day(), bulanNames[t.Month()], t.Year())
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

	// Build query for kegiatan
	var kegiatanList []models.KegiatanEkskul
	query := s.db.Model(&models.KegiatanEkskul{}).
		Preload("Ekstrakurikuler").
		Preload("TahunPelajaran").
		Preload("AbsensiPelatih.Pelatih").
		Where("tahun_pelajaran_id = ?", req.TahunPelajaranID)

	if req.EkstrakurikulerID != nil {
		query = query.Where("ekstrakurikuler_id = ?", *req.EkstrakurikulerID)
	}

	if req.Bulan != nil && req.Tahun != nil {
		startDate := time.Date(*req.Tahun, time.Month(*req.Bulan), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)
		query = query.Where("tanggal_kegiatan BETWEEN ? AND ?", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	} else if req.Tahun != nil {
		query = query.Where("YEAR(tanggal_kegiatan) = ?", *req.Tahun)
	}

	if err := query.Order("tanggal_kegiatan ASC, waktu_mulai ASC").Find(&kegiatanList).Error; err != nil {
		return nil, "", err
	}

	// Get all unique pelatih from kegiatan list
	pelatihMap := make(map[uint]PelatihInfo)
	for _, kegiatan := range kegiatanList {
		for _, absensi := range kegiatan.AbsensiPelatih {
			if absensi.Pelatih != nil {
				// Filter by pelatih_id if specified
				if req.PelatihID != nil && absensi.Pelatih.ID != *req.PelatihID {
					continue
				}
				pelatihMap[absensi.Pelatih.ID] = PelatihInfo{
					ID:      absensi.Pelatih.ID,
					Nama:    absensi.Pelatih.Nama,
					Telepon: absensi.Pelatih.Telepon,
				}
			}
		}
	}

	// Convert map to slice
	var pelatihList []PelatihInfo
	for _, pelatih := range pelatihMap {
		pelatihList = append(pelatihList, pelatih)
	}

	// Sort pelatih by name
	sort.Slice(pelatihList, func(i, j int) bool {
		return pelatihList[i].Nama < pelatihList[j].Nama
	})

	// Get ekstrakurikuler name for title
	ekskulName := "SEMUA EKSTRAKURIKULER"
	if req.EkstrakurikulerID != nil && len(kegiatanList) > 0 {
		ekskulName = kegiatanList[0].Ekstrakurikuler.Name
	}

	// Create Excel file
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheetName := "Absensi Pelatih"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, "", err
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 5)  // No
	f.SetColWidth(sheetName, "B", "B", 30) // Nama Pelatih
	f.SetColWidth(sheetName, "C", "C", 12) // Tanggal
	f.SetColWidth(sheetName, "D", "D", 10) // Waktu
	f.SetColWidth(sheetName, "E", "E", 40) // Materi Kegiatan
	f.SetColWidth(sheetName, "F", "F", 15) // Status Kehadiran
	f.SetColWidth(sheetName, "G", "G", 20) // TTD

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

	currentRow := 1

	// Title
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("DAFTAR HADIR PELATIH EKSTRAKURIKULER %s", strings.ToUpper(ekskulName)))
	f.MergeCell(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("G%d", currentRow))
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("A%d", currentRow), titleStyle)
	currentRow++

	// Subtitle
	if tahunText != "" {
		subtitle := fmt.Sprintf("BULAN %s TAHUN %s", bulanText, tahunText)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow), subtitle)
		f.MergeCell(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("G%d", currentRow))
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("A%d", currentRow), subtitleStyle)
		currentRow++
	}

	currentRow++ // Empty row

	// Header
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow), "No")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", currentRow), "Nama Pelatih")
	f.SetCellValue(sheetName, fmt.Sprintf("C%d", currentRow), "Tanggal")
	f.SetCellValue(sheetName, fmt.Sprintf("D%d", currentRow), "Waktu")
	f.SetCellValue(sheetName, fmt.Sprintf("E%d", currentRow), "Materi Kegiatan")
	f.SetCellValue(sheetName, fmt.Sprintf("F%d", currentRow), "Status Kehadiran")
	f.SetCellValue(sheetName, fmt.Sprintf("G%d", currentRow), "TTD")

	for col := 'A'; col <= 'G'; col++ {
		f.SetCellStyle(sheetName, fmt.Sprintf("%c%d", col, currentRow), fmt.Sprintf("%c%d", col, currentRow), headerStyle)
	}
	currentRow++

	// Data rows
	rowNum := 1
	for _, pelatih := range pelatihList {
		for _, kegiatan := range kegiatanList {
			tanggal := formatTanggalIndonesia(kegiatan.TanggalKegiatan)
			waktuMulai := ""
			if kegiatan.WaktuMulai != nil {
				waktuMulai = *kegiatan.WaktuMulai
			}
			waktuSelesai := ""
			if kegiatan.WaktuSelesai != nil {
				waktuSelesai = *kegiatan.WaktuSelesai
			}
			waktu := waktuMulai
			if waktuSelesai != "" {
				waktu += " - " + waktuSelesai
			}
			materi := kegiatan.MateriKegiatan

			// Check if pelatih hadir
			isHadir := false
			for _, absensi := range kegiatan.AbsensiPelatih {
				if absensi.PelatihID == pelatih.ID {
					isHadir = true
					break
				}
			}

			status := "Hadir"
			if !isHadir {
				status = "Tidak Hadir"
			}

			f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow), rowNum)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", currentRow), pelatih.Nama)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", currentRow), tanggal)
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", currentRow), waktu)
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", currentRow), materi)
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", currentRow), status)
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", currentRow), "")

			for col := 'A'; col <= 'G'; col++ {
				f.SetCellStyle(sheetName, fmt.Sprintf("%c%d", col, currentRow), fmt.Sprintf("%c%d", col, currentRow), dataStyle)
			}

			currentRow++
			rowNum++
		}
	}

	// Generate filename
	filename := fmt.Sprintf("Absensi_Pelatih_%s", ekskulName)
	if req.Bulan != nil && req.Tahun != nil {
		filename += fmt.Sprintf("_%s_%d", bulanText, *req.Tahun)
	}
	filename += ".xlsx"

	return f, filename, nil
}
