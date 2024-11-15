package middelware

import (
	"log"
	"net/http"

	"github.com/IsraelTeo/api-paw/auth"
	"github.com/IsraelTeo/api-paw/payload"
)

func Log(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request : %q, Method: %q", r.URL.Path, r.Method)
		f(w, r)
	}
}

func ValidateJWT(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := auth.ValidateToken(r)
		if err != nil {
			response := payload.NewResponse(payload.MessageTypeError, "Invalid token.", nil)
			payload.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}
		f(w, r)
	}
}

func ValidateJWTAdmin(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userData, err := auth.ValidateToken(r)
		if err != nil {
			response := payload.NewResponse(payload.MessageTypeError, "Invalid token", nil)
			payload.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}

		if !userData.IsAdmin {
			response := payload.NewResponse(payload.MessageTypeError, "Not admin", nil)
			payload.ResponseJSON(w, http.StatusForbidden, response)
			return
		}

		f(w, r)
	}
}
