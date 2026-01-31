# MBEL vs. The World (Engineers' Edition)

Most localization tools were built 20 years ago. They don't understand the reality of modern dev-cycles.

## 1. The "JSON Merge Hell"
In massive projects, `en.json` becomes a battleground. Adding one key at the end of a 5000-line file results in a git conflict every single time another dev does the same.

**The MBEL Solution**: MBEL encourages **Namespaced Files**. You work on `auth.mbel` or `billing.mbel`. Files are small, logical, and code-like. Merge conflicts are reduced by ~90% because you're touching specific domain files, not a global giant blob.

## 2. The "ICU Syntax" Trauma
Have you seen ICU plural rules in JSON?
`{count, plural, =0{no items} one{1 item} other{# items}}`

It's a "language-within-a-string". It's unreadable, fragile, and LLMs often break the nested brackets.

**The MBEL Solution**: MBEL uses **Native Block Logic**.
```mbel
items(count) {
    [0]     => "Empty"
    [one]   => "1 item"
    [other] => "{count} items"
}
```
It feels like code. It acts like code. It's deterministic.

## 3. The "AI Guessing" Game
When you send 5000 JSON keys to an AI for translation, it has no context. Is `Book` a noun or a verb? 

**The MBEL Solution**: First-class **AI Metadata**.
```mbel
@AI_Context: "Button to reserve a hotel room (Verb)"
book_btn = "Book Now"
```
Our CLI tools pass this context to LLMs, guaranteeing 99.9% translation accuracy without human correction.

---

## Technical Comparison

| Feature | MBEL | JSON / i18next | Fluent (Mozilla) |
| :--- | :--- | :--- | :--- |
| **Syntax** | **Clean DSL** | String-in-JSON | Complex DSL |
| **Logic** | **Native Blocks** | ICU (Fragile) | Functional |
| **AI Ready** | **✅ Yes (@AI_)** | ❌ No | ❌ No |
| **Merge Safety**| **High** | Low (Giant Blobs) | Medium |
| **Safety** | **Lint/Compile** | Runtime Only | Runtime |
