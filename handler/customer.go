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
