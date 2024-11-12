package model

import "gorm.io/gorm"

type Pet struct {
	gorm.Model
	ID         uint    `json:"id" gorm:"unique;primary_key;auto_increment"`
	Name       string  `json:"name" gorm:"size:70"`
	Specie     string  `json:"specie" gorm:"size:50"`
	Gender     string  `json:"gender" gorm:"size:10"`
	Race       string  `json:"race" gorm:"size:50"`
	Age        uint    `json:"age"`
	Weight     float64 `json:"weight"`
	CustomerID uint    `json:"customer_id" gorm:"not null;index"`
}
