# Architecture & Project Structure

## Overview

MBEL (Multilingual Expression Language) is a sophisticated i18n compiler designed for modern AI-driven applications. This document covers the codebase structure, compilation pipeline, data flow, and extension points for developers.

---

## Directory Structure

```
github.com/makkiattooo/MBEL
├── pkg/mbel/                    # Core library (production code)
│   ├── api.go                   # Public API: Init, T, GlobalT, TDefault, MustT
│   ├── ast.go                   # AST node definitions (Program, Statement, Expression)
│   ├── compiler.go              # AST → Runtime compilation
│   ├── http.go                  # HTTP middleware for locale extraction
│   ├── lexer.go                 # Text → Tokens
│   ├── manager.go               # Locale management, repository interface
│   ├── parser.go                # Tokens → AST
│   ├── plurals.go               # CLDR plural rules for 25+ languages
│   ├── runtime.go               # String resolution & interpolation
│   ├── sourcemap.go             # Source mapping for debugging
│   └── token.go                 # Token type definitions
│
├── cmd/mbel/                    # Command-line interface
│   └── main.go                  # mbel init/compile/lint/fmt/watch commands
│
├── tests/                       # Integration tests
│   ├── *_test.go                # Unit + integration tests
│   └── bench_test.go            # Benchmarks
│
├── examples/                    # Usage examples
│   ├── *.mbel                   # MBEL source files
│   └── server/                  # HTTP server example
│
├── docs/                        # Documentation
│   ├── en/                      # English docs
│   ├── pl/, de/, fr/, ...       # Localized docs
│   └── SUITE.md (per-language)  # Index of available docs
│
├── go.mod                       # Go module definition
├── README.md                    # Project overview
└── CHANGELOG.md                 # Version history
```

**Important**: All core MBEL code lives in `pkg/mbel/` to follow Go conventions. Users import: `github.com/makkiattooo/MBEL/pkg/mbel`

---

## Compilation Pipeline

The MBEL compiler transforms `.mbel` files through a multi-stage pipeline:

```
┌─────────────────────────────────────────────────────────────┐
│                    Input: .mbel file                         │
└────────────────────────────┬────────────────────────────────┘
                             ↓
                ┌────────────────────────┐
                │  LEXER (lexer.go)      │
                │  Text → Tokens         │
                │  Extracts:             │
                │  - Keywords (@lang)    │
                │  - Strings ("...")     │
                │  - Blocks ([...] => ) │
                │  - Comments (# ...)    │
                └────────────┬───────────┘
                             ↓
                ┌────────────────────────┐
                │  PARSER (parser.go)    │
                │  Tokens → AST          │
                │  Builds:               │
                │  - MetadataStatement   │
                │  - AssignStatement     │
                │  - SectionStatement    │
                │  - BlockExpression     │
                │  - AIAnnotation        │
                └────────────┬───────────┘
                             ↓
                ┌────────────────────────┐
                │  COMPILER (compiler.go)│
                │  AST → Runtime Map     │
                │  Resolves:             │
                │  - Pluralization       │
                │  - Nesting (dot.key)   │
                │  - Metadata (__meta)   │
                │  - Terms (__terms)     │
                └────────────┬───────────┘
                             ↓
            ┌────────────────────────────────┐
            │  Output: map[string]interface{}│
            │  {                             │
            │    "key": "value",             │
            │    "count(n)": RuntimeBlock,   │
            │    "section.key": "value",     │
            │    "__meta": {...},            │
            │    "__terms": {...}            │
            │  }                             │
            └────────────┬────────────────────┘
                         ↓
                ┌────────────────────────┐
                │  RUNTIME (runtime.go)  │
                │  Get(key, args)        │
                │  Performs:             │
                │  - Variable interpol.  │
                │  - Plural resolution   │
                │  - HTML escaping       │
                │  - Term substitution   │
                └────────────┬───────────┘
                             ↓
            ┌────────────────────────────┐
            │  Output: Translated String │
            └────────────────────────────┘
```

### Each Component's Role

- **Lexer**: Converts raw MBEL text into tokens (strings, variables, keywords, symbols).
- **Parser**: Groups tokens into logical blocks (assignments, metadata, patterns, AI annotations).
- **Compiler**: Transforms AST into executable Runtime structures with plural rules.
- **Runtime**: Resolves keys, interpolates variables, applies plural logic, handles escaping.
---

## Key Components (Deep Dive)

### 1. Lexer (`pkg/mbel/lexer.go`)
- **Input**: Raw MBEL text
- **Output**: Stream of tokens with line/column info
- **Time Complexity**: O(n) where n = file size
- **Key Methods**:
  - `NewLexer(input)` — Initialize scanner
  - `NextToken()` — Emit next token
  - `readIdentifier()`, `readString()`, `readComment()` — Type-specific scanners

**Token types tracked**:
- Operators: `=`, `=>`, `@`, `.`, `:`, `,`, `..`
- Delimiters: `{`, `}`, `[`, `]`, `(`, `)`
- Literals: `IDENT`, `STRING` (single/triple-quoted), `NUMBER`
- Special: `COMMENT`, `NEWLINE`, `EOF`

**Line/Column tracking**:
```go
// Multi-line strings update lexer.line counter
// Single-line comments update lexer.column
```

### 2. Parser (`pkg/mbel/parser.go`)
- **Input**: Token stream (from Lexer)
- **Output**: Abstract Syntax Tree (AST)
- **Time Complexity**: O(n) with error recovery
- **Key Methods**:
  - `ParseProgram()` — Entry point
  - `parseStatement()` — Top-level constructs
  - `parseBlockCases()` — Plural block parsing
  - `parseAIAnnotation()` — Extract AI metadata

**AST Nodes** (defined in `ast.go`):
- `Program` — Root node containing all statements
- `AssignStatement` — `key = value` or `key(arg) { block }`
- `BlockExpression` — Plural/gender/category rules
- `BlockCase` — Single condition `[category] => "output"`
- `MetadataStatement` — `@lang: pl` declarations
- `AIAnnotation` — Extracted from `# AI_Context:` comments

**Error recovery**:
```go
// If parsing fails, synchronize() skips to next statement boundary
p.synchronize()  // Skip tokens until NEWLINE or next keyword
```

### 3. Compiler (`pkg/mbel/compiler.go`)
- **Input**: AST from Parser
- **Output**: Compiled map (`map[string]interface{}`) ready for Runtime
- **Time Complexity**: O(k) where k = number of keys
- **Key Methods**:
  - `Compile(node)` — Dispatch based on node type
  - `compileProgram()` — Process all statements
  - `compileBlock()` — Create RuntimeBlock with plural rules

**Output structure**:
```go
map[string]interface{}{
  "greeting":           "Hello {name}",
  "items_count":        &RuntimeBlock{...},
  "section.nested":     "Nested value",
  "__meta":             map[string]string{"lang": "en"},
  "__terms":            map[string]string{"app-name": "MyApp"},
  "__ai":               map[string][]map[string]string{...},
  "__imports":          []string{"namespace1", "namespace2"},
}
```

### 4. Runtime (`pkg/mbel/runtime.go`)
- **Input**: Compiled map + key/arguments
- **Output**: Interpolated, pluralized string
- **Time Complexity**: O(1) lookup + O(m) interpolation (m = placeholders)
- **Key Methods**:
  - `Get(key, args...)` — Main resolution method
  - `interpolate(s, arg)` — Replace `{var}` placeholders
  - `ResolveWithLang(arg, lang)` — Apply language-specific plural rules

**Interpolation rules**:
```go
// Matches: {variable_name} or {var123}
// Supports: Vars (named type), map[string]interface{}, scalars
argRe := regexp.MustCompile(`\{([a-zA-Z_][a-zA-Z0-9_]*)\}`)

// Term substitution: {-term-name}
termRe := regexp.MustCompile(`\{-([a-zA-Z_][a-zA-Z0-9_-]*)\}`)
```

**Optional HTML escaping**:
```go
if r.escapeHTML {
    valStr = html.EscapeString(valStr)  // <>&" → entities
}
```

### 5. Manager (`pkg/mbel/manager.go`)
- **Input**: Repository (interface for any data source)
- **Output**: Language-specific Runtime instances
- **Configuration**:
  - `DefaultLocale` — Fallback language (default: "en")
  - `Watch` — Enable hot-reload (file-based only)
  - `LazyLoad` — Load languages on-demand vs. all upfront

**Key Methods**:
- `NewManager(rootPath, cfg)` — File-based initialization
- `NewManagerWithRepo(repo, cfg)` — Custom repository
- `Load()` — Reload all translations
- `Get(lang, key, args...)` — Retrieve + fallback chain
- `watchLoop()` — Poll for file changes (if Watch enabled)

**Fallback chain**:
```
1. Try requested language (e.g., "pl")
2. Try language prefix (e.g., "pl-PL" → "pl")
3. Try default language (e.g., "en")
4. Return key as fallback
```

**Lazy-loading** (if `LazyLoad=true`):
```go
// First time accessing "pl"? Compile on-demand
if _, exists := m.runtimes[lang]; !exists {
    m.runtimes[lang] = NewRuntime(m.allData[lang])
}
```

### 6. Plurals (`pkg/mbel/plurals.go`)
- **Responsibility**: CLDR plural categorization
- **Coverage**: 25+ languages with language-specific rules
- **Format**: Functions returning category (one/few/many/other/etc.)

**Supported languages & rules**:
```
English:   one, other
Polish:    one, few, many
Russian:   one, few, many
Arabic:    zero, one, two, few, many, other
French:    one (0, 1), other
Czech:     one, few, other
...and 19 more
```

**Usage**:
```go
category := mbel.ResolvePluralCategoryExtended("pl", 5)
// Returns: "many" for Polish, 5

cases := &RuntimeBlock.Cases
output := cases[category]  // Look up plural form
```

### 7. API (`pkg/mbel/api.go`)
- **Responsibility**: High-level convenience functions
- **Pattern**: Global manager + context-based locale

**Public functions**:
```go
// Global setup
Init(rootPath, cfg)                     // Initialize global manager

// Context-based translation
T(ctx, key, args...)                    // Main API
TWithLocale(ctx, lang, key, args...)    // Override locale
MustT(ctx, key, args...)                // Panic on miss

// Explicit language
TDefault(key, lang, args...)            // No context needed
GlobalT(key, args...)                   // Use default language

// Context utilities
WithLocale(ctx, lang)                   // Inject locale
LocaleFromContext(ctx)                  // Extract locale

// Metrics
GetMetrics()                            // Access call counters
ResetMetrics()                          // Clear counters
```

### 8. HTTP Middleware (`pkg/mbel/http.go`)
- **Responsibility**: Automatic locale extraction from HTTP requests
- **Method**: Parse `Accept-Language` header

**Usage**:
```go
router.Use(mbel.Middleware)  // Extract locale from request

func handler(w http.ResponseWriter, r *http.Request) {
    locale := mbel.LocaleFromContext(r.Context())
    msg := mbel.T(r.Context(), "greeting")
}
```

---

## Extension Points (Current & Future)

### 1. Custom Repository (Current ✓)

Implement `Repository` interface to load from database, API, or other sources:

```go
type Repository interface {
    LoadAll() (map[string]map[string]interface{}, error)
}

// Example: Database repository
type DatabaseRepository struct {
    db *sql.DB
}

func (r *DatabaseRepository) LoadAll() (map[string]map[string]interface{}, error) {
    results := make(map[string]map[string]interface{})
    
    // Query database for all translations
    rows, _ := r.db.Query("SELECT lang, key, value FROM translations")
    defer rows.Close()
    
    for rows.Next() {
        var lang, key string
        var value interface{}
        rows.Scan(&lang, &key, &value)
        
        if _, ok := results[lang]; !ok {
            results[lang] = make(map[string]interface{})
        }
        results[lang][key] = value
    }
    return results, nil
}

// Initialize with custom repo
m, _ := mbel.NewManagerWithRepo(&DatabaseRepository{db: db}, cfg)
```

### 2. Custom Expression Evaluators (Planned)

Future support for Lua, Jinja2, or custom expression languages:

```go
// Concept (not yet implemented)
// Enable in .mbel files:
msg = "Hello {lua:string.upper(name)}"
msg = "Price: {jinja:price | currency(locale)}"

// To be implemented:
mbel.RegisterExpressionHandler("lua", luaEvaluator)
mbel.RegisterExpressionHandler("jinja", jinjaEvaluator)
```

### 3. Export Formats (Planned)

Convert compiled MBEL to industry-standard formats:

```bash
# Convert to GNU gettext (.po/.pot)
mbel export --format=po locales/ -o translations.po

# Convert to XLIFF (localization standard)
mbel export --format=xliff locales/ -o translations.xliff

# Emit JSON for frontend bundling
mbel export --format=json locales/ -o messages/
```

### 4. Plugin System (Future Vision)

Extensible architecture for custom compilation hooks:

```go
// Future API concept
type Plugin interface {
    Name() string
    OnCompile(program *Program) error
    OnGet(key string, vars Vars) (string, error)
    OnMetrics() map[string]interface{}
}

// Register at startup
mbel.RegisterPlugin(&ValidationPlugin{})
mbel.RegisterPlugin(&CustomTransliterationPlugin{})
```

---

## Performance Characteristics

### Compilation Speed
| Operation | Time | Notes |
|-----------|------|-------|
| Lex 1KB file | ~50µs | Character-by-character scanning |
| Parse 1KB file | ~100µs | Token stream → AST |
| Compile 1KB file | ~50µs | AST → Runtime map |
| **Total 10KB file** | **~5ms** | End-to-end |

### Runtime Speed
| Operation | Time | Memory |
|-----------|------|--------|
| Get (no interp) | ~870 ns | 461 bytes |
| Get + interpolate 1 var | ~1.2 µs | 461 bytes + value |
| Plural resolve + get | ~950 ns | inline |
| LazyLoad 1 language | ~2ms | 50-200 KB |

### Memory Usage
- Per-language (typical): **50-200 KB**
- With lazy-loading: **~1-5 KB at startup**
- Manager overhead: **<10 KB**

### Caching
- **File cache**: mtime-based (checks modification time)
- **Runtime cache**: In-memory, keyed by language
- **Regex cache**: Compiled regexes (termRe, argRe) at module init

---

## Data Flow: Complete Example

**Scenario**: User in German calls `mbel.T(ctx, "user.items_count", 5)`

```
1. T(ctx, "user.items_count", 5)
   └─ Extract locale from context: "de"
   └─ Call Manager.Get("de", "user.items_count", 5)

2. Manager.Get("de", ...)
   ├─ Check runtimes["de"] exists?
   ├─ If lazy-loading and missing:
   │  └─ Compile allData["de"] → new Runtime
   └─ Call Runtime.Get("user.items_count", 5)

3. Runtime.Get("user.items_count", 5)
   ├─ Look up key in Data map
   ├─ Found: &RuntimeBlock
   ├─ Call ResolveWithLang(5, "de")
   │  └─ Apply German plural rule: 5 % 10 = 5 → "other"
   ├─ Look up cases["other"] → "{n} Elemente"
   ├─ Interpolate {n} → "5"
   └─ Return "5 Elemente"

4. Return to user
   └─ Output: "5 Elemente"

[Total time: ~900ns + interpolation ~300ns = ~1.2µs]
```

---

## Building & Testing

### Local Development
```bash
cd /path/to/mbel

# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Benchmark performance
go test -bench=. -benchmem ./tests/

# Build CLI
go build -o mbel ./cmd/mbel

# Format code
go fmt ./...

# Run linter
go vet ./...
```

### Adding New Features

**New token type**:
1. Add constant to `TokenType` in `token.go`
2. Handle in `Lexer.NextToken()`
3. Update parser if needed

**New AST node**:
1. Define struct + implement `Node` interface in `ast.go`
2. Handle in `Parser.parseStatement()` or `parseExpression()`
3. Handle in `Compiler.Compile()`

**New language plural rules**:
1. Add rule function in `plurals.go` (e.g., `pluralFinnish()`)
2. Register in `PluralRules` map with language code
3. Add test case in `tests/parser_test.go`

**Tests**:
```go
// Unit test template
func TestMyFeature(t *testing.T) {
    input := `key = "value"`
    l := mbel.NewLexer(input)
    p := mbel.NewParser(l)
    program := p.ParseProgram()
    
    // Assert expected behavior
    if len(program.Statements) != 1 {
        t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
    }
}
```

---

## Troubleshooting & Debugging

### Issue: Translation not found
```go
// Problem
msg := mbel.T(ctx, "unknown.key")  // Returns: "unknown.key"

// Solution
1. Check filename/structure match: locales/[lang].mbel
2. Verify @lang metadata matches directory
3. Use sourcemap to find original location:
   loc := sourcemap["unknown.key"]
   fmt.Printf("File: %s, Line: %d\n", loc.File, loc.Line)
```

### Issue: Plural not working
```go
// Problem
// count(n) block exists but returns wrong form

// Solution
1. Check language code: ResolvePluralCategoryExtended("pl", 5)
2. Verify case labels: [one], [few], [many] correct for language
3. Check argument passed: Get("count", 5) not Get("count", "5")
```

### Issue: Performance degradation
```go
// Solution
1. Enable lazy-loading: LazyLoad: true in Config
2. Profile: go test -bench=. -cpuprofile=cpu.prof ./tests/
3. Check cache hits: GetMetrics()["cache_hits"]
4. Monitor file watch: Disable Watch mode in production
```

---

## Future Roadmap

- **v1.3**: JavaScript/TypeScript SDK
- **v1.4**: Python SDK  
- **v1.5**: Plugin system
- **v2.0**: Custom expression evaluators (Lua, Jinja)
- **v2.1**: Export formats (PO, XLIFF, JSON)
- **v2.2**: Database integration examples

---

## Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure `go test ./...` and `go vet ./...` pass
5. Update docs/en/ARCHITECTURE.md if adding features
6. Submit a pull request

**Code style**: Follow `gofmt` and `golint` conventions.

