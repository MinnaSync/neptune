package httpx

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, data any, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	jsonb, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = w.Write(jsonb)
	return err
}
