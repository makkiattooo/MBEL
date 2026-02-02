# Security Best Practices: Protecting Your MBEL Applications

This guide covers security considerations when using MBEL in production environments, including XSS prevention, input validation, and secure deployment patterns.

---

## Overview

MBEL provides built-in security features for web applications, but security is a shared responsibility. This document outlines best practices for securing your MBEL deployments.

**Key security concerns:**
- **Cross-Site Scripting (XSS)**: User input in translations can execute JavaScript
- **Translation injection**: Compromised translation files could inject malicious content
- **Information disclosure**: Translation keys might reveal sensitive system information
- **Denial of Service (DoS)**: Large translations or nested structures could consume resources

---

## 1. XSS Prevention (HTML Escaping)

### Using the HTML Escape Option

By default, MBEL **does NOT escape HTML** when interpolating values. This is intentional to provide maximum flexibility. However, if you're rendering translations in a web/HTML context, you **must** enable HTML escaping to prevent XSS attacks.

**Enable HTML escaping for web applications:**

```go
import (
	"github.com/makkiattooo/MBEL/pkg/mbel"
)

func main() {
	// Safe for web contexts - enable HTML escaping
	m, err := mbel.NewManager("./dist/translations.json", mbel.Config{
		DefaultLocale: "en",
		HTMLEscape:    true, // Escape HTML entities in variables
	})
	if err != nil {
		panic(err)
	}

	// Variables are automatically escaped
	vars := mbel.Vars{"name": "<script>alert('xss')</script>"}
	msg := m.Get("en", "greeting", vars)
	// Output: "Welcome, &lt;script&gt;alert('xss')&lt;/script&gt;!"
}
```

**Without HTML escaping (use only for non-web contexts):**

```go
// ❌ DON'T USE FOR WEB APPLICATIONS
m, _ := mbel.NewManager("./dist/translations.json", mbel.Config{
	DefaultLocale: "en",
	HTMLEscape: false, // Raw interpolation - system messages only
})

// Only safe if variables are fully controlled and trusted
msg := m.Get("en", "greeting", mbel.Vars{"name": "Alice"})
```

### What Gets Escaped

When `HTMLEscape: true`:
- `<` → `&lt;`
- `>` → `&gt;`
- `&` → `&amp;`
- `"` → `&quot;`
- `'` → `&#x27;`

**This prevents execution:**
```html
<!-- Before (vulnerable) -->
<p>Welcome, <script>alert('xss')</script>!</p>

<!-- After (safe with HTMLEscape: true) -->
<p>Welcome, &lt;script&gt;alert('xss')&lt;/script&gt;!</p>
```

### HTTP Context Example

```go
package main

import (
	"net/http"
	"github.com/makkiattooo/MBEL/pkg/mbel"
)

var manager *mbel.Manager

func init() {
	var err error
	// ALWAYS enable HTML escaping for web
	manager, err = mbel.NewManager("./dist/translations.json", mbel.Config{
		DefaultLocale: "en",
		HTMLEscape:    true,
	})
	if err != nil {
		panic(err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Get user input from query parameter (untrusted)
	userName := r.URL.Query().Get("name")

	// MBEL automatically escapes it - safe to render
	msg := manager.Get("en", "greeting", mbel.Vars{"name": userName})
	
	// Safe to render in HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("<p>" + msg + "</p>"))
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":8080", nil)
}
```

---

## 2. Input Validation

### Validate Variables Before Interpolation

Even with HTML escaping, validate and sanitize input:

```go
import (
	"regexp"
	"strings"
	"fmt"
)

// Whitelist approach: only allow safe characters
func validateUserName(name string) (string, error) {
	// Only alphanumeric, spaces, hyphens, max 50 chars
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9\s\-]{1,50}$`, name)
	if !matched {
		return "", fmt.Errorf("invalid name format")
	}
	return strings.TrimSpace(name), nil
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	userName, err := validateUserName(r.URL.Query().Get("name"))
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Now safe to use with MBEL
	msg := manager.Get("en", "greeting", mbel.Vars{"name": userName})
	w.Write([]byte(msg))
}
```

### Validate Translation Keys

Prevent attackers from requesting arbitrary keys:

```go
// Only allow alphanumeric, dots, underscores (namespace.key format)
func validateKey(key string) (string, error) {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_.]{1,100}$`, key)
	if !matched {
		return "", fmt.Errorf("invalid key format")
	}
	return key, nil
}

func getTranslation(w http.ResponseWriter, r *http.Request) {
	key, err := validateKey(r.URL.Query().Get("key"))
	if err != nil {
		http.Error(w, "Invalid key", http.StatusBadRequest)
		return
	}

	msg := manager.Get("en", key, nil)
	w.Write([]byte(msg))
}
```

---

## 3. Translation File Security

### Secure File Permissions

Protect translation files with proper filesystem permissions:

```bash
# Compiled translations (readable by app, protected)
chmod 644 dist/translations.json
chown app:app dist/translations.json

# Source files (restrict access to translation team)
chmod 640 locales/*.mbel
chown translator:translators locales/

# Sourcemap (same as translations)
chmod 644 dist/translations.sourcemap.json
```

### Verify Translation Integrity

```bash
# Generate checksums during build
sha256sum dist/translations.json > dist/translations.json.sha256

# Verify before deployment
sha256sum -c dist/translations.json.sha256
```

**In Go:**

```go
import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
)

func verifyTranslations(filePath, checksumPath string) error {
	data, _ := ioutil.ReadFile(filePath)
	hash := sha256.Sum256(data)
	
	expected, _ := ioutil.ReadFile(checksumPath)
	if hex.EncodeToString(hash[:]) != strings.TrimSpace(string(expected)) {
		return fmt.Errorf("translation file integrity check failed")
	}
	return nil
}
```

### Secure Translation Pipeline

1. **Version Control**: Store translations in Git with signed commits
```bash
git commit -S locales/  # GPG-signed commits
git tag -s v1.0.0      # Signed tags
```

2. **Code Review**: Review all translation changes before merge
3. **Automated Validation**: Run linting and testing in CI/CD
```bash
mbel lint locales/
mbel compile locales/ -o dist/translations.json
go test ./...
```

4. **Secure Deployment**: Use HTTPS and verify checksums
```bash
curl --compressed https://api.example.com/translations.json \
  | sha256sum -c translations.json.sha256
```

---

## 4. Deployment Security

### Environment Variables

Never hardcode sensitive data in translations:

```mbel
# ❌ BAD - Don't include secrets in translations
api_key = "sk-1234567890abcdef"
database_url = "postgres://user:password@host/db"

# ✅ GOOD - Reference from environment
api_endpoint = "https://api.example.com/v1"
app_version = "1.2.3"
```

### Secure Configuration

```go
package main

import (
	"os"
	"github.com/makkiattooo/MBEL/pkg/mbel"
)

func init() {
	// Load translations from secure location
	translationPath := os.Getenv("TRANSLATION_PATH")
	if translationPath == "" {
		// Fallback to embedded or packaged translations
		translationPath = "./dist/translations.json"
	}

	manager, _ = mbel.NewManager(translationPath, mbel.Config{
		DefaultLocale: "en",
		HTMLEscape:    true,
	})
}
```

### Content Security Policy (CSP)

Set CSP headers to prevent inline script execution:

```go
func cspMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent inline scripts, restrict script sources
		w.Header().Set("Content-Security-Policy", 
			"default-src 'self'; script-src 'self'")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}
```

---

## 5. Runtime Security

### Rate Limiting

Protect against DoS attacks on translation lookups:

```go
import (
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(rate.Limit(1000), 100) // 1000 req/sec, burst 100

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	name := r.URL.Query().Get("name")
	msg := manager.Get("en", "greeting", mbel.Vars{"name": name})
	w.Write([]byte(msg))
}
```

### Resource Limits

```go
// Limit nested depth and translation size
const (
	MAX_NESTED_DEPTH       = 10
	MAX_TRANSLATION_SIZE   = 10_000 // bytes
	MAX_VARIABLE_SIZE      = 1_000  // bytes per variable
)

func safeGet(lang, key string, vars mbel.Vars) (string, error) {
	// Validate variable sizes
	for k, v := range vars {
		if len(fmt.Sprint(v)) > MAX_VARIABLE_SIZE {
			return "", fmt.Errorf("variable %s exceeds size limit", k)
		}
	}

	msg := manager.Get(lang, key, vars)
	if len(msg) > MAX_TRANSLATION_SIZE {
		return "", fmt.Errorf("translation exceeds size limit")
	}

	return msg, nil
}
```

### Logging and Monitoring

```go
import "log"

func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	lang := r.Header.Get("Accept-Language")

	// Log translation requests (be cautious with PII)
	log.Printf("translation_requested lang=%s key=%s client_ip=%s",
		lang, "greeting", r.RemoteAddr)

	msg := manager.Get(lang, "greeting", mbel.Vars{"name": name})
	w.Write([]byte(msg))
}
```

---

## 6. Common Vulnerabilities

### Template Injection

**Vulnerable:**
```go
// ❌ DON'T: Let users control format strings
format := r.URL.Query().Get("format")
msg := fmt.Sprintf(format, userData) // DANGEROUS
```

**Safe:**
```go
// ✅ DO: Use MBEL's controlled variable system
msg := manager.Get("en", "greeting", mbel.Vars{
	"name": userData, // Escaped automatically
})
```

### Directory Traversal

**Vulnerable:**
```go
// ❌ DON'T: Directly use user input as translation key
key := r.URL.Query().Get("key")
m.Get("en", key, nil) // User could pass "../../config"
```

**Safe:**
```go
// ✅ DO: Validate key format first
key, err := validateKey(r.URL.Query().Get("key"))
if err != nil {
	http.Error(w, "Invalid key", http.StatusBadRequest)
	return
}
m.Get("en", key, nil)
```

---

## 7. Data Protection

### GDPR Considerations

**Personal Data in Translations:**
- Names, email addresses in personalized messages
- User preferences stored in translation keys
- Language choices

**Mitigation:**
```go
// Pseudonymize in logs
func pseudonymizeEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts[0]) < 3 {
		return "***@" + parts[1]
	}
	return parts[0][:3] + "***@" + parts[1]
}

// Implement data retention
func deleteTranslationHistory(userID string, before time.Time) error {
	// Delete logs, caches, temporary data
	return database.Delete("translation_logs", 
		"user_id = ? AND created_at < ?", userID, before)
}
```

### PCI Compliance

Never include payment data in translations:

```mbel
# ❌ BAD - Exposes card data
confirmation = "Your payment {card_number} was processed"

# ✅ GOOD - Use masked values
confirmation = "Your payment ending in {card_last_4} was processed"
```

---

## 8. Security Checklist

**Development:**
- [ ] HTML escaping enabled for all web applications
- [ ] Input validation implemented for all variables
- [ ] No sensitive data in translation files
- [ ] Translation keys don't expose system internals
- [ ] No hardcoded secrets or passwords

**Testing:**
- [ ] XSS payload testing: `<script>alert('xss')</script>`
- [ ] SQL injection in variables: `'; DROP TABLE users; --`
- [ ] Large input testing (DoS prevention)
- [ ] Special character handling: quotes, backslashes, unicode
- [ ] Rate limiting tests
- [ ] File permission tests

**Deployment:**
- [ ] File permissions secured (644 for translations)
- [ ] HTTPS enabled for all endpoints
- [ ] CSP headers configured
- [ ] Security headers set (X-Frame-Options, X-Content-Type-Options)
- [ ] Rate limiting implemented
- [ ] Logging and monitoring active
- [ ] Translation file checksums verified
- [ ] Secrets stored in environment variables only
- [ ] Regular security updates applied

**Operations:**
- [ ] Regular security audits scheduled
- [ ] Keep MBEL and Go dependencies updated
- [ ] Monitor for suspicious translation access patterns
- [ ] Implement backup and recovery procedures
- [ ] Document incident response procedures
- [ ] Encryption at rest for sensitive translations

---

## 9. Resources

- **OWASP Top 10**: https://owasp.org/www-project-top-ten/
- **CWE-79 (XSS)**: https://cwe.mitre.org/data/definitions/79.html
- **Go Security**: https://golang.org/doc/security
- **MBEL GitHub**: https://github.com/makkiattooo/MBEL

---

## Reporting Security Issues

Found a security vulnerability in MBEL?

**Do NOT open a public GitHub issue.** Instead:

1. Email security report to: contact via GitHub profile
2. Include: Description, steps to reproduce, impact assessment
3. Allow 48-72 hours for acknowledgment and patch development
4. Coordinate responsible disclosure timeline

Researchers will be credited in security advisory.

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

