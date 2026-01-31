# MBEL: Suggerimenti e Best Practice

Segui questi pattern per ottenere il massimo da MBEL.

## 1. Convenzioni di Naming
Usa chiavi gerarchiche invece di nomi piatti come `label1`.
*   **Pattern**: `[feature].[screen].[component].[element]`

## 2. Organizzazione dei File
Dividi le tue localizzazioni in unità logiche. La struttura delle cartelle definisce automaticamente i namespace con l'opzione `--ns`.

## 3. Uso del Contesto IA
Non risparmiare su `@AI_Context`. È l'assicurazione sulla qualità delle tue traduzioni.

## 4. Blocchi Logici vs Interpolazione nel Codice
Evita di costruire frasi nel tuo codice Go.
*   **Male (Go)**: `fmt.Sprintf(mbel.T(ctx, "ciao") + " " + name)`
*   **Bene (MBEL)**:
    ```mbel
    welcome(name) {
        [other] => "Benvenuto, {name}!"
    }
    ```
