package middelware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/IsraelTeo/api-paw/payload"
	"github.com/golang-jwt/jwt"
)

func Log(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %q, Method %q", r.URL.Path, r.Method)
		f(w, r)
	}
}

func AuthenticationAdmin(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			response := payload.NewResponse(payload.MessageTypeError, "Token is null", nil)
			payload.ResponseJSON(w, http.StatusBadRequest, response)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("API_SECRET")), nil
		})
		if err != nil || !token.Valid {
			response := payload.NewResponse(payload.MessageTypeError, "Invalid token", nil)
			payload.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			role, roleOk := claims["role"].(string)
			if !roleOk || role != "ADMIN" {
				response := payload.NewResponse(payload.MessageTypeError, "Access denied: ADMIN role required", nil)
				payload.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			}
		} else {
			response := payload.NewResponse(payload.MessageTypeError, "Error to getting claims", nil)
			payload.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}

		f(w, r)
	}
}

func AuthenticationEmployee(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			response := payload.NewResponse(payload.MessageTypeError, "Token is null", nil)
			payload.ResponseJSON(w, http.StatusBadRequest, response)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("API_SECRET")), nil
		})
		if err != nil || !token.Valid {
			response := payload.NewResponse(payload.MessageTypeError, "Invalid token", nil)
			payload.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			role, roleOk := claims["role"].(string)
			if !roleOk || role != "EMPLOYEE	" {
				response := payload.NewResponse(payload.MessageTypeError, "Access denied: EMPLOYEE role required", nil)
				payload.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			}
		} else {
			response := payload.NewResponse(payload.MessageTypeError, "Error to getting claims", nil)
			payload.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}

		f(w, r)
	}
}
