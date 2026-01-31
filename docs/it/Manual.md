# MBEL: Manuale di Riferimento Completo

**Versione:** 1.2.0
**Data:** Gennaio 2026

---

## 📖 Indice

1.  [Introduzione](#1-introduzione)
2.  [Il Linguaggio MBEL](#2-il-linguaggio-mbel)
    *   [Struttura del File](#21-struttura-del-file)
    *   [Tipi di Dati](#22-tipi-di-dati)
    *   [Interpolazione & Variabili](#23-interpolazione--variabili)
    *   [Logica & Controllo](#24-logica--controllo)
    *   [Pluralizzazione](#25-pluralizzazione)
3.  [Strumenti CLI](#3-strumenti-cli)
4.  [Integrazione Go SDK](#4-integrazione-go-sdk)

---

## 1. Introduzione

Modern Billed-English Language (MBEL) è un formato di localizzazione progettato per l'era dell'Intelligenza Artificiale.

---

## 2. Il Linguaggio MBEL

### 2.1 Struttura del File

```mbel
@namespace: "features.auth"  # Metadati

[Schermata Login]            # Sezione

title = "Accedi"             # Assegnazione
```

### 2.2 Tipi di Dati

*   **Stringa Letterale**: Tra doppi apici.
*   **Stringa Multilinea**: Tra tripli apici `"""`.

### 2.3 Interpolazione & Variabili

Variabili tra parentesi graffe `{}`.

*   **Sintassi**: `Ciao, {user_name}!`

### 2.4 Logica & Controllo

**Sintassi:** `chiave(variabile) { casi }`

#### Corrispondenza Esatta
```mbel
theme(mode) {
    [dark]  => "Modalità Scura"
    [light] => "Modalità Chiara"
}
```

#### Corrispondenza Intervallo
```mbel
battery(percent) {
    [0]       => "Scarica"
    [1-19]    => "Bassa"
    [100]     => "Carica"
}
```

### 2.5 Pluralizzazione

**Esempio (Italiano):**
```mbel
items_it(n) {
    [one]   => "1 elemento"    # n = 1
    [other] => "{n} elementi"  # n != 1
}
```

---

## 1. Perché MBEL? (Confronto reale)

Pensi ancora che il JSON sia sufficiente? Confrontiamo una semplice regola del plurale + interpolazione.

#### ❌ Il metodo JSON (Disordinato)
```json
{
  "cart_items_one": "Hai 1 articolo nel carrello.",
  "cart_items_other": "Hai {{count}} articoli nel carrello.",
  "greeting_male": "Bentornato, Sig. {{name}}",
  "greeting_female": "Bentornata, Sig.ra {{name}}"
}
```
*La logica è dispersa su più chiavi. Il codice Go/JS deve decidere quale chiave recuperare.*

#### ✅ Il metodo MBEL (Pulito)
```mbel
cart_items(n) {
    [one]   => "Hai 1 articolo nel carrello."
    [other] => "Hai {n} articoli nel carrello."
}

greeting(gender) {
    [male]   => "Bentornato, Sig. {name}"
    [female] => "Bentornata, Sig.ra {name}"
}
```
*Una sola chiave, logica pulita. Il runtime si occupa del lavoro pesante.*

---

## 2. Guida alla Sintassi

### Chiavi Base
Semplici coppie chiave-valore.

```mbel
title = "La Mia Applicazione"
```

### Interpolazione vs Variabili Logiche
Distinzione importante:
1. **Variabile di Controllo**: Quella tra parentesi `key(var)`. Decide QUALE caso viene scelto.
2. **Variabili di Interpolazione**: Quelle tra graffe `{var}`. Vengono semplicemente sostituite dal testo.

```mbel
# 'gender' è la variabile di controllo
# '{name}' è la variabile di interpolazione
greeting(gender) {
    [male]   => "Ciao Sig. {name}"
    [female] => "Ciao Sig.ra {name}"
}
```
*Al runtime:* `mbel.T(ctx, "greeting", mbel.Vars{"gender": "male", "name": "Bob"})`

### Metadati IA
I metadati sono memorizzati nel campo `__ai` dell'oggetto compilato. Non influenzano il testo a runtime, ma potenziano massicciamente gli agenti di traduzione.

---

## 3. Strumenti CLI

*   `mbel init`: Creazione progetto.
*   `mbel lint`: Validazione sintassi.
*   `mbel compile`: Compilazione JSON.
*   `mbel watch`: Ricaricamento live.
*   **`stats`**: Statistiche.
*   **`fmt`**: Formattazione.

---

## 4. Integrazione Go SDK

### 4.1 Architettura

*   **Manager**: Punto di ingresso centrale.
*   **Runtime**: Istanza per lingua specifica.
*   **Repository**: Interfaccia sorgente dati.

### 4.2 Inizializzazione

```go
import "github.com/yourusername/mbel"

func init() {
    mbel.Init("./locales", mbel.Config{
        DefaultLocale: "it",
        Watch:         true,
    })
}
```

### 4.3 Utilizzo (Funzione T)

```go
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // 1. Chiave Semplice
    title := mbel.T(ctx, "page_title")

    // 2. Con Variabili
    msg := mbel.T(ctx, "welcome", mbel.Vars{"name": "Mario"})

    // 3. Con Plurale
    items := mbel.T(ctx, "cart_items", 5) 
}
```

### 4.4 Middleware HTTP

MBEL analizza automaticamente l'header `Accept-Language`.

```go
router.Use(mbel.Middleware)
```
