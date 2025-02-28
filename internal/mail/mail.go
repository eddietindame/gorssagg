package mail

import (
	"fmt"
	"strconv"

	"github.com/eddietindame/gorssagg/internal/config"
	"gopkg.in/gomail.v2"
)

func SendResetEmail(toEmail, resetLink string) error {
	fmt.Println(config.EMAIL_ADDRESS, config.SMTP_HOST, config.SMTP_PORT, config.SMTP_USERNAME, config.SMTP_PASSWORD, toEmail)

	fromEmail := config.EMAIL_ADDRESS
	smtpPort, _ := strconv.Atoi(config.SMTP_PORT)

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", fromEmail)
	mailer.SetHeader("To", toEmail)
	mailer.SetHeader("Subject", "Password Reset Request")
	mailer.SetBody("text/html", fmt.Sprintf(`
		<h2>Password Reset Request</h2>
		<p>Click the link below to reset your password:</p>
		<a href="%s">%s</a>
	`, resetLink, resetLink))

	dialer := gomail.NewDialer(config.SMTP_HOST, smtpPort, config.SMTP_USERNAME, config.SMTP_PASSWORD)

	if err := dialer.DialAndSend(mailer); err != nil {
		return fmt.Errorf("failed to send reset email: %w", err)
	}

	return nil
}
