# MBEL: Manuel de Référence Complet

**Version :** 1.2.0
**Date :** Janvier 2026

---

## 📖 Table des Matières

1.  [Introduction](#1-introduction)
2.  [Le Langage MBEL](#2-le-langage-mbel)
    *   [Structure du Fichier](#21-structure-du-fichier)
    *   [Types de Données](#22-types-de-données)
    *   [Interpolation & Variables](#23-interpolation--variables)
    *   [Logique & Contrôle](#24-logique--contrôle)
    *   [Règles de Pluralisation](#25-règles-de-pluralisation)
    *   [Métadonnées IA](#26-métadonnées-ia)
3.  [Outils CLI](#3-outils-cli)
4.  [Intégration Go SDK](#4-intégration-go-sdk)

---

## 1. Introduction

Modern Billed-English Language (MBEL) est un format de localisation conçu pour combler le fossé entre les développeurs et l'Intelligence Artificielle.

### Philosophie
1.  **Le Contexte est Roi**.
2.  **La Logique appartient aux Données**.
3.  **IA d'abord**.

---

## 2. Le Langage MBEL

### 2.1 Structure du Fichier

```mbel
@namespace: "features.auth"  # Métadonnées

[Écran de Connexion]         # Section

title = "Connexion"          # Assignation
```

### 2.2 Types de Données

*   **Chaîne Littérale** : Entre guillemets doubles.
*   **Chaîne Multiligne** : Entre triples guillemets `"""`.

### 2.3 Interpolation & Variables

Les variables sont entre accolades `{}`.

*   **Syntaxe** : `Bonjour, {user_name} !`

### 2.4 Logique & Contrôle

**Syntaxe :** `clé(variable) { cas }`

#### Correspondance Exacte
```mbel
theme(mode) {
    [dark]  => "Mode Sombre"
    [light] => "Mode Clair"
}
```

#### Correspondance de Plage
```mbel
battery(percent) {
    [0]       => "Vide"
    [1-19]    => "Faible"
    [100]     => "Pleine"
}
```

## 1. Pourquoi MBEL ? (Comparaison réelle)

Vous pensez toujours que le JSON suffit ? Comparons une règle de pluriel simple + interpolation.

#### ❌ La méthode JSON (Désordonnée)
```json
{
  "cart_items_one": "Vous avez 1 article dans votre panier.",
  "cart_items_other": "Vous avez {{count}} articles dans votre panier.",
  "greeting_male": "Bon retour, M. {{name}}",
  "greeting_female": "Bon retour, Mme {{name}}"
}
```
*La logique est dispersée sur plusieurs clés. Le code Go/JS doit décider quelle clé récupérer.*

#### ✅ La méthode MBEL (Propre)
```mbel
cart_items(n) {
    [one]   => "Vous avez 1 article dans votre panier."
    [other] => "Vous avez {n} articles dans votre panier."
}

greeting(gender) {
    [male]   => "Bon retour, M. {name}"
    [female] => "Bon retour, Mme {name}"
}
```
*Une seule clé, une logique propre. Le runtime s'occupe du reste.*

---

## 2. Guide de Syntaxe

### Clés de base
Paires clé-valeur simples. Utilisez `"""` pour les chaînes multilignes.

```mbel
title = "Mon Application"
```

### Interpolation vs Variables Logiques
Distinction importante :
1. **Variable de Contrôle** : Celle dans `key(var)`. Elle décide QUEL cas est choisi.
2. **Variables d'Interpolation** : Celles dans `{var}`. Elles sont simplement remplacées par du texte.

```mbel
# 'gender' est la variable de contrôle
# '{name}' est la variable d'interpolation
greeting(gender) {
    [male]   => "Bonjour M. {name}"
    [female] => "Bonjour Mme {name}"
}
```
*Au runtime :* `mbel.T(ctx, "greeting", mbel.Vars{"gender": "male", "name": "Bob"})`

### Métadonnées IA
Les métadonnées sont stockées dans le champ `__ai` de l'objet compilé. Elles n'affectent pas le texte au runtime mais aident massivement les agents de traduction.

### 2.5 Règles de Pluralisation

MBEL utilise les règles CLDR.

**Exemple (Français):**
En français, 0 et 1 sont au singulier ("one").
```mbel
items_fr(n) {
    [one]   => "{n} élément"   # n = 0 ou n = 1
    [other] => "{n} éléments"  # n >= 2
}
```

### 2.6 Métadonnées IA

Annotations commençant par `@AI_`.

```mbel
@AI_Tone: "Formel"
submit_btn = "Valider"
```

---

## 3. Outils CLI

### 3.1 Installation

```bash
go install github.com/yourusername/mbel/cmd/mbel@latest
```

### 3.2 Commandes

*   **`init`** : Assistant de configuration.
*   **`lint`** : Validation de syntaxe.
*   **`compile`** : Compilation vers JSON.
*   **`watch`** : Mode développement.
*   **`stats`** : Statistiques.
*   **`fmt`** : Formatage automatique.

---

## 4. Intégration Go SDK

### 4.1 Architecture

*   **Manager** : Point d'entrée central.
*   **Runtime** : Instance pour une langue spécifique.
*   **Repository** : Interface de source de données.

### 4.2 Initialisation

```go
import "github.com/yourusername/mbel"

func init() {
    mbel.Init("./locales", mbel.Config{
        DefaultLocale: "fr",
        Watch:         true,
    })
}
```

### 4.3 Utilisation (Fonction T)

`T` (Translate) résout la chaîne selon le contexte.

```go
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // 1. Clé Simple
    title := mbel.T(ctx, "page_title")

    // 2. Avec Variables
    msg := mbel.T(ctx, "welcome", mbel.Vars{"name": "Pierre"})

    // 3. Avec Pluriel
    items := mbel.T(ctx, "cart_items", 5) 
}
```

### 4.4 Middleware HTTP

MBEL analyse automatiquement l'en-tête `Accept-Language`.

```go
router.Use(mbel.Middleware)
```
