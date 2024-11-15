package model

import "gorm.io/gorm"

type EmployeeType struct {
	gorm.Model
	Name string `json:"name" gorm:"unique;not null;size:20"`
}
