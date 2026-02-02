# Szybki start: Twoja pierwsza aplikacja MBEL

Uruchom pełnoprawną aplikację wielojęzyczną z MBEL w **15 minut**. Ten przewodnik obejmuje inicjalizację projektu, pliki tłumaczeń, kompilację i wdrażanie.

---

## 1. Zainicjuj swój projekt

```bash
# Utwórz katalog projektu
mkdir -p hello-mbel && cd hello-mbel

# Zainicjuj moduł Go (używając MBEL jako bibliotekę)
go mod init hello-mbel
go get github.com/makkiattooo/MBEL@latest

# Utwórz strukturę projektu
mkdir -p locales cmd dist

# Zainstaluj CLI MBEL do kompilacji
go install github.com/makkiattooo/MBEL/cmd/mbel@latest
```

**Struktura projektu po konfiguracji:**
```
hello-mbel/
├── go.mod
├── go.sum
├── cmd/
│   └── main.go          # Twoja aplikacja
├── locales/             # Pliki tłumaczeń (.mbel)
│   ├── pl.mbel
│   ├── de.mbel
│   └── en.mbel
└── dist/                # Skompilowane tłumaczenia
```

---

## 2. Napisz pliki tłumaczeń

### Utwórz tłumaczenia polskie (`locales/pl.mbel`)

```mbel
# Polskie tłumaczenia dla Hello MBEL
# Sekcja metadanych
@namespace: hello
@lang: pl

# 1. Proste ciągi (bez zmiennych)
app_name = "Cześć MBEL"
app_version = "1.0.0"

# 2. Ciągi ze zmiennymi
greeting = "Witaj, {name}!"
goodbye = "Do widzenia, {name}! Czas: {time}"

# 3. Pluralizacja (automatycznie na podstawie reguł języka)
#    Polski: [one] (1), [few] (2-4), [many] (5+, 0, 21, 25, ...), [other]
items_count(n) {
    [one]  => "Masz 1 element"
    [few]  => "Masz {n} elementy"
    [many] => "Masz {n} elementów"
    [other] => "Masz {n} elementu"
}

# 4. Wiadomości multi-choice (płeć, status, itp.)
profile_updated(gender) {
    [male]   => "Zaktualizował swój profil"
    [female] => "Zaktualizowała swój profil"
    [other]  => "Zaktualizowali swój profil"
}

# 5. Zagnieżdżone klucze (użyj notacji z kropkami w kodzie)
ui.menu {
    home = "Strona główna"
    about = "O nas"
    contact = "Kontakt"
    settings = "Ustawienia"
}

# 6. Sformatowane wartości
order_total = "Razem: {price} (łącznie z podatkiem)"
```

### Dodaj tłumaczenia angielskie (`locales/en.mbel`)

```mbel
@namespace: hello
@lang: en

app_name = "Hello MBEL"
app_version = "1.0.0"

greeting = "Welcome, {name}!"
goodbye = "See you later, {name}! Time: {time}"

# English: [one], [other]
items_count(n) {
    [one]   => "You have 1 item"
    [other] => "You have {n} items"
}

profile_updated(gender) {
    [male]   => "He updated his profile"
    [female] => "She updated her profile"
    [other]  => "They updated their profile"
}

ui.menu {
    home = "Home"
    about = "About"
    contact = "Contact"
    settings = "Settings"
}

order_total = "Total: {price} (including tax)"
```

### Dodaj więcej języków

Utwórz `locales/de.mbel`, `locales/fr.mbel`, itp., postępując według tego samego wzoru.

---

## 3. Sprawdź poprawność i sformatuj tłumaczenia

```bash
# Sprawdź błędy składni
mbel lint locales/

# Automatycznie formatuj i poprawiaj problemy
mbel fmt locales/

# Pokaż statystyki
mbel stats locales/
```

**Przykładowe wyjście:**
```
✓ locales/pl.mbel: 15 keys, 1 plural, no errors
✓ locales/en.mbel: 15 keys, 1 plural, no errors
✓ locales/de.mbel: 15 keys, 1 plural, no errors
```

---

## 4. Skompiluj tłumaczenia

Kompilator MBEL konwertuje pliki `.mbel` na zoptymalizowany JSON:

```bash
# Skompiluj wszystkie tłumaczenia do jednego pliku JSON
mbel compile locales/ -o dist/translations.json

# Dołącz mapę źródeł do debugowania (zalecane w programowaniu)
mbel compile locales/ -o dist/translations.json -sourcemap

# Pliki wyjściowe:
#   dist/translations.json          (skompilowane tłumaczenia)
#   dist/translations.sourcemap.json (opcjonalne info do debugowania)
```

**Struktura skompilowanego wyjścia (`translations.json`):**
```json
{
  "pl": {
    "app_name": "Cześć MBEL",
    "greeting": "Witaj, {name}!",
    "items_count": {
      "__plural": "n",
      "one": "Masz 1 element",
      "few": "Masz {n} elementy",
      "many": "Masz {n} elementów",
      "other": "Masz {n} elementu"
    },
    "ui": {
      "menu": {
        "home": "Strona główna"
      }
    }
  },
  "en": { ... }
}
```

---

## 5. Utwórz swoją aplikację Go

### Podstawowy przykład (`cmd/main.go`)

```go
package main

import (
	"context"
	"fmt"
	"log"

	mbel "github.com/makkiattooo/MBEL/pkg/mbel"
)

func main() {
	// 1. Zainicjuj Manager (ładuje wszystkie locale ze skompilowanego JSON)
	m, err := mbel.NewManager("./dist/translations.json", mbel.Config{
		DefaultLocale: "pl",
		FallbackChain: []string{"pl", "en"}, // Powrót jeśli brakuje tłumaczenia
	})
	if err != nil {
		log.Fatal(err)
	}

	// 2. Proste wyszukiwanie ciągu
	fmt.Println(m.Get("pl", "app_name", nil))     // "Cześć MBEL"
	fmt.Println(m.Get("en", "app_name", nil))     // "Hello MBEL"

	// 3. Interpolacja zmiennych
	vars := mbel.Vars{"name": "Alicja"}
	greeting := m.Get("pl", "greeting", vars)
	fmt.Println(greeting)  // "Witaj, Alicja!"

	// 4. Pluralizacja (automatycznie rozwiązuje formę liczby mnogiej)
	for _, count := range []int{1, 2, 5, 10, 21} {
		vars := mbel.Vars{"n": count}
		msg := m.Get("pl", "items_count", vars)
		fmt.Printf("n=%d: %s\n", count, msg)
	}
	// Wyjście:
	// n=1: Masz 1 element
	// n=2: Masz 2 elementy
	// n=5: Masz 5 elementów
	// ...

	// 5. Multi-choice (płeć)
	for _, gender := range []string{"male", "female", "other"} {
		vars := mbel.Vars{"gender": gender}
		msg := m.Get("pl", "profile_updated", vars)
		fmt.Printf("%s: %s\n", gender, msg)
	}

	// 6. Zagnieżdżone klucze (notacja z kropkami)
	homeText := m.Get("pl", "ui.menu.home", nil)
	fmt.Println(homeText)  // "Strona główna"

	// 7. Globalne API (jedna domyślna locale)
	mbel.Init(m) // Użyj skonfigurowanego Managera globalnie
	ctx := context.Background()
	fmt.Println(mbel.T(ctx, "greeting", mbel.Vars{"name": "Piotr"}))
	// Lub użyj domyślnej locale:
	fmt.Println(mbel.TDefault("greeting", mbel.Vars{"name": "Piotr"}))
}
```

### Aplikacja webowa (`cmd/web/main.go`)

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
	// Zainicjuj przy starcie serwera
	var err error
	manager, err = mbel.NewManager("./dist/translations.json", mbel.Config{
		DefaultLocale: "pl",
	})
	if err != nil {
		panic(err)
	}
}

// Middleware wyodrębnia nagłówek Accept-Language i ustawia w kontekście
func localeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Przeanalizuj nagłówek Accept-Language (np. "pl-PL,pl;q=0.9,en;q=0.8")
		lang := r.Header.Get("Accept-Language")
		if lang == "" {
			lang = "pl"
		}
		
		// Rozwiąż na obsługiwanego języka
		ctx := context.WithValue(r.Context(), "lang", lang)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Przykład handlera
func helloHandler(w http.ResponseWriter, r *http.Request) {
	lang := r.Context().Value("lang").(string)
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Gość"
	}

	// Pobierz tłumaczenie dla języka użytkownika
	msg := manager.Get(lang, "greeting", mbel.Vars{"name": name})
	
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "%s"}`, msg)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	
	// Zastosuj middleware
	http.ListenAndServe(":8080", localeMiddleware(mux))
}
```

---

## 6. Dodaj testy

Utwórz `cmd/main_test.go`:

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
		DefaultLocale: "pl",
	})

	tests := []struct {
		name      string
		lang      string
		key       string
		vars      mbel.Vars
		want      string
		contains  bool // Jeśli true, sprawdź tylko podciąg
	}{
		{
			name: "simple lookup",
			lang: "pl",
			key:  "app_name",
			vars: nil,
			want: "Cześć MBEL",
		},
		{
			name:     "interpolation",
			lang:     "pl",
			key:      "greeting",
			vars:     mbel.Vars{"name": "Alicja"},
			want:     "Witaj, Alicja!",
			contains: false,
		},
		{
			name:     "pluralization singular",
			lang:     "pl",
			key:      "items_count",
			vars:     mbel.Vars{"n": 1},
			want:     "1 element",
			contains: true,
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
```

Uruchom testy:
```bash
go test -v ./...
go test -bench=. ./...
```

---

## 7. Zbuduj i uruchom

```bash
# Zbuduj binarny plik
go build -o hello-mbel ./cmd

# Uruchom
./hello-mbel

# Oczekiwane wyjście:
# Cześć MBEL
# Hello MBEL
# Witaj, Alicja!
# n=1: Masz 1 element
# ...
```

---

## 8. Wdrażanie produkcyjne

### Opcja A: Osadź skompilowany JSON

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
	// Przeanalizuj i zainicjuj w czasie kompilacji
	var data map[string]interface{}
	if err := json.Unmarshal(translationsJSON, &data); err != nil {
		panic(err)
	}
	// Użyj danych do utworzenia Managera
}
```

**Buduj z osadzaniem:**
```bash
go build -o hello-mbel ./cmd
# Rozmiar binarny się nieznacznie zwiększa, ale nie ma potrzeby plików zewnętrznych
```

### Opcja B: Wyślij ze skompilowanym JSON

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

Zbuduj i uruchom:
```bash
docker build -t hello-mbel .
docker run -p 8080:8080 hello-mbel
```

### Opcja C: Załaduj z URL

```go
// Załaduj tłumaczenia ze zdalnego serwera
m, err := mbel.NewManager("https://api.example.com/translations.json", config)
if err != nil {
	log.Fatal(err)
}
```

---

## 9. Częste wzorce

### Negocjacja języka w HTTP

```go
// Wyodrębnij preferowany język z nagłówka Accept-Language
func getPreferredLanguage(r *http.Request) string {
	// Middleware HTTP może obsłużyć to automatycznie
	// Patrz ARCHITECTURE.md dla szczegółów HTTP Middleware
	acceptLang := r.Header.Get("Accept-Language")
	// Przeanalizuj i zwróć najlepsze dopasowanie (np. "pl-PL" → "pl")
	return parseAcceptLanguage(acceptLang)
}
```

### Pamięć podręczna i wydajność

```go
// Manager automatycznie buforuje załadowane tłumaczenia
m, _ := mbel.NewManager("./dist/translations.json", mbel.Config{
	DefaultLocale: "pl",
	// Opcjonalnie: niestandardowy rozmiar pamięci podręcznej (domyślnie: 1000 wpisów)
})

// Powtarzane wyszukiwania są O(1) po pierwszym dostępie
for i := 0; i < 1000000; i++ {
	msg := m.Get("pl", "app_name", nil) // Sub-mikrosekund po pierwszym wywołaniu
}
```

### Obsługa błędów i powroty

```go
// Jeśli tłumaczenie nie znalezione, wraca do nazwy klucza lub domyślnych locale
msg := m.Get("xx", "nonexistent_key", nil)
// Zwraca: nazwa klucza jako powrót
// Lub jeśli używasz FallbackChain: powraca do angielskiej wersji
```

---

## 10. Następne kroki

1. **[Podręcznik](Manual.md)** — Kompleksowa dokumentacja składni MBEL
2. **[ARCHITECTURE.md](ARCHITECTURE.md)** — Głębokie zanurzenie w kompilator, runtime, manager
3. **[DEVELOPMENT.md](DEVELOPMENT.md)** — Rozszerzenie MBEL: dodaj niestandardowe funkcje, wtyczki
4. **[Security Best Practices](SECURITY.md)** — Zapobieganie XSS, walidacja wejścia
5. **[Przewodnik AI Annotations](AI_ANNOTATIONS.md)** — Dodaj kontekst AI dla lepszych tłumaczeń
6. **[Sourcemap Debugging](SOURCEMAP.md)** — Debuguj i śledzij źródła tłumaczeń

---

## Rozwiązywanie problemów

| Problem | Rozwiązanie |
|---------|-------------|
| `no such file or directory: locales/` | Uruchom `mkdir -p locales` w katalogu projektu |
| `Undefined: mbel` lub błąd importu | Sprawdź `go.mod` ma poprawną ścieżkę: `github.com/makkiattooo/MBEL` |
| `Syntax error at line 5` | Uruchom `mbel lint locales/` i napraw zgodnie z komunikatami błędów |
| `Manager is nil` | Upewnij się, że kompilacja się powiodła: `mbel compile locales/ -o dist/translations.json` |
| Tłumaczenie nie znalezione | Sprawdź nazwę klucza (notacja z kropkami), zweryfikuj kod języka, sprawdź łańcuch powrotu |
| Degradacja wydajności | Włącz buforowanie, użyj skompilowanego JSON, unikaj ponownego ładowania Managera dla każdego żądania |

---

## Przykładowe aplikacje

Sprawdź katalog `/examples` w repozytorium MBEL:
- **`ai_context.mbel`** — Używanie adnotacji AI w celu uzyskania lepszego kontekstu
- **`basic.mbel`** — Proste ciągi i interpolacja
- **`full_app.mbel`** — Kompletny przykład z pluralizacją, wyborem, zagnieżdżonymi kluczami
- **`gender.mbel`** — Wybór wiadomości na podstawie płci
- **`plural.mbel`** — Kompleksowe przykłady pluralizacji
- **`server/main.go`** — Serwer HTTP z negocjacją języka

Klonuj i odkrywaj:
```bash
git clone https://github.com/makkiattooo/MBEL.git
cd MBEL/examples/server
go run main.go
# Odwiedź http://localhost:8080
```
