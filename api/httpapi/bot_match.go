package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/yashbek/jotunheim/db/firebasedb"
	services "github.com/yashbek/jotunheim/services/engine"
)

func CreateBotGameHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		PlayerEmail string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	gameID := uuid.New().String()
	newBoard := services.NewBoard(11)

	game := map[string]interface{}{
		"player1":    request.PlayerEmail,
		"player2":    "bot",
		"status":     "ready",
		"created_at": time.Now().Format(time.RFC3339),
		"board":      newBoard,
		"moves":      []services.Move{},
	}

	if err := firebasedb.FirebaseClient.CreateGame("games", gameID, game); err != nil {
		http.Error(w, "Failed to create game", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"game_id": gameID,
	})
}
