# MBEL: The Complete Reference Manual

**Version:** 1.2.0
**Date:** January 2026

---

## 📖 Table of Contents

1.  [Introduction](#1-introduction)
2.  [The MBEL Language](#2-the-mbel-language)
    *   [File Structure](#21-file-structure)
    *   [Data Types](#22-data-types)
    *   [Interpolation & Variables](#23-interpolation--variables)
    *   [Logic & Control Flow](#24-logic--control-flow)
    *   [Pluralization Rules](#25-pluralization-rules)
    *   [AI Metadata](#26-ai-metadata)
3.  [CLI Toolchain](#3-cli-toolchain)
    *   [Installation](#31-installation)
    *   [Commands Reference](#32-commands-reference)
4.  [Go SDK Integration](#4-go-sdk-integration)
    *   [Architecture](#41-architecture)
    *   [Initialization](#42-initialization)
    *   [Runtime Usage](#43-runtime-usage)
    *   [Repository Pattern](#44-repository-pattern)
    *   [HTTP Middleware](#45-http-middleware)

---

## 1. Introduction

Modern Billed-English Language (MBEL) is a specialized localization format designed to bridge the gap between human developers and Artificial Intelligence. Unlike legacy formats like JSON, YAML, or PO, MBEL treats localization as **code**.

### Core Philosophy
1.  **Context is King**: Translations without context are guesses. MBEL enforces context via metadata.
2.  **Logic belongs in the Data**: Pluralization rules, gender switches, and ranges are defined in the file, not in the application code.
3.  **AI-First**: The syntax is optimized for LLM parsing and generation.

## 1. Why MBEL? (The Real Comparison)

Still think JSON is okay? Let's compare a simple plural rule + interpolation.

#### ❌ The JSON Way (Messy)
```json
{
  "cart_items_one": "You have 1 item in your cart.",
  "cart_items_other": "You have {{count}} items in your cart.",
  "greeting_male": "Welcome back, Mr. {{name}}",
  "greeting_female": "Welcome back, Ms. {{name}}"
}
```
*Logic is scattered across multiple keys. Go/JS code must decide which key to fetch.*

#### ✅ The MBEL Way (Clean)
```mbel
cart_items(n) {
    [one]   => "You have 1 item in your cart."
    [other] => "You have {n} items in your cart."
}

greeting(gender) {
    [male]   => "Welcome back, Mr. {name}"
    [female] => "Welcome back, Ms. {name}"
}
```
*One key, clean logic. The runtime handles the heavy lifting.*

---

## 2. Syntax Guide

### Basic Keys
Simple key-value pairs. Use `"""` for multiline strings.

```mbel
title = "My Application"
```

### Interpolation vs Logic Variables
Important distinction: 
1. **Control Variable**: The one in `key(var)`. It decides WHICH case is picked.
2. **Interpolation Variables**: The ones in `{var}`. They are just swapped as text.

```mbel
# 'gender' is the control variable
# '{name}' is the interpolation variable
greeting(gender) {
    [male]   => "Hello Mr. {name}"
    [female] => "Hello Ms. {name}"
}
```
*At runtime:* `mbel.T(ctx, "greeting", mbel.Vars{"gender": "male", "name": "Bob"})`

### AI Metadata
The metadata is stored in the `__ai` field of the compiled object. It does not affect runtime but empowers translation agents.

#### Range Match
Matches numerical ranges (inclusive).

```mbel
battery(percent) {
    [0]       => "Empty"
    [1-19]    => "Low Battery"
    [20-99]   => "Normal"
    [100]     => "Full"
}
```

### 2.5 Pluralization Rules

MBEL uses Common Locale Data Repository (CLDR) rules internally. You use standard categories as matching keys.

**Supported Categories:** `zero`, `one`, `two`, `few`, `many`, `other`.

**Example (English - 2 forms):**
```mbel
items(n) {
    [one]   => "1 item"
    [other] => "{n} items"
}
```

**Example (Polish - 3 forms):**
```mbel
items_pl(n) {
    [one]   => "1 plik"
    [few]   => "{n} pliki"    # 2, 3, 4...
    [many]  => "{n} plików"   # 5-21...
}
```

> **Note:** The variable `n` is purely conventional. You can name the counter variable whatever you like (e.g., `count`, `quantity`).

### 2.6 AI Metadata

Annotations starting with `@AI_` are treated specially. They are attached to the *next* assignment key.

```mbel
@AI_Tone: "Formal"
@AI_MaxLength: "20"
submit_btn = "Proceed"
```

These annotations do not affect the runtime string but are available in the compiled AST for translation tools.

---

## 3. CLI Toolchain

The `mbel` binary provides a suite of tools for the entire lifecycle.

### 3.1 Installation

```bash
go install github.com/makkiattooo/MBEL/cmd/mbel@latest
```

### 3.2 Commands Reference

#### `init`
Interactive wizard to set up a new project.
*   **Usage**: `mbel init`
*   **Action**: Creates `locales/` directory and `en.mbel` sample.

#### `lint`
Validates syntax and AI constraints.
*   **Usage**: `mbel lint ./locales`
*   **Flags**:
    *   `-j <int>`: Number of parallel workers (default: CPU count).
    *   `-v`: Verbose output.
*   **Checks**: Syntax errors, MaxLength violations.

#### `compile`
Compiles all `.mbel` files into a single JSON object. This is useful for front-end consumption or production bundling.
*   **Usage**: `mbel compile ./locales -o ./dist/locales.json`
*   **Flags**:
    *   `-o <file>`: Output file path.
    *   `--pretty`: Pretty-print JSON (default: true).
    *   `--ns`: Auto-derive namespace from folder structure (e.g. `locales/en/auth.mbel` -> `auth`).

#### `watch`
Development mode. Watches for file changes and (optionally) recompiles.
*   **Usage**: `mbel watch ./locales`
*   **Action**: Prints changed files to stdout. Useful when chained with other tools.

#### `stats`
Generates analytics about your localization coverage.
*   **Usage**: `mbel stats ./locales`
*   **Metrics**: Total keys, Logic block complexity, Duplicates.

#### `fmt`
Code formatter. Ensures consistent style (spacing, indentation).
*   **Usage**: `mbel fmt ./locales`

---

## 4. Go SDK Integration

The Go SDK is thread-safe, fast, and production-ready.

### 4.1 Architecture

*   **Manager**: The central entry point. Holds the configuration and the current state of loaded languages.
*   **Runtime**: An instance for a specific language (e.g., "en-US").
*   **Repository**: An interface for loading raw data (Files, DB, API).

### 4.2 Initialization

**Standard (File-based):**
```go
import "github.com/makkiattooo/MBEL"

func init() {
    // Loads all .mbel files from ./locales
    // recursive watching enabled
    err := mbel.Init("./locales", mbel.Config{
        DefaultLocale: "en",
        Watch:         true,
    })
    if err != nil {
        panic(err)
    }
}
```

### 4.3 Runtime Usage

**The `T` Function:**
`T` (Translate) is the main function. It resolves the string based on the Context.

```go
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // 1. Simple Key
    title := mbel.T(ctx, "page_title")

    // 2. With Vars
    msg := mbel.T(ctx, "welcome_message", mbel.Vars{
        "name": "John",
        "day":  "Monday",
    })

    // 3. With Plural
    // Passing a single int is shorthand for mbel.Vars{"n": count}
    // IF the logic block uses 'n' as argument.
    items := mbel.T(ctx, "cart_items", 5) 
    
    fmt.Fprintln(w, title, msg, items)
}
```

### 4.4 Repository Pattern

For enterprise scale, you might want to store translations in a database. Implement the `Repository` interface:

```go
type Repository interface {
    LoadAll() (map[string]map[string]interface{}, error)
}
```

**Using custom repo:**
```go
repo := NewPostgresRepository(dbConn)
manager, _ := mbel.NewManagerWithRepo(repo, mbel.Config{DefaultLocale: "en"})
```

### 4.5 HTTP Middleware

MBEL includes a robust middleware that parses the `Accept-Language` header (RFC 2616) with quality weights (q-factors).

```go
router := mux.NewRouter()
router.Use(mbel.Middleware)

// The middleware automatically injects the best matching language
// into the request Context.
```

---

*Documentation generated automatically by MBEL Team.*
