package websockets

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yashbek/jotunheim/db/firebasedb"
	"github.com/yashbek/jotunheim/models"
	services "github.com/yashbek/jotunheim/services/engine"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Connection struct {
	Socket *websocket.Conn
	UserID string
	Mu     sync.Mutex
}

type Server struct {
	Connections map[string]*Connection
	Mu          sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		Connections: make(map[string]*Connection),
	}
}

func (s *Server) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		return
	}
	log.Print("HERE")

	userID := r.URL.Query().Get("email")

	connection := &Connection{
		Socket: conn,
		UserID: userID,
	}

	s.Mu.Lock()
	s.Connections[userID] = connection
	s.Mu.Unlock()

	defer func() {
		s.Mu.Lock()
		delete(s.Connections, userID)
		s.Mu.Unlock()
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var msg struct {
			Type    string      `json:"type"`
			Content interface{} `json:"content"`
		}

		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		switch msg.Type {
		case "color":
			s.handleColorMessage(connection, msg.Content)
		case "chat":
			s.handleChatMessage(connection, msg.Content)
		case "move":
			s.handleMoveMessage(connection, msg.Content)
		case "game_end":
			s.handleGameEndMessage(connection, msg.Content)
		case "bot_move":
			s.handleBotMove(connection, msg.Content)
		default:
			log.Printf("Unknown message type: %s", msg.Type)
		}
	}
}

func (s *Server) handleChatMessage(_ *Connection, content interface{}) {
	log.Println(content)
}

func (s *Server) handleGameEndMessage(_ *Connection, content interface{}) {
	gameEndData, ok := content.(map[string]interface{})
	if !ok {
		log.Printf("Invalid game end data format")
		return
	}

	gameID := gameEndData["id"].(string)
	winner := gameEndData["winner"].(string)

	var gameInfo map[string]interface{}
	if err := firebasedb.FirebaseClient.ReadGame("games", gameID, &gameInfo); err != nil {
		log.Printf("Error reading game: %v", err)
		return
	}

	update := map[string]interface{}{
		"winner":   winner,
		"status":   "completed",
		"ended_at": time.Now().Format(time.RFC3339),
	}

	if err := firebasedb.FirebaseClient.UpdateGame(gameID, update); err != nil {
		log.Printf("Error updating game end state: %v", err)
		return
	}

	notification := map[string]interface{}{
		"type": "game_end",
		"content": map[string]interface{}{
			"winner": winner,
			"id":     gameID,
		},
	}

	s.Mu.RLock()
	if player1, ok := gameInfo["player1"].(string); ok {
		if conn1, exists := s.Connections[player1]; exists {
			conn1.Socket.WriteJSON(notification)
		}
	}

	if player2, ok := gameInfo["player2"].(string); ok && player2 != "bot" {
		if conn2, exists := s.Connections[player2]; exists {
			conn2.Socket.WriteJSON(notification)
		}
	}
	s.Mu.RUnlock()

	log.Printf("Game %s completed. Winner: %s", gameID, winner)
}

func (s *Server) handleColorMessage(conn *Connection, content interface{}) {
	color := ""

	gameInfo := map[string]interface{}{}
	gameID := content.(string)
	firebasedb.FirebaseClient.ReadGame("games", gameID, &gameInfo)
	player1, ok := gameInfo["player1"].(string)
	if !ok {
		log.Printf("Invalid game info format")
		return
	}
	player2, ok := gameInfo["player2"].(string)
	if !ok {
		log.Printf("Invalid game info format")
		return
	}

	switch conn.UserID {
	case player1:
		color = "w"
	case player2:
		color = "b"
	}

	colorJSON := map[string]interface{}{
		"type":  "color",
		"color": color,
	}

	s.Mu.RLock()
	conn.Socket.WriteJSON(colorJSON)
	s.Mu.RUnlock()
}

func (s *Server) handleMoveMessage(_ *Connection, content interface{}) {
	moveData, ok := content.(map[string]interface{})
	if !ok {
		log.Printf("Invalid move data format")
		return
	}

	gameID := moveData["id"].(string)

	game := services.Board{}

	firebasedb.FirebaseClient.ReadGameBoard("games", gameID, &game)

	gameInfo := map[string]interface{}{}
	firebasedb.FirebaseClient.ReadGame("games", gameID, &gameInfo)
	player1 := gameInfo["player1"].(string)
	player2 := gameInfo["player2"].(string)

	fromMap := moveData["from"].(map[string]interface{})
	toMap := moveData["to"].(map[string]interface{})

	fromX := int(fromMap["row"].(float64))
	fromY := int(fromMap["col"].(float64))
	toX := int(toMap["row"].(float64))
	toY := int(toMap["col"].(float64))

	move := services.Move{
		From: services.Position{
			X: fromX,
			Y: fromY,
		},
		To: services.Position{
			X: toX,
			Y: toY,
		},
	}

	game.MakeMove(move)
	var moves services.MovesOrder
	firebasedb.FirebaseClient.ReadGameMoves("games", gameID, &moves)

	if toX != -1 && toY != -1 {
		moves.Moves = append(moves.Moves, move)
	}

	update := map[string]interface{}{
		"board": game,
		"moves": moves,
	}

	firebasedb.FirebaseClient.UpdateGame(gameID, update)

	log.Println(content)

	moveUpdate := map[string]interface{}{
		"type": "move_update",
		"content": map[string]interface{}{
			"from": moveData["from"],
			"to":   moveData["to"],
		},
	}

	s.Mu.RLock()
	if conn1, exists := s.Connections[player1]; exists {
		conn1.Socket.WriteJSON(moveUpdate)
	}
	if conn2, exists := s.Connections[player2]; exists {
		conn2.Socket.WriteJSON(moveUpdate)
	}
	s.Mu.RUnlock()
}

func (s *Server) SendMatchNotification(userID string, match models.GameMatch) error {
	s.Mu.RLock()
	conn, exists := s.Connections[userID]
	s.Mu.RUnlock()

	if !exists {
		return fmt.Errorf("no connection found for user %s", userID)
	}

	notification := map[string]interface{}{
		"type":     "match_found",
		"game_id":  match.GameID,
		"opponent": match.Opponent.Email,
		"color":    match.Color,
	}

	conn.Mu.Lock()
	defer conn.Mu.Unlock()
	return conn.Socket.WriteJSON(notification)
}

func (s *Server) handleBotMove(conn *Connection, content interface{}) {
	moveData, ok := content.(map[string]interface{})
	if !ok {
		log.Printf("Invalid bot move data format")
		return
	}

	gameID := moveData["id"].(string)
	piecesData := moveData["pieces"].([]interface{})

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

	for _, piece := range piecesData {
		pieceArr := piece.([]interface{})
		row := int(pieceArr[0].(float64))
		col := int(pieceArr[1].(float64))
		pieceType := pieceTypeMap[pieceArr[2].(string)]
		board.Board[row][col] = pieceType
	}

	_, bestMove := services.Minimax(board, 3, math.Inf(-1), math.Inf(1), true)
	if bestMove == nil {
		log.Printf("No valid moves found for bot")
		return
	}

	botMoveData := map[string]interface{}{
		"type": "move_update",
		"content": map[string]interface{}{
			"from": map[string]interface{}{
				"row":   bestMove.From.X,
				"col":   bestMove.From.Y,
				"piece": "b", // Bot plays as black
			},
			"to": map[string]interface{}{
				"row": bestMove.To.X,
				"col": bestMove.To.Y,
			},
		},
	}

	moveJSON, err := json.Marshal(botMoveData)
	if err != nil {
		log.Printf("Error marshaling bot move: %v", err)
		return
	}

	conn.Mu.Lock()
	err = conn.Socket.WriteMessage(websocket.TextMessage, moveJSON)
	conn.Mu.Unlock()

	if err != nil {
		log.Printf("Error sending bot move: %v", err)
		return
	}

	var game map[string]interface{}
	if err := firebasedb.FirebaseClient.ReadGame("games", gameID, &game); err != nil {
		log.Printf("Error reading game: %v", err)
		return
	}

	moves, ok := game["moves"].([]services.Move)
	if !ok {
		moves = make([]services.Move, 0)
	}
	moves = append(moves, *bestMove)

	update := map[string]interface{}{
		"moves": moves,
	}
	if err := firebasedb.FirebaseClient.UpdateGame(gameID, update); err != nil {
		log.Printf("Error updating game: %v", err)
	}
}
