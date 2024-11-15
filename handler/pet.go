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

func GetPetById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := payload.NewResponse(payload.MessageTypeError, "Invalid Method", nil)
		payload.ResponseJSON(w, http.StatusMethodNotAllowed, response)
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	pet := model.Pet{}
	if err := db.GDB.First(&pet, id).Error; err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Pet was not found", nil)
		payload.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := payload.NewResponse(payload.MessageTypeSuccess, "Pet found", pet)
	payload.ResponseJSON(w, http.StatusOK, response)
}

func GetAllPets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := payload.NewResponse(payload.MessageTypeError, "Method get not permit", nil)
		payload.ResponseJSON(w, http.StatusMethodNotAllowed, response)
		return
	}

	var pets []model.Pet
	if err := db.GDB.Find(&pets).Error; err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Pets not found", nil)
		payload.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	if len(pets) == 0 {
		response := payload.NewResponse(payload.MessageTypeSuccess, "Pets List empty", nil)
		payload.ResponseJSON(w, http.StatusNoContent, response)
		return
	}

	response := payload.NewResponse(payload.MessageTypeSuccess, "Pets found", pets)
	payload.ResponseJSON(w, http.StatusNoContent, response)
}

func SavePet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := payload.NewResponse(payload.MessageTypeError, "Method post not permit", nil)
		payload.ResponseJSON(w, http.StatusMethodNotAllowed, response)
		return
	}

	pet := model.Pet{}
	if err := json.NewDecoder(r.Body).Decode(&pet); err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Bad request: invalid JSON data", nil)
		payload.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	if result := db.GDB.Create(&pet); result.Error != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Internal Server Error", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := payload.NewResponse(payload.MessageTypeSuccess, "Pet created successfusly", nil)
	payload.ResponseJSON(w, http.StatusCreated, response)
}

func UpdatePet(w http.ResponseWriter, r *http.Request) {
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
	pet := model.Pet{}
	if err := db.GDB.First(&pet, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := payload.NewResponse(payload.MessageTypeError, "Pet not found", nil)
			payload.ResponseJSON(w, http.StatusNotFound, response)
			log.Printf("pet not found: %v", err)
			return
		}

		response := payload.NewResponse(payload.MessageTypeError, "Database error", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		log.Printf("database error: %v", err)
		return
	}

	db.GDB.Save(&pet)
	response := payload.NewResponse(payload.MessageTypeSuccess, "Pet updated successfull", pet)
	payload.ResponseJSON(w, http.StatusOK, response)
}

func DeletePet(w http.ResponseWriter, r *http.Request) {
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
	pet := model.Pet{}
	if err := db.GDB.First(&pet, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := payload.NewResponse(payload.MessageTypeError, "Pet not found", nil)
			payload.ResponseJSON(w, http.StatusNotFound, response)
			log.Printf("pet not found: %v", err)
			return
		}

		response := payload.NewResponse(payload.MessageTypeError, "Database error", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		log.Printf("database error: %v", err)
		return
	}

	db.GDB.Delete(&pet)
	response := payload.NewResponse(payload.MessageTypeSuccess, "Pet deleted successfull", nil)
	payload.ResponseJSON(w, http.StatusOK, response)
}
