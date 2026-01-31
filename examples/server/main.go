package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/yourusername/mbel"
)

func main() {
	// 1. Initialize MBEL with Hot Reload enabled
	// We assume running from project root, so path is examples/server/locales
	err := mbel.Init("./examples/server/locales", mbel.Config{
		DefaultLocale: "en",
		Watch:         true,
	})
	if err != nil {
		fmt.Printf("⚠️  Warning: Failed to load locales: %v\n", err)
		fmt.Println("Did you run from project root? Trying local ./locales path...")
		// Fallback for running inside examples/server
		if err := mbel.Init("./locales", mbel.Config{DefaultLocale: "en", Watch: true}); err != nil {
			panic(err)
		}
	}

	// 2. Setup Server
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHome)

	// 3. Wrap with Middleware
	addr := ":8080"
	fmt.Printf("🚀 MBEL Demo Server running on http://localhost%s\n", addr)
	fmt.Println("   Try: http://localhost:8080?name=Alice&gender=female")
	fmt.Println("   Try: http://localhost:8080?name=Bob&gender=male&lang=pl")

	if err := http.ListenAndServe(addr, mbel.Middleware(mux)); err != nil {
		panic(err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Override lang from query param for easy testing in browser
	if queryLang := r.URL.Query().Get("lang"); queryLang != "" {
		ctx = mbel.WithLocale(ctx, queryLang)
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Stranger"
	}

	gender := r.URL.Query().Get("gender")

	// Simulate dynamic count
	rand.Seed(time.Now().UnixNano())
	files := rand.Intn(15) // 0-14

	// --- RENDER RESPONSE ---

	// 1. Simple key
	title := mbel.T(ctx, "title")

	// 2. Logic block with named parameters (gender -> logic, name -> interpolation)
	// Note: 'greeting' block in .mbel must have argument defined e.g. greeting(gender)
	greeting := mbel.T(ctx, "greeting", mbel.Vars{
		"gender": gender,
		"name":   name,
	})

	// 3. Plural logic
	status := mbel.T(ctx, "files_count", mbel.Vars{"n": files})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprintf(w, `
<html>
<body style="font-family: sans-serif; padding: 2rem;">
    <h1>%s</h1>
    <h2>%s</h2>
    <p>%s</p>
    <hr>
    <p><i>Current Locale: %s</i></p>
    <p><small>Try modifying .mbel files in locales/ folder and refresh!</small></p>
</body>
</html>
`, title, greeting, status, mbel.LocaleFromContext(ctx))
}
