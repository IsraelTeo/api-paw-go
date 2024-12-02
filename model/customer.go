package model

import "gorm.io/gorm"

type Customer struct {
	gorm.Model
	FirstName   string `json:"first_name" gorm:"size:70;not null"`
	LastName    string `json:"last_name" gorm:"size:70;not null"`
	DNI         string `json:"dni" gorm:"size:15;unique;not null"`
	Email       string `json:"email" gorm:"size:100;unique;not null"`
	PhoneNumber string `json:"phone_number" gorm:"unique;size:15"`
	PetID       uint   `json:"pet_id" gorm:"index" validate:"required"`
	Pet         Pet    `json:"pet" gorm:"foreignKey:PetID;constraint:OnDelete:CASCADE"`
}
