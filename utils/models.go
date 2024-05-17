package utils

// Duelist represents a participant in the rating system
// can be a Jedi or Sith
type Duelist struct {
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
	YouTube       string   `json:"youtube"`
}
