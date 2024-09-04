package matchmaking

import (
	"testing"

	"github.com/yashbek/jotunheim/models"
)

// @TODO: please add tests

func TestAddRecord(t *testing.T) {
	Init()
	/* Expected output
			1  
 		   / \ 
 		  2   4
 		     / \ 
 		    5   3          
	*/
	const expected = "2, 1, 5, 4, 3"
	records := []models.Profile{
		{
			Elo: 100,
			ID: "1",
		},
		{
			Elo: 0,
			ID: "2",
		},
		{
			Elo: 300,
			ID: "3",
		},
		{
			Elo: 500,
			ID: "4",
		},
		{
			Elo: 400,
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

func TestDeleteRecord(t *testing.T) {

}
