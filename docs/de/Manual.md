# MBEL: Das Vollständige Referenzhandbuch

**Version:** 1.2.0
**Datum:** Januar 2026

---

## 📖 Inhaltsverzeichnis

1.  [Einführung](#1-einführung)
2.  [Die MBEL-Sprache](#2-die-mbel-sprache)
    *   [Dateistruktur](#21-dateistruktur)
    *   [Datentypen](#22-datentypen)
    *   [Interpolation & Variablen](#23-interpolation--variablen)
    *   [Logik & Kontrollfluss](#24-logik--kontrollfluss)
    *   [Pluralisierungsregeln](#25-pluralisierungsregeln)
    *   [KI-Metadaten](#26-ki-metadaten)
3.  [CLI-Toolchain](#3-cli-toolchain)
4.  [Go SDK Integration](#4-go-sdk-integration)

---

## 1. Einführung

Modern Billed-English Language (MBEL) ist ein spezialisiertes Lokalisierungsformat, das entwickelt wurde, um die Lücke zwischen menschlichen Entwicklern und Künstlicher Intelligenz zu schließen.

### Kernphilosophie
1.  **Kontext ist König**: Übersetzungen ohne Kontext sind Ratespiele.
2.  **Logik gehört in die Daten**: Pluralregeln und Geschlechterwechsel werden in der Datei definiert.
3.  **AI-First**: Die Syntax ist für LLMs optimiert.

## 1. Warum MBEL? (Echter Vergleich)

Glauben Sie immer noch, dass JSON ausreicht? Vergleichen wir eine einfache Pluralregel + Interpolation.

#### ❌ Der JSON-Weg (Unordentlich)
```json
{
  "cart_items_one": "Sie haben 1 Artikel in Ihrem Warenkorb.",
  "cart_items_other": "Sie haben {{count}} Artikel in Ihrem Warenkorb.",
  "greeting_male": "Willkommen zurück, Herr {{name}}",
  "greeting_female": "Willkommen zurück, Frau {{name}}"
}
```
*Die Logik ist über mehrere Schlüssel verstreut. Der Go/JS-Code muss entscheiden, welcher Schlüssel abgerufen werden soll.*

#### ✅ Der MBEL-Weg (Sauber)
```mbel
cart_items(n) {
    [one]   => "Sie haben 1 Artikel in Ihrem Warenkorb."
    [other] => "Sie haben {n} Artikel in Ihrem Warenkorb."
}

greeting(gender) {
    [male]   => "Willkommen zurück, Herr {name}"
    [female] => "Willkommen zurück, Frau {name}"
}
```
*Ein Schlüssel, saubere Logik. Die Runtime übernimmt die schwere Arbeit.*

---

## 2. Syntax-Leitfaden

### Basisschlüssel
Einfache Schlüssel-Wert-Paare. Verwenden Sie `"""` für mehrzeilige Zeichenfolgen.

```mbel
title = "Meine Anwendung"
```

### Interpolation vs. Logik-Variablen
Wichtige Unterscheidung:
1. **Steuervariable**: Diejenige in `key(var)`. Sie entscheidet, WELCHER Fall ausgewählt wird.
2. **Interpolationsvariablen**: Diejenigen in `{var}`. Sie werden einfach als Text ausgetauscht.

```mbel
# 'gender' ist die Steuervariable
# '{name}' ist die Interpolationsvariable
greeting(gender) {
    [male]   => "Hallo Herr {name}"
    [female] => "Hallo Frau {name}"
}
```
*Zur Laufzeit:* `mbel.T(ctx, "greeting", mbel.Vars{"gender": "male", "name": "Bob"})`

### KI-Metadaten
Die Metadaten werden im Feld `__ai` des kompilierten Objekts gespeichert. Sie beeinflussen den Text zur Laufzeit nicht, unterstützen aber Übersetzungsagenten massiv.

#### Bereichsübereinstimmung
```mbel
battery(percent) {
    [0]       => "Leer"
    [1-19]    => "Niedrig"
    [100]     => "Voll"
}
```

### 2.5 Pluralisierungsregeln

MBEL verwendet CLDR-Regeln.

**Beispiel (Deutsch):**
```mbel
items_de(n) {
    [one]   => "1 Element"     # n = 1
    [other] => "{n} Elemente"  # n = 0, n > 1
}
```

### 2.6 KI-Metadaten

Annotationen, die mit `@AI_` beginnen.

```mbel
@AI_Tone: "Formell"
submit_btn = "Absenden"
```

---

## 3. CLI-Toolchain

### 3.1 Installation

```bash
go install github.com/makkiattooo/MBEL/cmd/mbel@latest
```

### 3.2 Befehlsreferenz

*   **`init`**: Interaktiver Einrichtungsassistent.
*   **`lint`**: Validiert Syntax.
*   **`compile`**: Kompiliert zu JSON.
*   **`watch`**: Entwicklungsmodus (Hot-Reload).
*   **`stats`**: Zeigt Statistiken an.
*   **`fmt`**: Code-Formatierer.

---

## 4. Go SDK Integration

Das Go SDK ist thread-safe und produktionsbereit.

### 4.1 Architektur

*   **Manager**: Zentraler Einstiegspunkt.
*   **Runtime**: Instanz für eine bestimmte Sprache.
*   **Repository**: Schnittstelle für Datenquellen (Dateien, DB).

### 4.2 Initialisierung

```go
import "github.com/makkiattooo/mbel"

func init() {
    mbel.Init("./locales", mbel.Config{
        DefaultLocale: "de",
        Watch:         true,
    })
}
```

### 4.3 Verwendung (T-Funktion)

`T` (Translate) löst den String basierend auf dem Kontext auf.

```go
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // 1. Einfacher Schlüssel
    title := mbel.T(ctx, "page_title")

    // 2. Mit Variablen
    msg := mbel.T(ctx, "welcome", mbel.Vars{"name": "Hans"})

    // 3. Mit Plural
    items := mbel.T(ctx, "cart_items", 5) 
}
```

### 4.4 Repository Pattern

Sie können Übersetzungen aus einer Datenbank laden, indem Sie `Repository` implementieren.

### 4.5 HTTP Middleware

MBEL analysiert automatisch den `Accept-Language`-Header.

```go
router.Use(mbel.Middleware)
```
