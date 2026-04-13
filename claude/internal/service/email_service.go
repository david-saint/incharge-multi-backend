package service

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log/slog"
	"net/smtp"
	"path/filepath"
	"runtime"

	"github.com/incharge/server/internal/config"
)

// EmailService handles sending transactional emails via SMTP.
type EmailService struct {
	cfg       *config.MailConfig
	templates map[string]*template.Template
}

// NewEmailService creates and initializes the email service.
func NewEmailService(cfg *config.MailConfig) *EmailService {
	svc := &EmailService{
		cfg:       cfg,
		templates: make(map[string]*template.Template),
	}
	svc.loadTemplates()
	return svc
}

func (s *EmailService) loadTemplates() {
	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Join(filepath.Dir(filename), "..", "..", "templates", "emails")

	for _, name := range []string{"verify", "reset"} {
		path := filepath.Join(baseDir, name+".html")
		t, err := template.ParseFiles(path)
		if err != nil {
			slog.Warn("email template not found, using fallback", "template", name, "error", err)
			// Create a minimal fallback template.
			t = template.Must(template.New(name).Parse(fallbackTemplate(name)))
		}
		s.templates[name] = t
	}
}

func fallbackTemplate(name string) string {
	switch name {
	case "verify":
		return `<html><body><h2>Verify Your Email</h2><p>Hi {{.Name}},</p><p>Click the link below to verify your email:</p><p><a href="{{.URL}}">Verify Email</a></p><p>This link expires in 72 hours.</p><p>— InCharge</p></body></html>`
	case "reset":
		return `<html><body><h2>Reset Your Password</h2><p>Hi,</p><p>Click the link below to reset your password:</p><p><a href="{{.URL}}">Reset Password</a></p><p>This link expires in 60 minutes.</p><p>— InCharge</p></body></html>`
	default:
		return `<html><body>{{.}}</body></html>`
	}
}

// VerifyEmailData holds data for the verification email template.
type VerifyEmailData struct {
	Name string
	URL  string
}

// ResetEmailData holds data for the password reset email template.
type ResetEmailData struct {
	URL string
}

// SendVerificationEmail sends a verification email to the user.
func (s *EmailService) SendVerificationEmail(to, name, verifyURL string) error {
	var body bytes.Buffer
	if err := s.templates["verify"].Execute(&body, VerifyEmailData{Name: name, URL: verifyURL}); err != nil {
		return fmt.Errorf("failed to render verify template: %w", err)
	}
	return s.sendHTML(to, "Verify Your Email - InCharge", body.String())
}

// SendPasswordResetEmail sends a password reset email.
func (s *EmailService) SendPasswordResetEmail(to, resetURL string) error {
	var body bytes.Buffer
	if err := s.templates["reset"].Execute(&body, ResetEmailData{URL: resetURL}); err != nil {
		return fmt.Errorf("failed to render reset template: %w", err)
	}
	return s.sendHTML(to, "Reset Your Password - InCharge", body.String())
}

// sendHTML sends an HTML email via SMTP.
func (s *EmailService) sendHTML(to, subject, htmlBody string) error {
	if s.cfg.Host == "" {
		slog.Warn("SMTP not configured, skipping email", "to", to, "subject", subject)
		return nil
	}

	from := s.cfg.FromAddr
	headers := fmt.Sprintf("From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n",
		s.cfg.FromName, from, to, subject)

	msg := []byte(headers + htmlBody)
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)

	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)

	if s.cfg.Encryption == "tls" {
		return s.sendTLS(addr, auth, from, to, msg)
	}

	return smtp.SendMail(addr, auth, from, []string{to}, msg)
}

func (s *EmailService) sendTLS(addr string, auth smtp.Auth, from, to string, msg []byte) error {
	tlsConfig := &tls.Config{
		ServerName: s.cfg.Host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		// Fall back to STARTTLS.
		return s.sendSTARTTLS(addr, auth, from, to, msg)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.cfg.Host)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.Auth(auth); err != nil {
		return err
	}
	if err := client.Mail(from); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		return err
	}
	return w.Close()
}

func (s *EmailService) sendSTARTTLS(addr string, auth smtp.Auth, from, to string, msg []byte) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()

	tlsConfig := &tls.Config{ServerName: s.cfg.Host}
	if err := client.StartTLS(tlsConfig); err != nil {
		return err
	}
	if err := client.Auth(auth); err != nil {
		return err
	}
	if err := client.Mail(from); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		return err
	}
	return w.Close()
}
