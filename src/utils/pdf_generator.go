package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// PDFGenerator handles PDF generation
type PDFGenerator struct {
	pdf *gofpdf.Fpdf
}

// NewPDFGenerator creates a new PDF generator with Arial font
func NewPDFGenerator() *PDFGenerator {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	
	// Add Arial font (using built-in Helvetica as alternative for Arial)
	// gofpdf uses Helvetica which is very similar to Arial
	pdf.SetFont("Arial", "", 12)
	
	return &PDFGenerator{pdf: pdf}
}

// AddHeader adds the report header with kop surat, title and academic year
func (p *PDFGenerator) AddHeader(title string, tahunPelajaran string) {
	// Add kop surat image at the top (full width, ignore margins)
	pageWidth := 210.0 // A4 width in mm
	
	// Kop image - full page width from edge to edge
	kopURL := "https://pintu-storage.sdnsukapura01.sch.id/kop.png"
	kopWidth := pageWidth // Full width
	kopX := 0.0           // Start from left edge (no margin)
	kopY := 3.0           // Very top, minimal margin
	
	// Download and add kop image from URL
	resp, err := http.Get(kopURL)
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		
		// Read image data
		imgData, err := io.ReadAll(resp.Body)
		if err == nil {
			// Register image from bytes
			imgName := "kop"
			imageOpts := gofpdf.ImageOptions{
				ImageType: "PNG",
				ReadDpi:   true,
			}
			p.pdf.RegisterImageOptionsReader(imgName, imageOpts, bytes.NewReader(imgData))
			
			// Add image to PDF (full width, auto height)
			p.pdf.ImageOptions(imgName, kopX, kopY, kopWidth, 0, false, imageOpts, 0, "")
		}
	}
	
	// Check if there was an error
	kopHeight := 35.0 // Approximate space for kop
	if p.pdf.Error() != nil {
		p.pdf.ClearError()
		kopHeight = 0
	}
	
	// Move down after kop with more spacing
	p.pdf.SetY(kopY + kopHeight + 15) // 15mm spacing after kop (increased from 12mm)
	
	// Title - Bold, 14pt, centered
	p.pdf.SetFont("Arial", "B", 14)
	p.pdf.CellFormat(0, 10, title, "", 1, "C", false, 0, "")
	
	// Academic Year - Regular, 12pt, centered
	p.pdf.SetFont("Arial", "", 12)
	p.pdf.CellFormat(0, 8, "TAHUN PELAJARAN "+tahunPelajaran, "", 1, "C", false, 0, "")
	
	// Add spacing (back to 5mm)
	p.pdf.Ln(5)
}

// AddStudentInfoSimple adds student name and NISN with colon
func (p *PDFGenerator) AddStudentInfoSimple(nama, nisn string) {
	p.pdf.SetFont("Arial", "", 11)
	
	// Nama with colon
	labelWidth := 30.0
	colonWidth := 5.0
	
	p.pdf.CellFormat(labelWidth, 7, "Nama", "", 0, "L", false, 0, "")
	p.pdf.CellFormat(colonWidth, 7, ":", "", 0, "L", false, 0, "")
	p.pdf.CellFormat(0, 7, nama, "", 1, "L", false, 0, "")
	
	// NISN with colon
	p.pdf.CellFormat(labelWidth, 7, "NISN", "", 0, "L", false, 0, "")
	p.pdf.CellFormat(colonWidth, 7, ":", "", 0, "L", false, 0, "")
	p.pdf.CellFormat(0, 7, nisn, "", 1, "L", false, 0, "")
	
	// Add spacing
	p.pdf.Ln(5)
}

// AddNilaiTableWithMergedAverage adds nilai table with merged rata-rata column
func (p *PDFGenerator) AddNilaiTableWithMergedAverage(nilaiList []interface{}, rataRata float64) {
	p.pdf.SetFont("Arial", "B", 11)
	
	// Define light gray color for header (RGB: 220, 220, 220)
	p.pdf.SetFillColor(220, 220, 220)
	
	// Table header
	colWidth1 := 20.0  // NO
	colWidth2 := 80.0  // MATA PELAJARAN
	colWidth3 := 35.0  // NILAI
	colWidth4 := 55.0  // NILAI RATA-RATA AKHIR
	headerHeight := 8.0   // Header row height
	dataRowHeight := 10.0 // Increased data row height for better spacing
	
	p.pdf.CellFormat(colWidth1, headerHeight, "NO", "1", 0, "C", true, 0, "")
	p.pdf.CellFormat(colWidth2, headerHeight, "MATA PELAJARAN", "1", 0, "C", true, 0, "")
	p.pdf.CellFormat(colWidth3, headerHeight, "NILAI", "1", 0, "C", true, 0, "")
	p.pdf.CellFormat(colWidth4, headerHeight, "NILAI RATA-RATA AKHIR", "1", 1, "C", true, 0, "")
	
	// Reset fill color to white for data rows
	p.pdf.SetFillColor(255, 255, 255)
	
	// Table rows
	p.pdf.SetFont("Arial", "", 11)
	
	// Calculate total height needed for all rows
	rowCount := len(nilaiList)
	
	for i, item := range nilaiList {
		if mapItem, ok := item.(map[string]interface{}); ok {
			mapel := ""
			var nilai interface{}
			
			if m, exists := mapItem["mapel"]; exists {
				mapel = fmt.Sprintf("%v", m)
			}
			if n, exists := mapItem["nilai"]; exists {
				nilai = n
			}
			
			// NO column - always has all borders for each row
			if i == 0 {
				// First row - top, left, right borders
				p.pdf.CellFormat(colWidth1, dataRowHeight, fmt.Sprintf("%d", i+1), "LTR", 0, "C", false, 0, "")
			} else if i == rowCount-1 {
				// Last row - all borders including bottom
				p.pdf.CellFormat(colWidth1, dataRowHeight, fmt.Sprintf("%d", i+1), "LBRT", 0, "C", false, 0, "")
			} else {
				// Middle rows - left and right only
				p.pdf.CellFormat(colWidth1, dataRowHeight, fmt.Sprintf("%d", i+1), "LR", 0, "C", false, 0, "")
			}
			
			// MATA PELAJARAN column
			if i == 0 {
				// First row - top, left, right borders
				p.pdf.CellFormat(colWidth2, dataRowHeight, mapel, "LTR", 0, "L", false, 0, "")
			} else if i == rowCount-1 {
				// Last row - all borders including bottom
				p.pdf.CellFormat(colWidth2, dataRowHeight, mapel, "LBRT", 0, "L", false, 0, "")
			} else {
				// Middle rows - left and right only
				p.pdf.CellFormat(colWidth2, dataRowHeight, mapel, "LR", 0, "L", false, 0, "")
			}
			
			// NILAI column
			var nilaiStr string
			switch v := nilai.(type) {
			case float64:
				nilaiStr = fmt.Sprintf("%.2f", v)
			case int:
				nilaiStr = fmt.Sprintf("%d", v)
			case int64:
				nilaiStr = fmt.Sprintf("%d", v)
			default:
				nilaiStr = fmt.Sprintf("%v", v)
			}
			
			if i == 0 {
				// First row - top, left, right borders
				p.pdf.CellFormat(colWidth3, dataRowHeight, nilaiStr, "LTR", 0, "C", false, 0, "")
			} else if i == rowCount-1 {
				// Last row - all borders including bottom
				p.pdf.CellFormat(colWidth3, dataRowHeight, nilaiStr, "LBRT", 0, "C", false, 0, "")
			} else {
				// Middle rows - left and right only
				p.pdf.CellFormat(colWidth3, dataRowHeight, nilaiStr, "LR", 0, "C", false, 0, "")
			}
			
			// NILAI RATA-RATA AKHIR column (merged for all rows)
			if i == 0 {
				// First row - draw top border only
				p.pdf.CellFormat(colWidth4, dataRowHeight, "", "LTR", 1, "C", false, 0, "")
			} else if i == rowCount-1 {
				// Last row - draw bottom border
				p.pdf.CellFormat(colWidth4, dataRowHeight, "", "LBR", 1, "C", false, 0, "")
			} else {
				// Middle rows - only side borders
				p.pdf.CellFormat(colWidth4, dataRowHeight, "", "LR", 1, "C", false, 0, "")
			}
		}
	}
	
	// Add rata-rata value in the center of merged cell
	x := p.pdf.GetX()
	y := p.pdf.GetY()
	
	// Calculate middle position
	totalTableHeight := dataRowHeight * float64(rowCount)
	middleRowY := y - totalTableHeight + (totalTableHeight / 2) - (dataRowHeight / 2)
	
	p.pdf.SetXY(x+colWidth1+colWidth2+colWidth3, middleRowY)
	p.pdf.SetFont("Arial", "B", 12)
	p.pdf.CellFormat(colWidth4, dataRowHeight, fmt.Sprintf("%.2f", rataRata), "", 1, "C", false, 0, "")
	
	// Move cursor to after the table
	p.pdf.SetXY(x, y)
	
	// Add spacing
	p.pdf.Ln(5)
}

// AddMotivationalText adds motivational message
func (p *PDFGenerator) AddMotivationalText() {
	p.pdf.SetFont("Arial", "", 11)
	
	text := "Teruslah belajar dengan giat dan penuh semangat agar prestasi yang kamu cita-citakan dapat tercapai sesuai dengan harapan. " +
		"Ingatlah bahwa setiap usaha yang kamu lakukan hari ini adalah investasi untuk masa depan yang lebih cerah."
	
	// Use MultiCell for text wrapping
	p.pdf.MultiCell(0, 6, text, "", "L", false)
	
	// Add spacing
	p.pdf.Ln(10)
}

// AddSignatures adds signature section for Orang Tua and Kepala Sekolah
func (p *PDFGenerator) AddSignatures(tanggalPengumuman time.Time, namaKepsek string, ttdKepsekURL string) {
	// Calculate column widths
	pageWidth := 210.0 // A4 width in mm
	margin := 10.0
	usableWidth := pageWidth - (2 * margin)
	colWidth := usableWidth / 2
	
	currentY := p.pdf.GetY()
	
	// Format date for right column
	months := []string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	formattedDate := fmt.Sprintf("Jakarta, %d %s %d", 
		tanggalPengumuman.Day(), 
		months[tanggalPengumuman.Month()-1], 
		tanggalPengumuman.Year())
	
	// Right column - Jakarta date and Kepala Sekolah (shift right by 20mm)
	rightColX := margin + colWidth + 20 // Shifted right by 20mm
	p.pdf.SetFont("Arial", "", 11)
	p.pdf.SetXY(rightColX, currentY)
	p.pdf.CellFormat(colWidth, 5, formattedDate, "", 1, "L", false, 0, "") // Reduced from 7 to 5
	
	// Left column - Orang Tua Murid (shift right by 20mm and aligned with Kepala SDN line)
	leftColX := margin + 20 // Shifted right by 20mm
	p.pdf.SetXY(leftColX, currentY+5) // Reduced from +7 to +5
	p.pdf.CellFormat(colWidth, 5, "Orang Tua Murid", "", 1, "L", false, 0, "")
	
	// Right column - Kepala Sekolah title
	p.pdf.SetXY(rightColX, currentY+5)
	p.pdf.CellFormat(colWidth, 5, "Kepala SDN Sukapura 01", "", 1, "L", false, 0, "")
	
	// Right column - Add signature image from URL (larger size)
	imgHeight := 28.0 // Increased from 20mm to 28mm
	imgWidth := 40.0  // Increased from 30mm to 40mm
	
	// Position image - not overlapping, below "Kepala SDN Sukapura 01" text
	imgY := currentY + 12 // Moved down (was currentY + 8 for overlapping)
	
	// Download and add TTD image from URL
	if ttdKepsekURL != "" {
		resp, err := http.Get(ttdKepsekURL)
		if err == nil && resp.StatusCode == 200 {
			defer resp.Body.Close()
			
			// Read image data
			imgData, err := io.ReadAll(resp.Body)
			if err == nil {
				// Register image from bytes
				imgName := "ttd_kepsek"
				imageOpts := gofpdf.ImageOptions{
					ImageType: "PNG",
					ReadDpi:   true,
				}
				p.pdf.RegisterImageOptionsReader(imgName, imageOpts, bytes.NewReader(imgData))
				
				// Add image to PDF
				p.pdf.ImageOptions(imgName, rightColX, imgY, imgWidth, imgHeight, false, imageOpts, 0, "")
			}
		}
	}
	
	// Check if there was an error adding the image
	if p.pdf.Error() != nil {
		// If image fails, just skip it
		p.pdf.ClearError()
	}
	
	// Position for kepala sekolah name (after image, reduced spacing)
	nameY := imgY + imgHeight + 1 // Reduced from +2 to +1
	
	// Right column - Kepala Sekolah name (bold)
	p.pdf.SetXY(rightColX, nameY)
	p.pdf.SetFont("Arial", "B", 11)
	p.pdf.CellFormat(colWidth-20, 5, namaKepsek, "", 1, "L", false, 0, "") // Reduced from 7 to 5
	
	// Left column - Dotted line for parent signature aligned with nama kepsek
	dots := "..................................."
	dotsWidth := p.pdf.GetStringWidth(dots)
	
	// Draw the dots at the same Y position as nama kepsek
	p.pdf.SetXY(leftColX, nameY)
	p.pdf.CellFormat(dotsWidth, 5, dots, "", 1, "L", false, 0, "")
	
	// Right column - NIP (not bold, reduced spacing)
	p.pdf.SetXY(rightColX, nameY+5) // Reduced from +7 to +5
	p.pdf.SetFont("Arial", "", 11)
	p.pdf.CellFormat(colWidth-20, 5, "NIP. 198805102014032004", "", 1, "L", false, 0, "")
}

// GetBytes returns the PDF as byte array
func (p *PDFGenerator) GetBytes() ([]byte, error) {
	var buf bytes.Buffer
	err := p.pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("error generating PDF: %v", err)
	}
	return buf.Bytes(), nil
}

