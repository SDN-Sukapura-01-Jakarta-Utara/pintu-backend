package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// DownloadPDFDokumentasiEkskul generates PDF document for ekstrakurikuler documentation
func (s *AbsensiEkskulService) DownloadPDFDokumentasiEkskul(req *dtos.KegiatanEkskulDownloadWordRequest) ([]byte, string, error) {
	// Get ekstrakurikuler name and tahun pelajaran
	var ekstrakurikuler models.Ekstrakurikuler
	if err := s.db.First(&ekstrakurikuler, req.EkstrakurikulerID).Error; err != nil {
		return nil, "", fmt.Errorf("ekstrakurikuler not found: %v", err)
	}

	var tahunPelajaran models.TahunPelajaran
	if err := s.db.First(&tahunPelajaran, req.TahunPelajaranID).Error; err != nil {
		return nil, "", fmt.Errorf("tahun pelajaran not found: %v", err)
	}

	// Get all kegiatan for the month
	var kegiatanList []models.KegiatanEkskul
	query := s.db.Where("ekstrakurikuler_id = ? AND tahun_pelajaran_id = ? AND deleted_at IS NULL",
		req.EkstrakurikulerID, req.TahunPelajaranID)
	
	// Filter by month and year
	query = query.Where("EXTRACT(MONTH FROM tanggal_kegiatan) = ? AND EXTRACT(YEAR FROM tanggal_kegiatan) = ?",
		req.Bulan, req.Tahun)
	
	if err := query.Order("tanggal_kegiatan ASC, waktu_mulai ASC").Find(&kegiatanList).Error; err != nil {
		return nil, "", fmt.Errorf("failed to get kegiatan: %v", err)
	}

	if len(kegiatanList) == 0 {
		return nil, "", fmt.Errorf("tidak ada kegiatan pada bulan dan tahun yang dipilih")
	}

	// Month names in Indonesian
	bulanNames := []string{"", "JANUARI", "FEBRUARI", "MARET", "APRIL", "MEI", "JUNI", "JULI", "AGUSTUS", "SEPTEMBER", "OKTOBER", "NOVEMBER", "DESEMBER"}
	bulanText := ""
	if req.Bulan >= 1 && req.Bulan <= 12 {
		bulanText = bulanNames[req.Bulan]
	}

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Title
	title := fmt.Sprintf("DOKUMENTASI EKSTRAKURIKULER %s", strings.ToUpper(ekstrakurikuler.Name))
	pdf.CellFormat(0, 10, title, "", 1, "C", false, 0, "")

	// Subtitle
	pdf.SetFont("Arial", "B", 14)
	subtitle := fmt.Sprintf("BULAN %s TAHUN %d", bulanText, req.Tahun)
	pdf.CellFormat(0, 6, subtitle, "", 1, "C", false, 0, "")
	
	// Space after subtitle (12pt = 4mm)
	pdf.Ln(4)

	// Loop through each kegiatan
	for idx, kegiatan := range kegiatanList {
		// Add new page for each activity (except the first one)
		if idx > 0 {
			pdf.AddPage()
		}

		// Format date
		tanggal := kegiatan.TanggalKegiatan.Format("02 January 2006")
		monthMap := map[string]string{
			"January": "Januari", "February": "Februari", "March": "Maret",
			"April": "April", "May": "Mei", "June": "Juni",
			"July": "Juli", "August": "Agustus", "September": "September",
			"October": "Oktober", "November": "November", "December": "Desember",
		}
		for en, id := range monthMap {
			tanggal = strings.Replace(tanggal, en, id, 1)
		}

		// Format time
		waktu := "-"
		if kegiatan.WaktuMulai != nil && kegiatan.WaktuSelesai != nil {
			waktuMulaiStr := strings.Replace((*kegiatan.WaktuMulai)[:5], ":", ".", 1)
			waktuSelesaiStr := strings.Replace((*kegiatan.WaktuSelesai)[:5], ":", ".", 1)
			waktu = fmt.Sprintf("%s - %s WIB", waktuMulaiStr, waktuSelesaiStr)
		} else if kegiatan.WaktuMulai != nil {
			waktuMulaiStr := strings.Replace((*kegiatan.WaktuMulai)[:5], ":", ".", 1)
			waktu = fmt.Sprintf("%s WIB", waktuMulaiStr)
		}

		// Get page width
		pageWidth, _ := pdf.GetPageSize()
		leftMargin, _, rightMargin, _ := pdf.GetMargins()
		tableWidth := pageWidth - leftMargin - rightMargin

		// Calculate column widths (30% and 70%)
		col1Width := tableWidth * 0.30
		col2Width := tableWidth * 0.70

		// Table header row 1: Tanggal
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(col1Width, 7, "Tanggal", "1", 0, "L", false, 0, "")
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(col2Width, 7, " "+tanggal, "1", 1, "L", false, 0, "")

		// Row 2: Waktu
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(col1Width, 7, "Waktu", "1", 0, "L", false, 0, "")
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(col2Width, 7, " "+waktu, "1", 1, "L", false, 0, "")

		// Row 3: Materi (multi-line if needed)
		pdf.SetFont("Arial", "B", 10)
		x := pdf.GetX()
		y := pdf.GetY()
		pdf.CellFormat(col1Width, 7, "Materi", "1", 0, "L", false, 0, "")
		
		pdf.SetFont("Arial", "", 10)
		// Calculate height needed for materi with safety check
		materiText := " " + kegiatan.MateriKegiatan
		if len(materiText) > 500 {
			materiText = materiText[:500] + "..." // Truncate if too long
		}
		
		materiLines := pdf.SplitLines([]byte(materiText), col2Width-2)
		lineHeight := 7.0
		cellHeight := float64(len(materiLines)) * lineHeight
		if cellHeight < 7 {
			cellHeight = 7
		}
		// Limit max cell height to prevent overflow
		if cellHeight > 100 {
			cellHeight = 100
		}
		
		// Draw cell border for label
		pdf.SetXY(x, y)
		pdf.CellFormat(col1Width, cellHeight, "", "1", 0, "", false, 0, "")
		
		// Draw multi-line text for materi
		pdf.SetXY(x+col1Width, y)
		pdf.MultiCell(col2Width, lineHeight, materiText, "1", "L", false)

		// Row 4: Dokumentasi (photos)
		pdf.SetFont("Arial", "B", 10)
		x = pdf.GetX()
		y = pdf.GetY()
		pdf.CellFormat(tableWidth, 7, "Dokumentasi", "LTR", 1, "L", false, 0, "")

		// Check for photos
		hasPhotos := false
		photoCount := 0
		isFirstPhoto := true
		
		if kegiatan.FotoKegiatan != nil && *kegiatan.FotoKegiatan != "" {
			var fotoUrls []string
			if err := json.Unmarshal([]byte(*kegiatan.FotoKegiatan), &fotoUrls); err == nil && len(fotoUrls) > 0 {
				for _, url := range fotoUrls {
					// Try to download image
					imgBytes, err := downloadImageForPDF(url)
					if err != nil {
						// Skip if image download fails
						fmt.Printf("Warning: Failed to download image %s: %v\n", url, err)
						continue
					}
					
					if len(imgBytes) == 0 {
						continue
					}
					
					// Try to register image with gofpdf - wrap in error recovery
					func() {
						defer func() {
							if r := recover(); r != nil {
								fmt.Printf("Warning: Failed to process image %s: %v\n", url, r)
							}
						}()
						
						// Determine image type from bytes signature
						imgType := ""
						if len(imgBytes) >= 4 {
							// Check for PNG signature (89 50 4E 47)
							if imgBytes[0] == 0x89 && imgBytes[1] == 0x50 && imgBytes[2] == 0x4E && imgBytes[3] == 0x47 {
								imgType = "png"
							// Check for JPEG signature (FF D8 FF)
							} else if imgBytes[0] == 0xFF && imgBytes[1] == 0xD8 && imgBytes[2] == 0xFF {
								imgType = "jpg"
							}
						}
						
						// Skip if unknown image type
						if imgType == "" {
							fmt.Printf("Warning: Unknown image type for %s\n", url)
							return
						}
						
						imgOpt := gofpdf.ImageOptions{ImageType: imgType, ReadDpi: false}
						
						// Try to register the image using bytes.NewReader to avoid corrupting binary data
						reader := bytes.NewReader(imgBytes)
						imgInfo := pdf.RegisterImageOptionsReader(url, imgOpt, reader)
						
						if imgInfo != nil {
							// Fixed dimensions for each photo: 120mm width, maintain aspect ratio
							fixedWidth := 120.0
							imgWidth := imgInfo.Width() / 2.83 // Convert pixels to mm (72 DPI)
							imgHeight := imgInfo.Height() / 2.83
							
							// Calculate height based on fixed width
							ratio := fixedWidth / imgWidth
							fixedHeight := imgHeight * ratio
							
							// Limit max height to 90mm
							if fixedHeight > 90.0 {
								fixedHeight = 90.0
								ratio = fixedHeight / imgHeight
								fixedWidth = imgWidth * ratio
							}
							
							// Check if we need new page (image + 4mm margin)
							_, pageHeight := pdf.GetPageSize()
							bottomMargin := 15.0 // Bottom margin
							needsNewPage := pdf.GetY()+fixedHeight+4 > pageHeight-bottomMargin
							
							if needsNewPage && !isFirstPhoto {
								// Close current dokumentasi section
								pdf.CellFormat(tableWidth, 0, "", "LBR", 1, "", false, 0, "")
								
								// Add new page
								pdf.AddPage()
								
								// Redraw dokumentasi header on new page
								pdf.SetFont("Arial", "B", 10)
								pdf.CellFormat(tableWidth, 7, "Dokumentasi (lanjutan)", "LTR", 1, "L", false, 0, "")
							}
							
							// Save current position
							currentY := pdf.GetY()
							
							// Calculate centered position for image
							xPos := leftMargin + (tableWidth-fixedWidth)/2
							
							// Draw left border only
							pdf.Line(leftMargin, currentY, leftMargin, currentY+fixedHeight+4)
							
							// Draw right border only
							pdf.Line(leftMargin+tableWidth, currentY, leftMargin+tableWidth, currentY+fixedHeight+4)
							
							// Place the image centered
							pdf.ImageOptions(url, xPos, currentY+2, fixedWidth, fixedHeight, false, imgOpt, 0, "")
							
							// Move to next line after image
							pdf.SetY(currentY + fixedHeight + 4)
							
							photoCount++
							hasPhotos = true
							isFirstPhoto = false
						}
					}()
				}
			}
		}
		
		// Show message if no valid photos were processed
		if !hasPhotos || photoCount == 0 {
			pdf.SetFont("Arial", "I", 10)
			pdf.CellFormat(tableWidth, 7, "(belum ada dokumentasi)", "LBR", 1, "C", false, 0, "")
		} else {
			// Close bottom border
			pdf.CellFormat(tableWidth, 0, "", "LBR", 1, "", false, 0, "")
		}
	}

	// Generate PDF bytes
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, "", fmt.Errorf("failed to generate PDF: %v", err)
	}

	// Generate filename
	filename := fmt.Sprintf("Dokumentasi_%s_%s_%d.pdf",
		strings.ReplaceAll(ekstrakurikuler.Name, " ", "_"),
		bulanText,
		req.Tahun)

	return buf.Bytes(), filename, nil
}

// downloadImageForPDF downloads image from URL for PDF document
func downloadImageForPDF(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download image: status %d", resp.StatusCode)
	}
	
	return io.ReadAll(resp.Body)
}
