package model

import (
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	ID          uint      `json:"id" gorm:"unique;primary_key;auto_increment"`
	FirstName   string    `json:"first_name" gorm:"size:70;not null"`
	LastName    string    `json:"last_name" gorm:"size:70;not null"`
	DNI         string    `json:"dni" gorm:"size:15;unique;not null"`
	Email       string    `json:"email" gorm:"size:100;unique;not null"`
	PhoneNumber string    `json:"phone_number" gorm:"unique;size:15"`
	Direction   string    `json:"direction" gorm:"size:100"`
	BirthDate   time.Time `json:"birth_date" gorm:"not null"`
	Role        Role      `json:"role" gorm:"not null"`
}
