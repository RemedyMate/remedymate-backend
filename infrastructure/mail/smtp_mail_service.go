package mail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
)

type SMTPMailService struct {
	host string
	port string
	user string
	pass string
	from string
}

func NewSMTPMailService() *SMTPMailService {
	return &SMTPMailService{
		host: os.Getenv("SMTP_HOST"),
		port: os.Getenv("SMTP_PORT"),
		user: os.Getenv("SMTP_USER"),
		pass: os.Getenv("SMTP_PASS"),
		from: os.Getenv("SMTP_FROM"),
	}
}

func (s *SMTPMailService) Send(to, subject, htmlBody string) error {
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	auth := smtp.PlainAuth("", s.user, s.pass, s.host)

	msg := "From: " + s.from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n" +
		htmlBody

	// Try TLS if port 465
	if s.port == "465" {
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         s.host,
		}
		conn, err := tls.Dial("tcp", addr, tlsconfig)
		if err != nil {
			return err
		}
		c, err := smtp.NewClient(conn, s.host)
		if err != nil {
			return err
		}
		if err = c.Auth(auth); err != nil {
			return err
		}
		if err = c.Mail(s.from); err != nil {
			return err
		}
		if err = c.Rcpt(to); err != nil {
			return err
		}
		wc, err := c.Data()
		if err != nil {
			return err
		}
		if _, err = wc.Write([]byte(msg)); err != nil {
			return err
		}
		if err = wc.Close(); err != nil {
			return err
		}
		return c.Quit()
	}

	// Fallback STARTTLS/normal
	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg))
}
