# MBEL: FAQ (Frequently Asked Questions)

### Q: Why not just use JSON?
**A:** JSON is a data format, not a language. MBEL is a programmable localization language. JSON doesn't support comments, logic blocks, or AI metadata natively. MBEL makes your localization files 10x more readable and "AI-Ready".

### Q: Does MBEL support plural rules for my language?
**A:** Yes. MBEL comes with built-in CLDR rules for over 15 major world languages including English, Polish, German, French, Russian, Spanish, Chinese, Japanese, and more.

### Q: How fast is the Go SDK?
**A:** Extremely fast. MBEL compiles files into an optimized AST. Resolution (the `T` function) is a simple O(1) map lookup or a fast condition match. In our benchmarks, it's comparable to or faster than `i18next` (JS) or `gettext`.

### Q: Can I use MBEL with React/Vue/Frontend?
**A:** Yes! Use the `mbel compile` command to export your files to JSON. You can then load this JSON into any JS-based i18n library, or wait for the native MBEL-JS package (coming soon).

### Q: What Happens if a key is missing?
**A:** The `T` function returns the key itself as a fallback. This ensures your app never shows empty spaces or crashes due to missing translations.

### Q: My translation has a `{n}` but I passed a string.
**A:** That's fine. MBEL is flexible. It will try to stringify any argument provided and swap the placeholder.
