package utils

import (
	"net/smtp"
	"os"
)

func SendEmail(toEmail, subject, body string) error {
	from := os.Getenv("EMAIL_ADDRESS")
	fromPass := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("EMAIL_HOST")

	smtpPort := "587"
	auth := smtp.PlainAuth("", from, fromPass, smtpHost)

	to := []string{toEmail}
	message := []byte(
		"Subject: Camagru - " + subject + "\r\n\r\n" +
			body + "\r\n",
	)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return nil
	}

	return nil
}
