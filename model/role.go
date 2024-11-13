package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	ID   uint   `json:"id" gorm:"primary_key;unique;auto_increment"`
	Name string `json:"name" gorm:"size:100;unique;not null;size:20"`
}
