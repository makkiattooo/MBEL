# MBEL: Kompletny Podręcznik Użytkownika

**Wersja:** 1.2.0
**Data:** Styczeń 2026

---

## 📖 Spis Treści

1.  [Wprowadzenie](#1-wprowadzenie)
2.  [Język MBEL](#2-język-mbel)
    *   [Struktura Pliku](#21-struktura-pliku)
    *   [Typy Danych](#22-typy-danych)
    *   [Interpolacja i Zmienne](#23-interpolacja-i-zmienne)
    *   [Logika i Sterowanie](#24-logika-i-sterowanie)
    *   [Zasady Liczby Mnogiej](#25-zasady-liczby-mnogiej)
    *   [Metadane AI](#26-metadane-ai)
3.  [Narzędzia CLI](#3-narzędzia-cli)
    *   [Instalacja](#31-instalacja)
    *   [Opis Poleceń](#32-opis-poleceń)
4.  [Integracja Go SDK](#4-integracja-go-sdk)
    *   [Architektura](#41-architektura)
    *   [Inicjalizacja](#42-inicjalizacja)
    *   [Użycie Runtime](#43-użycie-runtime)
    *   [Wzorzec Repository](#44-wzorzec-repository)
    *   [Middleware HTTP](#45-middleware-http)

---

## 1. Wprowadzenie

Modern Billed-English Language (MBEL) to wyspecjalizowany format lokalizacyjny zaprojektowany, aby wypełnić lukę między programistami a Sztuczną Inteligencją. W przeciwieństwie do przestarzałych formatów takich jak JSON czy YAML, MBEL traktuje lokalizację jako **kod**.

### Główna Filozofia
1.  **Kontekst to Król**: Tłumaczenie bez kontekstu to zgadywanie. MBEL wymusza kontekst poprzez metadane.
2.  **Logika należy do Danych**: Zasady liczby mnogiej, płcie i zakresy są definiowane w pliku, a nie w kodzie aplikacji.
3.  **AI-First**: Składnia jest zoptymalizowana pod kątem parsowania i generowania przez LLM.

## 1. Dlaczego MBEL? (Realne Porównanie)

Nadal uważasz, że JSON jest wystarczający? Porównajmy prostą regułę liczby mnogiej + interpolację.

#### ❌ Styl JSON (Bałagan)
```json
{
  "cart_items_one": "Masz 1 produkt w koszyku.",
  "cart_items_few": "Masz {{count}} produkty w koszyku.",
  "cart_items_many": "Masz {{count}} produktów w koszyku.",
  "greeting_male": "Witaj ponownie, Panie {{name}}",
  "greeting_female": "Witaj ponownie, Pani {{name}}"
}
```
*Logika jest rozproszona między wiele kluczy. Kod aplikacji (Go/JS) musi decydować, który klucz pobrać.*

#### ✅ Styl MBEL (Czysty)
```mbel
cart_items(n) {
    [one]   => "Masz 1 produkt w koszyku."
    [few]   => "Masz {n} produkty w koszyku."
    [many]  => "Masz {n} produktów w koszyku."
}

greeting(gender) {
    [male]   => "Witaj ponownie, Panie {name}"
    [female] => "Witaj ponownie, Pani {name}"
}
```
*Jeden klucz, czysta logika. Runtime zajmuje się resztą.*

---

## 2. Przewodnik po Składni

### Podstawowe Klucze
Proste pary klucz-wartość. Teksty wieloliniowe zaczynaj od `"""`.

```mbel
title = "Moja Aplikacja"
```

### Zmienne Sterujące vs Interpolacyjne
Kluczowe rozróżnienie:
1. **Zmienna Sterująca**: Ta w nawiasach `klucz(zmienna)`. Decyduje o tym, KTÓRY przypadek zostanie wybrany.
2. **Zmienne Interpolacyjne**: Te w klamrach `{zmienna}`. Są po prostu wstawiane jako tekst.

```mbel
# 'gender' to zmienna sterująca
# '{name}' to zmienna interpolacyjna
greeting(gender) {
    [male]   => "Witaj Panie {name}"
    [female] => "Witaj Pani {name}"
}
```
*Użycie w Go:* `mbel.T(ctx, "greeting", mbel.Vars{"gender": "male", "name": "Jan"})`

### Metadane AI
Metadane są przechowywane w polu `__ai`. Nie wpływają na tekst w czasie działania aplikacji, ale dają "nadludzkie" możliwości narzędziom do automatycznych tłumaczeń.

#### Dopasowanie Zakresu (Range Match)
```mbel
battery(percent) {
    [0]       => "Rozładowana"
    [1-19]    => "Słaba Bateria"
    [20-99]   => "W Normie"
    [100]     => "Pełna"
}
```

### 2.5 Zasady Liczby Mnogiej

MBEL używa zasad CLDR. Pola: `zero`, `one`, `two`, `few`, `many`, `other`.

**Przykład (Język Polski - Skomplikowany!):**
```mbel
items_pl(n) {
    [one]   => "1 plik"       # n = 1
    [few]   => "{n} pliki"    # n % 10 e {2,3,4} && n % 100 !e {12,13,14}
    [many]  => "{n} plików"   # n != 1 && (n % 10 e {0,1} || n % 10 e {5..9} ...)
    [other] => "{n} pliku"    # Ułamki itp.
}
```
MBEL oblicza to automatycznie! Ty tylko definiujesz przypadki.

### 2.6 Metadane AI

Adnotacje zaczynające się od `@AI_` są dołączane do *następnego* klucza.

```mbel
@AI_Tone: "Formalny"
@AI_MaxLength: "20"
submit_btn = "Zatwierdź"
```

---

## 3. Narzędzia CLI

### 3.1 Instalacja

```bash
go install github.com/makkiattooo/MBEL/cmd/mbel@latest
```

### 3.2 Opis Poleceń

#### `init`
Interaktywny kreator projektu. Tworzy folder `locales/` i przykładowy plik.

#### `lint`
Waliduje składnię i reguły AI.
*   `mbel lint ./locales`
*   Flagi: `-j` (wątki), `-v` (gadatliwość).

#### `compile`
Kompiluje pliki do JSON (dla produkcji).
*   `mbel compile ./locales -o out.json`
*   Flaga `--ns`: Automatycznie tworzy namespace z nazwy folderu.

#### `watch`
Tryb deweloperski. Nasłuchuje zmian i recompiluje w tle.

#### `fmt`
Formater kodu.

---

## 4. Integracja Go SDK

SDK jest bezpieczne wielowątkowo (thread-safe) i gotowe na produkcję.

### 4.1 Architektura

*   **Manager**: Centralny punkt wejścia.
*   **Runtime**: Instancja dla konkretnego języka.
*   **Repository**: Źródło danych (Pliki, Baza Danych).

### 4.2 Użycie (T Function)

Funkcja `T` (Translate) rozwiązuje tekst na podstawie kontekstu.

```go
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // 1. Prosty Klucz
    title := mbel.T(ctx, "page_title")

    // 2. Ze zmiennymi (Vars)
    msg := mbel.T(ctx, "welcome_message", mbel.Vars{
        "name": "Jan",
    })

    // 3. Z liczbą mnogą i logiką
    // Przekazanie int jako argumentu automatycznie jest traktowane jako 'n'
    items := mbel.T(ctx, "cart_items", 5) 
}
```

### 4.4 Wzorzec Repository

Dla systemów Enterprise, możesz ładować tłumaczenia z bazy danych. Zaimplementuj interfejs:

```go
type Repository interface {
    LoadAll() (map[string]map[string]interface{}, error)
}
```

### 4.5 Middleware HTTP

Automatycznie parsuje nagłówek `Accept-Language` z przeglądarki i ustawia odpowiedni locale w Context.

```go
router.Use(mbel.Middleware)
```

---

*Dokumentacja wygenerowana automatycznie przez MBEL Team.*
