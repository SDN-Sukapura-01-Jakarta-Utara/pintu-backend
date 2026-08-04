package services

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// ExportExcelPerEkskul exports ekstrakurikuler data per ekskul to Excel
func (s *PesertaDidikEkstrakurikulerServiceImpl) ExportExcelPerEkskul(req *dtos.ExportExcelPerEkskulRequest) ([]byte, string, error) {
	// Load WIB timezone
	wibLoc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		wibLoc = time.FixedZone("WIB", 7*60*60) // Fallback to UTC+7
	}

	// Get data from repository (no pagination - get all)
	registrations, _, err := s.repository.GetRekapPerEkskul(
		req.TahunPelajaranID,
		"",    // nama
		"",    // nis
		0,     // rombelID
		req.EkstrakurikulerID,
		10000, // large limit to get all
		0,     // offset
	)
	if err != nil {
		return nil, "", fmt.Errorf("gagal mengambil data rekapitulasi")
	}

	// Get tahun pelajaran info
	type dbAccessor interface {
		GetDB() *gorm.DB
	}
	
	repoImpl, ok := s.repository.(dbAccessor)
	if !ok {
		return nil, "", fmt.Errorf("cannot access database")
	}
	
	db := repoImpl.GetDB()
	
	var tahunPelajaran models.TahunPelajaran
	if err := db.First(&tahunPelajaran, req.TahunPelajaranID).Error; err != nil {
		return nil, "", fmt.Errorf("tahun pelajaran tidak ditemukan")
	}

	// Group by ekstrakurikuler
	ekskulMap := make(map[uint]*struct {
		ID       uint
		Nama     string
		Kategori string
		Siswa    []struct {
			Nama       string
			NIS        string
			NamaRombel string
		}
	})
	
	for _, reg := range registrations {
		if reg.Ekstrakurikuler == nil || reg.PesertaDidikRombel == nil || reg.PesertaDidikRombel.PesertaDidik == nil {
			continue
		}
		
		ekskulID := reg.EkstrakurikulerID
		
		// Initialize ekstrakurikuler if not exists
		if _, exists := ekskulMap[ekskulID]; !exists {
			ekskulMap[ekskulID] = &struct {
				ID       uint
				Nama     string
				Kategori string
				Siswa    []struct {
					Nama       string
					NIS        string
					NamaRombel string
				}
			}{
				ID:       ekskulID,
				Nama:     reg.Ekstrakurikuler.Name,
				Kategori: reg.Ekstrakurikuler.Kategori,
			}
		}
		
		// Add student
		namaRombel := ""
		if reg.PesertaDidikRombel.Rombel != nil {
			namaRombel = reg.PesertaDidikRombel.Rombel.Name
		}
		
		ekskulMap[ekskulID].Siswa = append(ekskulMap[ekskulID].Siswa, struct {
			Nama       string
			NIS        string
			NamaRombel string
		}{
			Nama:       reg.PesertaDidikRombel.PesertaDidik.Nama,
			NIS:        reg.PesertaDidikRombel.PesertaDidik.NIS,
			NamaRombel: namaRombel,
		})
	}

	// Create Excel file
	f := excelize.NewFile()
	defer f.Close()

	// Process each ekstrakurikuler as separate sheet
	sheetIndex := 0
	for _, ekskul := range ekskulMap {
		sheetName := fmt.Sprintf("Ekskul_%d", ekskul.ID)
		if sheetIndex == 0 {
			f.SetSheetName("Sheet1", sheetName)
		} else {
			_, err := f.NewSheet(sheetName)
			if err != nil {
				return nil, "", err
			}
		}
		
		// Set column widths
		f.SetColWidth(sheetName, "A", "A", 5)
		f.SetColWidth(sheetName, "B", "B", 35)
		f.SetColWidth(sheetName, "C", "C", 15)
		f.SetColWidth(sheetName, "D", "D", 15)

		// Title style
		titleStyle, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Bold: true,
				Size: 14,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
		})

		// Header style
		headerStyle, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Bold: true,
				Size: 11,
			},
			Fill: excelize.Fill{
				Type:    "pattern",
				Color:   []string{"#D3D3D3"},
				Pattern: 1,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
			},
		})

		// Data style
		dataStyle, _ := f.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: "left",
				Vertical:   "center",
				WrapText:   true,
			},
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
			},
		})

		// Center data style
		dataCenterStyle, _ := f.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
			},
		})

		// Title - Row 1: DATA EKSKUL [NAMA EKSKUL]
		title1 := fmt.Sprintf("DATA EKSKUL %s", strings.ToUpper(ekskul.Nama))
		f.SetCellValue(sheetName, "A1", title1)
		f.MergeCell(sheetName, "A1", "D1")
		f.SetCellStyle(sheetName, "A1", "D1", titleStyle)
		f.SetRowHeight(sheetName, 1, 25)

		// Title - Row 2: TAHUN AJARAN [TAHUN]
		title2 := fmt.Sprintf("TAHUN AJARAN %s", strings.ToUpper(tahunPelajaran.TahunPelajaran))
		f.SetCellValue(sheetName, "A2", title2)
		f.MergeCell(sheetName, "A2", "D2")
		f.SetCellStyle(sheetName, "A2", "D2", titleStyle)
		f.SetRowHeight(sheetName, 2, 25)

		// Headers
		headers := []string{"No", "Nama", "NIS", "Rombel"}
		for i, header := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 3)
			f.SetCellValue(sheetName, cell, header)
			f.SetCellStyle(sheetName, cell, cell, headerStyle)
		}
		f.SetRowHeight(sheetName, 3, 20)

		// Data rows
		row := 4
		for idx, siswa := range ekskul.Siswa {
			// No
			cell, _ := excelize.CoordinatesToCellName(1, row)
			f.SetCellValue(sheetName, cell, idx+1)
			f.SetCellStyle(sheetName, cell, cell, dataCenterStyle)

			// Nama
			cell, _ = excelize.CoordinatesToCellName(2, row)
			f.SetCellValue(sheetName, cell, siswa.Nama)
			f.SetCellStyle(sheetName, cell, cell, dataStyle)

			// NIS
			cell, _ = excelize.CoordinatesToCellName(3, row)
			f.SetCellValue(sheetName, cell, siswa.NIS)
			f.SetCellStyle(sheetName, cell, cell, dataCenterStyle)

			// Rombel
			cell, _ = excelize.CoordinatesToCellName(4, row)
			f.SetCellValue(sheetName, cell, siswa.NamaRombel)
			f.SetCellStyle(sheetName, cell, cell, dataCenterStyle)

			f.SetRowHeight(sheetName, row, 20)
			row++
		}

		sheetIndex++
	}

	// Save to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, "", err
	}

	// Generate filename
	filename := fmt.Sprintf("Data_Ekskul_Per_Ekstrakurikuler_%s.xlsx", 
		time.Now().In(wibLoc).Format("20060102_150405"))

	return buf.Bytes(), filename, nil
}

// ExportExcelPerRombel exports ekstrakurikuler data per rombel to Excel
func (s *PesertaDidikEkstrakurikulerServiceImpl) ExportExcelPerRombel(req *dtos.ExportExcelPerRombelRequest) ([]byte, string, error) {
	// Load WIB timezone
	wibLoc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		wibLoc = time.FixedZone("WIB", 7*60*60) // Fallback to UTC+7
	}

	// Get data from repository (no pagination - get all)
	registrations, _, err := s.repository.GetRekapPerRombel(
		req.RombelID,
		req.TahunPelajaranID,
		"",    // nama
		"",    // nis
		10000, // large limit to get all
		0,     // offset
	)
	if err != nil {
		return nil, "", fmt.Errorf("gagal mengambil data rekapitulasi")
	}

	// Get tahun pelajaran and rombel info
	type dbAccessor interface {
		GetDB() *gorm.DB
	}
	
	repoImpl, ok := s.repository.(dbAccessor)
	if !ok {
		return nil, "", fmt.Errorf("cannot access database")
	}
	
	db := repoImpl.GetDB()
	
	var tahunPelajaran models.TahunPelajaran
	if err := db.First(&tahunPelajaran, req.TahunPelajaranID).Error; err != nil {
		return nil, "", fmt.Errorf("tahun pelajaran tidak ditemukan")
	}

	var rombel models.Rombel
	if err := db.First(&rombel, req.RombelID).Error; err != nil {
		return nil, "", fmt.Errorf("rombel tidak ditemukan")
	}

	// Group by student
	siswaMap := make(map[uint]*struct {
		Nama           string
		NIS            string
		Ekstrakurikuler []string
		TotalEkskul    int
	})
	
	for _, reg := range registrations {
		if reg.PesertaDidikRombel == nil || reg.PesertaDidikRombel.PesertaDidik == nil {
			continue
		}
		
		studentID := reg.PesertaDidikRombelID
		
		// Initialize student if not exists
		if _, exists := siswaMap[studentID]; !exists {
			siswaMap[studentID] = &struct {
				Nama           string
				NIS            string
				Ekstrakurikuler []string
				TotalEkskul    int
			}{
				Nama:           reg.PesertaDidikRombel.PesertaDidik.Nama,
				NIS:            reg.PesertaDidikRombel.PesertaDidik.NIS,
				Ekstrakurikuler: []string{},
			}
		}
		
		// Add ekstrakurikuler if exists
		if reg.Ekstrakurikuler != nil {
			siswaMap[studentID].Ekstrakurikuler = append(siswaMap[studentID].Ekstrakurikuler, reg.Ekstrakurikuler.Name)
			siswaMap[studentID].TotalEkskul++
		}
	}

	// Convert to slice and sort by nama
	type SiswaData struct {
		Nama           string
		NIS            string
		Ekstrakurikuler string
		TotalEkskul    int
	}
	
	var siswaList []SiswaData
	for _, siswa := range siswaMap {
		ekskulStr := "-"
		if len(siswa.Ekstrakurikuler) > 0 {
			ekskulStr = ""
			for i, ekskul := range siswa.Ekstrakurikuler {
				if i > 0 {
					ekskulStr += ", "
				}
				ekskulStr += ekskul
			}
		}
		
		siswaList = append(siswaList, SiswaData{
			Nama:           siswa.Nama,
			NIS:            siswa.NIS,
			Ekstrakurikuler: ekskulStr,
			TotalEkskul:    siswa.TotalEkskul,
		})
	}

	// Sort by Nama A-Z
	sort.Slice(siswaList, func(i, j int) bool {
		return siswaList[i].Nama < siswaList[j].Nama
	})

	// Create Excel file
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Data Ekskul"
	f.SetSheetName("Sheet1", sheetName)

	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 5)
	f.SetColWidth(sheetName, "B", "B", 35)
	f.SetColWidth(sheetName, "C", "C", 15)
	f.SetColWidth(sheetName, "D", "D", 50)
	f.SetColWidth(sheetName, "E", "E", 15)

	// Title style
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 14,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})

	// Header style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 11,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#D3D3D3"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	// Data style
	dataStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	// Center data style
	dataCenterStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	// Title - Row 1: DATA EKSKUL KELAS [NAMA ROMBEL]
	title1 := fmt.Sprintf("DATA EKSKUL KELAS %s", strings.ToUpper(rombel.Name))
	f.SetCellValue(sheetName, "A1", title1)
	f.MergeCell(sheetName, "A1", "E1")
	f.SetCellStyle(sheetName, "A1", "E1", titleStyle)
	f.SetRowHeight(sheetName, 1, 25)

	// Title - Row 2: TAHUN AJARAN [TAHUN]
	title2 := fmt.Sprintf("TAHUN AJARAN %s", strings.ToUpper(tahunPelajaran.TahunPelajaran))
	f.SetCellValue(sheetName, "A2", title2)
	f.MergeCell(sheetName, "A2", "E2")
	f.SetCellStyle(sheetName, "A2", "E2", titleStyle)
	f.SetRowHeight(sheetName, 2, 25)

	// Headers
	headers := []string{"No", "Nama", "NIS", "Ekstrakurikuler", "Total Ekskul"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 3)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}
	f.SetRowHeight(sheetName, 3, 20)

	// Data rows
	row := 4
	for idx, siswa := range siswaList {
		// No
		cell, _ := excelize.CoordinatesToCellName(1, row)
		f.SetCellValue(sheetName, cell, idx+1)
		f.SetCellStyle(sheetName, cell, cell, dataCenterStyle)

		// Nama
		cell, _ = excelize.CoordinatesToCellName(2, row)
		f.SetCellValue(sheetName, cell, siswa.Nama)
		f.SetCellStyle(sheetName, cell, cell, dataStyle)

		// NIS
		cell, _ = excelize.CoordinatesToCellName(3, row)
		f.SetCellValue(sheetName, cell, siswa.NIS)
		f.SetCellStyle(sheetName, cell, cell, dataCenterStyle)

		// Ekstrakurikuler
		cell, _ = excelize.CoordinatesToCellName(4, row)
		f.SetCellValue(sheetName, cell, siswa.Ekstrakurikuler)
		f.SetCellStyle(sheetName, cell, cell, dataStyle)

		// Total Ekskul
		cell, _ = excelize.CoordinatesToCellName(5, row)
		f.SetCellValue(sheetName, cell, siswa.TotalEkskul)
		f.SetCellStyle(sheetName, cell, cell, dataCenterStyle)

		f.SetRowHeight(sheetName, row, 20)
		row++
	}

	// Save to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, "", err
	}

	// Generate filename
	filename := fmt.Sprintf("Data_Ekskul_Kelas_%s_%s.xlsx", 
		rombel.Name, time.Now().In(wibLoc).Format("20060102_150405"))

	return buf.Bytes(), filename, nil
}
