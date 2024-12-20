package services

import (
	"math"
)

const (
	Empty    = 0
	King     = 1
	Defender = 2
	Attacker = 3
)

type Position struct {
	X, Y int
}

type Move struct {
	From, To Position
}

type MovesOrder struct {
	Moves []Move `json:"moves"`
}

type Board struct {
	KingPosition *Position `json:"king_pos"`
	Size         int
	Board        [][]int `json:"board"`
	Winner       string  `json:"winner"`
}

func NewBoard(size int) *Board {
	board := &Board{
		Size:  size,
		Board: make([][]int, size),
	}
	for i := range board.Board {
		board.Board[i] = make([]int, size)
	}
	board.setupBoard()
	return board
}

func (b *Board) setupBoard() {
	center := b.Size / 2
	b.Board[center][center] = King
	b.KingPosition = &Position{X: center, Y: center}

	for i := -2; i <= 2; i++ {
		for j := -2; j <= 2; j++ {
			if abs(i)+abs(j) == 1 ||
				(abs(i) == 1 && abs(j) == 1) ||
				(abs(i) == 0 && abs(j) == 2) ||
				(abs(i) == 2 && abs(j) == 0) {
				b.Board[center+i][center+j] = Defender
			}
		}
	}

	for i := 3; i <= 7; i++ {
		b.Board[0][i] = Attacker
	}
	b.Board[1][5] = Attacker

	for i := 3; i <= 7; i++ {
		b.Board[10][i] = Attacker
	}
	b.Board[9][5] = Attacker

	for i := 3; i <= 7; i++ {
		b.Board[i][0] = Attacker
	}
	b.Board[5][1] = Attacker

	for i := 3; i <= 7; i++ {
		b.Board[i][10] = Attacker
	}
	b.Board[5][9] = Attacker
}

func (b *Board) Copy() *Board {
	newBoard := NewBoard(b.Size)
	for i := range b.Board {
		copy(newBoard.Board[i], b.Board[i])
	}
	return newBoard
}

func (b *Board) Evaluate1() float64 {
	if b.IsKingCaptured() {
		return math.Inf(-1)
	}
	if b.IsKingEscaped() {
		return math.Inf(1)
	}

	var score float64
	score = 1250
	kingPos := b.findKing()
	if kingPos != nil {
		distanceToEdge := min(
			kingPos.X,
			kingPos.Y,
			b.Size-1-kingPos.X,
			b.Size-1-kingPos.Y,
		)
		score -= float64(distanceToEdge) * 10
	}

	defenders, attackers := b.countPieces()
	score += float64(defenders)*100 - float64(attackers)*100

	return score
}

func (b *Board) GetMoves(isAttacker bool) []Move {
	moves := make([]Move, 0)
	pieceType := Attacker
	if !isAttacker {
		pieceType = Defender
	}

	for i := 0; i < b.Size; i++ {
		for j := 0; j < b.Size; j++ {
			if b.Board[i][j] == pieceType {
				dirs := []Position{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
				for _, dir := range dirs {
					newX, newY := i+dir.X, j+dir.Y
					for newX >= 0 && newX < b.Size && newY >= 0 && newY < b.Size && b.Board[newX][newY] == Empty {
						moves = append(moves, Move{
							From: Position{i, j},
							To:   Position{newX, newY},
						})
						newX += dir.X
						newY += dir.Y
					}
				}
			}
		}
	}
	return moves
}

func Minimax(board *Board, depth int, alpha, beta float64, isAttacker bool) (float64, *Move) {
	if depth == 0 {
		return board.Evaluate(), nil
	}

	moves := board.GetMoves(isAttacker)
	if len(moves) == 0 {
		return board.Evaluate(), nil
	}

	var bestMove *Move = &moves[0]

	if isAttacker {
		maxEval := math.Inf(-1)
		for _, move := range moves {
			boardCopy := board.Copy()
			boardCopy.MakeMove(move)
			eval, _ := Minimax(boardCopy, depth-1, alpha, beta, false)
			if eval > maxEval {
				maxEval = eval
				bestMove = &move
			}
			alpha = math.Max(alpha, eval)
			if beta <= alpha {
				break
			}
		}
		return maxEval, bestMove
	} else {
		minEval := math.Inf(1)
		for _, move := range moves {
			boardCopy := board.Copy()
			boardCopy.MakeMove(move)
			eval, _ := Minimax(boardCopy, depth-1, alpha, beta, true)
			if eval < minEval {
				minEval = eval
				bestMove = &move
			}
			beta = math.Min(beta, eval)
			if beta <= alpha {
				break
			}
		}
		return minEval, bestMove
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(values ...int) int {
	m := values[0]
	for _, v := range values[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

func (b *Board) findKing() *Position {
	for i := range b.Board {
		for j := range b.Board[i] {
			if b.Board[i][j] == King {
				return &Position{i, j}
			}
		}
	}
	return nil
}

func (b *Board) countPieces() (defenders, attackers int) {
	for i := range b.Board {
		for j := range b.Board[i] {
			switch b.Board[i][j] {
			case Defender:
				defenders++
			case Attacker:
				attackers++
			}
		}
	}
	return
}

func (b *Board) IsKingCaptured() bool {
	for i := 0; i < b.Size; i++ {
		for j := 0; j < b.Size; j++ {
			if b.Board[i][j] == King {
				surroundings := 0
				dirs := []Position{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
				for _, dir := range dirs {
					newX, newY := i+dir.X, j+dir.Y
					for newX >= 0 && newX < b.Size && newY >= 0 && newY < b.Size && b.Board[newX][newY] == Attacker {
						surroundings += 1
					}
				}
				if surroundings >= 3 {
					b.Winner = "b"
					return true
				}
				return false
			}
		}
	}
	return false
}

func (b *Board) IsKingEscaped() bool {
	winningPos := []Position{{0, b.Size - 1}, {b.Size - 1, 0}, {0, 0}, {b.Size - 1, b.Size - 1}}
	for _, pos := range winningPos {
		if b.Board[pos.X][pos.Y] == King {
			b.Winner = "w"
			return true
		}
	}
	return false
}

func (b *Board) MakeMove(move Move) {
	piece := b.Board[move.From.X][move.From.Y]

	b.Board[move.From.X][move.From.Y] = Empty

	if move.To.X != -1 && move.To.Y != -1 {
		b.Board[move.To.X][move.To.Y] = piece
	}

	b.IsKingCaptured()
	b.IsKingEscaped()

	if piece == King {
		b.KingPosition = &Position{
			X: move.To.X,
			Y: move.To.Y,
		}
	}
}

func (b *Board) Evaluate() float64 {
	if b.IsKingCaptured() {
		return math.Inf(-1)
	}
	if b.IsKingEscaped() {
		return math.Inf(1)
	}

	var score float64

	score += b.evaluateControl()

	score += b.evaluateBlockade()

	defenders, attackers := b.countPieces()
	score += float64(defenders)*50 - float64(attackers)*50

	return score
}

func (b *Board) evaluateControl() float64 {
	var score float64
	kingPos := b.findKing()
	if kingPos == nil {
		return 0
	}

	corners := []Position{
		{0, 0}, {0, b.Size - 1},
		{b.Size - 1, 0}, {b.Size - 1, b.Size - 1},
	}

	for _, corner := range corners {
		pathControl := b.evaluatePathToCorner(*kingPos, corner)
		score += pathControl
	}

	for i := 0; i < b.Size; i++ {
		for j := 0; j < b.Size; j++ {
			if b.Board[i][j] == Defender {
				if i == 0 || i == b.Size-1 || j == 0 || j == b.Size-1 {
					score += 15
				}

				if b.isPastEnemyLines(i, j) {
					score += 25
				}
			}
		}
	}

	return score
}

func (b *Board) evaluateBlockade() float64 {
	var score float64
	kingPos := b.findKing()
	if kingPos == nil {
		return 0
	}

	directions := []Position{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	for _, dir := range directions {
		blockadeScore := b.evaluateDirectionalBlockade(*kingPos, dir)
		score -= blockadeScore
	}

	score += b.evaluateGaps(*kingPos)

	return score
}

func (b *Board) evaluatePathToCorner(from, to Position) float64 {
	var score float64
	path := b.getPathCoordinates(from, to)

	for _, pos := range path {
		switch b.Board[pos.X][pos.Y] {
		case Defender:
			score += 10
		case Attacker:
			score -= 15
		case Empty:
			score += 2
		}
	}

	return score
}

func (b *Board) getPathCoordinates(from, to Position) []Position {
	path := make([]Position, 0)
	x, y := from.X, from.Y

	for x != to.X || y != to.Y {
		if x < to.X {
			x++
		} else if x > to.X {
			x--
		}
		if y < to.Y {
			y++
		} else if y > to.Y {
			y--
		}
		path = append(path, Position{x, y})
	}

	return path
}

func (b *Board) isPastEnemyLines(row, col int) bool {
	kingPos := b.findKing()
	if kingPos == nil {
		return false
	}

	minRow := min(row, kingPos.X)
	maxRow := max(row, kingPos.X)
	minCol := min(col, kingPos.Y)
	maxCol := max(col, kingPos.Y)

	foundAttacker := false
	for i := minRow; i <= maxRow; i++ {
		for j := minCol; j <= maxCol; j++ {
			if b.Board[i][j] == Attacker {
				foundAttacker = true
			}
		}
	}

	return foundAttacker
}

func (b *Board) evaluateGaps(kingPos Position) float64 {
	var score float64
	directions := []Position{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	for _, dir := range directions {
		x, y := kingPos.X+dir.X, kingPos.Y+dir.Y
		gapLength := 0

		for x >= 0 && x < b.Size && y >= 0 && y < b.Size {
			if b.Board[x][y] == Empty {
				gapLength++
			} else if b.Board[x][y] == Attacker {
				break
			}
			x += dir.X
			y += dir.Y
		}

		if gapLength > 0 {
			score += float64(gapLength) * 20
		}
	}

	return score
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (b *Board) evaluateDirectionalBlockade(kingPos Position, direction Position) float64 {
	var score float64
	x, y := kingPos.X+direction.X, kingPos.Y+direction.Y

	consecutiveAttackers := 0
	distanceFromKing := 1

	for x >= 0 && x < b.Size && y >= 0 && y < b.Size {
		switch b.Board[x][y] {
		case Attacker:
			consecutiveAttackers++
			score += 10.0 / float64(distanceFromKing)

			if consecutiveAttackers > 1 {
				score += 5
			}

		case Defender:
			consecutiveAttackers = 0
			score -= 15

		case Empty:
			consecutiveAttackers = 0
			score -= 5

		case King:
			return 0
		}

		if b.hasAdjacentAttacker(x, y) {
			score += 8
		}

		distanceFromKing++
		x += direction.X
		y += direction.Y
	}

	if distanceFromKing > 1 {
		score += 10
	}

	return score
}

func (b *Board) hasAdjacentAttacker(row, col int) bool {
	directions := []Position{
		{-1, 0}, {1, 0},
		{0, -1}, {0, 1},
	}

	for _, dir := range directions {
		newRow := row + dir.X
		newCol := col + dir.Y

		if newRow >= 0 && newRow < b.Size &&
			newCol >= 0 && newCol < b.Size &&
			b.Board[newRow][newCol] == Attacker {
			return true
		}
	}

	return false
}
