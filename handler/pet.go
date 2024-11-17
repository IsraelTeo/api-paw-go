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

// GetPetById maneja la solicitud HTTP GET para obtener una mascota por su ID.
// @Description Devuelve una mascota especificada por su ID.
// @Accept json
// @Produce json
// @Param id path int true "ID de la mascota"
// @Success 200 {object} payload.Response{MessageType=string, Message=string, Data=model.Pet} "Mascota encontrada"
// @Failure 400 {object} payload.Response{MessageType=string, Message=string} "Método no permitido"
// @Failure 404 {object} payload.Response{MessageType=string, Message=string} "Mascota no encontrada"
// @Failure 405 {object} payload.Response{MessageType=string, Message=string} "Método no permitido"
// @Router /api/v1/pet/{id} [get]
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

// GetAllPets maneja la solicitud HTTP GET para obtener todas las mascotas registradas.
// @Description Devuelve una lista con todas las mascotas en el sistema.
// @Accept json
// @Produce json
// @Success 200 {object} payload.Response{MessageType=string, Message=string, Data=[]model.Pet} "Mascotas encontradas"
// @Failure 404 {object} payload.Response{MessageType=string, Message=string} "Mascotas no encontradas"
// @Failure 405 {object} payload.Response{MessageType=string, Message=string} "Método no permitido"
// @Failure 204 {object} payload.Response{MessageType=string, Message=string} "Lista de mascotas vacía"
// @Router /api/v1/pets [get]
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

// SavePet maneja la solicitud HTTP POST para guardar una nueva mascota.
// @Description Registra una nueva mascota en el sistema.
// @Accept json
// @Produce json
// @Param pet body model.Pet true "Nueva mascota"
// @Success 201 {object} payload.Response{MessageType=string, Message=string} "Mascota creada exitosamente"
// @Failure 400 {object} payload.Response{MessageType=string, Message=string} "Solicitud incorrecta o JSON inválido"
// @Failure 500 {object} payload.Response{MessageType=string, Message=string} "Error interno del servidor"
// @Router /api/v1/pet [post]
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

// UpdatePet maneja la solicitud HTTP PUT para actualizar una mascota existente.
// @Description Actualiza una mascota existente por su ID.
// @Accept json
// @Produce json
// @Param id path int true "ID de la mascota"
// @Param pet body model.Pet true "Mascota actualizada"
// @Success 200 {object} payload.Response{MessageType=string, Message=string, Data=model.Pet} "Mascota actualizada"
// @Failure 400 {object} payload.Response{MessageType=string, Message=string} "ID inválido o formato incorrecto"
// @Failure 404 {object} payload.Response{MessageType=string, Message=string} "Mascota no encontrada"
// @Failure 500 {object} payload.Response{MessageType=string, Message=string} "Error interno del servidor"
// @Router /api/v1/pet/{id} [put]
func UpdatePet(w http.ResponseWriter, r *http.Request) {
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

// DeletePet maneja la solicitud HTTP DELETE para eliminar una mascota por su ID.
// @Description Elimina una mascota especificada por su ID.
// @Accept json
// @Produce json
// @Param id path int true "ID de la mascota"
// @Success 200 {object} payload.Response{MessageType=string, Message=string} "Mascota eliminada"
// @Failure 400 {object} payload.Response{MessageType=string, Message=string} "ID inválido"
// @Failure 404 {object} payload.Response{MessageType=string, Message=string} "Mascota no encontrada"
// @Failure 500 {object} payload.Response{MessageType=string, Message=string} "Error interno del servidor"
// @Router /api/v1/pet/{id} [delete]
func DeletePet(w http.ResponseWriter, r *http.Request) {
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
