package services

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"pintu-backend/src/modules/models"
	"github.com/jung-kurt/gofpdf"
	qrcode "github.com/skip2/go-qrcode"
	"golang.org/x/net/html"
)

// ============================================================================
// PALET WARNA & KONSTANTA UKURAN
// Ubah nilai-nilai di sini untuk menyesuaikan tema/warna kartu pelajar.
// ============================================================================

type cardColors struct {
	PrimaryR, PrimaryG, PrimaryB       int // Warna utama (header, border kartu)
	AccentR, AccentG, AccentB          int // Warna aksen (garis, badge, bullet)
	TextDarkR, TextDarkG, TextDarkB    int // Warna teks utama (nama, nilai data)
	TextMutedR, TextMutedG, TextMutedB int // Warna teks label/keterangan
	LightBgR, LightBgG, LightBgB       int // Warna latar terang (kotak foto, dll)
}

var colors = cardColors{
	PrimaryR: 168, PrimaryG: 28, PrimaryB: 32, // Merah tegas (warna utama)
	AccentR: 196, AccentG: 148, AccentB: 66, // Coklat keemasan (warna aksen)
	TextDarkR: 25, TextDarkG: 25, TextDarkB: 25,
	TextMutedR: 115, TextMutedG: 115, TextMutedB: 115,
	LightBgR: 246, LightBgG: 244, LightBgB: 241,
}

const (
	cardWidth  = 90.0 // mm, diperbesar dari 85.0 ke 90.0
	cardHeight = 57.0 // mm, diperbesar dari 54.0 ke 57.0
	headerH    = 12.0 // dirapatkan lagi supaya jarak alamat ke garis aksen tidak kejauhan
	accentBarH = 1.3
)

var bulanIndonesia = []string{
	"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
	"Juli", "Agustus", "September", "Oktober", "November", "Desember",
}

// formatTanggalIndonesia memformat tanggal dengan nama bulan berbahasa Indonesia,
// misal "05 Juli 2026", karena time.Format bawaan Go hanya punya nama bulan Inggris.
func formatTanggalIndonesia(t time.Time) string {
	return fmt.Sprintf("%02d %s %d", t.Day(), bulanIndonesia[int(t.Month())], t.Year())
}

// toTitleCase mengubah string menjadi Title Case (setiap kata diawali huruf kapital)
// Contoh: "PAGAR DEWA" -> "Pagar Dewa"
func toTitleCase(s string) string {
	words := strings.Fields(strings.ToLower(s))
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + word[1:]
		}
	}
	return strings.Join(words, " ")
}

// KartuPelajarData represents data needed for student card
type KartuPelajarData struct {
	Siswa         []models.PesertaDidik
	KepalaSekolah *models.Kepegawaian
	VisiMisi      *models.VisiMisi
}

// GenerateKartuPelajarPDF generates PDF for student cards
func GenerateKartuPelajarPDF(data *KartuPelajarData) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")

	logoDKI, err := downloadImage("https://pintu-storage.sdnsukapura01.sch.id/logo-dki.png")
	if err != nil {
		logoDKI = nil
	}
	logoSekolah, err := downloadImage("https://pintu-storage.sdnsukapura01.sch.id/logo-without-bg.png")
	if err != nil {
		logoSekolah = nil
	}
	ttdKepsek, err := downloadImage("https://pintu-storage.sdnsukapura01.sch.id/ttd-kepsek-tanpa-bg.png")
	if err != nil {
		ttdKepsek = nil
	}
	elemenDekor, err := downloadImage("https://pintu-storage.sdnsukapura01.sch.id/elemen-katpel.png")
	if err != nil {
		elemenDekor = nil
	}

	// 4 kartu per halaman (depan kiri, belakang kanan)
	for i := 0; i < len(data.Siswa); i += 4 {
		pdf.AddPage()

		end := i + 4
		if end > len(data.Siswa) {
			end = len(data.Siswa)
		}
		pageStudents := data.Siswa[i:end]

		for idx, siswa := range pageStudents {
			yPos := 10.0 + float64(idx)*68.0
			s := siswa
			drawFrontCard(pdf, &s, data.KepalaSekolah, yPos, 10.0, logoDKI, logoSekolah, ttdKepsek, elemenDekor)
			drawBackCard(pdf, data.VisiMisi, yPos, 110.0, logoSekolah, elemenDekor)
			drawCutMarks(pdf, 10.0, yPos, cardWidth, cardHeight)
			drawCutMarks(pdf, 110.0, yPos, cardWidth, cardHeight)
		}
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("gagal generate PDF: %s", err.Error())
	}
	return buf.Bytes(), nil
}

// ============================================================================
// KARTU DEPAN
// ============================================================================

func drawFrontCard(pdf *gofpdf.Fpdf, siswa *models.PesertaDidik, kepsek *models.Kepegawaian, y, x float64, logoDKI, logoSekolah, ttdKepsek, elemenDekor []byte) {
	// Base kartu dengan background cream muda
	pdf.SetFillColor(254, 253, 250) // Warna cream lebih muda dan lembut
	pdf.SetDrawColor(colors.PrimaryR, colors.PrimaryG, colors.PrimaryB)
	pdf.SetLineWidth(0.4)
	pdf.Rect(x, y, cardWidth, cardHeight, "FD")

	// Header flat modern + aksen geometris diagonal di pojok
	drawHeaderBlock(pdf, x, y, cardWidth, headerH)

	// Logo — lebar tetap tapi tinggi menyesuaikan rasio asli gambar (h: 0 = auto),
	// supaya logo tidak gepeng/terdistorsi seperti saat dipaksa persegi.
	logoDKIW := 7.5
	logoSekolahW := 9.0 // Dikecilkan dari 10.0 ke 9.0
	if logoDKI != nil {
		pdf.RegisterImageOptionsReader("logoDKI", gofpdf.ImageOptions{ImageType: "PNG"}, bytes.NewReader(logoDKI))
		pdf.Image("logoDKI", x+3, y+2, logoDKIW, 0, false, "", 0, "")
	}
	if logoSekolah != nil {
		pdf.RegisterImageOptionsReader("logoSekolah", gofpdf.ImageOptions{ImageType: "PNG"}, bytes.NewReader(logoSekolah))
		pdf.Image("logoSekolah", x+cardWidth-3-logoSekolahW, y+0.8, logoSekolahW, 0, false, "", 0, "")
	}

	// Teks header
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Helvetica", "B", 6.0) // Diubah ke Bold
	pdf.SetXY(x+14, y+1.2)
	pdf.CellFormat(cardWidth-28, 2.0, "PEMERINTAH PROVINSI DKI JAKARTA", "", 0, "C", false, 0, "")
	pdf.SetXY(x+14, y+3.2)
	pdf.CellFormat(cardWidth-28, 2.0, "DINAS PENDIDIKAN", "", 0, "C", false, 0, "")
	pdf.SetFont("Helvetica", "B", 8.5)
	pdf.SetXY(x+14, y+5.5)
	pdf.CellFormat(cardWidth-28, 3.2, "SDN SUKAPURA 01", "", 0, "C", false, 0, "")
	pdf.SetFont("Helvetica", "", 5.5) // Alamat diperbesar dari 5.0 ke 5.5
	pdf.SetXY(x+14, y+8.9)
	pdf.MultiCell(cardWidth-28, 1.4, "Jl. Beo No.15, Komp.Walikota No.2, Cilincing, Jakarta Utara", "", "C", false)
	pdf.SetTextColor(colors.TextDarkR, colors.TextDarkG, colors.TextDarkB)

	// Judul "KARTU PELAJAR" — dengan garis aksen tipis di bawahnya
	pdf.SetFont("Helvetica", "B", 9.0)
	pdf.SetTextColor(colors.PrimaryR, colors.PrimaryG, colors.PrimaryB)
	titleY := y + headerH + accentBarH + 2.0 // Diturunkan dari +1.5 ke +2.0
	pdf.SetXY(x, titleY)
	pdf.CellFormat(cardWidth, 3.4, "KARTU PELAJAR", "", 0, "C", false, 0, "")
	pdf.SetFillColor(colors.AccentR, colors.AccentG, colors.AccentB)
	pdf.Rect(x+cardWidth/2-8, titleY+3.9, 16, 0.6, "F") // Spacing diperlonggar dari 3.6 ke 3.9
	pdf.SetTextColor(colors.TextDarkR, colors.TextDarkG, colors.TextDarkB)

	// ---- Body ----
	bodyY := titleY + 5.5
	footerLineY := y + cardHeight - 16 // Diturunkan sedikit dari -17 ke -16

	// Background dekoratif transparan — dibatasi mulai dari bawah header,
	// supaya warna header tetap solid murni tanpa terkena opacity elemen ini
	drawBackgroundDecor(pdf, x, y, y+headerH+accentBarH, logoSekolah, elemenDekor)

	photoX := x + 3 // Digeser ke kiri dari x+4 ke x+3
	photoW, photoH := 17.0, 21.0 // Diperbesar lagi dari 16.0x20.0 ke 17.0x21.0
	photoY := bodyY + 1.0 // Diturunkan sedikit ke bawah

	// Kotak foto dengan bingkai kuning penuh (warna aksen)
	pdf.SetFillColor(colors.LightBgR, colors.LightBgG, colors.LightBgB)
	pdf.SetDrawColor(colors.AccentR, colors.AccentG, colors.AccentB) // Garis kuning
	pdf.SetLineWidth(0.8) // Garis lebih tebal
	pdf.Rect(photoX, photoY, photoW, photoH, "FD")
	
	// Coba download dan tampilkan foto siswa dari R2 storage
	photoDisplayed := false
	if siswa.Photo != "" {
		photoData, err := downloadImage(siswa.Photo)
		if err == nil {
			photoName := fmt.Sprintf("photo_%s", siswa.NIS)
			// Detect image type from URL or try both PNG and JPEG
			var imgOpts gofpdf.ImageOptions
			if strings.Contains(strings.ToLower(siswa.Photo), ".png") {
				imgOpts = gofpdf.ImageOptions{ImageType: "PNG"}
			} else if strings.Contains(strings.ToLower(siswa.Photo), ".jpg") || strings.Contains(strings.ToLower(siswa.Photo), ".jpeg") {
				imgOpts = gofpdf.ImageOptions{ImageType: "JPEG"}
			} else {
				// Default to auto-detect
				imgOpts = gofpdf.ImageOptions{ImageType: ""}
			}
			
			pdf.RegisterImageOptionsReader(photoName, imgOpts, bytes.NewReader(photoData))
			pdf.Image(photoName, photoX, photoY, photoW, photoH, false, "", 0, "")
			photoDisplayed = true
		}
	}
	
	// Jika foto tidak ada atau gagal download, tampilkan placeholder
	if !photoDisplayed {
		pdf.SetFont("Helvetica", "", 5.5)
		pdf.SetTextColor(160, 160, 160)
		pdf.SetXY(photoX, photoY+photoH/2-1.2)
		pdf.CellFormat(photoW, 2.4, "FOTO", "", 0, "C", false, 0, "")
		pdf.SetXY(photoX, photoY+photoH/2+1.2)
		pdf.CellFormat(photoW, 2.4, "3x4", "", 0, "C", false, 0, "")
	}
	pdf.SetTextColor(colors.TextDarkR, colors.TextDarkG, colors.TextDarkB)

	// Kolom info data siswa dengan QR code di kanan
	infoX := photoX + photoW + 0.8 // Dikurangi dari +1.0 ke +0.8 agar lebih dekat foto
	qrSize := 11.0
	infoWidth := (x + cardWidth - 3) - infoX - qrSize - 2 // Ada pengurangan untuk QR code
	labelW := 22.0
	rowY := bodyY + 1.0 // Diturunkan sejajar dengan foto

	// Nama - bisa multi-line, terbatas karena ada QR code
	namaHeight := drawInfoRowMultiLine(pdf, infoX, rowY, "Nama", siswa.Nama, labelW, infoWidth)
	
	namaRowY := rowY // Simpan posisi Y untuk Nama (untuk QR code sejajar)
	
	// Field lainnya - spacing sama semua
	rowY += namaHeight + 0.4
	
	// NIS / NISN - gabungkan kedua nilai, gunakan "-" jika kosong
	nis := siswa.NIS
	if nis == "" {
		nis = "-"
	}
	nisn := siswa.NISN
	if nisn == "" {
		nisn = "-"
	}
	nisNisn := nis + " / " + nisn
	drawInfoRow(pdf, infoX, rowY, "NIS / NISN", nisNisn, labelW, infoWidth)

	rowY += 3.6 // Disamakan dengan line height (sebelumnya 3.0)
	
	// Format TTL dengan tempat lahir title case dan tanggal format Indonesia
	tempTgl := toTitleCase(siswa.TempatLahir)
	if siswa.TanggalLahir != nil {
		tempTgl += ", " + formatTanggalIndonesia(*siswa.TanggalLahir)
	}
	drawInfoRowWithAutoResize(pdf, infoX, rowY, "TTL", tempTgl, labelW, infoWidth)

	rowY += 3.6 // Disamakan (sebelumnya 3.0)
	drawInfoRow(pdf, infoX, rowY, "Agama", siswa.Agama, labelW, infoWidth)

	rowY += 3.6 // Disamakan (sebelumnya 3.0)
	jenisKelamin := "Laki-laki"
	if siswa.JenisKelamin == "P" {
		jenisKelamin = "Perempuan"
	}
	drawInfoRow(pdf, infoX, rowY, "Jenis Kelamin", jenisKelamin, labelW, infoWidth)

	// QR code, sejajar dengan baris Nama di kanan
	if siswa.Barcode != "" {
		qrData, err := qrcode.Encode(siswa.Barcode, qrcode.Medium, 256)
		if err == nil {
			qrX := x + cardWidth - qrSize - 3
			qrY := namaRowY // Sejajar dengan Nama
			// Kotak background putih tanpa padding (persis ukuran barcode) dengan garis tebal
			pdf.SetFillColor(255, 255, 255)
			pdf.SetDrawColor(220, 220, 220)
			pdf.SetLineWidth(0.8) // Garis dipertebal dari 0.5 ke 0.8
			pdf.Rect(qrX, qrY, qrSize, qrSize, "FD") // Tanpa padding, persis ukuran barcode
			pdf.SetLineWidth(0.4) // Reset ke line width normal
			qrName := fmt.Sprintf("qr_%s", siswa.NIS)
			pdf.RegisterImageOptionsReader(qrName, gofpdf.ImageOptions{ImageType: "PNG"}, bytes.NewReader(qrData))
			pdf.Image(qrName, qrX, qrY, qrSize, qrSize, false, "", 0, "")
		}
	}

	// ---- Footer / tanda tangan ----
	footerY := footerLineY
	currentDate := formatTanggalIndonesia(time.Now())
	sigW := 34.0
	sigX := x + cardWidth - sigW - 3

	pdf.SetFont("Helvetica", "", 5.0)
	pdf.SetTextColor(colors.TextDarkR, colors.TextDarkG, colors.TextDarkB)
	pdf.SetXY(sigX, footerY+1)
	pdf.CellFormat(sigW, 2.2, fmt.Sprintf("Jakarta, %s", currentDate), "", 0, "C", false, 0, "")
	pdf.SetXY(sigX, footerY+2.8)
	pdf.CellFormat(sigW, 2.2, "Kepala SDN Sukapura 01", "", 0, "C", false, 0, "")

	if kepsek != nil {
		pdf.SetFont("Helvetica", "B", 5.5)
		pdf.SetTextColor(colors.TextDarkR, colors.TextDarkG, colors.TextDarkB)
		pdf.SetXY(sigX, footerY+10.5) // Diperbesar spacing dari 9.8 ke 10.5
		pdf.CellFormat(sigW, 2.2, kepsek.Nama, "", 0, "C", false, 0, "")
		pdf.SetFont("Helvetica", "", 4.8)
		pdf.SetTextColor(colors.TextDarkR, colors.TextDarkG, colors.TextDarkB)
		pdf.SetXY(sigX, footerY+12.5) // Disesuaikan dari 11.8 ke 12.5
		pdf.CellFormat(sigW, 2.2, fmt.Sprintf("NIP. %s", kepsek.NIP), "", 0, "C", false, 0, "")
	}

	// Tanda tangan asli kepala sekolah (PNG transparan)
	if ttdKepsek != nil {
		pdf.RegisterImageOptionsReader("ttdKepsek", gofpdf.ImageOptions{ImageType: "PNG"}, bytes.NewReader(ttdKepsek))
		ttdW := 12.5 // Diperbesar dari 11.0 ke 12.5
		ttdX := sigX + (sigW-ttdW)/2
		ttdY := footerY + 3.4
		pdf.Image("ttdKepsek", ttdX, ttdY, ttdW, 0, false, "", 0, "")
	}

	pdf.SetTextColor(colors.TextDarkR, colors.TextDarkG, colors.TextDarkB)

	// Catatan kecil masa berlaku kartu - mepet bawah
	noteW := 27.0
	noteY := y + cardHeight - 5.8 // Dinaikkan dari -4.8 ke -5.8
	noteH := 3.0
	pdf.SetAlpha(0.08, "Normal")
	pdf.SetFillColor(colors.AccentR, colors.AccentG, colors.AccentB)
	pdf.Rect(x+3, noteY, noteW, noteH, "F")
	pdf.SetAlpha(1, "Normal")
	pdf.SetFillColor(colors.AccentR, colors.AccentG, colors.AccentB)
	pdf.Rect(x+3, noteY, 0.7, noteH, "F")
	pdf.SetFont("Helvetica", "BI", 4.0)
	pdf.SetTextColor(colors.PrimaryR, colors.PrimaryG, colors.PrimaryB)
	pdf.SetXY(x+4.2, noteY+0.7)
	pdf.CellFormat(noteW-2, 2.0, "* Berlaku selama menjadi siswa aktif", "", 0, "L", false, 0, "")
}

// ============================================================================
// KARTU BELAKANG (VISI & MISI)
// ============================================================================

func drawBackCard(pdf *gofpdf.Fpdf, visiMisi *models.VisiMisi, y, x float64, logoSekolah, elemenDekor []byte) {
	// Base kartu dengan background cream muda
	pdf.SetFillColor(254, 253, 250) // Warna cream lebih muda dan lembut
	pdf.SetDrawColor(colors.PrimaryR, colors.PrimaryG, colors.PrimaryB)
	pdf.SetLineWidth(0.4)
	pdf.Rect(x, y, cardWidth, cardHeight, "FD")

	backHeaderH := headerH // disamakan dengan tinggi header kartu depan
	drawHeaderBlock(pdf, x, y, cardWidth, backHeaderH)

	pdf.SetFont("Helvetica", "B", 8.5)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(x, y+5)
	pdf.CellFormat(cardWidth, 4, "VISI & MISI SEKOLAH", "", 0, "C", false, 0, "")
	pdf.SetTextColor(colors.TextDarkR, colors.TextDarkG, colors.TextDarkB)

	// Latar dekoratif lembut, dibatasi mulai dari bawah header (sama seperti kartu depan)
	drawBackgroundDecor(pdf, x, y, y+backHeaderH+accentBarH, logoSekolah, elemenDekor)

	contentX := x + 4
	maxWidth := cardWidth - 8
	contentY := y + backHeaderH + accentBarH + 2.0 // Diturunkan dari 1.2 ke 2.0

	// Hard-coded Visi
	contentY = drawSectionBadge(pdf, contentX, contentY, "VISI")
	visiText := "Menciptakan lingkungan pendidikan yang mendukung peserta didik untuk menjadi cerdas secara holistik, kreatif, dan adaptif di era teknologi, serta memiliki karakter dan nilai-nilai Pancasila yang kuat"
	pdf.SetFont("Helvetica", "I", 6.0) // Diperbesar dari 5.5 ke 6.0
	pdf.SetTextColor(colors.TextDarkR, colors.TextDarkG, colors.TextDarkB)
	pdf.SetXY(contentX+2, contentY)
	pdf.MultiCell(maxWidth-2, 2.4, visiText, "", "L", false) // Line height dari 2.3 ke 2.4
	contentY = pdf.GetY() + 1.2 // Dikurangi dari 1.5 ke 1.2

	// Hard-coded Misi
	contentY = drawSectionBadge(pdf, contentX, contentY, "MISI")
	misiItems := []string{
		"Mengembangkan berbagai kegiatan untuk menguatkan keimanan dan ketakwaan kepada Tuhan Yang Maha Esa.",
		"Menyediakan lingkungan belajar yang mendukung pengembangan kecerdasan intelektual, emosional, dan sosial peserta didik.",
		"Mengintegrasikan teknologi pendidikan dalam proses pembelajaran untuk meningkatkan kreativitas dan inovasi peserta didik.",
		"Mengembangkan kurikulum yang adaptif dan relevan dengan perkembangan zaman dan kebutuhan peserta didik.",
	}
	pdf.SetFont("Helvetica", "", 5.6) // Diperbesar dari 5.2 ke 5.6
	for _, item := range misiItems {
		pdf.SetFillColor(colors.AccentR, colors.AccentG, colors.AccentB)
		pdf.Rect(contentX+1, contentY+0.8, 1.2, 1.2, "F") // Bullet diperbesar dari 1.1 ke 1.2
		pdf.SetXY(contentX+3.2, contentY) // Disesuaikan dari 3 ke 3.2
		pdf.MultiCell(maxWidth-3.2, 2.4, item, "", "L", false) // Line height dari 2.3 ke 2.4
		contentY = pdf.GetY() + 0.4 // Dikurangi dari 0.5 ke 0.4
		if contentY > y+cardHeight-2 {
			break
		}
	}
}

func drawSectionBadge(pdf *gofpdf.Fpdf, x, y float64, label string) float64 {
	pdf.SetFont("Helvetica", "B", 5.0)
	w := pdf.GetStringWidth(label) + 4
	h := 3.2
	pdf.SetFillColor(colors.AccentR, colors.AccentG, colors.AccentB)
	pdf.Rect(x, y, w, h, "F")
	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(x, y)
	pdf.CellFormat(w, h, label, "", 0, "C", false, 0, "")
	pdf.SetTextColor(colors.TextDarkR, colors.TextDarkG, colors.TextDarkB)
	return y + h + 0.8
}

// ============================================================================
// HELPER — HEADER & LATAR DEKORATIF (FLAT MODERN, TANPA GELOMBANG)
// ============================================================================

// drawHeaderBlock menggambar header flat solid + aksen geometris diagonal
// di pojok kanan atas + bar aksen tipis sebagai pemisah ke body.
func drawHeaderBlock(pdf *gofpdf.Fpdf, x, y, w, h float64) {
	pdf.SetFillColor(colors.PrimaryR, colors.PrimaryG, colors.PrimaryB)
	pdf.Rect(x, y, w, h, "F")

	// Aksen geometris diagonal di pojok kanan atas (potongan segitiga, bukan lengkung)
	pdf.SetAlpha(0.18, "Normal")
	pdf.SetFillColor(colors.AccentR, colors.AccentG, colors.AccentB)
	pdf.Polygon([]gofpdf.PointType{
		{X: x + w - 24, Y: y},
		{X: x + w, Y: y},
		{X: x + w, Y: y + h*0.75},
	}, "F")
	pdf.SetAlpha(1, "Normal")

	// Bar aksen tipis solid sebagai pemisah header-body
	pdf.SetFillColor(colors.AccentR, colors.AccentG, colors.AccentB)
	pdf.Rect(x, y+h, w, accentBarH, "F")
}

// drawBackgroundDecor menggambar elemen dekoratif transparan HANYA di area body
// (di bawah header), supaya warna header tetap solid murni tanpa terkena efek
// opacity dari elemen dekoratif ini. Hanya logo sekolah samar di tengah area body
// dan elemen dekoratif di kanan bawah.
func drawBackgroundDecor(pdf *gofpdf.Fpdf, x, y, bodyTop float64, logoSekolah, elemenDekor []byte) {
	bodyBottom := y + cardHeight
	bodyH := bodyBottom - bodyTop
	pdf.ClipRect(x, bodyTop, cardWidth, bodyH, false)

	// Logo sekolah, samar, di tengah area body (bukan tengah kartu penuh,
	// supaya tidak menyentuh/menimpa header sama sekali).
	if logoSekolah != nil {
		pdf.SetAlpha(0.09, "Normal")
		pdf.RegisterImageOptionsReader("logoSekolahWatermark", gofpdf.ImageOptions{ImageType: "PNG"}, bytes.NewReader(logoSekolah))
		wmW, wmH := 38.0, 40.0 // Width dikecilkan dari 44 ke 38, height tetap 40
		centerY := bodyTop + bodyH/2
		pdf.Image("logoSekolahWatermark", x+cardWidth/2-wmW/2, centerY-wmH/2, wmW, wmH, false, "", 0, "")
		pdf.SetAlpha(1, "Normal")
	}

	// Elemen dekoratif di kanan bawah dengan opacity
	if elemenDekor != nil {
		pdf.SetAlpha(0.15, "Normal")
		elemenName := fmt.Sprintf("elemenDekor_%f_%f", x, y)
		pdf.RegisterImageOptionsReader(elemenName, gofpdf.ImageOptions{ImageType: "PNG"}, bytes.NewReader(elemenDekor))
		// Ukuran diperbesar lagi (42 -> 55) dan posisi dinaikkan + digeser
		// lebih ke kiri supaya tidak tenggelam/terpotong di batas bawah kartu.
		elemenW := 55.0
		elemenX := x + cardWidth - elemenW + 2 // digeser lebih ke kiri (dari +3 ke -4)
		elemenY := bodyBottom - 46            // dinaikkan lebih tinggi (dari -30 ke -42)
		pdf.Image(elemenName, elemenX, elemenY, elemenW, 0, false, "", 0, "")
		pdf.SetAlpha(1, "Normal")
	}

	pdf.ClipEnd()
}

// ============================================================================
// HELPER — ELEMEN KECIL UI KARTU
// ============================================================================

// drawInfoRow menggambar satu baris "label : nilai" — label dan nilai sama-sama
// warna hitam, bedanya nilai dibuat bold supaya lebih menonjol.
func drawInfoRow(pdf *gofpdf.Fpdf, x, y float64, label, value string, labelW, maxWidth float64) {
	const baseFont = 7.0 // Diperbesar dari 6.5 ke 7.0

	pdf.SetFont("Helvetica", "", baseFont)
	pdf.SetTextColor(colors.TextDarkR, colors.TextDarkG, colors.TextDarkB)
	pdf.SetXY(x, y)
	pdf.CellFormat(labelW, 3.6, label, "", 0, "L", false, 0, "")
	pdf.SetXY(x+labelW-6.5, y) // Digeser lebih jauh ke kiri dari -5.5 ke -6.5
	pdf.CellFormat(2, 3.6, ":", "", 0, "L", false, 0, "")

	valueW := maxWidth - labelW + 4.5 // Disesuaikan dari +3.5 ke +4.5
	fontSize := baseFont
	pdf.SetFont("Helvetica", "B", fontSize)
	for pdf.GetStringWidth(value) > valueW && fontSize > 4.8 {
		fontSize -= 0.2
		pdf.SetFont("Helvetica", "B", fontSize)
	}
	pdf.SetXY(x+labelW-4.5, y) // Disesuaikan dari -3.5 ke -4.5
	pdf.CellFormat(valueW, 3.6, value, "", 0, "L", false, 0, "")
}

// drawInfoRowWithAutoResize menggambar satu baris dengan auto-resize yang lebih agresif
// untuk menghindari tabrakan dengan barcode
func drawInfoRowWithAutoResize(pdf *gofpdf.Fpdf, x, y float64, label, value string, labelW, maxWidth float64) {
	const baseFont = 7.0 // Diperbesar dari 6.5 ke 7.0

	pdf.SetFont("Helvetica", "", baseFont)
	pdf.SetTextColor(colors.TextDarkR, colors.TextDarkG, colors.TextDarkB)
	pdf.SetXY(x, y)
	pdf.CellFormat(labelW, 3.6, label, "", 0, "L", false, 0, "")
	pdf.SetXY(x+labelW-6.5, y) // Digeser lebih jauh ke kiri dari -5.5 ke -6.5
	pdf.CellFormat(2, 3.6, ":", "", 0, "L", false, 0, "")

	valueW := maxWidth - labelW + 4.5 // Disesuaikan dari +3.5 ke +4.5
	fontSize := baseFont
	pdf.SetFont("Helvetica", "B", fontSize)
	// Resize lebih agresif untuk menghindari barcode, minimal 4.5
	for pdf.GetStringWidth(value) > valueW && fontSize > 4.5 {
		fontSize -= 0.15
		pdf.SetFont("Helvetica", "B", fontSize)
	}
	pdf.SetXY(x+labelW-4.5, y) // Disesuaikan dari -3.5 ke -4.5
	pdf.CellFormat(valueW, 3.6, value, "", 0, "L", false, 0, "")
}

// drawInfoRowMultiLine menggambar satu baris "label : nilai" dengan dukungan multi-line untuk value
// dan mengembalikan tinggi total yang digunakan
func drawInfoRowMultiLine(pdf *gofpdf.Fpdf, x, y float64, label, value string, labelW, maxWidth float64) float64 {
	const baseFont = 7.0 // Diperbesar dari 6.5 ke 7.0
	const lineHeight = 3.6

	pdf.SetFont("Helvetica", "", baseFont)
	pdf.SetTextColor(colors.TextDarkR, colors.TextDarkG, colors.TextDarkB)
	pdf.SetXY(x, y)
	pdf.CellFormat(labelW, lineHeight, label, "", 0, "L", false, 0, "")
	pdf.SetXY(x+labelW-6.5, y) // Digeser lebih jauh ke kiri dari -5.5 ke -6.5
	pdf.CellFormat(2, lineHeight, ":", "", 0, "L", false, 0, "")

	valueW := maxWidth - labelW + 4.5 // Disesuaikan dari +3.5 ke +4.5
	pdf.SetFont("Helvetica", "B", baseFont)
	
	// Cek apakah value muat dalam satu baris
	if pdf.GetStringWidth(value) <= valueW {
		// Muat dalam satu baris, gunakan CellFormat biasa
		pdf.SetXY(x+labelW-4.5, y) // Disesuaikan dari -3.5 ke -4.5
		pdf.CellFormat(valueW, lineHeight, value, "", 0, "L", false, 0, "")
		return lineHeight
	}
	
	// Jika tidak muat, gunakan MultiCell untuk multi-line
	pdf.SetXY(x+labelW-4.5, y) // Disesuaikan dari -3.5 ke -4.5
	startY := pdf.GetY()
	pdf.MultiCell(valueW, lineHeight, value, "", "L", false)
	endY := pdf.GetY()
	
	return endY - startY
}

// drawCornerMarks menggambar aksen sudut di kotak foto, ala bingkai kamera studio.
func drawCornerMarks(pdf *gofpdf.Fpdf, x, y, w, h, size float64) {
	pdf.SetDrawColor(colors.AccentR, colors.AccentG, colors.AccentB)
	pdf.SetLineWidth(0.5)
	corners := [][2]float64{{x, y}, {x + w, y}, {x, y + h}, {x + w, y + h}}
	dirs := [][2]float64{{1, 1}, {-1, 1}, {1, -1}, {-1, -1}}
	for i, c := range corners {
		dx, dy := dirs[i][0], dirs[i][1]
		pdf.Line(c[0], c[1], c[0]+dx*size, c[1])
		pdf.Line(c[0], c[1], c[0], c[1]+dy*size)
	}
}

// drawCutMarks menggambar tanda potong tipis di luar kartu, memudahkan proses pemotongan cetak.
func drawCutMarks(pdf *gofpdf.Fpdf, x, y, w, h float64) {
	pdf.SetDrawColor(180, 180, 180)
	pdf.SetLineWidth(0.15)
	m := 2.0 // panjang tanda potong
	corners := [][2]float64{{x, y}, {x + w, y}, {x, y + h}, {x + w, y + h}}
	dirs := [][2]float64{{-1, -1}, {1, -1}, {-1, 1}, {1, 1}}
	for i, c := range corners {
		dx, dy := dirs[i][0], dirs[i][1]
		pdf.Line(c[0]+dx*0.5, c[1], c[0]+dx*m, c[1])
		pdf.Line(c[0], c[1]+dy*0.5, c[0], c[1]+dy*m)
	}
}

// ============================================================================
// HELPER — DOWNLOAD GAMBAR & PARSING HTML
// ============================================================================

func downloadImage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download image: status %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

// stripHTML menghapus tag HTML dan mengembalikan teks polos (dipakai untuk VISI).
func stripHTML(htmlStr string) string {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		re := regexp.MustCompile(`<[^>]*>`)
		return strings.TrimSpace(re.ReplaceAllString(htmlStr, ""))
	}

	var result strings.Builder
	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			if text != "" {
				result.WriteString(text)
				result.WriteString(" ")
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}
	extract(doc)
	return strings.TrimSpace(result.String())
}

// extractListItems mengambil tiap poin <li> dari HTML MISI agar bisa ditampilkan
// sebagai bullet list terpisah, bukan satu paragraf panjang. Kalau tidak ada
// tag <li> (misalnya teks polos dengan pemisah titik), fallback ke pemisahan kalimat.
func extractListItems(htmlStr string) []string {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return splitSentences(stripHTML(htmlStr))
	}

	var items []string
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "li" {
			var sb strings.Builder
			var extractText func(*html.Node)
			extractText = func(nn *html.Node) {
				if nn.Type == html.TextNode {
					sb.WriteString(nn.Data)
				}
				for c := nn.FirstChild; c != nil; c = c.NextSibling {
					extractText(c)
				}
			}
			extractText(n)
			if text := strings.TrimSpace(sb.String()); text != "" {
				items = append(items, text)
			}
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)

	if len(items) == 0 {
		return splitSentences(stripHTML(htmlStr))
	}
	return items
}

func splitSentences(text string) []string {
	raw := strings.Split(text, ".")
	var out []string
	for _, s := range raw {
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s+".")
		}
	}
	return out
}