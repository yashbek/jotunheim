package tests

import (
	"math"
	"reflect"
	"testing"

	services "github.com/yashbek/jotunheim/services/engine"
)

func TestBoardSetup(t *testing.T) {
	board := services.NewBoard(11)

	expectedBoard := [][]int{{0, 0, 0, 3, 3, 3, 3, 3, 0, 0, 0}, {0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, {3, 0, 0, 0, 0, 2, 0, 0, 0, 0, 3}, {3, 0, 0, 0, 2, 2, 2, 0, 0, 0, 3}, {3, 3, 0, 2, 2, 1, 2, 2, 0, 3, 3}, {3, 0, 0, 0, 2, 2, 2, 0, 0, 0, 3}, {3, 0, 0, 0, 0, 2, 0, 0, 0, 0, 3}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, {0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0}, {0, 0, 0, 3, 3, 3, 3, 3, 0, 0, 0}}

	if !reflect.DeepEqual(board.Board, expectedBoard) {
		t.Error("Boards are not equal")
	}
}

func TestBoardEval(t *testing.T) {
	board := services.NewBoard(11)

	eval, _ := services.Minimax(board, 2, math.Inf(-1), math.Inf(-1), true)
	if eval != 472 {
		t.Error("eval is not expected val not equal")
	}
}
