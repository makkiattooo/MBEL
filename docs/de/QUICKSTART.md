# Schnelleinstieg: Ihre erste MBEL-Anwendung

Führen Sie eine produktionsreife mehrsprachige Anwendung mit MBEL in **15 Minuten** aus. Dieser Leitfaden behandelt Projekteinrichtung, Übersetzungsdateien, Kompilierung und Bereitstellung.

---

## 1. Initialisieren Sie Ihr Projekt

```bash
# Erstellen Sie das Projektverzeichnis
mkdir -p hello-mbel && cd hello-mbel

# Initialisieren Sie das Go-Modul (MBEL als Bibliothek verwenden)
go mod init hello-mbel
go get github.com/makkiattooo/MBEL@latest

# Erstellen Sie die Projektstruktur
mkdir -p locales cmd dist

# Installieren Sie MBEL CLI zum Kompilieren
go install github.com/makkiattooo/MBEL/cmd/mbel@latest
```

**Projektstruktur nach Setup:**
```
hello-mbel/
├── go.mod
├── go.sum
├── cmd/
│   └── main.go          # Ihre Anwendung
├── locales/             # Übersetzungsdateien (.mbel)
│   ├── de.mbel
│   ├── en.mbel
│   └── fr.mbel
└── dist/                # Kompilierte Übersetzungen
```

---

## 2. Schreiben Sie Übersetzungsdateien

### Erstellen Sie deutsche Übersetzungen (`locales/de.mbel`)

```mbel
# Deutsche Übersetzungen für Hello MBEL
# Metadaten-Abschnitt
@namespace: hello
@lang: de

# 1. Einfache Strings (ohne Variablen)
app_name = "Hallo MBEL"
app_version = "1.0.0"

# 2. Strings mit Variableninterpolation
greeting = "Willkommen, {name}!"
goodbye = "Auf Wiedersehen, {name}! Zeit: {time}"

# 3. Pluralisierung (automatisch basierend auf Sprachregeln)
#    Deutsch: [one] (1), [other] (alles andere)
items_count(n) {
    [one]   => "Du hast 1 Element"
    [other] => "Du hast {n} Elemente"
}

# 4. Mehrzweck-Nachrichten (Geschlecht, Status, usw.)
profile_updated(gender) {
    [male]   => "Er hat sein Profil aktualisiert"
    [female] => "Sie hat ihr Profil aktualisiert"
    [other]  => "Sie haben ihr Profil aktualisiert"
}

# 5. Verschachtelte Schlüssel (Punktnotation im Code verwenden)
ui.menu {
    home = "Startseite"
    about = "Über uns"
    contact = "Kontakt"
    settings = "Einstellungen"
}

# 6. Formatierte Werte
order_total = "Gesamtsumme: {price} (inklusive Steuern)"
```

### Englische Übersetzungen hinzufügen (`locales/en.mbel`)

```mbel
@namespace: hello
@lang: en

app_name = "Hello MBEL"
app_version = "1.0.0"

greeting = "Welcome, {name}!"
goodbye = "See you later, {name}! Time: {time}"

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

### Weitere Sprachen hinzufügen

Erstellen Sie `locales/fr.mbel`, `locales/es.mbel`, usw., und folgen Sie dem gleichen Muster.

---

## 3. Überprüfen und formatieren Sie Übersetzungen

```bash
# Syntaxfehler prüfen
mbel lint locales/

# Automatisch formatieren und Probleme beheben
mbel fmt locales/

# Statistiken anzeigen
mbel stats locales/
```

**Beispielausgabe:**
```
✓ locales/de.mbel: 15 keys, 1 plural, no errors
✓ locales/en.mbel: 15 keys, 1 plural, no errors
✓ locales/fr.mbel: 15 keys, 2 plurals, no errors
```

---

## 4. Kompilieren Sie Übersetzungen

Der MBEL-Compiler konvertiert `.mbel`-Dateien in optimiertes JSON:

```bash
# Kompilieren Sie alle Übersetzungen in eine einzelne JSON-Datei
mbel compile locales/ -o dist/translations.json

# Quellenmap zum Debuggen einschließen (für die Entwicklung empfohlen)
mbel compile locales/ -o dist/translations.json -sourcemap

# Ausgabedateien:
#   dist/translations.json          (Übersetzungen kompiliert)
#   dist/translations.sourcemap.json (Optionale Debugginginformationen)
```

**Kompilierte Ausgabestruktur (`translations.json`):**
```json
{
  "de": {
    "app_name": "Hallo MBEL",
    "greeting": "Willkommen, {name}!",
    "items_count": {
      "__plural": "n",
      "one": "Du hast 1 Element",
      "other": "Du hast {n} Elemente"
    },
    "ui": {
      "menu": {
        "home": "Startseite"
      }
    }
  },
  "en": { ... }
}
```

---

## 5. Erstellen Sie Ihre Go-Anwendung

### Einfaches Beispiel (`cmd/main.go`)

```go
package main

import (
	"context"
	"fmt"
	"log"

	mbel "github.com/makkiattooo/MBEL/pkg/mbel"
)

func main() {
	// 1. Manager initialisieren (lädt alle Locales aus kompiliertem JSON)
	m, err := mbel.NewManager("./dist/translations.json", mbel.Config{
		DefaultLocale: "de",
		FallbackChain: []string{"de", "en"}, // Fallback wenn Übersetzung fehlt
	})
	if err != nil {
		log.Fatal(err)
	}

	// 2. Einfaches String-Nachschlagen
	fmt.Println(m.Get("de", "app_name", nil))     // "Hallo MBEL"
	fmt.Println(m.Get("en", "app_name", nil))     // "Hello MBEL"

	// 3. Variableninterpolation
	vars := mbel.Vars{"name": "Alice"}
	greeting := m.Get("de", "greeting", vars)
	fmt.Println(greeting)  // "Willkommen, Alice!"

	// 4. Pluralisierung (löst automatisch die Pluralform auf)
	for _, count := range []int{0, 1, 2, 5, 10} {
		vars := mbel.Vars{"n": count}
		msg := m.Get("de", "items_count", vars)
		fmt.Printf("n=%d: %s\n", count, msg)
	}

	// 5. Multi-Choice (Geschlecht)
	for _, gender := range []string{"male", "female", "other"} {
		vars := mbel.Vars{"gender": gender}
		msg := m.Get("de", "profile_updated", vars)
		fmt.Printf("%s: %s\n", gender, msg)
	}

	// 6. Verschachtelte Schlüssel (Punktnotation)
	homeText := m.Get("de", "ui.menu.home", nil)
	fmt.Println(homeText)  // "Startseite"

	// 7. Globales API (einzelne Standardlocale)
	mbel.Init(m) // Verwende Manager global
	ctx := context.Background()
	fmt.Println(mbel.T(ctx, "greeting", mbel.Vars{"name": "Klaus"}))
	fmt.Println(mbel.TDefault("greeting", mbel.Vars{"name": "Klaus"}))
}
```

### Webanwendung (`cmd/web/main.go`)

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
	var err error
	manager, err = mbel.NewManager("./dist/translations.json", mbel.Config{
		DefaultLocale: "de",
	})
	if err != nil {
		panic(err)
	}
}

func localeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lang := r.Header.Get("Accept-Language")
		if lang == "" {
			lang = "de"
		}
		ctx := context.WithValue(r.Context(), "lang", lang)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	lang := r.Context().Value("lang").(string)
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Gast"
	}

	msg := manager.Get(lang, "greeting", mbel.Vars{"name": name})
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "%s"}`, msg)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":8080", localeMiddleware(mux))
}
```

---

## 6. Testen hinzufügen

Erstellen Sie `cmd/main_test.go`:

```go
package main

import (
	"testing"
	"strings"

	mbel "github.com/makkiattooo/MBEL/pkg/mbel"
)

func TestTranslations(t *testing.T) {
	m, _ := mbel.NewManager("./dist/translations.json", mbel.Config{
		DefaultLocale: "de",
	})

	tests := []struct {
		name     string
		lang     string
		key      string
		vars     mbel.Vars
		want     string
		contains bool
	}{
		{
			name: "simple lookup",
			lang: "de",
			key:  "app_name",
			vars: nil,
			want: "Hallo MBEL",
		},
		{
			name:     "interpolation",
			lang:     "de",
			key:      "greeting",
			vars:     mbel.Vars{"name": "Alice"},
			want:     "Willkommen, Alice!",
			contains: false,
		},
		{
			name:     "pluralization",
			lang:     "de",
			key:      "items_count",
			vars:     mbel.Vars{"n": 5},
			want:     "5 Elemente",
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

Tests ausführen:
```bash
go test -v ./...
go test -bench=. ./...
```

---

## 7. Bauen und Ausführen

```bash
# Binary bauen
go build -o hello-mbel ./cmd

# Ausführen
./hello-mbel

# Erwartete Ausgabe:
# Hallo MBEL
# Hello MBEL
# Willkommen, Alice!
# ...
```

---

## 8. Produktionsbereitstellung

### Option A: JSON einbetten

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
	var data map[string]interface{}
	if err := json.Unmarshal(translationsJSON, &data); err != nil {
		panic(err)
	}
}
```

---

## 9. Häufige Muster

### Sprachaushandlung in HTTP

```go
func getPreferredLanguage(r *http.Request) string {
	acceptLang := r.Header.Get("Accept-Language")
	return parseAcceptLanguage(acceptLang)
}
```

### Zwischenspeicherung und Leistung

```go
// Manager speichert Übersetzungen automatisch zwischen
m, _ := mbel.NewManager("./dist/translations.json", mbel.Config{
	DefaultLocale: "de",
})

// Wiederholte Zugriffe sind O(1) nach dem ersten Zugriff
for i := 0; i < 1000000; i++ {
	msg := m.Get("de", "app_name", nil)
}
```

---

## 10. Nächste Schritte

1. **[Handbuch](Manual.md)** — Umfassende MBEL-Dokumentation
2. **[ARCHITECTURE.md](ARCHITECTURE.md)** — Tiefe technische Analyse
3. **[DEVELOPMENT.md](DEVELOPMENT.md)** — MBEL erweitern
4. **[Sicherheitsbest Practices](SECURITY.md)** — XSS-Prävention
5. **[AI Annotations Guide](AI_ANNOTATIONS.md)** — AI-Kontext für bessere Übersetzungen

---

## Fehlerbehebung

| Problem | Lösung |
|---------|--------|
| `no such file or directory: locales/` | `mkdir -p locales` ausführen |
| `Undefined: mbel` | Überprüfen Sie `go.mod` auf richtige Pfade |
| `Syntax error at line 5` | `mbel lint locales/` ausführen |
| Manager is nil | Kompilierung überprüfen: `mbel compile locales/ -o dist/translations.json` |
| Übersetzung nicht gefunden | Schlüsselname, Sprachcode und Fallback-Kette überprüfen |
| Leistungsabfall | Zwischenspeicherung aktivieren, kompiliertes JSON verwenden |

Beispiele können Sie im Verzeichnis `/examples` des MBEL-Repositorys finden!
