package utils

// Jedi represents a participant in the rating system
type Jedi struct {
	Name   string
	Rating float64
}

// Duel represents a duel between Jedi/Sith
type Duel struct {
	Title         string   `json:"title"`
	Duelist       string   `json:"duelist"`
	Versus        string   `json:"versus"`
	Location      string   `json:"location"`
	Winner        string   `json:"winner"`
	IsMulti       bool     `json:"isMulti"`
	MultiDuelists []string `json:"multiDuelists"`
	Youtube       string   `json:"youtube"`
}
