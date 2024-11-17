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

// GetCustomerById maneja la solicitud HTTP GET para obtener un cliente por su ID.
// @Description Obtiene un cliente especificado por su ID, incluyendo la información de sus mascotas asociadas.
// @Accept json
// @Produce json
// @Param id path int true "ID del cliente"
// @Success 200 {object} payload.Response{MessageType=string, Message=string, Data=model.Customer} "Cliente encontrado"
// @Failure 400 {object} payload.Response{MessageType=string, Message=string} "Método no permitido"
// @Failure 404 {object} payload.Response{MessageType=string, Message=string} "Cliente no encontrado"
// @Failure 405 {object} payload.Response{MessageType=string, Message=string} "Método no permitido"
// @Router /api/v1/customer/{id} [get]
func GetCustomerById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := payload.NewResponse(payload.MessageTypeError, "Invalid Method", nil)
		payload.ResponseJSON(w, http.StatusMethodNotAllowed, response)
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	customer := model.Customer{}
	if err := db.GDB.Preload("Pets").First(&customer, id).Error; err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Customer was not found", nil)
		payload.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := payload.NewResponse(payload.MessageTypeSuccess, "Customer found", customer)
	payload.ResponseJSON(w, http.StatusOK, response)
}

// GetAllCustomers maneja la solicitud HTTP GET para obtener todos los clientes.
// @Description Obtiene una lista de todos los clientes registrados, incluyendo la información de sus mascotas asociadas.
// @Accept json
// @Produce json
// @Success 200 {object} payload.Response{MessageType=string, Message=string, Data=[]model.Customer} "Clientes encontrados"
// @Failure 400 {object} payload.Response{MessageType=string, Message=string} "Método no permitido"
// @Failure 404 {object} payload.Response{MessageType=string, Message=string} "Clientes no encontrados"
// @Failure 405 {object} payload.Response{MessageType=string, Message=string} "Método no permitido"
// @Failure 204 {object} payload.Response{MessageType=string, Message=string} "Lista de clientes vacía"
// @Router /api/v1/customers [get]
func GetAllCustomers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := payload.NewResponse(payload.MessageTypeError, "Method get not permit", nil)
		payload.ResponseJSON(w, http.StatusMethodNotAllowed, response)
		return
	}

	var customers []model.Customer
	if err := db.GDB.Preload("Pets").Find(&customers).Error; err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Customer was not found", nil)
		payload.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	if len(customers) == 0 {
		response := payload.NewResponse(payload.MessageTypeSuccess, "Customers List empty", nil)
		payload.ResponseJSON(w, http.StatusNoContent, response)
		return
	}

	response := payload.NewResponse(payload.MessageTypeSuccess, "Customers found", customers)
	payload.ResponseJSON(w, http.StatusOK, response)
}

// SaveCustomer maneja la solicitud HTTP POST para registrar un nuevo cliente.
// @Description Crea un nuevo cliente en el sistema.
// @Accept json
// @Produce json
// @Param customer body model.Customer true "Nuevo cliente"
// @Success 201 {object} payload.Response{MessageType=string, Message=string} "Cliente creado exitosamente"
// @Failure 400 {object} payload.Response{MessageType=string, Message=string} "Solicitud incorrecta o JSON inválido"
// @Failure 409 {object} payload.Response{MessageType=string, Message=string} "Email ya está en uso"
// @Failure 500 {object} payload.Response{MessageType=string, Message=string} "Error interno del servidor"
// @Router /api/v1/customer [post]
func SaveCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := payload.NewResponse(payload.MessageTypeError, "Method post not permit", nil)
		payload.ResponseJSON(w, http.StatusMethodNotAllowed, response)
		return
	}

	customer := model.Customer{}
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Bad request: invalid JSON data", nil)
		payload.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	if err := db.GDB.Where("email = ?", customer.Email).First(&customer).Error; err == nil {
		response := payload.NewResponse(payload.MessageTypeError, "Email already in use", nil)
		payload.ResponseJSON(w, http.StatusConflict, response)
		return
	}

	if result := db.GDB.Create(&customer); result.Error != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Internal Server Error", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := payload.NewResponse(payload.MessageTypeSuccess, "Customer created successfusly", nil)
	payload.ResponseJSON(w, http.StatusCreated, response)
}

// UpdateCustomer maneja la solicitud HTTP PUT para actualizar un cliente existente.
// @Description Actualiza los datos de un cliente en el sistema.
// @Accept json
// @Produce json
// @Param id path int true "ID del cliente"
// @Param customer body model.Customer true "Cliente actualizado"
// @Success 200 {object} payload.Response{MessageType=string, Message=string, Data=model.Customer} "Cliente actualizado"
// @Failure 400 {object} payload.Response{MessageType=string, Message=string} "ID inválido o formato incorrecto"
// @Failure 404 {object} payload.Response{MessageType=string, Message=string} "Cliente no encontrado"
// @Failure 500 {object} payload.Response{MessageType=string, Message=string} "Error interno del servidor"
// @Router /api/v1/customer/{id} [put]
func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
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

	customer := model.Customer{}
	if err := db.GDB.First(&customer, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := payload.NewResponse(payload.MessageTypeError, "Customer not found", nil)
			payload.ResponseJSON(w, http.StatusNotFound, response)
			log.Printf("customer not found: %v", err)
			return
		}

		response := payload.NewResponse(payload.MessageTypeError, "Database error", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		log.Printf("database error: %v", err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Bad request", nil)
		payload.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	db.GDB.Save(&customer)
	response := payload.NewResponse(payload.MessageTypeSuccess, "Customer updated successfull", customer)
	payload.ResponseJSON(w, http.StatusOK, response)
}

// DeleteCustomer maneja la solicitud HTTP DELETE para eliminar un cliente por su ID.
// @Description Elimina un cliente especificado por su ID, y también elimina sus mascotas asociadas.
// @Accept json
// @Produce json
// @Param id path int true "ID del cliente"
// @Success 200 {object} payload.Response{MessageType=string, Message=string} "Cliente eliminado"
// @Failure 400 {object} payload.Response{MessageType=string, Message=string} "ID inválido"
// @Failure 404 {object} payload.Response{MessageType=string, Message=string} "Cliente no encontrado"
// @Failure 500 {object} payload.Response{MessageType=string, Message=string} "Error interno del servidor"
// @Router /api/v1/customer/{id} [delete]
func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
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
	customer := model.Customer{}
	if err := db.GDB.First(&customer, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := payload.NewResponse(payload.MessageTypeError, "Customer not found", nil)
			payload.ResponseJSON(w, http.StatusNotFound, response)
			log.Printf("customer not found: %v", err)
			return
		}

		response := payload.NewResponse(payload.MessageTypeError, "Database error", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		log.Printf("database error: %v", err)
		return
	}

	if err := db.GDB.Where("customer_id = ?", customer.ID).Delete(&model.Pet{}).Error; err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Failed to delete associated pets", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		log.Printf("failed to delete associated pets: %v", err)
		return
	}

	db.GDB.Delete(&customer)
	response := payload.NewResponse(payload.MessageTypeSuccess, "Customer deleted successfull", nil)
	payload.ResponseJSON(w, http.StatusOK, response)
}
