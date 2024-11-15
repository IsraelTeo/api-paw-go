package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/IsraelTeo/api-paw/model"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(user model.User) (string, error) {
	payload := jwt.MapClaims{
		"email":      user.Email,
		"authorized": true,
		"is_admin":   user.IsAdmin,
		"iat":        time.Now().Unix(),
		"exp":        time.Now().Add(time.Hour * 2).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString([]byte(os.Getenv("API_SECRET")))
	if err != nil {
		log.Println("Error signing the token:", err)
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(r *http.Request) (model.User, error) {
	// Obtiene el token desde la solicitud
	token := GetToken(r)
	if token == "" {
		log.Println("No token found in request")
		return model.User{}, fmt.Errorf("no token found in request")
	}

	// Valida el token JWT
	jwtToken, err := jwt.Parse(token, validateMethodAndGetSecret)
	if err != nil {
		log.Printf("Token not valid: %v\n", err)
		return model.User{}, fmt.Errorf("invalid token: %w", err)
	}

	// Verifica si el token tiene el tipo de reclamaciones adecuado
	userData, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		log.Println("Unable to retrieve payload information or token is invalid")
		return model.User{}, fmt.Errorf("invalid token claims")
	}

	// Intenta extraer el campo "email" de las reclamaciones
	_, ok = userData["email"].(string)
	if !ok {
		log.Println("Email field missing or not a string in token claims")
		return model.User{}, fmt.Errorf("email field is missing or invalid in token claims")
	}

	// Crea la respuesta de usuario, ahora con el rol
	response := model.User{
		Email:   userData["email"].(string),
		IsAdmin: userData["is_admin"].(bool),
	}

	return response, nil
}

func GetToken(r *http.Request) string {
	params := r.URL.Query()
	if token := params.Get("token"); token != "" {
		return token
	}

	if tokenString := r.Header.Get("Authorization"); len(strings.Split(tokenString, " ")) == 2 {
		return strings.Split(tokenString, " ")[1]
	}

	return ""
}

func validateMethodAndGetSecret(token *jwt.Token) (any, error) {
	_, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return nil, fmt.Errorf("method not valid")
	}

	return []byte(os.Getenv("API_SECRET")), nil
}
