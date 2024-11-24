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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetUserById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := payload.NewResponse(payload.MessageTypeError, "Invalid Method", nil)
		payload.ResponseJSON(w, http.StatusMethodNotAllowed, response)
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	user := model.User{}
	if err := db.GDB.First(&user, id).Error; err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "User was not found", nil)
		payload.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := payload.NewResponse(payload.MessageTypeSuccess, "User found", user)
	payload.ResponseJSON(w, http.StatusOK, response)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := payload.NewResponse(payload.MessageTypeError, "Method get not permit", nil)
		payload.ResponseJSON(w, http.StatusMethodNotAllowed, response)
		return
	}

	var users []model.User
	if err := db.GDB.Find(&users).Error; err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Users not found", nil)
		payload.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	if len(users) == 0 {
		response := payload.NewResponse(payload.MessageTypeSuccess, "Users List empty", nil)
		payload.ResponseJSON(w, http.StatusNoContent, response)
		return
	}

	response := payload.NewResponse(payload.MessageTypeSuccess, "Users found", users)
	payload.ResponseJSON(w, http.StatusOK, response)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := payload.NewResponse(payload.MessageTypeError, "Method post not permit", nil)
		payload.ResponseJSON(w, http.StatusMethodNotAllowed, response)
		return
	}

	user := model.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Bad request: invalid JSON data", nil)
		payload.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	if err := service.ValidateEntity(&user); err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Internal server error", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	if exists, err := service.ValidateUniqueField("email", user.Email, &model.User{}); err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Internal server error", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	} else if exists {
		response := payload.NewResponse(payload.MessageTypeError, "Email already in use", nil)
		payload.ResponseJSON(w, http.StatusConflict, response)
		return

	}

	empty := service.IsEmpty(user.Password)
	if empty {
		response := payload.NewResponse(payload.MessageTypeError, "Password cannot be empty", nil)
		payload.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Error hashed password", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	user.Password = string(hashedPassword)
	if result := db.GDB.Create(&user); result.Error != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Internal Server Error", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := payload.NewResponse(payload.MessageTypeSuccess, "User created successfusly", nil)
	payload.ResponseJSON(w, http.StatusCreated, response)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
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

	user := model.User{}
	if err := db.GDB.First(&user, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := payload.NewResponse(payload.MessageTypeError, "User not found", nil)
			payload.ResponseJSON(w, http.StatusNotFound, response)
			log.Printf("user not found: %v", err)
			return
		}

		response := payload.NewResponse(payload.MessageTypeError, "Database error", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		log.Printf("database error: %v", err)
		return
	}

	var input model.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Invalid request body", nil)
		payload.ResponseJSON(w, http.StatusBadRequest, response)
		log.Printf("error decoding request body: %v", err)
		return
	}

	if err := service.ValidateEntity(&user); err != nil {
		log.Printf("validation error: %v", err)
		response := payload.NewResponse(payload.MessageTypeError, "Bad request.", nil)
		payload.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	user.Email = input.Email
	user.Password = input.Password
	if err := db.GDB.Save(&user).Error; err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Error saving user", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		log.Printf("error saving user: %v", err)
		return
	}

	response := payload.NewResponse(payload.MessageTypeSuccess, "User updated successfully", user)
	payload.ResponseJSON(w, http.StatusOK, response)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
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

	user := model.User{}
	if err := db.GDB.First(&user, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := payload.NewResponse(payload.MessageTypeError, "User not found", nil)
			payload.ResponseJSON(w, http.StatusNotFound, response)
			log.Printf("user not found: %v", err)
			return
		}

		response := payload.NewResponse(payload.MessageTypeError, "Database error", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		log.Printf("database error: %v", err)
		return
	}

	db.GDB.Delete(&user)
	response := payload.NewResponse(payload.MessageTypeSuccess, "User deleted successfull", nil)
	payload.ResponseJSON(w, http.StatusOK, response)
}
