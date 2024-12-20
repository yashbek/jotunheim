package httpapi

import (
	"encoding/json"
	"math"
	"net/http"

	"github.com/yashbek/jotunheim/db/firebasedb"
	services "github.com/yashbek/jotunheim/services/engine"
)

func EvaluateHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Pieces [][]interface{} `json:"pieces"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	board := services.NewBoard(11)
	board.Board = make([][]int, 11)
	for i := range board.Board {
		board.Board[i] = make([]int, 11)
	}

	pieceTypeMap := map[string]int{
		"k": services.King,
		"w": services.Defender,
		"b": services.Attacker,
	}

	for _, piece := range request.Pieces {
		row := int(piece[0].(float64))
		col := int(piece[1].(float64))
		pieceType := pieceTypeMap[piece[2].(string)]
		board.Board[row][col] = pieceType
	}

	// evaluation := board.Evaluate()
	evaluation, _ := services.Minimax(board, 2, math.Inf(-1), math.Inf(-1), true)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"evaluation": evaluation,
	})
}

func GetGameHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	gameID := r.URL.Query().Get("id")
	if gameID == "" {
		http.Error(w, "need game id for query", http.StatusBadRequest)
		return
	}

	var gameData map[string]interface{}
	if err := firebasedb.FirebaseClient.ReadGame("games", gameID, &gameData); err != nil {
		http.Error(w, "Failed to fetch game", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gameData)
}
