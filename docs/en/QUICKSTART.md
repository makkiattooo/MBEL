# Getting Started: Your First MBEL Application

This guide walks you through creating a complete multi-language application with MBEL in 15 minutes.

---

## Step 1: Initialize Project

```bash
mkdir hello-mbel
cd hello-mbel

# Initialize Go module
go mod init hello-mbel
go get github.com/makkiattooo/MBEL

# Create directory structure
mkdir -p locales cmd

# Copy mbel CLI (or install from releases)
go install github.com/makkiattooo/MBEL/cmd/mbel@latest
```

---

## Step 2: Create Base Translations

Create `locales/en.mbel`:

```mbel
# English translations for Hello MBEL

# Metadata
@namespace: hello
@lang: en

# Simple messages
app_name = "Hello MBEL"
greeting = "Welcome, {name}!"
goodbye = "See you later, {name}!"

# Pluralization example
items_count(n) {
    [one]   => "You have 1 item"
    [other] => "You have {n} items"
}

# Gender-based messages
profile_updated(gender) {
    [male]   => "He updated his profile"
    [female] => "She updated her profile"
    [other]  => "They updated their profile"
}

# Nested keys
menu {
    home = "Home"
    about = "About"
    contact = "Contact"
}
```

---

## Step 3: Add Translations

Create `locales/de.mbel`:

```mbel
@namespace: hello
@lang: de

app_name = "Hallo MBEL"
greeting = "Willkommen, {name}!"
goodbye = "Bis später, {name}!"

items_count(n) {
    [one]   => "Du hast 1 Element"
    [other] => "Du hast {n} Elemente"
}

profile_updated(gender) {
    [male]   => "Er hat sein Profil aktualisiert"
    [female] => "Sie hat ihr Profil aktualisiert"
    [other]  => "Sie haben ihr Profil aktualisiert"
}

menu {
    home = "Startseite"
    about = "Über"
    contact = "Kontakt"
}
```

---

## Step 4: Validate & Format

```bash
# Check for syntax errors
mbel lint locales/

# Auto-format all files
mbel fmt locales/
```

---

## Step 5: Create Your Application

Create `cmd/main.go`:

```go
package main

import (
	"fmt"
	"log"

	mbel "hello-mbel"
)

func main() {
	// 1. Initialize manager
	m, err := mbel.NewManager("./locales", mbel.Config{
		DefaultLocale: "en",
	})
	if err != nil {
		log.Fatal(err)
	}

	// 2. Translate simple strings
	fmt.Println(m.Get("en", "app_name", nil))
	fmt.Println(m.Get("de", "app_name", nil))

	// 3. Interpolate variables
	greeting := m.Get("en", "greeting", mbel.Vars{"name": "Alice"})
	fmt.Println(greeting) // "Welcome, Alice!"

	// 4. Pluralize
	for _, count := range []int{0, 1, 2, 5} {
		msg := m.Get("en", "items_count", mbel.Vars{"n": count})
		fmt.Printf("n=%d: %s\n", count, msg)
	}

	// 5. Gender-based translation
	for _, gender := range []string{"male", "female", "other"} {
		msg := m.Get("en", "profile_updated", mbel.Vars{"gender": gender})
		fmt.Printf("%s: %s\n", gender, msg)
	}

	// 6. Use global API
	mbel.Init("./locales", mbel.Config{DefaultLocale: "en"})
	fmt.Println(mbel.T(context.Background(), "greeting", mbel.Vars{"name": "Bob"}))
}
```

---

## Step 6: Compile & Build

```bash
# Compile translations to JSON
mbel compile locales/ -o dist/translations.json -sourcemap

# Build Go binary
go build -o hello-mbel ./cmd

# Run
./hello-mbel
```

---

## Step 7: Test Everything

Create `main_test.go`:

```go
package main

import (
	"testing"
	mbel "hello-mbel"
)

func TestTranslations(t *testing.T) {
	m, _ := mbel.NewManager("./locales", mbel.Config{
		DefaultLocale: "en",
	})

	// Test simple lookup
	greeting := m.Get("en", "greeting", mbel.Vars{"name": "Test"})
	if greeting != "Welcome, Test!" {
		t.Errorf("got %q, want %q", greeting, "Welcome, Test!")
	}

	// Test pluralization
	plural := m.Get("en", "items_count", mbel.Vars{"n": 1})
	if !strings.Contains(plural, "1 item") {
		t.Errorf("got %q", plural)
	}

	// Test German
	de_greeting := m.Get("de", "greeting", mbel.Vars{"name": "Klaus"})
	if de_greeting != "Willkommen, Klaus!" {
		t.Errorf("got %q", de_greeting)
	}
}
```

```bash
go test -v
```

---

## Step 8: Deploy

### Option A: Embed Compiled JSON

```go
//go:embed dist/translations.json
var translationsJSON string

func init() {
	var data map[string]interface{}
	json.Unmarshal([]byte(translationsJSON), &data)
	// Use data for your Manager
}
```

### Option B: Ship with Locales

```dockerfile
FROM golang:1.21
WORKDIR /app
COPY . .
RUN go build -o hello-mbel ./cmd

FROM alpine
COPY --from=0 /app/hello-mbel .
COPY --from=0 /app/locales ./locales
CMD ["./hello-mbel"]
```

---

## Next Steps

1. **Add More Languages** — Copy `locales/en.mbel` to new language files
2. **Add AI Context** — Annotate strings with `# AI_Context: ...`
3. **Enable XSS Protection** — Use `NewRuntimeWithOptions(data, true)` for web
4. **Monitor** — Log metrics with `mbel.GetMetrics()`
5. **Setup CI/CD** — Use GitHub Actions template in `.github/workflows/ci.yml`

---

## Full Example Structure

```
hello-mbel/
├── go.mod
├── go.sum
├── cmd/
│   └── main.go
├── main_test.go
├── locales/
│   ├── en.mbel
│   ├── de.mbel
│   └── fr.mbel
├── dist/
│   ├── translations.json
│   └── translations.sourcemap.json
└── Dockerfile
```

---

## Troubleshooting

### "No such file or directory: locales"
```bash
# Make sure locales directory exists
mkdir -p locales
```

### "Undefined: mbel"
```bash
# Run from project root
go get github.com/makkiattooo/MBEL
go mod tidy
```

### "Syntax error at line 5"
```bash
# Check MBEL syntax
mbel lint locales/
# Fix and re-run
```

---

## Learn More

- **[Full Manual](Manual.md)** — Comprehensive documentation
- **[AI Annotations](AI_ANNOTATIONS.md)** — Add context for AI translation
- **[Sourcemap Debugging](SOURCEMAP.md)** — Track translation sources
- **[Security Best Practices](SECURITY.md)** — XSS prevention & validation

