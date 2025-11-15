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
	<html lang="en">
	<head>
	<meta charset="UTF-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
	<title>%s</title>
	</head>

	<body style="
		margin:0; 
		padding:0; 
		background:#0f0f17; 
		font-family:'Inter', Arial, sans-serif;
	">

		<!-- Container -->
		<div style="
			max-width:600px; 
			margin:40px auto; 
			padding:0;
			border-radius:20px;
			background:rgba(20,20,35,0.85);
			backdrop-filter:blur(6px);
			border:1px solid rgba(255,255,255,0.08);
			box-shadow:
				0 0 25px rgba(0,150,255,0.25),
				0 0 45px rgba(180,70,255,0.15);
			overflow:hidden;
		">

			<!-- Header -->
			<div style="
				padding:32px 20px; 
				text-align:center; 
				color:#fff;
				background:linear-gradient(135deg,#6b5bff,#7a2cff,#ff2cdf);
				box-shadow:inset 0 0 30px rgba(255,255,255,0.15);
			">
				<div style="
					font-size:28px; 
					font-weight:800; 
					letter-spacing:-0.6px;
					text-shadow:0 0 8px rgba(255,255,255,0.35);
				">
					%s
				</div>
			</div>

			<!-- Body -->
			<div style="padding:36px 30px; color:#e6e6e6; text-align:center;">
				
				<h2 style="
					color:#ffffff; 
					font-size:26px; 
					margin:0 0 18px; 
					font-weight:700; 
					text-shadow:0 0 12px rgba(255, 0, 221, 0.35);
				">
					%s
				</h2>

				<p style="
					font-size:16px; 
					line-height:1.7; 
					color:#c7c7d1; 
					margin:0 0 28px;
				">
					%s
				</p>

				<!-- CTA -->
				<a href="%s" style="
					display:inline-block;
					padding:15px 40px;
					font-size:17px;
					font-weight:700;
					color:#fff;
					text-decoration:none;
					border-radius:10px;
					background:linear-gradient(135deg,#6b5bff,#8b33ff,#ff2cdf);
					box-shadow:
						0 6px 18px rgba(138, 43, 226, 0.45),
						0 0 15px rgba(255,0,150,0.45);
					transition:all 0.3s ease;
				">
					%s
				</a>

				<p style="
					font-size:13px; 
					color:#8a8a99; 
					margin-top:30px; 
					line-height:1.6;
				">
					If you didn’t request this, you may ignore this email.<br>
					This link will expire shortly for security reasons.
				</p>
			</div>

			<!-- Footer -->
			<div style="
				background:#11111a; 
				padding:16px; 
				text-align:center; 
				font-size:13px; 
				color:#666;
				border-top:1px solid rgba(255,255,255,0.08);
			">
				© %d %s — All rights reserved.
			</div>

		</div>
	</body>
	</html>
	`, title, appName, title, message, buttonLink, buttonText, 2025, appName)
}
