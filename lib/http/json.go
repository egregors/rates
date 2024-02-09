package lib

import (
	"encoding/json"
	"net/http"
)

// RespJSON is a utility function to write a JSON response
func RespJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
