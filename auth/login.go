package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/IsraelTeo/api-paw-go/db"
	"github.com/IsraelTeo/api-paw-go/model"
	"github.com/IsraelTeo/api-paw-go/payload"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Bad request", nil)
		payload.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	userData, err := userByEmailAndPassword(credentials.Email, credentials.Password)
	if err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Invalid email or password", nil)
		payload.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

	token, err := GenerateToken(userData)
	if err != nil {
		response := payload.NewResponse(payload.MessageTypeError, "Error generating token", nil)
		payload.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	responseMap := map[string]interface{}{
		"role":  userData.IsAdmin,
		"token": token,
	}

	response := payload.NewResponse(payload.MessageTypeSuccess, "Login successfully", responseMap)
	payload.ResponseJSON(w, http.StatusOK, response)
}

func userByEmailAndPassword(email, password string) (model.User, error) {
	user := model.User{}
	if err := db.GDB.Where("email = ?", email).First(&user).Error; err != nil {
		log.Printf("email invalid: %v", err)
		return user, err
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Password invalid: %v", err)
		return user, err
	}

	return user, nil
}
