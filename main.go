package main

import (
	"fmt"
	"math"
)

// Jedi represents a participant in the rating system
type Jedi struct {
	Name  string
	Rating float64
}

// K is the K-factor
const K = 32

// CalculateExpectedScore calculates the expected score for a participant
func CalculateExpectedScore(ratingA, ratingB float64) float64 {
	return 1 / (1 + math.Pow(10, (ratingB-ratingA)/400))
}

// UpdateRatings updates the ratings of two participants after a battle
func UpdateRatings(winner, loser *Jedi) {
	expectedScoreWinner := CalculateExpectedScore(winner.Rating, loser.Rating)
	expectedScoreLoser := CalculateExpectedScore(loser.Rating, winner.Rating)

	// Winner's new rating
	winner.Rating = winner.Rating + K*(1-expectedScoreWinner)

	// Loser's new rating
	loser.Rating = loser.Rating + K*(0-expectedScoreLoser)
}

func main() {
	jedi1 := &Jedi{Name: "Yoda", Rating: 1000}
	jedi2 := &Jedi{Name: "Darth Vader", Rating: 1000}

	fmt.Printf("%s vs %s\n", jedi1.Name, jedi2.Name)
	fmt.Printf("Initial Ratings: %s: %.2f, %s: %.2f\n", jedi1.Name, jedi1.Rating, jedi2.Name, jedi2.Rating)

	// Simulate a battle where jedi1 wins
	UpdateRatings(jedi1, jedi2)

	fmt.Printf("Updated Ratings: %s: %.2f, %s: %.2f\n", jedi1.Name, jedi1.Rating, jedi2.Name, jedi2.Rating)
}
