package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/IsraelTeo/api-paw-go/db"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var validate *validator.Validate

func InitValidator() {
	if validate == nil {
		validate = validator.New()
	}
}

func ValidateEntity[T any](model *T) error {
	return validate.Struct(model)
}

func VerifyListEmpty[T any](list []T) bool {
	return len(list) == 0
}

func checkIfFieldExists[T any](field string, value interface{}, model *T) (bool, error) {
	err := db.GDB.Where(field+" = ?", value).First(model).Error
	if err == nil {
		return true, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return false, err
}

func ValidateUniqueField[T any](field, value string, model *T) (bool, error) {
	exists, err := checkIfFieldExists(field, value, model)
	return exists, err
}

func IsEmpty(s string) bool {
	return s == ""
}

func ValidateBirthDate(birthDate string) error {
	_, err := time.Parse("2006-01-02", birthDate)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYY-MM-DD")
	}
	return nil
}
