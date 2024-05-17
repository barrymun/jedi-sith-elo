package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/barrymun/jedi-sith-elo/utils"
)

func main() {
	// load duels from data.json
	duels, err := utils.LoadDuels("data.json")
	if err != nil {
		log.Fatalf("Failed to load duels: %v", err)
	}

	duelistNames := utils.GetUniqueNames(duels)

	// initialize Jedi/Sith map with starting ratings
	duelists := make(map[string]*utils.Duelist)
	for _, name := range duelistNames {
		duelists[name] = &utils.Duelist{Name: name, Rating: 1000}
	}

	for _, duel := range duels {
		utils.UpdateRatings(duelists, duel)
	}

	// prepare for sorting
	duelistList := make([]*utils.Duelist, 0, len(duelists))
	for _, duelist := range duelists {
		duelistList = append(duelistList, duelist)
	}

	// sort in descending order
	sort.Slice(duelistList, func(i, j int) bool {
		return duelistList[i].Rating > duelistList[j].Rating
	})

	for _, duelist := range duelistList {
		fmt.Printf("%s: %.2f\n", duelist.Name, duelist.Rating)
	}
}
