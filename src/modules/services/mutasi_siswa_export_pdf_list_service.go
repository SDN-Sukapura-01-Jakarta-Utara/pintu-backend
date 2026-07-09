package services

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"

	"github.com/jung-kurt/gofpdf"
	"gorm.io/gorm"
)

// ExportListPDF exports mutasi siswa list to PDF with table format
func (s *MutasiSiswaServiceImpl) ExportListPDF(req *dtos.MutasiSiswaExportExcelRequest) ([]byte, error) {
	// Get tahun pelajaran
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

	// Create PDF - Portrait
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.SetAutoPageBreak(true, 10)
	pdf.AddPage()

	// Download and add kop
	kopURL := "https://pintu-storage.sdnsukapura01.sch.id/kop.png"
	resp, err := http.Get(kopURL)
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		kopBytes, _ := io.ReadAll(resp.Body)
		
		// Register kop image
		opt := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}
		pdf.RegisterImageOptionsReader("kop", opt, bytes.NewReader(kopBytes))
		
		// Add kop image - full width
		pdf.ImageOptions("kop", 10, 5, 190, 0, false, opt, 0, "")
		pdf.Ln(45) // Increased space to push title down
	} else {
		pdf.Ln(5)
	}

	// Title
	pdf.SetFont("Times", "B", 12) // Slightly smaller font for portrait
	title := fmt.Sprintf("DATA CALON MURID BARU SEMESTER %d\nTAHUN PELAJARAN %s", 
		req.Semester, tahunPelajaran.TahunPelajaran)
	pdf.MultiCell(190, 6, title, "", "C", false)
	pdf.Ln(3)

	// Table column widths - adjusted for portrait (210mm width - 20mm margins = 190mm)
	colWidths := []float64{8, 40, 20, 25, 22, 18, 18, 15} // Total: 166mm
	leftMargin := (210.0 - 166.0) / 2.0 // Center the table

	// Header style
	pdf.SetFillColor(211, 211, 211) // Gray background
	pdf.SetFont("Times", "B", 8) // Smaller font for portrait
	pdf.SetTextColor(0, 0, 0)

	// Table headers
	headers := []string{"No", "Nama CMB", "NISN", "Tempat Lahir", "Tgl Lahir", "JK", "Agama", "Kelas"}
	
	pdf.SetX(leftMargin)
	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 7, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Data rows
	pdf.SetFont("Times", "", 7) // Smaller font for data
	pdf.SetFillColor(255, 255, 255)

	for idx, item := range data {
		// Check if new page is needed
		if pdf.GetY() > 250 { // Adjusted for portrait
			pdf.AddPage()
			
			// Reprint header on new page
			pdf.SetFont("Times", "B", 8)
			pdf.SetFillColor(211, 211, 211)
			pdf.SetX(leftMargin)
			for i, header := range headers {
				pdf.CellFormat(colWidths[i], 7, header, "1", 0, "C", true, 0, "")
			}
			pdf.Ln(-1)
			pdf.SetFont("Times", "", 7)
			pdf.SetFillColor(255, 255, 255)
		}

		// No
		pdf.SetX(leftMargin)
		pdf.CellFormat(colWidths[0], 6, fmt.Sprintf("%d", idx+1), "1", 0, "C", false, 0, "")

		// Nama CMB
		pdf.CellFormat(colWidths[1], 6, item.NamaLengkap, "1", 0, "L", false, 0, "")

		// NISN
		nisn := ""
		if item.NISN != nil {
			nisn = *item.NISN
		}
		pdf.CellFormat(colWidths[2], 6, nisn, "1", 0, "C", false, 0, "")

		// Tempat Lahir
		pdf.CellFormat(colWidths[3], 6, item.TempatLahir, "1", 0, "L", false, 0, "")

		// Tanggal Lahir
		tanggalLahir := item.TanggalLahir.Format("02-01-2006")
		pdf.CellFormat(colWidths[4], 6, tanggalLahir, "1", 0, "C", false, 0, "")

		// Jenis Kelamin
		pdf.CellFormat(colWidths[5], 6, item.JenisKelamin, "1", 0, "C", false, 0, "")

		// Agama
		pdf.CellFormat(colWidths[6], 6, item.Agama, "1", 0, "C", false, 0, "")

		// Pindahan Kelas
		pindahanKelas := ""
		if item.PindahanKelas != nil {
			pindahanKelas = fmt.Sprintf("%d", *item.PindahanKelas)
		}
		pdf.CellFormat(colWidths[7], 6, pindahanKelas, "1", 0, "C", false, 0, "")

		pdf.Ln(-1)
	}

	// Output PDF
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
