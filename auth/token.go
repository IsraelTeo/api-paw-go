package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/IsraelTeo/api-paw/model"
	"github.com/golang-jwt/jwt/v4"
)

var JwtKey = []byte(os.Getenv("API_SECRET"))

type Claims struct {
	Email string     `json:"email"`
	Role  model.Role `json:"string"`
	jwt.RegisteredClaims
}

func GenerateToken(email string, role model.Role) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}

func ValidateToken(r *http.Request) error {
	jwtToken := GetToken(r)
	if jwtToken == "" {
		log.Println("No token found in request")
		return fmt.Errorf("no token found")
	}

	token, err := jwt.Parse(jwtToken, func(f *jwt.Token) (interface{}, error) {
		if _, ok := f.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", f.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})

	if err != nil {
		log.Println("Error parsing token:", err)
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		Pretty(claims)
		return nil
	}

	log.Println("Invalid token")
	return fmt.Errorf("invalid token")
}

func Pretty(data interface{}) {
	pretty, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(pretty))
}

func GetToken(r *http.Request) string {
	params := r.URL.Query()
	token := params.Get("token")
	if token != "" {
		return token
	}
	tokenString := r.Header.Get("Authorization")
	if len(strings.Split(tokenString, " ")) == 2 {
		return strings.Split(tokenString, " ")[1]
	}
	return ""
}
