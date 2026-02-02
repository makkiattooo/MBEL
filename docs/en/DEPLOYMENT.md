# Production Deployment Guide

## Overview

This guide covers deploying MBEL-based localization systems to production, including CI/CD integration, performance optimization, and monitoring.

---

## Pre-Deployment Checklist

- [ ] All `.mbel` files pass `mbel lint` validation
- [ ] No syntax errors: `go test ./...` passes
- [ ] Benchmark baseline recorded: `go test -bench . -benchmem`
- [ ] Security review: HTML escape enabled for web contexts
- [ ] Translations reviewed for all target languages
- [ ] Sourcemap generated for debugging: `mbel compile -sourcemap`
- [ ] CI/CD pipeline configured (GitHub Actions, GitLab CI, etc.)
- [ ] Monitoring/alerting configured for translation failures

---

## Build & Compile Strategy

### Option 1: Compile at Build Time (Recommended)

Compile translations once during build, embed the JSON:

```bash
#!/bin/bash
# build.sh

# 1. Validate
mbel lint locales/ || exit 1

# 2. Format & cleanup
mbel fmt locales/

# 3. Compile to JSON
mbel compile locales/ \
  -o dist/translations.json \
  -sourcemap \
  -pretty=false

# 4. Create sourcemap index for debugging
cp dist/translations.sourcemap.json dist/sourcemap.idx

# 5. Embed in binary (Go example)
go embed dist/translations.json
go build -o app
```

### Option 2: Lazy-Load on Startup

For large translation repositories, load languages on-demand:

```go
m, _ := mbel.NewManager("./locales", mbel.Config{
    DefaultLocale: "en",
    LazyLoad: true,      // Load runtimes only when needed
    HotReload: false,    // Disable watch mode in production
})
```

**Benefits:**
- Faster startup time
- Lower memory footprint
- Only load languages users actually request

### Option 3: Remote Translation Server

Fetch translations from a centralized service:

```go
type RemoteRepository struct {
    baseURL string
    client  *http.Client
}

func (r *RemoteRepository) Get(lang string) (mbel.Data, error) {
    resp, _ := r.client.Get(fmt.Sprintf("%s/api/translations/%s", r.baseURL, lang))
    var data mbel.Data
    json.NewDecoder(resp.Body).Decode(&data)
    return data, nil
}

m := mbel.NewManagerWithRepository(&RemoteRepository{
    baseURL: "https://api.example.com",
    client: httpClient,
})
```

**Benefits:**
- Update translations without redeploying
- Centralized management
- A/B testing capabilities

---

## Docker Deployment

### Build Stage

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /build
COPY . .

# Validate translations
RUN go run ./cmd/mbel lint locales/

# Compile translations
RUN go run ./cmd/mbel compile locales/ \
    -o dist/translations.json \
    -sourcemap

# Build binary
RUN go build -o app ./cmd/myapp
```

### Runtime Stage

```dockerfile
FROM alpine:3.18

WORKDIR /app
COPY --from=builder /build/app .
COPY --from=builder /build/dist/translations*.json ./dist/

EXPOSE 8080
CMD ["./app"]
```

### Kubernetes ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: translations
data:
  en.json: |
    {
      "__meta": {"lang": "en"},
      "greeting": "Hello, {name}!"
    }
  de.json: |
    {
      "__meta": {"lang": "de"},
      "greeting": "Hallo, {name}!"
    }
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    spec:
      containers:
      - name: app
        volumeMounts:
        - name: translations
          mountPath: /etc/translations
      volumes:
      - name: translations
        configMap:
          name: translations
```

---

## Performance Optimization

### 1. Enable HTML Escape Only When Needed

```go
// Production: Only escape for web contexts
isWeb := os.Getenv("APP_CONTEXT") == "web"
rt := mbel.NewRuntimeWithOptions(data, isWeb)
```

### 2. Use Lazy-Loading for Large Repos

```go
config := mbel.Config{
    DefaultLocale: "en",
    LazyLoad: true,  // Load runtimes on first use
    HotReload: false, // Disable in production
}
m, _ := mbel.NewManager("./locales", config)
```

### 3. Cache Compiled Translations

```go
// Serve precompiled JSON to frontend
http.HandleFunc("/api/translations/:lang", func(w http.ResponseWriter, r *http.Request) {
    lang := r.PathValue("lang")
    
    // Set cache headers
    w.Header().Set("Cache-Control", "public, max-age=86400") // 1 day
    w.Header().Set("Content-Type", "application/json")
    
    // Stream compiled JSON
    json.NewEncoder(w).Encode(m.Data(lang))
})
```

### 4. Monitoring Metrics

```go
// Log metrics periodically
ticker := time.NewTicker(1 * time.Minute)
go func() {
    for range ticker.C {
        metrics := mbel.GetMetrics()
        log.Printf("Translations: %d gets, %d interpolations",
            metrics["get_calls"],
            metrics["interpolate_ops"])
    }
}()
```

---

## Error Handling & Fallbacks

### Graceful Degradation

```go
func (m *Manager) GetSafe(lang, key string, vars mbel.Vars) string {
    // Try primary language
    if msg, err := m.Get(lang, key, vars); err == nil {
        return msg
    }
    
    // Fallback to default
    if msg, err := m.Get("en", key, vars); err == nil {
        return msg
    }
    
    // Final fallback: return key itself
    return key
}
```

### Logging Translation Failures

```go
func logMissing(key, lang string) {
    log.WithFields(log.Fields{
        "key": key,
        "lang": lang,
        "severity": "warning",
    }).Warning("Translation not found")
    
    // Alert if critical keys missing
    if isCritical(key) {
        alertOn Slack(fmt.Sprintf("Critical translation missing: %s", key))
    }
}
```

---

## Monitoring & Alerting

### Prometheus Metrics

```go
import "github.com/prometheus/client_golang/prometheus"

var (
    translationGets = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "mbel_translation_gets_total",
            Help: "Total number of translation lookups",
        },
    )
    translationMisses = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "mbel_translation_misses_total",
            Help: "Total number of missing translations",
        },
    )
)

prometheus.MustRegister(translationGets, translationMisses)
```

### Health Checks

```go
func healthCheck(m *mbel.Manager) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Check all languages load
        for _, lang := range []string{"en", "de", "fr"} {
            if _, err := m.Get(lang, "health_check_key", nil); err != nil {
                w.WriteHeader(http.StatusServiceUnavailable)
                return
            }
        }
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
    }
}
```

---

## Database Integration

### Store Translations in PostgreSQL

```sql
CREATE TABLE translations (
    id SERIAL PRIMARY KEY,
    language VARCHAR(10) NOT NULL,
    key VARCHAR(255) NOT NULL,
    value TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(language, key)
);
```

```go
type DBRepository struct {
    db *sql.DB
}

func (r *DBRepository) Get(lang string) (mbel.Data, error) {
    rows, _ := r.db.Query(
        "SELECT key, value FROM translations WHERE language = $1",
        lang,
    )
    defer rows.Close()
    
    data := make(map[string]interface{})
    for rows.Next() {
        var key, value string
        rows.Scan(&key, &value)
        data[key] = value
    }
    return data, nil
}
```

---

## Versioning Strategy

### Semantic Versioning for Translations

```json
{
  "__meta": {
    "lang": "en",
    "version": "1.2.0",
    "build_date": "2024-02-01T12:00:00Z",
    "git_commit": "abc123def456"
  }
}
```

### Changelog Management

```markdown
## Translations v1.2.0

### Added
- New greeting messages for onboarding flow
- AI annotations for e-commerce section

### Changed
- Updated error messages for clarity
- Improved tone for success messages

### Fixed
- Missing Polish translation for "delete_confirm"
```

---

## Rollback Procedures

```bash
#!/bin/bash
# rollback.sh

VERSION=$1
BACKUP_DATE=$(date -d "1 day ago" +%Y%m%d)

# Restore from backup
aws s3 cp s3://backups/translations-$BACKUP_DATE.json \
    dist/translations.json

# Validate
go run ./cmd/mbel lint dist/translations.json || exit 1

# Restart service
systemctl restart myapp
```

---

## Load Testing

```bash
# Test high concurrency
ab -c 100 -n 10000 \
    "http://localhost:8080/api/translations/en"

# Monitor under load
watch -n 1 'curl http://localhost:8080/health'
```

---

## Compliance & Auditing

### GDPR-Compliant Access Logging

```go
func logTranslationAccess(userID, lang, key string) {
    log.WithFields(log.Fields{
        "user_id": userID,
        "lang": lang,
        "key": key,
        "timestamp": time.Now().Unix(),
    }).Info("Translation accessed")
}
```

### Encryption at Rest

```bash
# Encrypt compiled JSON before deployment
openssl enc -aes-256-cbc -in dist/translations.json \
    -out dist/translations.json.enc \
    -k "$ENCRYPTION_KEY"
```

