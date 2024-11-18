package service

import (
	"errors"

	"github.com/IsraelTeo/api-paw-go/db"
	"gorm.io/gorm"
)

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
