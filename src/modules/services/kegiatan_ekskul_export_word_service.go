package services

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"strings"
	"time"
)

// DownloadWordDokumentasiEkskul generates Word document for ekstrakurikuler documentation
func (s *AbsensiEkskulService) DownloadWordDokumentasiEkskul(req *dtos.KegiatanEkskulDownloadWordRequest) ([]byte, string, error) {
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

	// Build document XML
	var docContent strings.Builder
	docContent.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	docContent.WriteString(`<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">`)
	docContent.WriteString(`<w:body>`)

	// Image counter for relationship IDs
	imageCounter := 1
	imageData := make(map[string][]byte) // Store downloaded images

	// Add title
	title := fmt.Sprintf("DOKUMENTASI EKSTRAKURIKULER %s", strings.ToUpper(ekstrakurikuler.Name))
	docContent.WriteString(createParagraphWithSpacing(title, true, true, 32, 120)) // Bold, centered, size 16pt, spacing after 120 twips (6pt)

	// Add subtitle
	subtitle := fmt.Sprintf("BULAN %s TAHUN %d", bulanText, req.Tahun)
	docContent.WriteString(createParagraphWithSpacing(subtitle, true, true, 28, 240)) // Bold, centered, size 14pt, spacing after 240 twips (12pt)

	// Loop through each kegiatan
	for idx, kegiatan := range kegiatanList {
		// Add page break for each activity (except the first one)
		if idx > 0 {
			docContent.WriteString(`<w:p><w:r><w:br w:type="page"/></w:r></w:p>`)
		}

		// Format date
		tanggal := kegiatan.TanggalKegiatan.Format("02 January 2006")
		// Translate month names to Indonesian
		monthMap := map[string]string{
			"January": "Januari", "February": "Februari", "March": "Maret",
			"April": "April", "May": "Mei", "June": "Juni",
			"July": "Juli", "August": "Agustus", "September": "September",
			"October": "Oktober", "November": "November", "December": "Desember",
		}
		for en, id := range monthMap {
			tanggal = strings.Replace(tanggal, en, id, 1)
		}

		// Add time
		waktu := "-"
		if kegiatan.WaktuMulai != nil && kegiatan.WaktuSelesai != nil {
			// Parse time and format to HH.MM
			waktuMulaiStr := strings.Replace((*kegiatan.WaktuMulai)[:5], ":", ".", 1)
			waktuSelesaiStr := strings.Replace((*kegiatan.WaktuSelesai)[:5], ":", ".", 1)
			waktu = fmt.Sprintf("%s - %s WIB", waktuMulaiStr, waktuSelesaiStr)
		} else if kegiatan.WaktuMulai != nil {
			waktuMulaiStr := strings.Replace((*kegiatan.WaktuMulai)[:5], ":", ".", 1)
			waktu = fmt.Sprintf("%s WIB", waktuMulaiStr)
		}

		// Create table with border for this kegiatan
		docContent.WriteString(`<w:tbl>`)
		docContent.WriteString(`<w:tblPr>`)
		docContent.WriteString(`<w:tblW w:w="5000" w:type="pct"/>`)
		docContent.WriteString(`<w:tblBorders>`)
		docContent.WriteString(`<w:top w:val="single" w:sz="12" w:space="0" w:color="000000"/>`)
		docContent.WriteString(`<w:left w:val="single" w:sz="12" w:space="0" w:color="000000"/>`)
		docContent.WriteString(`<w:bottom w:val="single" w:sz="12" w:space="0" w:color="000000"/>`)
		docContent.WriteString(`<w:right w:val="single" w:sz="12" w:space="0" w:color="000000"/>`)
		docContent.WriteString(`<w:insideH w:val="single" w:sz="6" w:space="0" w:color="000000"/>`)
		docContent.WriteString(`<w:insideV w:val="single" w:sz="6" w:space="0" w:color="000000"/>`)
		docContent.WriteString(`</w:tblBorders>`)
		docContent.WriteString(`</w:tblPr>`)

		// Row 1: Tanggal
		docContent.WriteString(createTableRow("Tanggal", tanggal))

		// Row 2: Waktu
		docContent.WriteString(createTableRow("Waktu", waktu))

		// Row 3: Materi
		docContent.WriteString(createTableRow("Materi", kegiatan.MateriKegiatan))

		// Row 4: Dokumentasi (photos) - full width row
		docContent.WriteString(`<w:tr>`)
		docContent.WriteString(`<w:tc>`)
		docContent.WriteString(`<w:tcPr><w:gridSpan w:val="2"/><w:tcW w:w="5000" w:type="pct"/><w:tcMar><w:left w:w="100" w:type="dxa"/></w:tcMar></w:tcPr>`)
		
		// Add "Dokumentasi" header
		docContent.WriteString(createSimpleParagraph("Dokumentasi", true))
		
		hasPhotos := false
		if kegiatan.FotoKegiatan != nil && *kegiatan.FotoKegiatan != "" {
			var fotoUrls []string
			if err := json.Unmarshal([]byte(*kegiatan.FotoKegiatan), &fotoUrls); err == nil && len(fotoUrls) > 0 {
				hasPhotos = true
				
				// Download and add images
				for _, url := range fotoUrls {
					// Try to download image
					imgBytes, err := downloadImageForWord(url)
					if err == nil && len(imgBytes) > 0 {
						imgKey := fmt.Sprintf("image%d.jpg", imageCounter)
						imageData[imgKey] = imgBytes
						
						// Add image paragraph (centered)
						docContent.WriteString(createImageParagraphCentered(imageCounter, imgKey))
						imageCounter++
					}
				}
			}
		}
		
		// If no photos, show message
		if !hasPhotos {
			docContent.WriteString(createSimpleParagraphCentered("(belum ada dokumentasi)", false))
		}
		
		docContent.WriteString(`</w:tc>`)
		docContent.WriteString(`</w:tr>`)

		docContent.WriteString(`</w:tbl>`)
	}

	docContent.WriteString(`</w:body>`)
	docContent.WriteString(`</w:document>`)

	// Create Word document (docx is a zip file)
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Add document.xml
	docXML, err := zipWriter.Create("word/document.xml")
	if err != nil {
		return nil, "", err
	}
	docXML.Write([]byte(docContent.String()))

	// Add images to word/media/ folder
	for imgName, imgBytes := range imageData {
		imgFile, err := zipWriter.Create("word/media/" + imgName)
		if err != nil {
			continue
		}
		imgFile.Write(imgBytes)
	}

	// Add [Content_Types].xml
	contentTypes := createContentTypesXML(imageCounter)
	contentTypesFile, err := zipWriter.Create("[Content_Types].xml")
	if err != nil {
		return nil, "", err
	}
	contentTypesFile.Write([]byte(contentTypes))

	// Add _rels/.rels
	rels := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`

	relsFile, err := zipWriter.Create("_rels/.rels")
	if err != nil {
		return nil, "", err
	}
	relsFile.Write([]byte(rels))

	// Add word/_rels/document.xml.rels for images
	docRels := createDocumentRelsXML(imageCounter)
	docRelsFile, err := zipWriter.Create("word/_rels/document.xml.rels")
	if err != nil {
		return nil, "", err
	}
	docRelsFile.Write([]byte(docRels))

	zipWriter.Close()

	// Generate filename
	filename := fmt.Sprintf("Dokumentasi_%s_%s_%d.docx",
		strings.ReplaceAll(ekstrakurikuler.Name, " ", "_"),
		bulanText,
		req.Tahun)

	return buf.Bytes(), filename, nil
}

// createParagraph creates a Word paragraph XML element
func createParagraph(text string, bold, centered bool, fontSize int) string {
	var p strings.Builder
	
	p.WriteString(`<w:p>`)
	
	// Paragraph properties
	if centered || fontSize > 0 {
		p.WriteString(`<w:pPr>`)
		if centered {
			p.WriteString(`<w:jc w:val="center"/>`)
		}
		p.WriteString(`</w:pPr>`)
	}
	
	p.WriteString(`<w:r>`)
	
	// Run properties
	if bold || fontSize > 0 {
		p.WriteString(`<w:rPr>`)
		if bold {
			p.WriteString(`<w:b/>`)
		}
		if fontSize > 0 {
			p.WriteString(fmt.Sprintf(`<w:sz w:val="%d"/>`, fontSize))
		}
		p.WriteString(`</w:rPr>`)
	}
	
	// Escape XML special characters
	escapedText := html.EscapeString(text)
	p.WriteString(fmt.Sprintf(`<w:t xml:space="preserve">%s</w:t>`, escapedText))
	
	p.WriteString(`</w:r>`)
	p.WriteString(`</w:p>`)
	
	return p.String()
}

// createParagraphWithSpacing creates a Word paragraph XML element with spacing after
func createParagraphWithSpacing(text string, bold, centered bool, fontSize int, spacingAfter int) string {
	var p strings.Builder
	
	p.WriteString(`<w:p>`)
	
	// Paragraph properties
	p.WriteString(`<w:pPr>`)
	if centered {
		p.WriteString(`<w:jc w:val="center"/>`)
	}
	if spacingAfter > 0 {
		p.WriteString(fmt.Sprintf(`<w:spacing w:after="%d"/>`, spacingAfter))
	}
	p.WriteString(`</w:pPr>`)
	
	if text != "" {
		p.WriteString(`<w:r>`)
		
		// Run properties
		if bold || fontSize > 0 {
			p.WriteString(`<w:rPr>`)
			if bold {
				p.WriteString(`<w:b/>`)
			}
			if fontSize > 0 {
				p.WriteString(fmt.Sprintf(`<w:sz w:val="%d"/>`, fontSize))
			}
			p.WriteString(`</w:rPr>`)
		}
		
		// Escape XML special characters
		escapedText := html.EscapeString(text)
		p.WriteString(fmt.Sprintf(`<w:t xml:space="preserve">%s</w:t>`, escapedText))
		
		p.WriteString(`</w:r>`)
	}
	
	p.WriteString(`</w:p>`)
	
	return p.String()
}

// createSimpleParagraph creates a simple paragraph without formatting options
func createSimpleParagraph(text string, bold bool) string {
	var p strings.Builder
	p.WriteString(`<w:p><w:r>`)
	if bold {
		p.WriteString(`<w:rPr><w:b/></w:rPr>`)
	}
	escapedText := html.EscapeString(text)
	p.WriteString(fmt.Sprintf(`<w:t xml:space="preserve">%s</w:t>`, escapedText))
	p.WriteString(`</w:r></w:p>`)
	return p.String()
}

// createSimpleParagraphCentered creates a simple centered paragraph
func createSimpleParagraphCentered(text string, bold bool) string {
	var p strings.Builder
	p.WriteString(`<w:p>`)
	p.WriteString(`<w:pPr><w:jc w:val="center"/></w:pPr>`)
	p.WriteString(`<w:r>`)
	if bold {
		p.WriteString(`<w:rPr><w:b/></w:rPr>`)
	}
	escapedText := html.EscapeString(text)
	p.WriteString(fmt.Sprintf(`<w:t xml:space="preserve">%s</w:t>`, escapedText))
	p.WriteString(`</w:r></w:p>`)
	return p.String()
}

// createTableRow creates a table row with label and value (with left padding in both cells)
func createTableRow(label, value string) string {
	var row strings.Builder
	row.WriteString(`<w:tr>`)
	
	// Cell 1: Label (30% width) with left padding
	row.WriteString(`<w:tc>`)
	row.WriteString(`<w:tcPr><w:tcW w:w="1500" w:type="pct"/><w:tcMar><w:left w:w="100" w:type="dxa"/></w:tcMar></w:tcPr>`)
	row.WriteString(createSimpleParagraph(label, true))
	row.WriteString(`</w:tc>`)
	
	// Cell 2: Value (70% width) with left padding
	row.WriteString(`<w:tc>`)
	row.WriteString(`<w:tcPr><w:tcW w:w="3500" w:type="pct"/><w:tcMar><w:left w:w="100" w:type="dxa"/></w:tcMar></w:tcPr>`)
	row.WriteString(createSimpleParagraph(value, false))
	row.WriteString(`</w:tc>`)
	
	row.WriteString(`</w:tr>`)
	return row.String()
}

// downloadImageForWord downloads image from URL for Word document
func downloadImageForWord(url string) ([]byte, error) {
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

// createImageParagraph creates a paragraph with embedded image
func createImageParagraph(imageID int, imageName string) string {
	rIdStr := fmt.Sprintf("rId%d", imageID)
	
	// Image dimensions (width: 4 inches = 3657600 EMUs, height: 3 inches = 2743200 EMUs)
	width := 3657600
	height := 2743200
	
	var p strings.Builder
	p.WriteString(`<w:p>`)
	p.WriteString(`<w:r>`)
	p.WriteString(`<w:drawing>`)
	p.WriteString(`<wp:inline distT="0" distB="0" distL="0" distR="0">`)
	p.WriteString(fmt.Sprintf(`<wp:extent cx="%d" cy="%d"/>`, width, height))
	p.WriteString(`<wp:effectExtent l="0" t="0" r="0" b="0"/>`)
	p.WriteString(`<wp:docPr id="1" name="Picture"/>`)
	p.WriteString(`<a:graphic>`)
	p.WriteString(`<a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/picture">`)
	p.WriteString(`<pic:pic xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture">`)
	p.WriteString(`<pic:nvPicPr>`)
	p.WriteString(fmt.Sprintf(`<pic:cNvPr id="%d" name="%s"/>`, imageID, imageName))
	p.WriteString(`<pic:cNvPicPr/>`)
	p.WriteString(`</pic:nvPicPr>`)
	p.WriteString(`<pic:blipFill>`)
	p.WriteString(fmt.Sprintf(`<a:blip r:embed="%s"/>`, rIdStr))
	p.WriteString(`<a:stretch><a:fillRect/></a:stretch>`)
	p.WriteString(`</pic:blipFill>`)
	p.WriteString(`<pic:spPr>`)
	p.WriteString(`<a:xfrm>`)
	p.WriteString(`<a:off x="0" y="0"/>`)
	p.WriteString(fmt.Sprintf(`<a:ext cx="%d" cy="%d"/>`, width, height))
	p.WriteString(`</a:xfrm>`)
	p.WriteString(`<a:prstGeom prst="rect"><a:avLst/></a:prstGeom>`)
	p.WriteString(`</pic:spPr>`)
	p.WriteString(`</pic:pic>`)
	p.WriteString(`</a:graphicData>`)
	p.WriteString(`</a:graphic>`)
	p.WriteString(`</wp:inline>`)
	p.WriteString(`</w:drawing>`)
	p.WriteString(`</w:r>`)
	p.WriteString(`</w:p>`)
	
	return p.String()
}

// createImageParagraphCentered creates a centered paragraph with embedded image
func createImageParagraphCentered(imageID int, imageName string) string {
	rIdStr := fmt.Sprintf("rId%d", imageID)
	
	// Image dimensions (width: 4 inches = 3657600 EMUs, height: 3 inches = 2743200 EMUs)
	width := 3657600
	height := 2743200
	
	var p strings.Builder
	p.WriteString(`<w:p>`)
	p.WriteString(`<w:pPr><w:jc w:val="center"/></w:pPr>`)
	p.WriteString(`<w:r>`)
	p.WriteString(`<w:drawing>`)
	p.WriteString(`<wp:inline distT="0" distB="0" distL="0" distR="0">`)
	p.WriteString(fmt.Sprintf(`<wp:extent cx="%d" cy="%d"/>`, width, height))
	p.WriteString(`<wp:effectExtent l="0" t="0" r="0" b="0"/>`)
	p.WriteString(`<wp:docPr id="1" name="Picture"/>`)
	p.WriteString(`<a:graphic>`)
	p.WriteString(`<a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/picture">`)
	p.WriteString(`<pic:pic xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture">`)
	p.WriteString(`<pic:nvPicPr>`)
	p.WriteString(fmt.Sprintf(`<pic:cNvPr id="%d" name="%s"/>`, imageID, imageName))
	p.WriteString(`<pic:cNvPicPr/>`)
	p.WriteString(`</pic:nvPicPr>`)
	p.WriteString(`<pic:blipFill>`)
	p.WriteString(fmt.Sprintf(`<a:blip r:embed="%s"/>`, rIdStr))
	p.WriteString(`<a:stretch><a:fillRect/></a:stretch>`)
	p.WriteString(`</pic:blipFill>`)
	p.WriteString(`<pic:spPr>`)
	p.WriteString(`<a:xfrm>`)
	p.WriteString(`<a:off x="0" y="0"/>`)
	p.WriteString(fmt.Sprintf(`<a:ext cx="%d" cy="%d"/>`, width, height))
	p.WriteString(`</a:xfrm>`)
	p.WriteString(`<a:prstGeom prst="rect"><a:avLst/></a:prstGeom>`)
	p.WriteString(`</pic:spPr>`)
	p.WriteString(`</pic:pic>`)
	p.WriteString(`</a:graphicData>`)
	p.WriteString(`</a:graphic>`)
	p.WriteString(`</wp:inline>`)
	p.WriteString(`</w:drawing>`)
	p.WriteString(`</w:r>`)
	p.WriteString(`</w:p>`)
	
	return p.String()
}

// createContentTypesXML creates Content_Types.xml with image support
func createContentTypesXML(imageCount int) string {
	var ct strings.Builder
	ct.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	ct.WriteString(`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">`)
	ct.WriteString(`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>`)
	ct.WriteString(`<Default Extension="xml" ContentType="application/xml"/>`)
	ct.WriteString(`<Default Extension="jpg" ContentType="image/jpeg"/>`)
	ct.WriteString(`<Default Extension="jpeg" ContentType="image/jpeg"/>`)
	ct.WriteString(`<Default Extension="png" ContentType="image/png"/>`)
	ct.WriteString(`<Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>`)
	ct.WriteString(`</Types>`)
	return ct.String()
}

// createDocumentRelsXML creates document.xml.rels for image relationships
func createDocumentRelsXML(imageCount int) string {
	var rels strings.Builder
	rels.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	rels.WriteString(`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`)
	
	// Add relationship for each image
	for i := 1; i < imageCount; i++ {
		rels.WriteString(fmt.Sprintf(`<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="media/image%d.jpg"/>`, i, i))
	}
	
	rels.WriteString(`</Relationships>`)
	return rels.String()
}
