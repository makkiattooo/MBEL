# MBEL: Tipps & Best Practices

Um MBEL optimal zu nutzen, folgen Sie diesen bewährten Mustern.

## 1. Namenskonventionen
Nutzen Sie hierarchische Schlüssel anstelle von flachen Namen wie `label1`.

*   **Muster**: `[feature].[screen].[component].[element]`
*   **Gut**: `auth.login.form.email_placeholder`

## 2. Datei-Organisation
Teilen Sie Ihre Lokalisierungen in logische Einheiten auf.

**Empfohlene Struktur:**
```text
locales/
  en/
    common.mbel      (Globale Strings)
    auth/
      login.mbel     (Login-Bildschirm)
  de/
    ... (gespiegelte Struktur)
```
Die Flagge `--ns` von MBEL wandelt `auth/login.mbel` automatisch in den Präfix `auth.login` um.

## 3. KI-Kontext nutzen
Sparen Sie nicht an `@AI_Context`. Es ist die Versicherung für Ihre Übersetzungsqualität.

```mbel
# Button unten auf der Checkout-Seite
@AI_Context: "Handlungsaufforderung zur Zahlung"
@AI_Tone: "Dringend, aber professionell"
pay_now = "Jetzt bezahlen"
```

## 4. Logik-Blöcke vs. Code-Interpolation
Vermeiden Sie es, Sätze im Go-Code zusammenzubauen.
*   **Schlecht (Go)**: `fmt.Sprintf(mbel.T(ctx, "hallo") + " " + name)`
*   **Gut (MBEL)**:
    ```mbel
    welcome(name) {
        [other] => "Willkommen, {name}!"
    }
    ```
