package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"gopkg.in/gomail.v2"
)

// EmailService handles email operations
type EmailService struct {
	smtpHost     string
	smtpPort     int
	smtpUsername string
	smtpPassword string
	fromName     string
	fromEmail    string
}

// NewEmailService creates a new email service
func NewEmailService() *EmailService {
	return &EmailService{
		smtpHost:     os.Getenv("SMTP_HOST"),
		smtpPort:     587, // Default SMTP port
		smtpUsername: os.Getenv("SMTP_USERNAME"),
		smtpPassword: os.Getenv("SMTP_PASSWORD"),
		fromName:     os.Getenv("SMTP_FROM_NAME"),
		fromEmail:    os.Getenv("SMTP_FROM_EMAIL"),
	}
}

// EmailData represents data for email template
type EmailData struct {
	// Informasi Pengirim
	IDTiket string
	Nama    string
	Email   string
	Telepon string

	// Informasi Pertanyaan
	TanggalPengajuan string
	Kategori         string
	Prioritas        string
	JudulPertanyaan  string
	DeskripsiPertanyaan string
	FilePertanyaan   []FileLink

	// Informasi Jawaban
	JudulJawaban     string
	DeskripsiJawaban string
	FileJawaban      []FileLink
}

// FileLink represents a file with name and URL
type FileLink struct {
	Name string
	URL  string
}

// SendPertanyaanReply sends email reply for pertanyaan
func (e *EmailService) SendPertanyaanReply(to string, data EmailData) error {
	// Create email message
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", e.fromName, e.fromEmail))
	m.SetHeader("To", to)
	m.SetHeader("Subject", fmt.Sprintf("Jawaban Pertanyaan - %s", data.IDTiket))

	// Generate HTML body
	htmlBody, err := e.generateEmailHTML(data)
	if err != nil {
		return fmt.Errorf("failed to generate email HTML: %w", err)
	}

	m.SetBody("text/html", htmlBody)

	// Attach files if any
	// Note: gomail doesn't support attaching from URL directly
	// Files are shown as hyperlinks in email body instead

	// Send email
	d := gomail.NewDialer(e.smtpHost, e.smtpPort, e.smtpUsername, e.smtpPassword)
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// generateEmailHTML generates HTML email body
func (e *EmailService) generateEmailHTML(data EmailData) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
            line-height: 1.6; 
            color: #1f2937; 
            background-color: #f3f4f6;
            padding: 20px;
        }
        .email-wrapper { 
            max-width: 650px; 
            margin: 0 auto; 
            background-color: #ffffff;
            border-radius: 12px;
            overflow: hidden;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        }
        .header { 
            background: linear-gradient(135deg, #DC2626 0%, #991B1B 100%);
            color: white; 
            padding: 40px 30px;
            text-align: center;
            position: relative;
        }
        .header::after {
            content: '';
            position: absolute;
            bottom: 0;
            left: 0;
            right: 0;
            height: 4px;
            background: linear-gradient(90deg, #FCA5A5, #DC2626, #FCA5A5);
        }
        .header h1 { 
            font-size: 28px; 
            margin-bottom: 8px;
            font-weight: 700;
            text-shadow: 0 2px 4px rgba(0,0,0,0.2);
        }
        .header p { 
            font-size: 16px; 
            opacity: 0.95;
            font-weight: 300;
        }
        .content { padding: 30px; }
        .section { 
            margin-bottom: 25px; 
            padding: 20px; 
            background-color: #fef2f2;
            border-left: 4px solid #DC2626;
            border-radius: 8px;
            transition: transform 0.2s;
        }
        .section:hover {
            transform: translateX(5px);
        }
        .section-title { 
            font-weight: 700; 
            color: #991B1B; 
            margin-bottom: 15px; 
            font-size: 18px;
        }
        .section-icon {
            display: inline-flex;
            align-items: center;
            justify-content: center;
            width: 32px;
            height: 32px;
            background-color: #DC2626;
            color: white;
            border-radius: 50%;
            font-size: 16px;
            line-height: 1;
        }
        .info-row { 
            margin: 12px 0;
            padding: 8px 0;
            border-bottom: 1px solid #fee2e2;
        }
        .info-row:last-child {
            border-bottom: none;
        }
        .label { 
            font-weight: 600; 
            color: #374151;
            display: block;
            margin-bottom: 4px;
        }
        .value { 
            color: #1f2937;
            word-wrap: break-word;
        }
        .description-row {
            margin: 12px 0;
            padding: 8px 0;
        }
        .description-label {
            font-weight: 600; 
            color: #374151;
            display: block;
            margin-bottom: 8px;
        }
        .description-value {
            color: #1f2937;
            word-wrap: break-word;
            line-height: 1.6;
            padding: 12px;
            background-color: rgba(255, 255, 255, 0.5);
            border-radius: 6px;
            border-left: 3px solid #e5e7eb;
        }
        .rich-text {
            /* Rich text editor content styling */
        }
        .rich-text p {
            margin: 8px 0;
            line-height: 1.6;
        }
        .rich-text p:first-child {
            margin-top: 0;
        }
        .rich-text p:last-child {
            margin-bottom: 0;
        }
        .rich-text span {
            /* Inherit color from parent instead of inline styles */
            color: inherit !important;
        }
        .rich-text strong, .rich-text b {
            font-weight: 700;
        }
        .rich-text em, .rich-text i {
            font-style: italic;
        }
        .rich-text ul, .rich-text ol {
            margin: 8px 0;
            padding-left: 20px;
        }
        .rich-text li {
            margin: 4px 0;
        }
        .file-link { 
            display: inline-block; 
            margin: 8px 12px 8px 0; 
            padding: 10px 18px; 
            background: linear-gradient(135deg, #DC2626 0%, #B91C1C 100%);
            color: white !important; 
            text-decoration: none; 
            border-radius: 6px;
            font-weight: 600;
            font-size: 14px;
            transition: all 0.3s;
            box-shadow: 0 2px 4px rgba(220, 38, 38, 0.3);
        }
        .file-link:hover { 
            background: linear-gradient(135deg, #B91C1C 0%, #991B1B 100%);
            transform: translateY(-2px);
            box-shadow: 0 4px 8px rgba(220, 38, 38, 0.4);
            color: white !important;
        }
        
        /* Question Section - Blue Theme */
        .question-section {
            background: linear-gradient(135deg, #eff6ff 0%, #dbeafe 100%);
            border-left: 4px solid #2563eb;
        }
        .question-section .section-title {
            color: #1d4ed8;
        }
        .question-section .section-icon {
            background-color: #2563eb;
        }
        .question-section .info-row {
            border-bottom: 1px solid #bfdbfe;
        }
        .question-section .description-value {
            border-left: 3px solid #2563eb;
            background-color: rgba(37, 99, 235, 0.05);
        }
        .question-section .file-link {
            background: linear-gradient(135deg, #2563eb 0%, #1d4ed8 100%);
            box-shadow: 0 2px 4px rgba(37, 99, 235, 0.3);
            color: white !important;
        }
        .question-section .file-link:hover {
            background: linear-gradient(135deg, #1d4ed8 0%, #1e40af 100%);
            box-shadow: 0 4px 8px rgba(37, 99, 235, 0.4);
            color: white !important;
        }

        /* Answer Section - Green Theme */
        .answer-section {
            background: linear-gradient(135deg, #f0fdf4 0%, #dcfce7 100%);
            border-left: 4px solid #16a34a;
            padding: 25px;
            border-radius: 8px;
            margin: 25px 0;
        }
        .answer-section .section-title {
            color: #15803d;
        }
        .answer-section .section-icon {
            background-color: #16a34a;
        }
        .answer-section .info-row {
            border-bottom: 1px solid #bbf7d0;
        }
        .answer-section .description-value {
            border-left: 3px solid #16a34a;
            background-color: rgba(22, 163, 74, 0.05);
        }
        .answer-section .file-link {
            background: linear-gradient(135deg, #16a34a 0%, #15803d 100%);
            box-shadow: 0 2px 4px rgba(22, 163, 74, 0.3);
            color: white !important;
        }
        .answer-section .file-link:hover {
            background: linear-gradient(135deg, #15803d 0%, #166534 100%);
            box-shadow: 0 4px 8px rgba(22, 163, 74, 0.4);
            color: white !important;
        }
        .confirmation { 
            background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
            padding: 20px; 
            border-left: 4px solid #f59e0b; 
            margin: 25px 0;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(245, 158, 11, 0.2);
        }
        .confirmation strong {
            color: #92400e;
            font-size: 16px;
            display: block;
            margin-bottom: 8px;
        }
        .footer { 
            background: linear-gradient(135deg, #1f2937 0%, #111827 100%);
            color: #e5e7eb;
            padding: 30px;
            margin-top: 30px;
        }
        .footer-title { 
            font-weight: 700; 
            margin-bottom: 15px;
            color: #DC2626;
            font-size: 18px;
            display: flex;
            align-items: center;
            gap: 8px;
        }
        .footer p { 
            margin: 8px 0;
            color: #d1d5db;
            line-height: 1.8;
        }
        .footer strong {
            color: #f3f4f6;
        }
        .divider {
            height: 2px;
            background: linear-gradient(90deg, transparent, #DC2626, transparent);
            margin: 20px 0;
        }
        @media only screen and (max-width: 600px) {
            .content { padding: 20px; }
            .header { padding: 30px 20px; }
            .header h1 { font-size: 24px; }
            .label { min-width: 100px; display: block; margin-bottom: 4px; }
        }
    </style>
</head>
<body>
    <div class="email-wrapper">
        <div class="header">
            <h1>PINTU SDN Sukapura 01</h1>
            <p>Sistem Informasi Terpadu</p>
        </div>

        <div class="content">
            <div class="section">
                <div class="section-title">
                    Informasi Pengirim
                </div>
                <div class="info-row">
                    <span class="label">ID Tiket:</span> 
                    <span class="value" style="font-weight: 700; color: #DC2626;">{{.IDTiket}}</span>
                </div>
                <div class="info-row">
                    <span class="label">Nama:</span> 
                    <span class="value">{{.Nama}}</span>
                </div>
                <div class="info-row">
                    <span class="label">Email:</span> 
                    <span class="value">{{.Email}}</span>
                </div>
                <div class="info-row">
                    <span class="label">Telepon:</span> 
                    <span class="value">{{.Telepon}}</span>
                </div>
            </div>

            <div class="section question-section">
                <div class="section-title">
                    Pertanyaan Anda
                </div>
                <div class="info-row">
                    <span class="label">Tanggal Pengajuan:</span> 
                    <span class="value">{{.TanggalPengajuan}}</span>
                </div>
                <div class="info-row">
                    <span class="label">Kategori:</span> 
                    <span class="value">{{.Kategori}}</span>
                </div>
                <div class="info-row">
                    <span class="label">Prioritas:</span> 
                    <span class="value" style="font-weight: 600;">{{.Prioritas}}</span>
                </div>
                <div class="info-row">
                    <span class="label">Judul:</span> 
                    <span class="value" style="font-weight: 600;">{{.JudulPertanyaan}}</span>
                </div>
                <div class="description-row">
                    <span class="description-label">Deskripsi:</span>
                    <div class="description-value rich-text">{{safeHTML .DeskripsiPertanyaan}}</div>
                </div>
                {{if .FilePertanyaan}}
                <div class="info-row">
                    <span class="label">File Lampiran:</span><br>
                    <div style="margin-top: 10px;">
                        {{range $index, $file := .FilePertanyaan}}
                        <a href="{{$file.URL}}" class="file-link" target="_blank">📎 File {{add $index 1}}</a>
                        {{end}}
                    </div>
                </div>
                {{end}}
            </div>

            <div class="divider"></div>

            <div class="answer-section">
                <div class="section-title">
                    Jawaban Kami
                </div>
                <div class="info-row">
                    <span class="label">Judul Jawaban:</span> 
                    <span class="value" style="font-weight: 700;">{{.JudulJawaban}}</span>
                </div>
                <div class="description-row">
                    <span class="description-label">Deskripsi:</span>
                    <div class="description-value rich-text">{{safeHTML .DeskripsiJawaban}}</div>
                </div>
                {{if .FileJawaban}}
                <div class="info-row">
                    <span class="label">File Jawaban:</span><br>
                    <div style="margin-top: 10px;">
                        {{range $index, $file := .FileJawaban}}
                        <a href="{{$file.URL}}" class="file-link" target="_blank">📄 {{$file.Name}}</a>
                        {{end}}
                    </div>
                </div>
                {{end}}
            </div>

            <div class="confirmation">
                <strong>⚠️ Konfirmasi Diperlukan</strong>
                Mohon balas email ini untuk mengkonfirmasi apakah jawaban kami sudah memenuhi kebutuhan Anda atau jika Anda memerlukan klarifikasi lebih lanjut.<br><br>
                
                <strong>Catatan Penting:</strong> Jika tidak ada balasan dalam waktu 3 (tiga) hari kerja sejak email ini dikirim, maka kami akan menganggap bahwa jawaban yang diberikan telah memenuhi kebutuhan Anda dan pertanyaan ini akan ditutup secara otomatis. Terima kasih atas perhatian dan kerjasamanya.
            </div>
        </div>

        <div class="footer">
            <div class="footer-title">
                Informasi Kontak
            </div>
            <p><strong>Alamat:</strong><br>
            Jl. Beo No.15, Komp.Walikota No.2, RT.12/RW.6<br>
            Sukapura, Kec. Cilincing, Jakarta Utara<br>
            DKI Jakarta 14140</p>
            
            <p><strong>Telepon:</strong> 021-4411729</p>
            <p><strong>Email:</strong> sdnsukapuraa01@gmail.com</p>
            
            <p><strong>Jam Operasional:</strong><br>
            Senin - Jumat: 06.30 - 15.00<br>
            Sabtu - Minggu: Tutup</p>
            
            <div style="margin-top: 20px; padding-top: 20px; border-top: 1px solid #374151; text-align: center; color: #9ca3af; font-size: 12px;">
                © 2026 SDN Sukapura 01. All rights reserved.
            </div>
        </div>
    </div>
</body>
</html>
`

	// Parse template
	t, err := template.New("email").Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
	}).Parse(tmpl)
	if err != nil {
		return "", err
	}

	// Execute template
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
