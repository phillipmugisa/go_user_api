package api

import (
	"fmt"
	"net/smtp"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// utility function for send verification code to email
func SendEmail(username, email, code string) error {
	subject := "Activation Code"

	// Message content
	message := []byte("To: " + email + "\r\n" +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"\r\n" +
		fmt.Sprintf("hello %s.\r\nThis is your activation code: %s.", username, code))

	// Connect to the SMTP server
	auth := smtp.PlainAuth("", os.Getenv("SENDER_EMAIL"), os.Getenv("SENDER_EMAIL_PASSWORD"), os.Getenv("SMTP_HOST"))

	err := smtp.SendMail(os.Getenv("SMTP_HOST")+":"+os.Getenv("SMTP_PORT"), auth, os.Getenv("SENDER_EMAIL"), []string{email}, message)

	return err
}

func HashPassword(password string) (string, error) {
    // Generate a salted hash for the password.
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}