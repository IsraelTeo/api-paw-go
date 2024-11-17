package config

import "os"

type EmailConfig struct {
	User     string
	Password string
}

var AppConfig EmailConfig

func GetEmailConfig() EmailConfig {
	return EmailConfig{
		User:     os.Getenv("EMAIL_USER"),
		Password: os.Getenv("EMAIL_PASSWORD"),
	}
}
