package json

import (
	"encoding/json"
	"net/http"
)

func Write(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func Read(w http.ResponseWriter, r *http.Request, data any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Throws an error if an unknown field is passed to the body
	return decoder.Decode(data)
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	Write(w, status, ErrorResponse{Error: msg})
}
