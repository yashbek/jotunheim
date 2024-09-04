package matchmaking

import (
	"testing"

	"github.com/yashbek/jotunheim/models"
)

// @TODO: please add tests

func TestAddRecord(t *testing.T) {
	/* 
	- No Rotation Expected output
		1  
	   / \ 
	  2   3
	     / \ 
		5   4
	- Left Left Rotation Expected output
		1  
	   / \ 
	  2   4
	     / \ 
		5   3     
	- Left Right Rotation Expected output
		1  
	   / \ 
	  2   5
	     / \ 
		4   3   
	- Right Right Rotation Expected output
		1  
	   / \ 
	  2   4
	     / \ 
		3   5   
	- Right Left Rotation Expected output
		1  
	   / \ 
	  2   5
	     / \ 
		3   4   
	*/
	expectedOutputs := map[string][]int {
		"2, 1, 5, 3, 4" : {100, 0, 300, 400, 200},
		"2, 1, 5, 4, 3" : {100, 0, 500, 400, 300},  
		"2, 1, 4, 5, 3" : {100, 0, 500, 300, 400},
		"2, 1, 3, 4, 5" : {100, 0, 300, 400, 500},
		"2, 1, 3, 5, 4" : {100, 0, 300, 500, 400},
	}
	for expected, inputs := range expectedOutputs {
		Init()
		records := []models.Profile{
			{
				Elo: inputs[0],
				ID: "1",
			},
			{
				Elo: inputs[1],
				ID: "2",
			},
			{
				Elo: inputs[2],
				ID: "3",
			},
			{
				Elo: inputs[3],
				ID: "4",
			},
			{
				Elo: inputs[4],
				ID: "5",
			},
		}
		for _, record := range records {
			AddRecord(record)
		}
		if exp := SprintInOrder("1"); exp != expected {
			t.Error("expected", exp)
		}
	}
}

func TestDeleteRecord(t *testing.T) {

}
