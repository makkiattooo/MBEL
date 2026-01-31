# MBEL : Conseils & Bonnes Pratiques

Pour tirer le meilleur parti de MBEL, suivez ces modèles éprouvés.

## 1. Conventions de nommage
Utilisez des clés hiérarchiques au lieu de noms plats comme `label1`.
*   **Modèle** : `[feature].[screen].[component].[element]`

## 2. Organisation des fichiers
Divisez vos localisations en unités logiques. La structure des dossiers définit automatiquement les espaces de noms avec l'option `--ns`.

## 3. Utilisation du contexte IA
N'économisez pas sur `@AI_Context`. C'est l'assurance qualité de vos traductions.

## 4. Blocs logiques vs Interpolation dans le code
Évitez de construire des phrases dans votre code Go.
*   **Mauvais (Go)** : `fmt.Sprintf(mbel.T(ctx, "bonjour") + " " + name)`
*   **Bon (MBEL)** :
    ```mbel
    welcome(name) {
        [other] => "Bienvenue, {name} !"
    }
    ```
