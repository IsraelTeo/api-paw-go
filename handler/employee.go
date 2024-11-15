package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/IsraelTeo/api-paw/db"
	"github.com/IsraelTeo/api-paw/model"
	"github.com/IsraelTeo/api-paw/payload"
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

	if len(employees) == 0 {
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

	if err := db.GDB.Where("dni = ?", employee.Email).First(&employee).Error; err == nil {
		response := payload.NewResponse(payload.MessageTypeError, "DNI already exists", nil)
		payload.ResponseJSON(w, http.StatusConflict, response)
		return
	}

	if err := db.GDB.Where("email = ?", employee.Email).First(&employee).Error; err == nil {
		response := payload.NewResponse(payload.MessageTypeError, "Email already in use", nil)
		payload.ResponseJSON(w, http.StatusConflict, response)
		return
	}

	if err := db.GDB.Where("phone_number = ?", employee.Email).First(&employee).Error; err == nil {
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
		response := payload.NewResponse(payload.MessageTypeError, "Method put not permit", nil)
		payload.ResponseJSON(w, http.StatusMethodNotAllowed, response)
		return
	}

	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr) // Convertir `id` a uint
	if err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Invalid ID format", nil)
		payload.ResponseJSON(w, http.StatusBadRequest, response)
		log.Printf("invalid ID format: %v", err)
		return
	}

	// Buscar el empleado en la base de datos
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

	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Bad request", nil)
		payload.ResponseJSON(w, http.StatusBadRequest, response)
		log.Printf("bad request %v:", err)
		return
	}

	db.GDB.Save(&employee)
	response := payload.NewResponse(payload.MessageTypeSuccess, "Employee updated successfull", employee)
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
	id, err := strconv.Atoi(idStr) // Convertir `id` a uint
	if err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Invalid ID format", nil)
		payload.ResponseJSON(w, http.StatusBadRequest, response)
		log.Printf("invalid ID format: %v", err)
		return
	}

	// Buscar el empleado en la base de datos
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
