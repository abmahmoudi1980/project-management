package services

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

// EmailService handles sending emails
type EmailService struct {
	smtpHost string
	smtpPort string
	smtpUser string
	smtpPass string
	appURL   string
}

// NewEmailService creates a new email service
func NewEmailService() *EmailService {
	return &EmailService{
		smtpHost: getEnv("SMTP_HOST", "smtp.gmail.com"),
		smtpPort: getEnv("SMTP_PORT", "587"),
		smtpUser: getEnv("SMTP_USER", ""),
		smtpPass: getEnv("SMTP_PASSWORD", ""),
		appURL:   getEnv("APP_URL", "http://localhost:5173"),
	}
}

// SendEmail sends a plain text email
func (s *EmailService) SendEmail(to, subject, body string) error {
	from := s.smtpUser
	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", from, to, subject, body))

	auth := smtp.PlainAuth("", s.smtpUser, s.smtpPass, s.smtpHost)
	addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)

	err := smtp.SendMail(addr, auth, from, []string{to}, msg)
	if err != nil {
		// In development, log the email instead of failing
		if s.smtpUser == "" {
			fmt.Printf("\n=== EMAIL (Development Mode) ===\n")
			fmt.Printf("To: %s\n", to)
			fmt.Printf("Subject: %s\n", subject)
			fmt.Printf("Body:\n%s\n", body)
			fmt.Printf("=================================\n\n")
			return nil
		}
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendPasswordResetEmail sends password reset email with Persian template
func (s *EmailService) SendPasswordResetEmail(to, token string) error {
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.appURL, token)

	subject := "بازیابی رمز عبور"
	body := s.getPasswordResetEmailTemplate(resetLink)

	return s.SendEmail(to, subject, body)
}

// getPasswordResetEmailTemplate returns Persian RTL email template
func (s *EmailService) getPasswordResetEmailTemplate(resetLink string) string {
	template := `
<!DOCTYPE html>
<html dir="rtl" lang="fa">
<head>
    <meta charset="UTF-8">
    <style>
        body {
            font-family: Tahoma, Arial, sans-serif;
            direction: rtl;
            text-align: right;
            background-color: #f4f4f4;
            padding: 20px;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #ffffff;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
        }
        h2 {
            color: #333333;
            margin-bottom: 20px;
        }
        p {
            color: #555555;
            line-height: 1.6;
            margin-bottom: 15px;
        }
        .button {
            display: inline-block;
            padding: 12px 30px;
            background-color: #3b82f6;
            color: #ffffff;
            text-decoration: none;
            border-radius: 5px;
            margin: 20px 0;
        }
        .button:hover {
            background-color: #2563eb;
        }
        .footer {
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #eeeeee;
            color: #999999;
            font-size: 12px;
        }
        .warning {
            background-color: #fef3c7;
            border-right: 4px solid #f59e0b;
            padding: 15px;
            margin: 20px 0;
            border-radius: 5px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h2>بازیابی رمز عبور</h2>
        <p>سلام،</p>
        <p>برای بازیابی رمز عبور خود، روی دکمه زیر کلیک کنید:</p>
        <p style="text-align: center;">
            <a href="%s" class="button">بازیابی رمز عبور</a>
        </p>
        <p>یا می‌توانید لینک زیر را کپی کرده و در مرورگر خود باز کنید:</p>
        <p style="word-break: break-all; background-color: #f9fafb; padding: 10px; border-radius: 5px;">%s</p>
        <div class="warning">
            <strong>توجه:</strong> این لینک تنها برای 1 ساعت معتبر است و فقط یک بار قابل استفاده می‌باشد.
        </div>
        <p>اگر شما درخواست بازیابی رمز عبور نداده‌اید، این ایمیل را نادیده بگیرید.</p>
        <div class="footer">
            <p>این ایمیل به صورت خودکار ارسال شده است. لطفاً به آن پاسخ ندهید.</p>
        </div>
    </div>
</body>
</html>
`
	return fmt.Sprintf(template, resetLink, resetLink)
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}
	return value
}
