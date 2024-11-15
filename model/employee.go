package model

import (
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	FirstName   string       `json:"first_name" gorm:"size:70;not null"`
	LastName    string       `json:"last_name" gorm:"size:70;not null"`
	DNI         string       `json:"dni" gorm:"size:15;unique;not null"`
	Email       string       `json:"email" gorm:"size:100;unique;not null"`
	PhoneNumber string       `json:"phone_number" gorm:"unique;size:15"`
	Direction   string       `json:"direction" gorm:"size:100"`
	BirthDate   time.Time    `json:"birth_date" gorm:"not null"`
	TypeID      uint         `json:"type_id" gorm:"not null;index"`
	Type        EmployeeType `json:"type" gorm:"foreignKey:TypeID"`
}
