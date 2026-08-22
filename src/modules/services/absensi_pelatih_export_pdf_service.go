package services

import (
	"bytes"
	"fmt"
	"pintu-backend/src/dtos"
	"pintu-backend/src/modules/models"
	"sort"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// DownloadPDFAbsensiPelatih generates PDF file for pelatih attendance
func (s *AbsensiEkskulService) DownloadPDFAbsensiPelatih(req *dtos.AbsensiPelatihExportRequest) ([]byte, string, error) {
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

	// Create PDF - Landscape for wider table and TTD space
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.SetAutoPageBreak(true, 10)
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 14)
	title := fmt.Sprintf("DAFTAR HADIR PELATIH EKSTRAKURIKULER %s", strings.ToUpper(ekskulName))
	pdf.CellFormat(277, 7, title, "", 1, "C", false, 0, "")

	// Subtitle
	if tahunText != "" {
		pdf.SetFont("Arial", "B", 12)
		subtitle := fmt.Sprintf("BULAN %s TAHUN %s", bulanText, tahunText)
		pdf.CellFormat(277, 6, subtitle, "", 1, "C", false, 0, "")
	}

	pdf.Ln(3)

	// Column widths - adjusted for landscape and full width (total ~277mm)
	noWidth := 10.0
	namaWidth := 50.0
	tanggalWidth := 30.0
	waktuWidth := 25.0
	materiWidth := 80.0
	statusWidth := 30.0
	ttdWidth := 52.0 // Larger for signature space

	// Header
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(211, 211, 211)

	pdf.CellFormat(noWidth, 6, "No", "1", 0, "C", true, 0, "")
	pdf.CellFormat(namaWidth, 6, "Nama Pelatih", "1", 0, "C", true, 0, "")
	pdf.CellFormat(tanggalWidth, 6, "Tanggal", "1", 0, "C", true, 0, "")
	pdf.CellFormat(waktuWidth, 6, "Waktu", "1", 0, "C", true, 0, "")
	pdf.CellFormat(materiWidth, 6, "Materi Kegiatan", "1", 0, "C", true, 0, "")
	pdf.CellFormat(statusWidth, 6, "Status Kehadiran", "1", 0, "C", true, 0, "")
	pdf.CellFormat(ttdWidth, 6, "TTD", "1", 1, "C", true, 0, "")

	// Data rows - increased height for TTD space
	pdf.SetFont("Arial", "", 8)
	pdf.SetFillColor(255, 255, 255)

	rowNum := 1
	rowHeight := 15.0 // Increased from 6 to 15 for TTD space
	for _, pelatih := range pelatihList {
		for _, kegiatan := range kegiatanList {
			tanggal := formatTanggalIndonesia(kegiatan.TanggalKegiatan)
			waktuMulai := ""
			if kegiatan.WaktuMulai != nil && len(*kegiatan.WaktuMulai) >= 5 {
				waktuMulai = (*kegiatan.WaktuMulai)[:5] // HH:MM only
			}
			waktuSelesai := ""
			if kegiatan.WaktuSelesai != nil && len(*kegiatan.WaktuSelesai) >= 5 {
				waktuSelesai = (*kegiatan.WaktuSelesai)[:5]
			}
			waktu := waktuMulai
			if waktuSelesai != "" {
				waktu += "-" + waktuSelesai
			}
			materi := kegiatan.MateriKegiatan
			if len(materi) > 60 {
				materi = materi[:57] + "..."
			}

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

			pdf.CellFormat(noWidth, rowHeight, fmt.Sprintf("%d", rowNum), "1", 0, "C", false, 0, "")
			pdf.CellFormat(namaWidth, rowHeight, pelatih.Nama, "1", 0, "L", false, 0, "")
			pdf.CellFormat(tanggalWidth, rowHeight, tanggal, "1", 0, "C", false, 0, "")
			pdf.CellFormat(waktuWidth, rowHeight, waktu, "1", 0, "C", false, 0, "")
			pdf.CellFormat(materiWidth, rowHeight, materi, "1", 0, "L", false, 0, "")
			pdf.CellFormat(statusWidth, rowHeight, status, "1", 0, "C", false, 0, "")
			pdf.CellFormat(ttdWidth, rowHeight, "", "1", 1, "C", false, 0, "")

			rowNum++
		}
	}

	// Generate filename
	filename := fmt.Sprintf("Absensi_Pelatih_%s", ekskulName)
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
