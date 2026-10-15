package handler

import (
	"encoding/json"
	"net/http"
)

func writeError(writer http.ResponseWriter, status int, message string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(map[string]string{"error": message})
}
