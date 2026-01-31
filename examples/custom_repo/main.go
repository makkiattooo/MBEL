package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yourusername/mbel"
)

// InMemoryRepository simulates loading translations from a database or API
// This demonstrates the power of the Repository Pattern in MBEL.
type InMemoryRepository struct{}

// LoadAll implements mbel.Repository
func (r *InMemoryRepository) LoadAll() (map[string]map[string]interface{}, error) {
	// In a real app, you would SELECT * FROM translations
	return map[string]map[string]interface{}{
		"en": {
			"title":   "Enterprise Dashboard",
			"welcome": "Welcome from InMemory Repository",
		},
		"pl": {
			"title":   "Panel Enterprise",
			"welcome": "Witaj z Repozytorium Pamięci",
		},
	}, nil
}

func main() {
	// 1. Create your custom repository (Database, Redis, API, etc.)
	repo := &InMemoryRepository{}

	// 2. Initialize Manager with the repo
	// Note: Watch is disabled as this repo doesn't support file watching
	m, err := mbel.NewManagerWithRepo(repo, mbel.Config{DefaultLocale: "en"})
	if err != nil {
		log.Fatal(err)
	}

	// 3. Use it!
	ctx := context.Background()

	// Simulate middleware setting locale
	ctxPL := mbel.WithLocale(ctx, "pl")

	fmt.Println("--- Enterprise MBEL Demo ---")
	fmt.Println("[EN] Title:", m.Get("en", "title"))
	// Note: We need to use m.Get directly if we are not setting the global std instance
	// Or we can set it:
	// mbel.InitWithRepo(repo, config...) -> if we exposed such helper.
	// For this demo, let's use the explicit manager instance.

	fmt.Println("[PL] Title:", m.Get("pl", "title"))
}
