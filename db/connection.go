package db

import (
	"fmt"
	"os"

	"github.com/IsraelTeo/api-paw-go/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var GDB *gorm.DB

func Connection() error {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)

	var err error
	if GDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
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
