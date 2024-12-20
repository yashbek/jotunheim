package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/yashbek/jotunheim/db/firebasedb"
	authjt "github.com/yashbek/jotunheim/services/auth"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	UID      string `json:"uid"`
	Username string `json:"username"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := firebasedb.FirebaseClient.SignInUser(req.Email, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := authjt.GenerateJWT(req.Email)
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(LoginResponse{
		Token:    token,
		Username: user.UserInfo.DisplayName,
		UID:      user.UID,
	})
}
