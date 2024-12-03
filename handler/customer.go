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
	"github.com/IsraelTeo/api-paw-go/service"
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
	if err := db.GDB.Preload("pet_id").First(&customer, id).Error; err != nil {
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
	if err := db.GDB.Preload("Pet").Find(&customers).Error; err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Customers were not found", nil)
		payload.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	empty := service.VerifyListEmpty(customers)
	if empty {
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

	if err := service.ValidateEntity(&customer); err != nil {
		log.Printf("validation error: %v", err)
		response := payload.NewResponse(payload.MessageTypeError, "Bad request.", nil)
		payload.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	if exists, err := service.ValidateUniqueField("email", customer.Email, &model.Customer{}); err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Internal server error", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	} else if exists {
		response := payload.NewResponse(payload.MessageTypeError, "Email already exists", nil)
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
	if err := db.GDB.Preload("Pet").First(&customer, uint(id)).Error; err != nil {
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

	var input model.Customer
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Bad request", nil)
		payload.ResponseJSON(w, http.StatusBadRequest, response)
		log.Printf("bad request %v:", err)
		return
	}

	if err := service.ValidateEntity(&customer); err != nil {
		log.Printf("validation error: %v", err)
		response := payload.NewResponse(payload.MessageTypeError, "Bad request", nil)
		payload.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	customer.FirstName = input.FirstName
	customer.LastName = input.LastName
	customer.DNI = input.DNI
	customer.Email = input.Email
	customer.PhoneNumber = input.PhoneNumber

	if err := db.GDB.Save(&customer).Error; err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Error updating employee", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		log.Printf("error updating employee: %v", err)
		return
	}

	if err := db.GDB.Preload("Pet").First(&customer, uint(customer.ID)).Error; err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Error loading pet data", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		log.Printf("error loading pet data: %v", err)
		return
	}

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

	if customer.PetID != 0 {
		pet := model.Pet{}
		if err := db.GDB.Delete(&pet, customer.PetID).Error; err != nil {
			response := payload.NewResponse(payload.MessageTypeError, "Error deleting pet", nil)
			payload.ResponseJSON(w, http.StatusInternalServerError, response)
			log.Printf("error deleting pet: %v", err)
			return
		}
	}

	if err := db.GDB.Delete(&customer).Error; err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Error deleting customer", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		log.Printf("error deleting customer: %v", err)
		return
	}

	response := payload.NewResponse(payload.MessageTypeSuccess, "Customer deleted successfull", nil)
	payload.ResponseJSON(w, http.StatusOK, response)
}
