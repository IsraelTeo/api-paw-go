package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/IsraelTeo/api-paw-go/db"
	"github.com/IsraelTeo/api-paw-go/model"
	"github.com/IsraelTeo/api-paw-go/payload"
	"github.com/IsraelTeo/api-paw-go/service"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetEmployeeById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := payload.NewResponse(payload.MessageTypeError, "Invalid Method", nil)
		payload.ResponseJSON(w, http.StatusMethodNotAllowed, response)
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	employee := model.Employee{}
	if err := db.GDB.First(&employee, id).Error; err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Employee was not found", nil)
		payload.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := payload.NewResponse(payload.MessageTypeSuccess, "Employee found", employee)
	payload.ResponseJSON(w, http.StatusOK, response)
}

func GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := payload.NewResponse(payload.MessageTypeError, "Method get not permit", nil)
		payload.ResponseJSON(w, http.StatusMethodNotAllowed, response)
		return
	}

	var employees []model.Employee
	if err := db.GDB.Find(&employees).Error; err != nil {
		log.Printf("employees list not found %v:", err)
		response := payload.NewResponse(payload.MessageTypeError, "Employees not found", nil)
		payload.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	empty := service.VerifyListEmpty(employees)
	if empty {
		response := payload.NewResponse(payload.MessageTypeSuccess, "Employees List empty", nil)
		payload.ResponseJSON(w, http.StatusNoContent, response)
		return
	}

	response := payload.NewResponse(payload.MessageTypeSuccess, "Employees found", employees)
	payload.ResponseJSON(w, http.StatusOK, response)
}

func SaveEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := payload.NewResponse(payload.MessageTypeError, "Method post not permit", nil)
		payload.ResponseJSON(w, http.StatusMethodNotAllowed, response)
		return
	}

	employee := model.Employee{}
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Bad request: invalid JSON data", nil)
		payload.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	parsedDate, err := time.Parse("2006-01-02", employee.BirthDateRaw)
	if err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Invalid date format, expected YYYY-MM-DD", nil)
		payload.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	employee.BirthDate = parsedDate

	if err := service.ValidateEntity(&employee); err != nil {
		log.Printf("validation error: %v", err)
		response := payload.NewResponse(payload.MessageTypeError, "Bad request.", nil)
		payload.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	if exists, err := service.ValidateUniqueField("dni", employee.DNI, &model.Employee{}); err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Internal server error.", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	} else if exists {
		response := payload.NewResponse(payload.MessageTypeError, "DNI already exists", nil)
		payload.ResponseJSON(w, http.StatusConflict, response)
		return
	}

	if exists, err := service.ValidateUniqueField("email", employee.Email, &model.Employee{}); err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Internal server error", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	} else if exists {
		response := payload.NewResponse(payload.MessageTypeError, "Email already exists", nil)
		payload.ResponseJSON(w, http.StatusConflict, response)
		return
	}

	if exists, err := service.ValidateUniqueField("phone_number", employee.PhoneNumber, &model.Employee{}); err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Internal server error", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	} else if exists {
		response := payload.NewResponse(payload.MessageTypeError, "Phone number already exists", nil)
		payload.ResponseJSON(w, http.StatusConflict, response)
		return
	}

	if result := db.GDB.Create(&employee); result.Error != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Internal Server Error", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := payload.NewResponse(payload.MessageTypeSuccess, "Employee created successfusly", nil)
	payload.ResponseJSON(w, http.StatusCreated, response)
}

func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		response := payload.NewResponse(payload.MessageTypeError, "Method PUT not allowed", nil)
		payload.ResponseJSON(w, http.StatusMethodNotAllowed, response)
		return
	}

	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Invalid ID format", nil)
		payload.ResponseJSON(w, http.StatusBadRequest, response)
		log.Printf("invalid ID format: %v", err)
		return
	}

	employee := model.Employee{}
	if err := db.GDB.First(&employee, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := payload.NewResponse(payload.MessageTypeError, "Employee not found", nil)
			payload.ResponseJSON(w, http.StatusNotFound, response)
			log.Printf("employee not found: %v", err)
			return
		}

		response := payload.NewResponse(payload.MessageTypeError, "Database error", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		log.Printf("database error: %v", err)
		return
	}

	var input model.Employee
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Bad request", nil)
		payload.ResponseJSON(w, http.StatusBadRequest, response)
		log.Printf("bad request %v:", err)
		return
	}

	if err := service.ValidateEntity(&employee); err != nil {
		log.Printf("validation error: %v", err)
		response := payload.NewResponse(payload.MessageTypeError, "Bad request", nil)
		payload.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	employee.FirstName = input.FirstName
	employee.LastName = input.LastName
	employee.BirthDate = input.BirthDate
	employee.DNI = input.DNI
	employee.Direction = input.Direction
	employee.PhoneNumber = input.PhoneNumber
	employee.Email = input.Email
	employee.TypeID = input.TypeID

	if err := db.GDB.Save(&employee).Error; err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Error saving employee", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		log.Printf("error saving employee: %v", err)
		return
	}

	response := payload.NewResponse(payload.MessageTypeSuccess, "Employee updated successfully", employee)
	payload.ResponseJSON(w, http.StatusOK, response)
}

func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		response := payload.NewResponse(payload.MessageTypeError, "Method delete not permit", nil)
		payload.ResponseJSON(w, http.StatusMethodNotAllowed, response)
		return
	}

	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Invalid ID format", nil)
		payload.ResponseJSON(w, http.StatusBadRequest, response)
		log.Printf("invalid ID format: %v", err)
		return
	}

	employee := model.Employee{}
	if err := db.GDB.First(&employee, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := payload.NewResponse(payload.MessageTypeError, "Employee not found", nil)
			payload.ResponseJSON(w, http.StatusNotFound, response)
			log.Printf("employee not found: %v", err)
			return
		}

		response := payload.NewResponse(payload.MessageTypeError, "Database error", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		log.Printf("database error: %v", err)
		return
	}

	db.GDB.Delete(&employee)
	response := payload.NewResponse(payload.MessageTypeSuccess, "Employee deleted successfull", nil)
	payload.ResponseJSON(w, http.StatusOK, response)
}
