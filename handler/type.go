package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/IsraelTeo/api-paw-go/db"
	"github.com/IsraelTeo/api-paw-go/model"
	"github.com/IsraelTeo/api-paw-go/payload"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// GetEmployeeTypeById maneja la solicitud HTTP GET para obtener un tipo de empleado por su ID.
// @Description Devuelve un tipo de empleado especificado por su ID.
// @Accept json
// @Produce json
// @Param id path int true "ID del tipo de empleado"
// @Success 200 {object} payload.Response{MessageType=string, Message=string, Data=model.EmployeeType} "Tipo de empleado encontrado"
// @Failure 400 {object} payload.Response{MessageType=string, Message=string} "Método no permitido"
// @Failure 404 {object} payload.Response{MessageType=string, Message=string} "Tipo de empleado no encontrado"
// @Failure 405 {object} payload.Response{MessageType=string, Message=string} "Método no permitido"
// @Router /api/v1/type/{id} [get]
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
	payload.ResponseJSON(w, http.StatusOK, response)
}

// GetAllEmployeeTypes maneja la solicitud HTTP GET para obtener todos los tipos de empleados.
// @Description Devuelve todos los tipos de empleados registrados en el sistema.
// @Accept json
// @Produce json
// @Success 200 {object} payload.Response{MessageType=string, Message=string, Data=[]model.EmployeeType} "Tipos de empleados encontrados"
// @Failure 404 {object} payload.Response{MessageType=string, Message=string} "Tipos de empleados no encontrados"
// @Failure 405 {object} payload.Response{MessageType=string, Message=string} "Método no permitido"
// @Failure 204 {object} payload.Response{MessageType=string, Message=string} "Lista de tipos de empleados vacía"
// @Router /api/v1/type [get]
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

// SaveEmployeeType maneja la solicitud HTTP POST para guardar un nuevo tipo de empleado.
// @Description Registra un nuevo tipo de empleado en el sistema.
// @Accept json
// @Produce json
// @Param employeeType body model.EmployeeType true "Nuevo tipo de empleado"
// @Success 201 {object} payload.Response{MessageType=string, Message=string} "Tipo de empleado creado exitosamente"
// @Failure 400 {object} payload.Response{MessageType=string, Message=string} "Solicitud incorrecta o JSON inválido"
// @Failure 409 {object} payload.Response{MessageType=string, Message=string} "El tipo de empleado ya existe"
// @Failure 500 {object} payload.Response{MessageType=string, Message=string} "Error interno del servidor"
// @Router /api/v1/type/{1} [post]
func UpdateEmployeeType(w http.ResponseWriter, r *http.Request) {
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

	var input model.EmployeeType
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Invalid request body", nil)
		payload.ResponseJSON(w, http.StatusBadRequest, response)
		log.Printf("error decoding request body: %v", err)
		return
	}

	employeeType.Name = input.Name

	db.GDB.Save(&employeeType)
	response := payload.NewResponse(payload.MessageTypeSuccess, "EmployeeType updated successfull", employeeType)
	payload.ResponseJSON(w, http.StatusOK, response)
}

// UpdateEmployeeType maneja la solicitud HTTP PUT para actualizar un tipo de empleado existente.
// @Description Actualiza un tipo de empleado existente por su ID.
// @Accept json
// @Produce json
// @Param id path int true "ID del tipo de empleado"
// @Param employeeType body model.EmployeeType true "Tipo de empleado actualizado"
// @Success 200 {object} payload.Response{MessageType=string, Message=string, Data=model.EmployeeType} "Tipo de empleado actualizado"
// @Failure 400 {object} payload.Response{MessageType=string, Message=string} "ID inválido o formato incorrecto"
// @Failure 404 {object} payload.Response{MessageType=string, Message=string} "Tipo de empleado no encontrado"
// @Failure 500 {object} payload.Response{MessageType=string, Message=string} "Error interno del servidor"
// @Router /api/v1/type/{id} [put]
func DeleteEmployeeType(w http.ResponseWriter, r *http.Request) {
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
