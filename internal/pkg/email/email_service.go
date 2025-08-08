package email

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"

	"github.com/jokosaputro95/cms-go/config"
)

type EmailService interface {
	SendVerificationEmail(to, token, username string) error
	SendWelcomeEmail(to, username string) error
}

type emailService struct {
	cfg *config.Config
}

func NewEmailService(cfg *config.Config) *emailService {
	return &emailService{cfg: cfg}
}

// SendVerificationEmail mengirimkan email verifikasi
func (s *emailService) SendVerificationEmail(to, token, username string) error {
	var body bytes.Buffer

	// Data yang akan dimasukkan ke template
	data := EmailData{
		AppName:         s.cfg.Server.AppName, 
		FirstName:       username, 
		VerificationURL: fmt.Sprintf("http://localhost:%s/auth/verify-email?token=%s", s.cfg.Server.ServerPort, token),
		AppURL:          fmt.Sprintf("http://localhost:%s", s.cfg.Server.ServerPort),
		SupportURL:      "http://localhost/support", // Ganti dengan URL dukungan Anda
		ExpiresIn:       "30 minutes",
	}

	// Persiapkan pesan email dengan header
	subject := "Verifikasi Email Anda"
	headers := map[string]string{
		"From":         s.cfg.Email.EmailSMTPUsername,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=utf-8",
	}

	for k, v := range headers {
		body.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}

	body.WriteString("\r\n")

	// Eksekusi template HTML
	err := emailTemplates["verification_html"].Execute(&body, data)
	if err != nil {
		return fmt.Errorf("gagal mengeksekusi template email: %w", err)
	}

	// Autentikasi dan kirim email
	auth := smtp.PlainAuth("", s.cfg.Email.EmailSMTPUsername, s.cfg.Email.EmailSMTPPassword, s.cfg.Email.EmailSMTPHost)
	addr := fmt.Sprintf("%s:%s", s.cfg.Email.EmailSMTPHost, s.cfg.Email.EmailSMTPPort)
	
	err = smtp.SendMail(addr, auth, s.cfg.Email.EmailSMTPUsername, []string{to}, body.Bytes())
	if err != nil {
		return fmt.Errorf("gagal mengirim email: %w", err)
	}

	return nil
}

func (s *emailService) SendWelcomeEmail(to, username string) error {
	log.Printf("Mengirim email selamat datang ke %s", to)
	return nil
}