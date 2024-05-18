package utils

import (
	"encoding/json"
	"io"
	"math"
	"os"
)

// calculateExpectedScore calculates the expected score for a participant
func CalculateExpectedScore(ratingA, ratingB float64) float64 {
	return 1 / (1 + math.Pow(10, (ratingB-ratingA)/400))
}

// check if a slice contains a string
func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// get unique names from the duels
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

// load duels from a JSON file
func LoadDuels(filename string) ([]Duel, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
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

// updates the ratings of two participants after a duel
func UpdateRatings(duelists map[string]*Duelist, duel Duel) {
	duelist := duelists[duel.Duelist]
	versus := duelists[duel.Versus]

	expectedScoreDuelist := CalculateExpectedScore(duelist.Rating, versus.Rating)
	expectedScoreVersus := CalculateExpectedScore(versus.Rating, duelist.Rating)

	// handle single duel
	// draw does not result in rating change for either participant in a single duel
	if !duel.IsMulti {
		if duel.Winner == duelist.Name {
			duelist.Rating += K * (1 - expectedScoreDuelist)
			versus.Rating += K * (0 - expectedScoreVersus)
		} else if duel.Winner == versus.Name {
			duelist.Rating += K * (0 - expectedScoreDuelist)
			versus.Rating += K * (1 - expectedScoreVersus)
		}
		return
	}

	// handle multi duel
	if duel.Winner == duelist.Name {
		if Contains(duel.MultiDuelists, duelist.Name) {
			duelist.Rating += K * (1 - expectedScoreDuelist) * 0.75
			versus.Rating += K * (0 - expectedScoreVersus) * 1.25
		} else {
			duelist.Rating += K * (1 - expectedScoreDuelist) * 1.25
			versus.Rating += K * (0 - expectedScoreVersus) * 0.75
		}
	} else if duel.Winner == versus.Name {
		if Contains(duel.MultiDuelists, versus.Name) {
			versus.Rating += K * (1 - expectedScoreVersus) * 0.75
			duelist.Rating += K * (0 - expectedScoreDuelist) * 1.25
		} else {
			versus.Rating += K * (1 - expectedScoreVersus) * 1.25
			duelist.Rating += K * (0 - expectedScoreDuelist) * 0.75
		}
	} else {
		// handle draw
		// rating change is 10% of the normal rating change
		// if a single participant is in the multi duel, they get a 10% bonus rating change
		// for a draw, and the other participant gets a 10% penalty
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
