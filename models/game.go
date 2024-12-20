package models

type GameMatch struct {
	Opponent User
	GameID   string
	Color    string
}
