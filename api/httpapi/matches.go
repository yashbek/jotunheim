package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/yashbek/jotunheim/db/firebasedb"
	services "github.com/yashbek/jotunheim/services/engine"
)

type MatchHistoryRequest struct {
	Email string `json:"email"`
}

type MatchHistory struct {
	ID          string          `json:"id"`
	PlayerWhite string          `json:"player_white"`
	PlayerBlack string          `json:"player_black"`
	Winner      string          `json:"winner"`
	CreatedAt   time.Time       `json:"created_at"`
	Moves       []services.Move `json:"moves"`
}

func MatchHistoryHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req MatchHistoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := req.Email
	var matches []MatchHistory
	ref := firebasedb.FirebaseClient.Client.NewRef("games")

	var results map[string]interface{}
	if err := ref.Get(r.Context(), &results); err != nil {
		http.Error(w, "Failed to fetch matches", http.StatusInternalServerError)
		return
	}

	for id, data := range results {
		matchData := data.(map[string]interface{})

		player1 := matchData["player1"].(string)
		player2 := matchData["player2"].(string)
		if player1 != email && player2 != email {
			continue
		}

		var moves []services.Move
		if movesData, ok := matchData["moves"]; ok {
			movesJSON, _ := json.Marshal(movesData)
			json.Unmarshal(movesJSON, &moves)
		}

		createdAtStr, ok := matchData["created_at"].(string)
		if !ok {
			log.Printf("Error: created_at is not a string")
		}

		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			log.Printf("Error parsing timestamp: %v", err)
			createdAt = time.Now()
		}

		var board services.Board
		firebasedb.FirebaseClient.ReadGameBoard("games", id, &board)

		winner := board.Winner
		if wplayer, ok := matchData["winner"].(string); ok {
			winner = wplayer
		}

		match := MatchHistory{
			ID:          id,
			PlayerWhite: player1,
			PlayerBlack: player2,
			Moves:       moves,
			CreatedAt:   createdAt,
			Winner:      winner,
		}

		matches = append(matches, match)
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].CreatedAt.After(matches[j].CreatedAt)
	})

	json.NewEncoder(w).Encode(matches)
}
