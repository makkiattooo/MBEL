# MBEL: Best Practices & Design Patterns

To get the most out of MBEL, follow these patterns used by professional localization teams.

## 1. Key Naming Conventions
Don't just use `label1`, `label2`. Use hierarchical keys.

*   **Pattern**: `[feature].[screen].[component].[element]`
*   **Good**: `auth.login.form.email_placeholder`
*   **Bad**: `login_email`

## 2. Organizing Files
For small projects, one file per language is fine. For large apps, use the folder structure.

**Recommended Structure:**
```text
locales/
  en/
    common.mbel      (Global app strings)
    auth/
      login.mbel     (Login screen)
      register.mbel  (Register screen)
    settings.mbel    (Settings panel)
  pl/
    ... (mirror structure)
```
MBEL's `--ns` (namespace) flag will automatically turn `auth/login.mbel` into the `auth.login` prefix.

## 3. Using Context for Humans and AI
Always use `@AI_Context` and comments. It helps translators (and LLMs) understand where the text appears.

```mbel
# Button on the bottom of the checkout page
@AI_Context: "Call to action for payment"
@AI_Tone: "Urgent but professional"
pay_now = "Complete Payment"
```

## 4. Logical Blocks vs. String Interpolation
**Avoid** building sentences in Go code.
*   **Bad (Go)**: `fmt.Sprintf(mbel.T(ctx, "hello") + " " + name)`
*   **Good (MBEL)**:
    ```mbel
    welcome(name) {
        [other] => "Welcome, {name}!"
    }
    ```

## 5. Handling Terms (Glossary)
Use the `-term` syntax for branding or words that never change.
```mbel
@brand_name = "SpaceCorp 🚀"

# Usage:
footer = "Copyright 2026 by {-brand_name}"
```
If you change `@brand_name` once, it updates everywhere.
