# MBEL: Najlepsze Praktyki i Wzorce Projektowe

Aby wycisnąć z MBEL 100%, stosuj sprawdzone wzorce używane przez profesjonalne zespoły lokalizacyjne.

## 1. Nazewnictwo Kluczy
Nie używaj `label1`, `label2`. Stosuj hierarchię.

*   **Wzorzec**: `[funkcja].[ekran].[komponent].[element]`
*   **Dobrze**: `auth.login.form.email_placeholder`
*   **Źle**: `login_email`

## 2. Organizacja Plików
W małych projektach jeden plik na język wystarczy. W dużych - używaj struktury folderów.

**Zalecana Struktura:**
```text
locales/
  en/
    common.mbel      (Globalne teksty)
    auth/
      login.mbel     (Ekran logowania)
  pl/
    ... (lustrzana struktura)
```
Flaga `--ns` automatycznie zamieni `auth/login.mbel` na prefiks `auth.login`.

## 3. Kontekst dla Ludzi i AI
Zawsze używaj `@AI_Context` i komentarzy. Pomaga to tłumaczom (i modelom LLM) zrozumieć, gdzie pojawia się tekst.

```mbel
# Przycisk na dole strony płatności
@AI_Context: "Wezwanie do działania przy płatności"
@AI_Tone: "Pilny, ale profesjonalny"
pay_now = "Zapłać teraz"
```

## 4. Bloki Logiczne vs. Interpolacja w Kodzie
**Unikaj** składania zdań w kodzie Go.
*   **Źle (Go)**: `fmt.Sprintf(mbel.T(ctx, "hello") + " " + name)`
*   **Dobrze (MBEL)**:
    ```mbel
    welcome(name) {
        [other] => "Witaj, {name}!"
    }
    ```
