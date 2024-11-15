package db

import (
	"os"

	"github.com/IsraelTeo/api-paw/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var GDB *gorm.DB

func Connection() error {
	var err error
	if GDB, err = gorm.Open(mysql.Open(os.Getenv("CONNECTION_STRING")), &gorm.Config{}); err != nil {
		return err
	}
	return nil
}

func MigrateDataBase() error {
	err := GDB.AutoMigrate(
		&model.Customer{},
		&model.Employee{},
		&model.Pet{},
		&model.User{},
		&model.EmployeeType{},
	)

	if err != nil {
		return err
	}
	return nil
}
