package utils

import (
	"fmt"
	"net/smtp"
)

func SendVerificationCodeEmail(toEmail, verificationCode string) error {
	from := "your-email@gmail.com"
	password := "your-email-password"

	// SMTP серверийн тохиргоо
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Имэйл агуулга
	subject := "Таны баталгаажуулах код"
	body := fmt.Sprintf("Таны баталгаажуулах код: %s", verificationCode)

	message := []byte("Subject: " + subject + "\r\n" + "To: " + toEmail + "\r\n" + "\r\n" + body)

	// SMTP аутентификац
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Имэйл илгээх
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, message)
	if err != nil {
		return fmt.Errorf("имэйл илгээхэд алдаа гарлаа: %v", err)
	}
	return nil
}

// 📌 Түр нууц үг илгээх функц

func SendPasswordResetEmail(toEmail, tempPassword string) error {
	from := "orgiloorgil16@gmail.com"
	password := "wxja rcue qzee ymwl"

	// SMTP тохиргоо
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Имэйл контент
	subject := "Таны түр нууц үг"
	body := fmt.Sprintf("Таны түр нууц үг: %s", tempPassword)

	message := []byte("Subject: " + subject + "\r\n" + "To: " + toEmail + "\r\n" + "\r\n" + body)

	// Имэйл илгээх
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, message)
	if err != nil {
		return fmt.Errorf("имэйл илгээхэд алдаа гарлаа: %v", err)
	}
	return nil
}
func SendVerificationEmail(toEmail, verificationToken string) error {
	from := "orgiloorgil16@gmail.com"
	password := "wxja rcue qzee ymwl" // Таны имэйл хаягийн нууц үг

	// SMTP серверийн тохиргоо
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Имэйл агуулга
	subject := "Имэйл баталгаажуулалт"
	body := fmt.Sprintf("Таны баталгаажуулах код: %s", verificationToken)

	message := []byte("Subject: " + subject + "\r\n" + "To: " + toEmail + "\r\n" + "\r\n" + body)

	// SMTP аутентификац
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Имэйл илгээх
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, message)
	if err != nil {
		return fmt.Errorf("имэйл илгээхэд алдаа гарлаа: %v", err)
	}
	return nil
}
