package payload

import (
	"encoding/json"
	"net/http"
)

const (
	MessageTypeError   = "error"
	MessageTypeSuccess = "success"
	Message            = "message"
)

type response struct {
	MessageType string      `json:"message_type"`
	Message     string      `json:"message"`
	Data        interface{} `json:"data"`
}

func NewResponse(messageType, message string, data interface{}) response {
	return response{
		MessageType: messageType,
		Message:     message,
		Data:        data,
	}
}

func ResponseJSON(w http.ResponseWriter, statusCode int, rep response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(&rep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
