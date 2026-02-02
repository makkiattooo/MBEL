# Security & Best Practices

## XSS Prevention (HTML Escaping)

By default, MBEL **does NOT escape HTML** when interpolating values. This is intentional to provide maximum flexibility. However, if you're rendering translations in a web/HTML context, you **must** enable HTML escaping to prevent XSS attacks.

### Enable HTML Escape:

```go
import "github.com/makkiattooo/MBEL"

// Safe for web contexts
r := mbel.NewRuntimeWithOptions(data, true) // true = enable HTML escape
msg := r.Get("greeting", mbel.Vars{"name": "<script>alert('xss')</script>"})
// Result: "Hello &lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"
```

### Default (No Escaping):

```go
// Raw interpolation - suitable for system/internal messages
r := mbel.NewRuntime(data) // or NewRuntimeWithOptions(data, false)
msg := r.Get("greeting", mbel.Vars{"name": "<b>Alice</b>"})
// Result: "Hello <b>Alice</b>"
```

## Principle: Data vs. Template

MBEL treats translations as **data + logic**, not templates. Keys contain hardcoded structure; only placeholders (`{var}`) are substituted. This limits injection attacks by design.

---

## Performance Best Practices

### 1. Use Lazy Loading for Large Projects

For applications with many languages:

```go
m, err := mbel.NewManager("./locales", mbel.Config{
    DefaultLocale: "en",
    LazyLoad: true, // Load runtimes on first use
})
```

This reduces startup time and memory footprint when you have 20+ languages.

### 2. Enable Compilation Caching

The `FileRepository` automatically caches compiled `.mbel` files by modification time. When you reload, unchanged files are reused.

```go
// During watch mode or reload cycles, cache is active by default
mgr.Load() // Only recompiles changed files
```

### 3. Monitor Metrics

Track runtime performance:

```go
metrics := mbel.GetMetrics()
fmt.Printf("Get calls: %d\n", metrics["get_calls"])
fmt.Printf("Interpolations: %d\n", metrics["interpolate_ops"])
```

---

## Input Validation

### Validate User-Provided Locale Codes

```go
func isValidLocale(lang string) bool {
    // Only allow 2-5 character locale codes
    if len(lang) < 2 || len(lang) > 5 {
        return false
    }
    // Optionally validate against known languages
    return strings.Contains("en,de,fr,pl,ru", lang[:2])
}
```

---

## CI/CD Integration

### 1. Validate All `.mbel` Files

```bash
mbel lint locales/ --v
```

### 2. Check Format Consistency

```bash
mbel fmt locales/ -n # Dry-run
```

### 3. Compile Before Deployment

```bash
mbel compile locales/ -o translations.json
```

---

## Common Pitfalls

### ❌ Don't:
- Interpolate untrusted input without `escapeHTML` in web contexts
- Trust user-provided locale codes without validation
- Assume all translation strings are simple — use patterns and logic blocks

### ✅ Do:
- Enable HTML escaping for user-facing web output
- Validate locale codes against a whitelist
- Test plural rules with edge cases (0, 1, 2, negative numbers)
- Use `mbel lint` in CI to catch syntax errors early

