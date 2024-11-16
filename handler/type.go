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

// GetEmployeeTypeById godoc
// @Summary Get an employee type by ID
// @Description Get details of an employee type by its ID
// @Tags EmployeeTypes
// @Accept  json
// @Produce  json
// @Param id path int true "Employee Type ID"
// @Success 200 {object} payload.Response{data=model.EmployeeType} "Employee Type found"
// @Failure 400 {object} payload.Response "Invalid ID format"
// @Failure 404 {object} payload.Response "Employee Type not found"
// @Failure 405 {object} payload.Response "Method not allowed"
// @Router /type/{id} [get]
func GetEmployeeTypeById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := payload.NewResponse(payload.MessageTypeError, "Invalid Method", nil)
		payload.ResponseJSON(w, http.StatusMethodNotAllowed, response)
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	role := model.EmployeeType{}
	if err := db.GDB.First(&role, id).Error; err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Employee Type was not found", nil)
		payload.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := payload.NewResponse(payload.MessageTypeSuccess, "Employee Type found", role)
	payload.ResponseJSON(w, http.StatusOK, response)
}

// GetAllEmployeeTypes godoc
// @Summary Get all employee types
// @Description Get a list of all employee types
// @Tags EmployeeTypes
// @Accept  json
// @Produce  json
// @Success 200 {object} payload.Response{data=[]model.EmployeeType} "Employee Types found"
// @Failure 204 {object} payload.Response "Employee Types List empty"
// @Failure 405 {object} payload.Response "Method not allowed"
// @Router /types [get]
func GetAllEmployeeTypes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := payload.NewResponse(payload.MessageTypeError, "Method get not permit", nil)
		payload.ResponseJSON(w, http.StatusMethodNotAllowed, response)
		return
	}

	var roles []model.EmployeeType
	if err := db.GDB.Find(&roles).Error; err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Employee Types not found", nil)
		payload.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	if len(roles) == 0 {
		response := payload.NewResponse(payload.MessageTypeSuccess, "Employee Types List empty", nil)
		payload.ResponseJSON(w, http.StatusNoContent, response)
		return
	}

	response := payload.NewResponse(payload.MessageTypeSuccess, "Employee Types found", roles)
	payload.ResponseJSON(w, http.StatusNoContent, response)
}

// SaveEmployeeType godoc
// @Summary Create a new employee type
// @Description Save a new employee type
// @Tags EmployeeTypes
// @Accept  json
// @Produce  json
// @Param body body model.EmployeeType true "Employee Type data"
// @Success 201 {object} payload.Response "Employee Type created successfully"
// @Failure 400 {object} payload.Response "Bad request: invalid JSON data"
// @Failure 409 {object} payload.Response "Employee Type already exists"
// @Failure 405 {object} payload.Response "Method not allowed"
// @Failure 500 {object} payload.Response "Internal Server Error"
// @Router /type [post]
func SaveEmployeeType(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := payload.NewResponse(payload.MessageTypeError, "Method post not permit", nil)
		payload.ResponseJSON(w, http.StatusMethodNotAllowed, response)
		return
	}

	role := model.EmployeeType{}
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Bad request: invalid JSON data", nil)
		payload.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	if err := db.GDB.Where("name = ?", role.Name).First(&role).Error; err == nil {
		response := payload.NewResponse(payload.MessageTypeError, "Employee Type already exists", nil)
		payload.ResponseJSON(w, http.StatusConflict, response)
		return
	}

	if result := db.GDB.Create(&role); result.Error != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Internal Server Error", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := payload.NewResponse(payload.MessageTypeSuccess, "Employee Type created successfusly", nil)
	payload.ResponseJSON(w, http.StatusCreated, response)
}

// UpdateEmployeeType godoc
// @Summary Update an existing employee type
// @Description Update the details of an existing employee type
// @Tags EmployeeTypes
// @Accept  json
// @Produce  json
// @Param id path int true "Employee Type ID"
// @Param body body model.EmployeeType true "Updated Employee Type data"
// @Success 200 {object} payload.Response{data=model.EmployeeType} "Employee Type updated successfully"
// @Failure 400 {object} payload.Response "Invalid ID format"
// @Failure 404 {object} payload.Response "Employee Type not found"
// @Failure 405 {object} payload.Response "Method not allowed"
// @Failure 500 {object} payload.Response "Internal Server Error"
// @Router /type/{id} [put]
func UpdateEmployeeType(w http.ResponseWriter, r *http.Request) {
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

	employeeType := model.EmployeeType{}
	if err := db.GDB.First(&employeeType, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := payload.NewResponse(payload.MessageTypeError, "Employee Type not found", nil)
			payload.ResponseJSON(w, http.StatusNotFound, response)
			log.Printf("customer not found: %v", err)
			return
		}

		response := payload.NewResponse(payload.MessageTypeError, "Database error", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		log.Printf("database error: %v", err)
		return
	}

	db.GDB.Save(&employeeType)
	response := payload.NewResponse(payload.MessageTypeSuccess, "EmployeeType updated successfull", employeeType)
	payload.ResponseJSON(w, http.StatusOK, response)
}

// DeleteEmployeeType godoc
// @Summary Delete an employee type by ID
// @Description Delete an employee type from the system by its ID
// @Tags EmployeeTypes
// @Accept  json
// @Produce  json
// @Param id path int true "Employee Type ID"
// @Success 200 {object} payload.Response "Employee Type deleted successfully"
// @Failure 400 {object} payload.Response "Invalid ID format"
// @Failure 404 {object} payload.Response "Employee Type not found"
// @Failure 405 {object} payload.Response "Method not allowed"
// @Failure 500 {object} payload.Response "Internal Server Error"
// @Router /type/{id} [delete]
func DeleteEmployeeType(w http.ResponseWriter, r *http.Request) {
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

	employeeType := model.EmployeeType{}
	if err := db.GDB.First(&employeeType, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := payload.NewResponse(payload.MessageTypeError, "Employee Type not found", nil)
			payload.ResponseJSON(w, http.StatusNotFound, response)
			log.Printf("customer not found: %v", err)
			return
		}

		response := payload.NewResponse(payload.MessageTypeError, "Database error", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		log.Printf("database error: %v", err)
		return
	}

	db.GDB.Delete(&employeeType)
	response := payload.NewResponse(payload.MessageTypeSuccess, "Employee type deleted successfull", nil)
	payload.ResponseJSON(w, http.StatusOK, response)
}
