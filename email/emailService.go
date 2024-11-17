package email

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/IsraelTeo/api-paw/config"
)

type EmailService struct {
	User     string
	Password string
	Host     string
	Port     string
}

func NewEmailService() *EmailService {
	return &EmailService{
		User:     config.AppConfig.User,
		Password: config.AppConfig.Password,
		Host:     "smtp.gmail.com",
		Port:     "587",
	}
}

func (e *EmailService) SendEmail(toUser, subject, messageBody string) error {
	from := e.User
	auth := smtp.PlainAuth("", e.User, e.Password, e.Host)

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		from, toUser, subject, messageBody,
	))

	if err := smtp.SendMail(
		fmt.Sprintf("%s:%s", e.Host, e.Port),
		auth,
		from,
		[]string{toUser},
		msg,
	); err != nil {
		log.Printf("Error: %v\n", err)
		return fmt.Errorf("error sending email: %w", err)
	}

	log.Println("Correo enviado exitosamente a", toUser)
	return nil
}
