package auth

import (
	"encoding/json"
	"net/http"

	"github.com/IsraelTeo/api-paw/db"
	"github.com/IsraelTeo/api-paw/model"
	"github.com/IsraelTeo/api-paw/payload"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type Cretendials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials Cretendials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Bad request", nil)
		payload.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	user := model.User{}
	if err := db.GDB.Where("email = ?", credentials.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := payload.NewResponse(payload.MessageTypeError, "User not found", nil)
			payload.ResponseJSON(w, http.StatusNotFound, response)
		} else {
			response := payload.NewResponse(payload.MessageTypeError, "Error when querying user", nil)
			payload.ResponseJSON(w, http.StatusInternalServerError, response)
		}
		return
	}

	if err := model.VerifyPassword(user.Password, credentials.Password); err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Invalid password", nil)
		payload.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

	token, err := GenerateToken(user)
	if err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Error generating token", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := payload.NewResponse(payload.MessageTypeSuccess, "Authentication success", token)
	payload.ResponseJSON(w, http.StatusOK, response)
}
