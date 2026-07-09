package services

import (
	"bytes"
	"fmt"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// ExportExcel exports mutasi siswa data to Excel
func (s *MutasiSiswaServiceImpl) ExportExcel(req *dtos.MutasiSiswaExportExcelRequest) ([]byte, error) {
	// Get tahun pelajaran - use type assertion to access DB
	type dbAccessor interface {
		GetDB() *gorm.DB
	}
	
	repoImpl, ok := s.repository.(dbAccessor)
	if !ok {
		return nil, fmt.Errorf("cannot access database")
	}
	
	db := repoImpl.GetDB()
	
	var tahunPelajaran models.TahunPelajaran
	if err := db.First(&tahunPelajaran, req.TahunPelajaranID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("tahun pelajaran tidak ditemukan")
		}
		return nil, err
	}

	// Get mutasi siswa data
	data, err := s.repository.GetByTahunPelajaranAndSemester(req.TahunPelajaranID, req.Semester)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data: %w", err)
	}

	// Create new Excel file
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Data Calon Murid Baru"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 5)
	f.SetColWidth(sheetName, "B", "B", 30)
	f.SetColWidth(sheetName, "C", "C", 15)
	f.SetColWidth(sheetName, "D", "D", 20)
	f.SetColWidth(sheetName, "E", "E", 15)
	f.SetColWidth(sheetName, "F", "F", 15)
	f.SetColWidth(sheetName, "G", "G", 15)
	f.SetColWidth(sheetName, "H", "H", 15)

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

	// Header style (gray background)
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

	// Title
	title := fmt.Sprintf("DATA CALON MURID BARU SEMESTER %d TAHUN PELAJARAN %s", 
		req.Semester, tahunPelajaran.TahunPelajaran)
	f.SetCellValue(sheetName, "A1", title)
	f.MergeCell(sheetName, "A1", "H1")
	f.SetCellStyle(sheetName, "A1", "H1", titleStyle)
	f.SetRowHeight(sheetName, 1, 25)

	// Headers
	headers := []string{"No", "Nama CMB", "NISN", "Tempat Lahir", "Tanggal Lahir", "Jenis Kelamin", "Agama", "Pindahan Kelas"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 2)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}
	f.SetRowHeight(sheetName, 2, 20)

	// Data rows
	row := 3
	for idx, item := range data {
		// No
		cell, _ := excelize.CoordinatesToCellName(1, row)
		f.SetCellValue(sheetName, cell, idx+1)
		f.SetCellStyle(sheetName, cell, cell, dataCenterStyle)

		// Nama CMB
		cell, _ = excelize.CoordinatesToCellName(2, row)
		f.SetCellValue(sheetName, cell, item.NamaLengkap)
		f.SetCellStyle(sheetName, cell, cell, dataStyle)

		// NISN
		cell, _ = excelize.CoordinatesToCellName(3, row)
		nisn := ""
		if item.NISN != nil {
			nisn = *item.NISN
		}
		f.SetCellValue(sheetName, cell, nisn)
		f.SetCellStyle(sheetName, cell, cell, dataCenterStyle)

		// Tempat Lahir
		cell, _ = excelize.CoordinatesToCellName(4, row)
		f.SetCellValue(sheetName, cell, item.TempatLahir)
		f.SetCellStyle(sheetName, cell, cell, dataStyle)

		// Tanggal Lahir
		cell, _ = excelize.CoordinatesToCellName(5, row)
		f.SetCellValue(sheetName, cell, item.TanggalLahir.Format("02-01-2006"))
		f.SetCellStyle(sheetName, cell, cell, dataCenterStyle)

		// Jenis Kelamin
		cell, _ = excelize.CoordinatesToCellName(6, row)
		f.SetCellValue(sheetName, cell, item.JenisKelamin)
		f.SetCellStyle(sheetName, cell, cell, dataCenterStyle)

		// Agama
		cell, _ = excelize.CoordinatesToCellName(7, row)
		f.SetCellValue(sheetName, cell, item.Agama)
		f.SetCellStyle(sheetName, cell, cell, dataCenterStyle)

		// Pindahan Kelas
		cell, _ = excelize.CoordinatesToCellName(8, row)
		if item.PindahanKelas != nil {
			f.SetCellValue(sheetName, cell, *item.PindahanKelas)
		} else {
			f.SetCellValue(sheetName, cell, "")
		}
		f.SetCellStyle(sheetName, cell, cell, dataCenterStyle)

		f.SetRowHeight(sheetName, row, 20)
		row++
	}

	// Save to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
