package utils

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/saifwork/mock-service/internal/core/config"
)

// SendEmail sends an HTML email using Gmail SMTP.
func SendEmail(cfg *config.Config, to, subject, htmlBody string) error {
	from := cfg.GmailUser
	password := cfg.GmailPassKey

	if from == "" || password == "" {
		return fmt.Errorf("GMAIL_USER or GMAIL_PASS not set in env")
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Proper MIME headers for HTML
	msg := strings.Builder{}
	msg.WriteString(fmt.Sprintf("From: MockNode <%s>\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n")
	msg.WriteString(htmlBody)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg.String()))
}

func BuildEmailTemplate(appName, title, message, buttonText, buttonLink string) string {
	return fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>%s</title>
	</head>
	<body style="margin:0;padding:0;background:#f4f4f7;font-family:Arial,sans-serif;">
		<div style="max-width:600px;margin:40px auto;background:#ffffff;border-radius:12px;box-shadow:0 4px 12px rgba(0,0,0,0.1);overflow:hidden;">
			<div style="background:linear-gradient(90deg,#6C63FF,#5C6BC0);padding:20px;text-align:center;color:#fff;font-size:24px;font-weight:bold;">
				%s
			</div>
			<div style="padding:30px 25px;text-align:center;color:#333;">
				<h2 style="color:#444;">%s</h2>
				<p style="font-size:16px;color:#555;">%s</p>
				<a href="%s" style="display:inline-block;margin-top:25px;padding:12px 28px;background:#6C63FF;color:#fff;text-decoration:none;font-weight:600;border-radius:6px;">
					%s
				</a>
			</div>
			<div style="background:#fafafa;padding:15px;text-align:center;font-size:13px;color:#999;">
				© %d %s. All rights reserved.
			</div>
		</div>
	</body>
	</html>`, title, appName, title, message, buttonLink, buttonText, 2025, appName)
}
