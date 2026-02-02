# Quick Start Guide: Your First MBEL Application

Get a production-ready multi-language application running in **15 minutes**. This guide covers project setup, translation files, compilation, and deployment.

---

## 1. Initialize Your Project

```bash
# Create project directory
mkdir -p hello-mbel && cd hello-mbel

# Initialize Go module (using MBEL as a library)
go mod init hello-mbel
go get github.com/makkiattooo/MBEL@latest

# Create project structure
mkdir -p locales cmd dist

# Install MBEL CLI for compilation
go install github.com/makkiattooo/MBEL/cmd/mbel@latest
```

**Project structure after setup:**
```
hello-mbel/
├── go.mod
├── go.sum
├── cmd/
│   └── main.go          # Your application
├── locales/             # Translation files (.mbel)
│   ├── en.mbel
│   ├── de.mbel
│   └── fr.mbel
└── dist/                # Compiled output
```

---

## 2. Write Translation Files

### Create English translations (`locales/en.mbel`)

```mbel
# English translations for Hello MBEL
# Metadata section
@namespace: hello
@lang: en

# 1. Simple strings (no variables)
app_name = "Hello MBEL"
app_version = "1.0.0"

# 2. Strings with variable interpolation
greeting = "Welcome, {name}!"
goodbye = "See you later, {name}! Time was: {time}"

# 3. Pluralization (automatic based on language rules)
#    English: [one] (singular), [other] (plural/zero)
items_count(n) {
    [one]   => "You have 1 item"
    [other] => "You have {n} items"
}

# 4. Multi-choice messages (gender, status, etc.)
profile_updated(gender) {
    [male]   => "He updated his profile"
    [female] => "She updated her profile"
    [other]  => "They updated their profile"
}

# 5. Nested keys (use dot notation in code)
#    Access as: "ui.menu.home", "ui.menu.about", etc.
ui.menu {
    home = "Home"
    about = "About"
    contact = "Contact"
    settings = "Settings"
}

# 6. Formatted values (numbers, dates, currencies)
#    Use within variables: {price:currency}, {date:short}
order_total = "Total: {price} (including tax)"
```

### Add German translations (`locales/de.mbel`)

```mbel
@namespace: hello
@lang: de

app_name = "Hallo MBEL"
app_version = "1.0.0"

greeting = "Willkommen, {name}!"
goodbye = "Auf Wiedersehen, {name}! Zeit war: {time}"

# German has 2 plural forms: [one], [other]
items_count(n) {
    [one]   => "Du hast 1 Element"
    [other] => "Du hast {n} Elemente"
}

profile_updated(gender) {
    [male]   => "Er hat sein Profil aktualisiert"
    [female] => "Sie hat ihr Profil aktualisiert"
    [other]  => "Sie haben ihr Profil aktualisiert"
}

ui.menu {
    home = "Startseite"
    about = "Über uns"
    contact = "Kontakt"
    settings = "Einstellungen"
}

order_total = "Gesamtsumme: {price} (inklusive Steuern)"
```

### Add more languages as needed

Create `locales/fr.mbel`, `locales/es.mbel`, etc., following the same pattern.

---

## 3. Validate and Format Translations

```bash
# Check for syntax errors
mbel lint locales/

# Auto-format and fix common issues
mbel fmt locales/

# Show statistics
mbel stats locales/
```

**Output example:**
```
✓ locales/en.mbel: 15 keys, 2 plurals, no errors
✓ locales/de.mbel: 15 keys, 2 plurals, no errors
✓ locales/fr.mbel: 15 keys, 3 plurals, no errors
```

---

## 4. Compile Translations

The MBEL compiler converts `.mbel` files to optimized JSON:

```bash
# Compile all translations to single JSON file
mbel compile locales/ -o dist/translations.json

# Include source map for debugging (recommended for development)
mbel compile locales/ -o dist/translations.json -sourcemap

# Output files:
#   dist/translations.json          (compiled translations)
#   dist/translations.sourcemap.json (optional debugging info)
```

**Compiled output structure (`translations.json`):**
```json
{
  "en": {
    "app_name": "Hello MBEL",
    "greeting": "Welcome, {name}!",
    "items_count": {
      "__plural": "n",
      "one": "You have 1 item",
      "other": "You have {n} items"
    },
    "ui": {
      "menu": {
        "home": "Home",
        "about": "About"
      }
    }
  },
  "de": { ... },
  "fr": { ... }
}
```

---

## 5. Create Your Go Application

### Basic example (`cmd/main.go`)

```go
package main

import (
	"context"
	"fmt"
	"log"

	mbel "github.com/makkiattooo/MBEL/pkg/mbel"
)

func main() {
	// 1. Initialize Manager (loads all locales from compiled JSON)
	m, err := mbel.NewManager("./dist/translations.json", mbel.Config{
		DefaultLocale: "en",
		FallbackChain: []string{"en"}, // Fallback if translation missing
	})
	if err != nil {
		log.Fatal(err)
	}

	// 2. Simple string lookup
	fmt.Println(m.Get("en", "app_name", nil))     // "Hello MBEL"
	fmt.Println(m.Get("de", "app_name", nil))     // "Hallo MBEL"

	// 3. Interpolate variables
	vars := mbel.Vars{"name": "Alice"}
	greeting := m.Get("en", "greeting", vars)
	fmt.Println(greeting)  // "Welcome, Alice!"

	// 4. Pluralization (automatically resolves plural form)
	for _, count := range []int{0, 1, 2, 5, 10} {
		vars := mbel.Vars{"n": count}
		msg := m.Get("en", "items_count", vars)
		fmt.Printf("n=%d: %s\n", count, msg)
	}
	// Output:
	// n=0: You have 0 items
	// n=1: You have 1 item
	// n=2: You have 2 items
	// n=5: You have 5 items
	// n=10: You have 10 items

	// 5. Multi-choice (gender)
	for _, gender := range []string{"male", "female", "other"} {
		vars := mbel.Vars{"gender": gender}
		msg := m.Get("en", "profile_updated", vars)
		fmt.Printf("%s: %s\n", gender, msg)
	}

	// 6. Nested keys (dot notation)
	homeText := m.Get("en", "ui.menu.home", nil)
	fmt.Println(homeText)  // "Home"

	// 7. Global API (single default locale)
	mbel.Init(m) // Use configured Manager globally
	ctx := context.Background()
	fmt.Println(mbel.T(ctx, "greeting", mbel.Vars{"name": "Bob"}))
	// Or use default locale:
	fmt.Println(mbel.TDefault("greeting", mbel.Vars{"name": "Bob"}))
}
```

### Web application example (`cmd/web/main.go`)

```go
package main

import (
	"context"
	"fmt"
	"net/http"

	mbel "github.com/makkiattooo/MBEL/pkg/mbel"
)

var manager *mbel.Manager

func init() {
	// Initialize on server startup
	var err error
	manager, err = mbel.NewManager("./dist/translations.json", mbel.Config{
		DefaultLocale: "en",
	})
	if err != nil {
		panic(err)
	}
}

// Middleware extracts Accept-Language header and sets in request context
func localeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse Accept-Language header (e.g., "de-DE,de;q=0.9,en;q=0.8")
		lang := r.Header.Get("Accept-Language")
		if lang == "" {
			lang = "en"
		}
		
		// Resolve to supported language
		ctx := context.WithValue(r.Context(), "lang", lang)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Handler example
func helloHandler(w http.ResponseWriter, r *http.Request) {
	lang := r.Context().Value("lang").(string)
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Guest"
	}

	// Get translation for user's language
	msg := manager.Get(lang, "greeting", mbel.Vars{"name": name})
	
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "%s"}`, msg)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	
	// Apply middleware
	http.ListenAndServe(":8080", localeMiddleware(mux))
}
```

---

## 6. Add Tests

Create `cmd/main_test.go`:

```go
package main

import (
	"testing"
	"strings"

	mbel "github.com/makkiattooo/MBEL/pkg/mbel"
)

func TestTranslations(t *testing.T) {
	// Setup
	m, _ := mbel.NewManager("./dist/translations.json", mbel.Config{
		DefaultLocale: "en",
	})

	tests := []struct {
		name      string
		lang      string
		key       string
		vars      mbel.Vars
		want      string
		contains  bool // If true, just check substring
	}{
		{
			name: "simple lookup",
			lang: "en",
			key:  "app_name",
			vars: nil,
			want: "Hello MBEL",
		},
		{
			name:     "interpolation",
			lang:     "en",
			key:      "greeting",
			vars:     mbel.Vars{"name": "Alice"},
			want:     "Welcome, Alice!",
			contains: false,
		},
		{
			name:     "pluralization singular",
			lang:     "en",
			key:      "items_count",
			vars:     mbel.Vars{"n": 1},
			want:     "1 item",
			contains: true,
		},
		{
			name:     "pluralization plural",
			lang:     "en",
			key:      "items_count",
			vars:     mbel.Vars{"n": 5},
			want:     "5 items",
			contains: true,
		},
		{
			name:     "german translation",
			lang:     "de",
			key:      "greeting",
			vars:     mbel.Vars{"name": "Klaus"},
			want:     "Willkommen, Klaus!",
			contains: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := m.Get(tt.lang, tt.key, tt.vars)
			
			if tt.contains {
				if !strings.Contains(result, tt.want) {
					t.Errorf("got %q, want to contain %q", result, tt.want)
				}
			} else {
				if result != tt.want {
					t.Errorf("got %q, want %q", result, tt.want)
				}
			}
		})
	}
}

func BenchmarkTranslationLookup(b *testing.B) {
	m, _ := mbel.NewManager("./dist/translations.json", mbel.Config{
		DefaultLocale: "en",
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get("en", "greeting", mbel.Vars{"name": "Test"})
	}
}
```

Run tests:
```bash
go test -v ./...
go test -bench=. ./...
```

---

## 7. Build and Run

```bash
# Build the binary
go build -o hello-mbel ./cmd

# Run
./hello-mbel

# Expected output:
# Hello MBEL
# Hallo MBEL
# Welcome, Alice!
# n=0: You have 0 items
# ...
```

---

## 8. Production Deployment

### Option A: Embed Compiled JSON

```go
package main

import (
	"embed"
	"encoding/json"

	mbel "github.com/makkiattooo/MBEL/pkg/mbel"
)

//go:embed dist/translations.json
var translationsJSON []byte

func init() {
	// Parse and initialize at compile time
	var data map[string]interface{}
	if err := json.Unmarshal(translationsJSON, &data); err != nil {
		panic(err)
	}
	// Use data to create Manager
}
```

**Build with embedding:**
```bash
go build -o hello-mbel ./cmd
# Binary size increases slightly, but no external files needed
```

### Option B: Ship with Compiled JSON

```dockerfile
FROM golang:1.21 AS builder
WORKDIR /src
COPY . .
RUN go mod download
RUN go build -o /tmp/hello-mbel ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /tmp/hello-mbel .
COPY dist/translations.json ./dist/
EXPOSE 8080
CMD ["./hello-mbel"]
```

Build and run:
```bash
docker build -t hello-mbel .
docker run -p 8080:8080 hello-mbel
```

### Option C: Load from URL

```go
// Load translations from remote server
m, err := mbel.NewManager("https://api.example.com/translations.json", config)
if err != nil {
	log.Fatal(err)
}
```

---

## 9. Common Patterns

### Language negotiation in HTTP

```go
// Extract preferred language from Accept-Language header
func getPreferredLanguage(r *http.Request) string {
	// HTTP middleware can handle this automatically
	// See ARCHITECTURE.md for HTTP Middleware details
	acceptLang := r.Header.Get("Accept-Language")
	// Parse and return best match (e.g., "de-DE" → "de")
	return parseAcceptLanguage(acceptLang)
}
```

### Caching and performance

```go
// Manager automatically caches loaded translations
m, _ := mbel.NewManager("./dist/translations.json", mbel.Config{
	DefaultLocale: "en",
	// Optional: custom cache size (default: 1000 entries)
})

// Repeated lookups are O(1) after first access
for i := 0; i < 1000000; i++ {
	msg := m.Get("en", "app_name", nil) // Sub-microsecond after first call
}
```

### Error handling and fallbacks

```go
// If translation not found, falls back to key name or default locale
msg := m.Get("xx", "nonexistent_key", nil)
// Returns: key name as fallback
// Or if using FallbackChain: falls back to English version

// Check if key exists (before displaying to user)
m.Get("en", "key_name", nil)
// Inspect manager's internal data to check existence first
```

---

## 10. Next Steps

1. **[Full Manual](Manual.md)** — Comprehensive MBEL syntax and features
2. **[ARCHITECTURE.md](ARCHITECTURE.md)** — Deep dive into compiler, runtime, manager
3. **[DEVELOPMENT.md](DEVELOPMENT.md)** — Extend MBEL: add custom features, plugins
4. **[Security Best Practices](SECURITY.md)** — XSS prevention, input validation
5. **[AI Annotations Guide](AI_ANNOTATIONS.md)** — Add AI context for better translations
6. **[Sourcemap Debugging](SOURCEMAP.md)** — Debug and track translation sources

---

## Troubleshooting

| Problem | Solution |
|---------|----------|
| `no such file or directory: locales/` | Run `mkdir -p locales` in project root |
| `Undefined: mbel` or import error | Check `go.mod` has correct path: `github.com/makkiattooo/MBEL` |
| `Syntax error at line 5` | Run `mbel lint locales/` and fix according to error messages |
| `Manager is nil` | Ensure compilation succeeded: `mbel compile locales/ -o dist/translations.json` |
| Translation not found | Check key name (dot notation), verify language code, check fallback chain |
| Performance degradation | Enable caching, use compiled JSON, avoid reloading Manager per-request |

---

## Example Applications

Check `/examples` directory in MBEL repository:
- **`ai_context.mbel`** — Using AI annotations for better context
- **`basic.mbel`** — Simple strings and interpolation
- **`full_app.mbel`** — Complete example with plurals, choices, nested keys
- **`gender.mbel`** — Gender-based message selection
- **`plural.mbel`** — Comprehensive pluralization examples
- **`server/main.go`** — HTTP server with language negotiation

Clone and explore:
```bash
git clone https://github.com/makkiattooo/MBEL.git
cd MBEL/examples/server
go run main.go
# Visit http://localhost:8080
```

