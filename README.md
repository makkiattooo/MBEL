# MBEL (Modern Billed-English Language)

> **"Stop fighting with JSON. Start coding your localization."**

Traditional i18n (JSON, YAML, .po) is broken. It was built 20 years ago for static strings. In today's era of **dynamic apps** and **AI-driven development**, it creates more problems than it solves:

*   💥 **Merge Conflicts**: Two devs add a key to `en.json`? Enjoy the git hell.
*   🕯️ **Zero Context**: What does `Register` mean? Is it a button? A header? A verb? A noun?
*   🐞 **Runtime Crashes**: One missing comma in your 5000-line JSON and your production build fails.
*   🤖 **AI Hallucinations**: Standard tools give zero guidance to LLMs, leading to terrible automated translations.

**MBEL is different.** It’s a programmable localization DSL that treats translations as code, not just data.

---

## 📚 Documentation / Dokumentacja

| Language | Core Manual | Technical Suite (FAQ/Tips/AI) |
| :--- | :--- | :--- |
| 🇬🇧 **English** | ✅ [Read](docs/en/Manual.md) | ✅ [Full Suite](docs/en/FAQ.md) |
| 🇵🇱 **Polski** | ✅ [Czytaj](docs/pl/Manual.md) | ✅ [Pełny Pakiet](docs/pl/FAQ.md) |
| 🇩🇪 **Deutsch** | ✅ [Handbuch](docs/de/Manual.md) | ✅ [Komplett](docs/de/FAQ.md) |
| 🇫🇷 **Français** | ✅ [Manuel](docs/fr/Manual.md) | ✅ [Complet](docs/fr/FAQ.md) |
| 🇪🇸 **Español** | ✅ [Manual](docs/es/Manual.md) | ✅ [Completo](docs/es/FAQ.md) |
| 🇮🇹 **Italiano** | ✅ [Manuale](docs/it/Manual.md) | ✅ [Completo](docs/it/FAQ.md) |
| 🇷🇺 **Русский** | ✅ [Руководство](docs/ru/Manual.md) | ✅ [Полный пакет](docs/ru/FAQ.md) |
| 🇨🇳 **中文** | ✅ [官方手册](docs/zh/Manual.md) | ✅ [完整版](docs/zh/FAQ.md) |
| 🇯🇵 **日本語** | ✅ [マニュアル](docs/ja/Manual.md) | ✅ [完全版](docs/ja/FAQ.md) |

---

## 🦾 The "Killer" Feature: AI Context

Most tools treat AI translation as a "black box". MBEL makes it **deterministic**. We attach metadata directly to your keys, which our CLI tools feed into LLMs to guarantee perfect translations.

**The MBEL Reality:**
```mbel
# This tells the AI precisely what to do
@AI_Context: "Button in the header for user registration (Verb)"
@AI_Tone: "Action-oriented, short"
@AI_MaxLength: "12"
register_btn = "Sign Up"
```

**The result?** No more "Register" (Noun) when you needed a "Register" (Verb).

---

## 🤯 The Syntax: Logic in Data

Stop building sentences in your Go/JS code. Define the logic once in MBEL, and let the runtime handle the complexity.

### 1. Simple but Powerful Plurals
```mbel
# No more key_one, key_other mess. One key, all rules.
cart_items(n) {
    [one]   => "You have 1 item"
    [other] => "You have {n} items"
}
```

### 2. Contextual Logic (Gender/Roles/Enums)
```mbel
# Control logic by 'gender', interpolate dynamic 'name'
greeting(gender) {
    [male]   => "Welcome back, Mr. {name}"
    [female] => "Welcome back, Ms. {name}"
    [other]  => "Hi {name}!"
}
```

### 3. Smart Ranges
```mbel
# Perfect for XP, Battery, or progress levels
xp_level(points) {
    [0-99]   => "Novice"
    [100-499]=> "Warrior"
    [other]  => "Legend"
}
```

---

## 🛠 Developer UX (Go SDK)

MBEL is built for Go developers by Go developers.

### Installation
```bash
go get github.com/makkiattooo/MBEL
```

### Professional Integration
```go
import "github.com/makkiattooo/MBEL"

func main() {
    // 1. Production-ready setup with Hot-Reload
    mbel.Init("./locales", mbel.Config{
        DefaultLocale: "en",
        Watch:         true, // Hit save in .mbel, see changes in UI instantly
    })
}

func handler(w http.ResponseWriter, r *http.Request) {
    // 2. Context-aware translation
    // Uses Accept-Language middleware automatically
    title := mbel.T(r.Context(), "title")
    
    // 3. Complex logic call
    msg := mbel.T(r.Context(), "greeting", mbel.Vars{
        "gender": "female",
        "name":   "Anna",
    })
}
```

---

## 👔 Enterprise Grade: The "Bebechy"

We don't just talk about Enterprise; we provide the interfaces. 

### 1. Repository Pattern (No Files? No Problem!)
If your architecture demands translations from **PostgreSQL**, **Redis**, or an external **API**, just implement the interface:

```go
type Repository interface {
    // LoadAll returns map[language] -> map[key]value
    LoadAll() (map[string]map[string]interface{}, error)
}

// Then swap it in one line:
repo := NewPostgresRepository(db)
mbel.InitWithRepo(repo, mbel.Config{})
```

### 2. CI-Ready Toolchain
*   `mbel lint`: Integrate into your GitHub Actions. Fail builds on syntax errors or MaxLength violations.
*   **Designed for CI**: CLI exits with non-zero codes on errors.
*   `mbel fmt`: Consistency across the team. Like `gofmt`, but for translations.
*   `mbel compile`: High-speed generation of JSON for your frontend. Perfect for CD pipelines.

---

*Build with ❤️ for developers who value their sanity and production stability.*