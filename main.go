package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sort"
)

// Jedi represents a participant in the rating system
type Jedi struct {
	Name   string
	Rating float64
}

// Duel represents a duel between Jedi/Sith
type Duel struct {
	Title         string   `json:"title"`
	Battler       string   `json:"battler"`
	Versus        string   `json:"versus"`
	Location      string   `json:"location"`
	Winner        string   `json:"winner"`
	IsMulti       bool     `json:"isMulti"`
	MultiDuelists []string `json:"multiDuelists"`
	Youtube       string   `json:"youtube"`
}

// K is the K-factor
const K = 32

// CalculateExpectedScore calculates the expected score for a participant
func CalculateExpectedScore(ratingA, ratingB float64) float64 {
	return 1 / (1 + math.Pow(10, (ratingB-ratingA)/400))
}

// UpdateRatings updates the ratings of two participants after a duel
func UpdateRatings(jedis map[string]*Jedi, duel Duel) {
	battler := jedis[duel.Battler]
	versus := jedis[duel.Versus]

	expectedScoreBattler := CalculateExpectedScore(battler.Rating, versus.Rating)
	expectedScoreVersus := CalculateExpectedScore(versus.Rating, battler.Rating)

	// If it's a single duel
	if !duel.IsMulti {
		if duel.Winner == battler.Name {
			battler.Rating += K * (1 - expectedScoreBattler)
			versus.Rating += K * (0 - expectedScoreVersus)
		} else if duel.Winner == versus.Name {
			battler.Rating += K * (0 - expectedScoreBattler)
			versus.Rating += K * (1 - expectedScoreVersus)
		}
		return
	}

	// Multi duel adjustments
	if duel.Winner == battler.Name {
		if contains(duel.MultiDuelists, battler.Name) {
			battler.Rating += K * (1 - expectedScoreBattler) * 0.75
			versus.Rating += K * (0 - expectedScoreVersus) * 1.25
		} else {
			battler.Rating += K * (1 - expectedScoreBattler) * 1.25
			versus.Rating += K * (0 - expectedScoreVersus) * 0.75
		}
	} else if duel.Winner == versus.Name {
		if contains(duel.MultiDuelists, versus.Name) {
			versus.Rating += K * (1 - expectedScoreVersus) * 0.75
			battler.Rating += K * (0 - expectedScoreBattler) * 1.25
		} else {
			versus.Rating += K * (1 - expectedScoreVersus) * 1.25
			battler.Rating += K * (0 - expectedScoreBattler) * 0.75
		}
	} else { // Draw
		if contains(duel.MultiDuelists, battler.Name) {
			battler.Rating -= K * 0.1
		} else {
			battler.Rating += K * 0.1
		}
		if contains(duel.MultiDuelists, versus.Name) {
			versus.Rating -= K * 0.1
		} else {
			versus.Rating += K * 0.1
		}
	}
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func main() {
	// Load duels from data.json
	duels, err := loadDuels("data.json")
	if err != nil {
		log.Fatalf("Failed to load duels: %v", err)
	}

	// Get unique names from duels
	jediNames := getUniqueNames(duels)

	// Initialize Jedi/Sith map with starting ratings
	jedis := make(map[string]*Jedi)
	for _, name := range jediNames {
		jedis[name] = &Jedi{Name: name, Rating: 1000}
	}

	// Update ratings based on duels
	for _, duel := range duels {
		UpdateRatings(jedis, duel)
	}

	// Convert map to slice for sorting
	jediList := make([]*Jedi, 0, len(jedis))
	for _, jedi := range jedis {
		jediList = append(jediList, jedi)
	}

	// Sort the Jedi/Sith by rating in descending order
	sort.Slice(jediList, func(i, j int) bool {
		return jediList[i].Rating > jediList[j].Rating
	})

	// Print the sorted ratings
	for _, jedi := range jediList {
		fmt.Printf("%s: %.2f\n", jedi.Name, jedi.Rating)
	}
}

// Load duels from a JSON file
func loadDuels(filename string) ([]Duel, error) {
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

// Get unique names from the duels
func getUniqueNames(duels []Duel) []string {
	nameSet := make(map[string]struct{})
	for _, duel := range duels {
		nameSet[duel.Battler] = struct{}{}
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
