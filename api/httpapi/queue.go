package httpapi

import (
	"encoding/json"
	"net/http"
)

type QueueUpRequest struct {
}

type QueueUpResponse struct {
	Message string `json:"message"`
}

func QueueUpHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var req QueueUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "",
	})
}
