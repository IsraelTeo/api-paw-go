package model

import "gorm.io/gorm"

type Customer struct {
	gorm.Model
	ID          uint   `json:"id" gorm:"unique;primary_key;auto_increment"`
	FirstName   string `json:"first_name" gorm:"size:70;not null"`
	LastName    string `json:"last_name" gorm:"size:70;not null"`
	DNI         string `json:"dni" gorm:"size:15;unique;not null"`
	Email       string `json:"email" gorm:"size:100;unique;not null"`
	PhoneNumber string `json:"phone_number" gorm:"unique;size:15"`
	Pets        []Pet  `json:"pets" gorm:"foreignKey:CustomerID"`
}
