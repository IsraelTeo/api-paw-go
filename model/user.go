package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"size:100;unique;not_null"`
	Password string `json:"password" gorm:"size:100"`
	IsAdmin  bool   `json:"is_admin" gorm:"dafault:false"`
}

func VerifyPassword(passwordHashed string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHashed), []byte(password))
}
