package httpx

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`            // machine code
	Message string `json:"message,omitempty"` // optional human message
}

type Envelope map[string]any

func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func OK(w http.ResponseWriter, payload any) {
	JSON(w, http.StatusOK, payload)
}

func Created(w http.ResponseWriter, payload any) {
	JSON(w, http.StatusCreated, payload)
}

func Fail(w http.ResponseWriter, status int, code string, msg ...string) {
	res := ErrorResponse{Error: code}
	if len(msg) > 0 {
		res.Message = msg[0]
	}
	JSON(w, status, res)
}
