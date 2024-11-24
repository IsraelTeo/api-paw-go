package model

import (
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	FirstName   string       `json:"first_name" gorm:"size:70;not null" validate:"required,min=2,max=70"`
	LastName    string       `json:"last_name" gorm:"size:70;not null" validate:"required,min=3,max=90"`
	DNI         string       `json:"dni" gorm:"size:15;unique;not null" validate:"required,len=15"`
	Email       string       `json:"email" gorm:"size:100;unique;not null" validate:"required,email"`
	PhoneNumber string       `json:"phone_number" gorm:"unique;size:15" validate:"required,len=15,numeric"`
	Direction   string       `json:"direction" gorm:"size:100" validate:"required,max=100"`
	BirthDate   time.Time    `json:"birth_date" gorm:"not null" validate:"required"`
	TypeID      uint         `json:"type_id" gorm:"not null;index" validate:"required"`
	Type        EmployeeType `json:"type" gorm:"foreignKey:TypeID"`
}



