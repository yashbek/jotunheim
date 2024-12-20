package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/yashbek/jotunheim/db/firebasedb"
	authjt "github.com/yashbek/jotunheim/services/auth"
)

type SignupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupResponse struct {
	Token string `json:"token"`
}

func SignupHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := firebasedb.FirebaseClient.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := authjt.GenerateJWT(req.Email)
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
		"uid":   user.UID,
	})
}
