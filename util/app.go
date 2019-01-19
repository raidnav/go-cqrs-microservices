package util

import (
	"encoding/json"
	"net/http"
)

func responseOK(w http.ResponseWriter, body interface{}) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(body)

	if err != nil {
		panic("Unable to parse json body")
	}
}

func responseERROR(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	body := map[string]string{
		"error": message,
	}

	err := json.NewEncoder(w).Encode(body)

	if err != nil {
		panic("Unable to parse json body")
	}
}
