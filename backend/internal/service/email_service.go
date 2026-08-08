package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"github.com/kelseyhightower/envconfig"

	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/program"
)

type EmailMessage struct {
	To       string
	Subject  string
	HTMLBody string
}

type TemplateMessage struct {
	To       string
	Subject  string
	Template string
	Data     any
}

type EmailService interface {
	Send(ctx context.Context, msg EmailMessage) error
	SendTemplate(ctx context.Context, msg TemplateMessage) error
}

type smtpConfig struct {
	Host     string `envconfig:"host" required:"true"`
	Port     int    `envconfig:"port" default:"587"`
	Username string `envconfig:"username" required:"true"`
	Password string `envconfig:"password" required:"true"`
	From     string `envconfig:"from" required:"true"`
	FromName string `envconfig:"from_name" default:"Selectify"`
}

type emailService struct {
	cfg     smtpConfig
	logoURL string
}

func NewEmailService() EmailService {
	var wrapper struct {
		SMTP   smtpConfig `envconfig:"smtp"`
		CDNURL string     `envconfig:"cdn_url" required:"true"`
	}

	prefix := program.AppPrefix
	if prefix == "" {
		prefix = "EVT"
	}

	if err := envconfig.Process(prefix, &wrapper); err != nil {
		panic(logger.Error(context.Background(), err, "failed to process SMTP env vars"))
	}

	base := strings.TrimRight(wrapper.CDNURL, "/")
	return &emailService{
		cfg:     wrapper.SMTP,
		logoURL: base + "/logos/logo.svg",
	}
}

func (s *emailService) SendTemplate(ctx context.Context, msg TemplateMessage) error {
	if msg.Template == "" {
		return fmt.Errorf("template is required")
	}

	// Todo:: Uncomment this when ready to send emails
	//htmlBody, err := email.Render(msg.Template, msg.Data, s.logoURL)
	//if err != nil {
	//	return logger.Error(ctx, err, "failed to render email template")
	//}

	//return s.Send(ctx, EmailMessage{
	//	To:       msg.To,
	//	Subject:  msg.Subject,
	//	HTMLBody: htmlBody,
	//})

	return nil
}

func (s *emailService) Send(ctx context.Context, msg EmailMessage) error {
	if msg.To == "" {
		return fmt.Errorf("email recipient is required")
	}
	if msg.Subject == "" {
		return fmt.Errorf("email subject is required")
	}
	if msg.HTMLBody == "" {
		return fmt.Errorf("email body is required")
	}

	fromHeader := s.cfg.From
	if s.cfg.FromName != "" {
		fromHeader = fmt.Sprintf("%s <%s>", s.cfg.FromName, s.cfg.From)
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("From: %s\r\n", fromHeader))
	builder.WriteString(fmt.Sprintf("To: %s\r\n", msg.To))
	builder.WriteString(fmt.Sprintf("Subject: %s\r\n", msg.Subject))
	builder.WriteString("MIME-Version: 1.0\r\n")
	builder.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	builder.WriteString("\r\n")
	builder.WriteString(msg.HTMLBody)

	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)

	if err := sendMail(addr, auth, s.cfg.From, []string{msg.To}, []byte(builder.String())); err != nil {
		return logger.Errorf(ctx, err, "failed to send email to %s", msg.To)
	}

	return nil
}

// sendMail dials the SMTP server, upgrades with STARTTLS when available, then sends the message.
// Uses an explicit STARTTLS path so Google SMTP (port 587) works reliably.
func sendMail(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		_ = conn.Close()
		return err
	}
	defer func() { _ = client.Close() }()

	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{ServerName: host}
		if err = client.StartTLS(tlsConfig); err != nil {
			return err
		}
	}

	if auth != nil {
		if ok, _ := client.Extension("AUTH"); ok {
			if err = client.Auth(auth); err != nil {
				return err
			}
		}
	}

	if err = client.Mail(from); err != nil {
		return err
	}
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return err
		}
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}
	if _, err = writer.Write(msg); err != nil {
		_ = writer.Close()
		return err
	}
	if err = writer.Close(); err != nil {
		return err
	}

	return client.Quit()
}
