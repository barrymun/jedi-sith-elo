package utils

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"os"
)

// CalculateExpectedScore calculates the expected score for a participant
func CalculateExpectedScore(ratingA, ratingB float64) float64 {
	return 1 / (1 + math.Pow(10, (ratingB-ratingA)/400))
}

// Helper function to check if a slice contains a string
func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// Get unique names from the duels
func GetUniqueNames(duels []Duel) []string {
	nameSet := make(map[string]struct{})
	for _, duel := range duels {
		nameSet[duel.Duelist] = struct{}{}
		nameSet[duel.Versus] = struct{}{}
		for _, duelist := range duel.MultiDuelists {
			nameSet[duelist] = struct{}{}
		}
	}

	names := make([]string, 0, len(nameSet))
	for name := range nameSet {
		names = append(names, name)
	}
	return names
}

// Load duels from a JSON file
func LoadDuels(filename string) ([]Duel, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var duels []Duel
	err = json.Unmarshal(byteValue, &duels)
	if err != nil {
		return nil, err
	}

	return duels, nil
}

// UpdateRatings updates the ratings of two participants after a duel
func UpdateRatings(jedis map[string]*Jedi, duel Duel) {
	duelist := jedis[duel.Duelist]
	versus := jedis[duel.Versus]

	expectedScoreBattler := CalculateExpectedScore(duelist.Rating, versus.Rating)
	expectedScoreVersus := CalculateExpectedScore(versus.Rating, duelist.Rating)

	// If it's a single duel
	if !duel.IsMulti {
		if duel.Winner == duelist.Name {
			duelist.Rating += K * (1 - expectedScoreBattler)
			versus.Rating += K * (0 - expectedScoreVersus)
		} else if duel.Winner == versus.Name {
			duelist.Rating += K * (0 - expectedScoreBattler)
			versus.Rating += K * (1 - expectedScoreVersus)
		}
		return
	}

	// Multi duel adjustments
	if duel.Winner == duelist.Name {
		if Contains(duel.MultiDuelists, duelist.Name) {
			duelist.Rating += K * (1 - expectedScoreBattler) * 0.75
			versus.Rating += K * (0 - expectedScoreVersus) * 1.25
		} else {
			duelist.Rating += K * (1 - expectedScoreBattler) * 1.25
			versus.Rating += K * (0 - expectedScoreVersus) * 0.75
		}
	} else if duel.Winner == versus.Name {
		if Contains(duel.MultiDuelists, versus.Name) {
			versus.Rating += K * (1 - expectedScoreVersus) * 0.75
			duelist.Rating += K * (0 - expectedScoreBattler) * 1.25
		} else {
			versus.Rating += K * (1 - expectedScoreVersus) * 1.25
			duelist.Rating += K * (0 - expectedScoreBattler) * 0.75
		}
	} else { // Draw
		if Contains(duel.MultiDuelists, duelist.Name) {
			duelist.Rating -= K * 0.1
		} else {
			duelist.Rating += K * 0.1
		}
		if Contains(duel.MultiDuelists, versus.Name) {
			versus.Rating -= K * 0.1
		} else {
			versus.Rating += K * 0.1
		}
	}
}
