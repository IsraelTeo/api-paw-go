package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       uint   `json:"id" gorm:"primary_key;auto_increment"`
	Email    string `json:"email" gorm:"size:100;unique;not_null"`
	Password string `json:"password" gorm:"size:100"`
	Role     Role   `json:"role"`
	
}
