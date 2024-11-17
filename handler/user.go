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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// GetUserById maneja la solicitud HTTP GET para obtener un usuario por su ID.
// @Description Devuelve un usuario especificado por su ID.
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Success 200 {object} payload.Response{MessageType=string, Message=string, Data=model.User} "Usuario encontrado"
// @Failure 400 {object} payload.Response{MessageType=string, Message=string} "Método no permitido"
// @Failure 404 {object} payload.Response{MessageType=string, Message=string} "Usuario no encontrado"
// @Failure 405 {object} payload.Response{MessageType=string, Message=string} "Método no permitido"
// @Router /api/v1/user/{id} [get]
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

// GetAllUsers maneja la solicitud HTTP GET para obtener todos los usuarios.
// @Description Devuelve todos los usuarios registrados en el sistema.
// @Accept json
// @Produce json
// @Success 200 {object} payload.Response{MessageType=string, Message=string, Data=[]model.User} "Usuarios encontrados"
// @Failure 404 {object} payload.Response{MessageType=string, Message=string} "Usuarios no encontrados"
// @Failure 405 {object} payload.Response{MessageType=string, Message=string} "Método no permitido"
// @Failure 204 {object} payload.Response{MessageType=string, Message=string} "Lista de usuarios vacía"
// @Router /api/v1/users [get]
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

// RegisterUser maneja la solicitud HTTP POST para registrar un nuevo usuario.
// @Description Registra un nuevo usuario en el sistema.
// @Accept json
// @Produce json
// @Param user body model.User true "Nuevo usuario"
// @Success 201 {object} payload.Response{MessageType=string, Message=string} "Usuario creado exitosamente"
// @Failure 400 {object} payload.Response{MessageType=string, Message=string} "Solicitud incorrecta o JSON inválido"
// @Failure 409 {object} payload.Response{MessageType=string, Message=string} "Correo electrónico ya en uso"
// @Failure 500 {object} payload.Response{MessageType=string, Message=string} "Error interno del servidor"
// @Router /auth/sign-up [post]
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

	if err := db.GDB.Where("email = ?", user.Email).First(&user).Error; err == nil {
		response := payload.NewResponse(payload.MessageTypeError, "Email already in use", nil)
		payload.ResponseJSON(w, http.StatusConflict, response)
		return
	}

	if len(user.Password) == 0 {
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

// UpdateUser maneja la solicitud HTTP PUT para actualizar la información de un usuario.
// @Description Actualiza la información de un usuario existente por su ID.
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Param user body model.User true "Nuevo usuario con datos actualizados"
// @Success 200 {object} payload.Response{MessageType=string, Message=string, Data=model.User} "Usuario actualizado"
// @Failure 400 {object} payload.Response{MessageType=string, Message=string} "ID inválido o formato incorrecto"
// @Failure 404 {object} payload.Response{MessageType=string, Message=string} "Usuario no encontrado"
// @Failure 500 {object} payload.Response{MessageType=string, Message=string} "Error interno del servidor"
// @Router /api/v1/user/{id} [put]
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

	db.GDB.Save(&user)
	response := payload.NewResponse(payload.MessageTypeSuccess, "User updated successfull", user)
	payload.ResponseJSON(w, http.StatusOK, response)
}

// DeleteUser maneja la solicitud HTTP DELETE para eliminar un usuario por su ID.
// @Description Elimina un usuario especificado por su ID.
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Success 200 {object} payload.Response{MessageType=string, Message=string} "Usuario eliminado"
// @Failure 400 {object} payload.Response{MessageType=string, Message=string} "ID inválido"
// @Failure 404 {object} payload.Response{MessageType=string, Message=string} "Usuario no encontrado"
// @Failure 500 {object} payload.Response{MessageType=string, Message=string} "Error interno del servidor"
// @Router /api/v1/user/{id} [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
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
