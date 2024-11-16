package middelware

import (
	"log"
	"net/http"

	"github.com/IsraelTeo/api-paw/auth"
	"github.com/IsraelTeo/api-paw/payload"
)

// Log es un middleware que registra cada solicitud HTTP con su método y URL.
// @Summary Registra la solicitud HTTP entrante
// @Description Registra el método y la URL de cada solicitud entrante
// @Tags Middleware
// @Accept json
// @Produce json
// @Success 200 {string} string "Solicitud registrada"
// @Router / [get]
func Log(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %q, Method: %q", r.URL.Path, r.Method)
		f(w, r)
	}
}

// ValidateJWT es un middleware que valida el token JWT en las solicitudes.
// @Summary Valida el token JWT en la solicitud
// @Description Valida el token JWT en el encabezado de la solicitud y procede si es válido. Devuelve un error 401 si no lo es.
// @Tags Middleware
// @Accept json
// @Produce json
// @Failure 401 {object} payload.Response "Token inválido"
// @Success 200 {string} string "Solicitud permitida"
// @Router / [get]
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

// ValidateJWTAdmin es un middleware que valida el token JWT y asegura que el usuario sea un administrador.
// @Summary Valida el token JWT y verifica si el usuario es administrador
// @Description Valida el token JWT en la solicitud y verifica el rol del usuario. Devuelve un error 403 si no es administrador, o 401 si el token es inválido.
// @Tags Middleware
// @Accept json
// @Produce json
// @Failure 401 {object} payload.Response "Token inválido"
// @Failure 403 {object} payload.Response "Usuario no administrador"
// @Success 200 {string} string "Solicitud permitida"
// @Router /admin [get]
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
