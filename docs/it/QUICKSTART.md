# Guida rapida: La tua prima applicazione MBEL

Avvia un'applicazione multilingue pronta per la produzione con MBEL in **15 minuti**. Questa guida copre la configurazione del progetto, i file di traduzione, la compilazione e la distribuzione.

---

## 1. Inizializza il tuo progetto

```bash
# Crea la directory del progetto
mkdir -p hello-mbel && cd hello-mbel

# Inizializza il modulo Go
go mod init hello-mbel
go get github.com/makkiattooo/MBEL@latest

# Crea la struttura del progetto
mkdir -p locales cmd dist

# Installa l'CLI di MBEL
go install github.com/makkiattooo/MBEL/cmd/mbel@latest
```

---

## 2. Scrivi i file di traduzione

### Crea le traduzioni italiane (`locales/it.mbel`)

```mbel
@namespace: hello
@lang: it

app_name = "Ciao MBEL"
app_version = "1.0.0"

greeting = "Benvenuto, {name}!"
goodbye = "Arrivederci, {name}! Tempo: {time}"

# Italiano: [one] (1), [other]
items_count(n) {
    [one]   => "Hai 1 articolo"
    [other] => "Hai {n} articoli"
}

profile_updated(gender) {
    [male]   => "Ha aggiornato il suo profilo"
    [female] => "Ha aggiornato il suo profilo"
    [other]  => "Hanno aggiornato il loro profilo"
}

ui.menu {
    home = "Home"
    about = "Chi siamo"
    contact = "Contatti"
    settings = "Impostazioni"
}

order_total = "Totale: {price} (tasse incluse)"
```

### Aggiungi altri linguaggi

Crea `locales/en.mbel`, `locales/de.mbel`, etc.

---

## 3. Valida e formatta le traduzioni

```bash
# Verifica errori di sintassi
mbel lint locales/

# Formatta automaticamente
mbel fmt locales/

# Mostra statistiche
mbel stats locales/
```

---

## 4. Compila le traduzioni

```bash
# Compila in un unico file JSON
mbel compile locales/ -o dist/translations.json

# Includi mappa sorgente per il debug
mbel compile locales/ -o dist/translations.json -sourcemap
```

---

## 5. Crea la tua applicazione Go

### Esempio di base (`cmd/main.go`)

```go
package main

import (
	"context"
	"fmt"
	"log"

	mbel "github.com/makkiattooo/MBEL/pkg/mbel"
)

func main() {
	// 1. Inizializza il Manager
	m, err := mbel.NewManager("./dist/translations.json", mbel.Config{
		DefaultLocale: "it",
		FallbackChain: []string{"it", "en"},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 2. Ricerca semplice
	fmt.Println(m.Get("it", "app_name", nil))
	fmt.Println(m.Get("en", "app_name", nil))

	// 3. Interpolazione di variabili
	vars := mbel.Vars{"name": "Alice"}
	greeting := m.Get("it", "greeting", vars)
	fmt.Println(greeting)

	// 4. Pluralizzazione
	for _, count := range []int{1, 2, 5} {
		vars := mbel.Vars{"n": count}
		msg := m.Get("it", "items_count", vars)
		fmt.Printf("n=%d: %s\n", count, msg)
	}

	// 5. API globale
	mbel.Init(m)
	ctx := context.Background()
	fmt.Println(mbel.T(ctx, "greeting", mbel.Vars{"name": "Mario"}))
}
```

---

## 6. Aggiungi test

```bash
go test -v ./...
```

---

## 7. Compila ed esegui

```bash
go build -o hello-mbel ./cmd
./hello-mbel
```

---

## 8. Distribuzione in produzione

### Opzione A: Incorpora JSON

```go
//go:embed dist/translations.json
var translationsJSON []byte
```

### Opzione B: Distribuisci con Docker

```dockerfile
FROM golang:1.21 AS builder
WORKDIR /src
COPY . .
RUN go build -o /tmp/hello-mbel ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /tmp/hello-mbel .
COPY dist/translations.json ./dist/
CMD ["./hello-mbel"]
```

---

## 9. Pattern comuni

### Negoziazione della lingua

```go
lang := r.Header.Get("Accept-Language")
```

### Caching e performance

```go
// Il Manager cache automaticamente
for i := 0; i < 1000000; i++ {
	msg := m.Get("it", "app_name", nil)
}
```

---

## 10. Prossimi passi

1. **[Manuale](Manual.md)** — Documentazione completa
2. **[ARCHITECTURE.md](ARCHITECTURE.md)** — Analisi tecnica approfondita
3. **[DEVELOPMENT.md](DEVELOPMENT.md)** — Estendi MBEL
4. **[Best practices di sicurezza](SECURITY.md)** — Prevenzione XSS

---

## Risoluzione dei problemi

| Problema | Soluzione |
|----------|----------|
| `no such file or directory: locales/` | `mkdir -p locales` |
| `Undefined: mbel` | Controlla `go.mod` |
| `Syntax error at line 5` | `mbel lint locales/` |
| Traduzione non trovata | Controlla chiave, lingua e catena fallback |
