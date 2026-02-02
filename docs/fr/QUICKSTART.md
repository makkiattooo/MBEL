# Démarrage rapide : Votre première application MBEL

Lancez une application multilingue prête pour la production avec MBEL en **15 minutes**. Ce guide couvre la configuration du projet, les fichiers de traduction, la compilation et le déploiement.

---

## 1. Initialisez votre projet

```bash
# Créez le répertoire du projet
mkdir -p hello-mbel && cd hello-mbel

# Initialisez le module Go (utiliser MBEL comme bibliothèque)
go mod init hello-mbel
go get github.com/makkiattooo/MBEL@latest

# Créez la structure du projet
mkdir -p locales cmd dist

# Installez l'CLI MBEL pour la compilation
go install github.com/makkiattooo/MBEL/cmd/mbel@latest
```

---

## 2. Écrivez les fichiers de traduction

### Créez des traductions françaises (`locales/fr.mbel`)

```mbel
@namespace: hello
@lang: fr

app_name = "Bonjour MBEL"
app_version = "1.0.0"

greeting = "Bienvenue, {name}!"
goodbye = "À bientôt, {name}! Temps: {time}"

# Français: [one] (1), [other]
items_count(n) {
    [one]   => "Vous avez 1 article"
    [other] => "Vous avez {n} articles"
}

profile_updated(gender) {
    [male]   => "Il a mis à jour son profil"
    [female] => "Elle a mis à jour son profil"
    [other]  => "Ils ont mis à jour leur profil"
}

ui.menu {
    home = "Accueil"
    about = "À propos"
    contact = "Contact"
    settings = "Paramètres"
}

order_total = "Total: {price} (taxes comprises)"
```

### Ajoutez d'autres langues

Créez `locales/en.mbel`, `locales/de.mbel`, etc., en suivant le même modèle.

---

## 3. Validez et formatez les traductions

```bash
# Vérifiez les erreurs de syntaxe
mbel lint locales/

# Formatez automatiquement et corrigez les problèmes
mbel fmt locales/

# Afficher les statistiques
mbel stats locales/
```

---

## 4. Compilez les traductions

```bash
# Compilez toutes les traductions en un seul fichier JSON
mbel compile locales/ -o dist/translations.json

# Incluez la carte source pour le débogage
mbel compile locales/ -o dist/translations.json -sourcemap
```

---

## 5. Créez votre application Go

### Exemple basique (`cmd/main.go`)

```go
package main

import (
	"context"
	"fmt"
	"log"

	mbel "github.com/makkiattooo/MBEL/pkg/mbel"
)

func main() {
	// 1. Initialisez le Manager
	m, err := mbel.NewManager("./dist/translations.json", mbel.Config{
		DefaultLocale: "fr",
		FallbackChain: []string{"fr", "en"},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 2. Recherche simple de chaîne
	fmt.Println(m.Get("fr", "app_name", nil))
	fmt.Println(m.Get("en", "app_name", nil))

	// 3. Interpolation de variables
	vars := mbel.Vars{"name": "Alice"}
	greeting := m.Get("fr", "greeting", vars)
	fmt.Println(greeting)

	// 4. Pluralisation
	for _, count := range []int{1, 2, 5} {
		vars := mbel.Vars{"n": count}
		msg := m.Get("fr", "items_count", vars)
		fmt.Printf("n=%d: %s\n", count, msg)
	}

	// 5. API global
	mbel.Init(m)
	ctx := context.Background()
	fmt.Println(mbel.T(ctx, "greeting", mbel.Vars{"name": "Bob"}))
}
```

---

## 6. Ajoutez des tests

```bash
go test -v ./...
go test -bench=. ./...
```

---

## 7. Construisez et exécutez

```bash
go build -o hello-mbel ./cmd
./hello-mbel
```

---

## 8. Déploiement en production

### Option A: Intégrez JSON compilé

```go
//go:embed dist/translations.json
var translationsJSON []byte
```

### Option B: Expédiez avec traductions

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

## 9. Modèles courants

### Négociation de langue en HTTP

```go
lang := r.Header.Get("Accept-Language")
```

### Mise en cache et performance

```go
// Le Manager met automatiquement en cache les traductions
for i := 0; i < 1000000; i++ {
	msg := m.Get("fr", "app_name", nil)
}
```

---

## 10. Prochaines étapes

1. **[Manuel](Manual.md)** — Documentation complète de MBEL
2. **[ARCHITECTURE.md](ARCHITECTURE.md)** — Analyse technique approfondie
3. **[DEVELOPMENT.md](DEVELOPMENT.md)** — Étendre MBEL
4. **[Meilleures pratiques de sécurité](SECURITY.md)** — Prévention XSS

---

## Dépannage

| Problème | Solution |
|----------|----------|
| `no such file or directory: locales/` | `mkdir -p locales` |
| `Undefined: mbel` | Vérifier `go.mod` |
| `Syntax error at line 5` | `mbel lint locales/` |
| Traduction non trouvée | Vérifier la clé, le code de langue, la chaîne de retour |
