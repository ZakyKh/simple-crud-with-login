package util

import (
	"encoding/json"
	"net/http"
)

func WriteErrorResponse(w http.ResponseWriter, message string, status int) {
	WriteJSONResponse(w, ErrorResponse{Message: message}, status)
}

func WriteJSONResponse(w http.ResponseWriter, body interface{}, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	bodyBytes, _ := json.Marshal(body)
	w.Write(bodyBytes)
}

type ErrorResponse struct {
	Message string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}