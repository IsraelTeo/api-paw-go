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

// GetEmployeeById maneja la solicitud HTTP GET para obtener un empleado por su ID.
// @Description Obtiene un empleado especificado por su ID.
// @Accept json
// @Produce json
// @Param id path int true "ID del empleado"
// @Success 200 {object} payload.Response{MessageType=string, Message=string, Data=model.Employee} "Empleado encontrado"
// @Failure 400 {object} payload.Response{MessageType=string, Message=string} "Método no permitido"
// @Failure 404 {object} payload.Response{MessageType=string, Message=string} "Empleado no encontrado"
// @Failure 405 {object} payload.Response{MessageType=string, Message=string} "Método no permitido"
// @Router /api/v1/employee/{id} [get]
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

// GetAllEmployees maneja la solicitud HTTP GET para obtener todos los empleados.
// @Description Obtiene una lista con todos los empleados registrados en el sistema.
// @Accept json
// @Produce json
// @Success 200 {object} payload.Response{MessageType=string, Message=string, Data=[]model.Employee} "Empleados encontrados"
// @Failure 400 {object} payload.Response{MessageType=string, Message=string} "Método no permitido"
// @Failure 404 {object} payload.Response{MessageType=string, Message=string} "Empleados no encontrados"
// @Failure 405 {object} payload.Response{MessageType=string, Message=string} "Método no permitido"
// @Failure 204 {object} payload.Response{MessageType=string, Message=string} "Lista de empleados vacía"
// @Router /api/v1/employees [get]
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

// SaveEmployee maneja la solicitud HTTP POST para registrar un nuevo empleado.
// @Description Crea un nuevo empleado en el sistema.
// @Accept json
// @Produce json
// @Param employee body model.Employee true "Nuevo empleado"
// @Success 201 {object} payload.Response{MessageType=string, Message=string} "Empleado creado exitosamente"
// @Failure 400 {object} payload.Response{MessageType=string, Message=string} "Solicitud incorrecta o JSON inválido"
// @Failure 409 {object} payload.Response{MessageType=string, Message=string} "DNI, Email o Número de teléfono ya existen"
// @Failure 500 {object} payload.Response{MessageType=string, Message=string} "Error interno del servidor"
// @Router /api/v1/employee [post]
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
