package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/barrymun/jedi-sith-elo/utils"
)

func main() {
	// Load duels from data.json
	duels, err := utils.LoadDuels("data.json")
	if err != nil {
		log.Fatalf("Failed to load duels: %v", err)
	}

	// Get unique names from duels
	jediNames := utils.GetUniqueNames(duels)

	// Initialize Jedi/Sith map with starting ratings
	jedis := make(map[string]*utils.Jedi)
	for _, name := range jediNames {
		jedis[name] = &utils.Jedi{Name: name, Rating: 1000}
	}

	// Update ratings based on duels
	for _, duel := range duels {
		utils.UpdateRatings(jedis, duel)
	}

	// Convert map to slice for sorting
	jediList := make([]*utils.Jedi, 0, len(jedis))
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
