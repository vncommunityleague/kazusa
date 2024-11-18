package internal

import (
	"encoding/json"
	"net/http"
)

func Json(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		ErrorJson(w, http.StatusInternalServerError, "serialization", err)
	}
}

type ErrorContainer struct {
	Error    string `json:"error"`
	Messages string `json:"description"`
}

func ErrorJson(w http.ResponseWriter, status int, id string, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(ErrorContainer{
		Error:    id,
		Messages: err.Error(),
	}); err != nil {
		panic(err)
	}
}
