package lib

import (
	"encoding/json"
	"net/http"
)

// RespJSON sends a JSON response
func RespJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// DecodeJSON decodes a JSON request
func DecodeJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}
