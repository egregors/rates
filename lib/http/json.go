package lib

import (
	"encoding/json"
	"io"
	"net/http"
)

// RespJSON sends a JSON response
func RespJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// DecodeJSON decodes a JSON request
func DecodeJSON(b io.Reader, v any) error {
	return json.NewDecoder(b).Decode(v)
}
