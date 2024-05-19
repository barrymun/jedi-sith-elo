package utils

import (
    "math"
    "testing"
)

func TestGetUniqueNames(t *testing.T) {
    duels := []Duel{
        {Duelist: "Alice", Versus: "Bob", MultiDuelists: []string{"Charlie", "David"}},
        {Duelist: "Bob", Versus: "Alice", MultiDuelists: []string{"Charlie", "David"}},
        {Duelist: "Charlie", Versus: "David", MultiDuelists: []string{"Alice", "Bob"}},
    }
    expected := []string{"Alice", "Bob", "Charlie", "David"}
    names := GetUniqueNames(duels)
    if len(names) != len(expected) {
        t.Errorf("got %v; expected %v", names, expected)
    }
    for i, name := range names {
        if name != expected[i] {
            t.Errorf("got %v; expected %v", name, expected[i])
        }
    }
}

func TestUpdateRatings(t *testing.T) {
    duelistNames := []string{"Alice", "Bob", "Charlie", "David"}
    duelists := make(map[string]*Duelist)
	for _, name := range duelistNames {
		duelists[name] = &Duelist{Name: name, Rating: 1000}
	}

    duels := []Duel{
        {Duelist: "Alice", Versus: "Bob", Winner: "Alice"},
        {Duelist: "Charlie", Versus: "David", Winner: "David"},
        {Duelist: "Bob", Versus: "Alice", Winner: "Bob"},
    }

    for _, duel := range duels {
        UpdateRatings(duelists, duel)
    }

    expected := map[string]float64{
        "Alice": 998.53,
        "Bob": 1001.47,
        "Charlie": 984,
        "David": 1016,
    }
    
    for name, rating := range expected {
        if (math.Round(duelists[name].Rating * 100) / 100) != rating {
            t.Errorf("got %v; expected %v", duelists[name].Rating, rating)
        }
    }
}
