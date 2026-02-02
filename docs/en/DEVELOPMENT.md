# Developer's Guide: Contributing to MBEL

This guide covers local development, building, testing, debugging, and extending MBEL. Whether you're fixing bugs or building new features, start here.

---

## Prerequisites

- **Go 1.21+** ([install](https://go.dev/doc/install))
- **Git** for version control
- **Make** (optional, for running commands)
- Text editor or IDE (VS Code, GoLand, etc.)

Verify your setup:
```bash
go version     # Should be 1.21 or higher
git --version
```

---

## Getting the Source Code

### Clone the repository

```bash
git clone https://github.com/makkiattooo/MBEL.git
cd MBEL
```

### Understand the directory structure

```
MBEL/
├── pkg/mbel/              # Core library (all .go files)
│   ├── api.go             # Public API (Init, T, GlobalT, etc.)
│   ├── ast.go             # AST node definitions
│   ├── compiler.go        # Compilation logic
│   ├── lexer.go           # Tokenization
│   ├── parser.go          # Parsing to AST
│   ├── runtime.go         # Execution & interpolation
│   ├── manager.go         # Locale management & caching
│   ├── token.go           # Token types
│   ├── plurals.go         # CLDR plural rules (25+ languages)
│   ├── sourcemap.go       # Source tracking for debugging
│   ├── http.go            # HTTP middleware
│   └── lexer_test.go      # Lexer unit tests
│
├── cmd/mbel/              # CLI tool
│   └── main.go            # Commands: init, compile, lint, fmt, watch, stats, diff, import, translate
│
├── tests/                 # Integration & unit tests
│   ├── api_test.go
│   ├── parser_test.go
│   ├── runtime_test.go
│   ├── bench_test.go
│   └── integration_test.go
│
├── examples/              # Example applications
│   ├── *.mbel             # Example MBEL files
│   └── server/            # HTTP server example
│
├── docs/                  # Documentation (9 languages)
│   ├── en/
│   │   ├── QUICKSTART.md      # Quick start guide
│   │   ├── ARCHITECTURE.md    # Deep technical dive
│   │   ├── DEVELOPMENT.md     # This file
│   │   ├── Manual.md
│   │   ├── FAQ.md
│   │   ├── TIPS.md
│   │   ├── COMPARISON.md
│   │   ├── SECURITY.md
│   │   └── SUITE.md           # Documentation index
│   └── [de/fr/es/it/ja/pl/ru/zh]/...
│
├── go.mod                 # Module definition
├── go.sum                 # Dependency checksums
├── README.md              # Project overview
├── CHANGELOG.md           # Version history
└── .github/workflows/     # CI/CD pipelines
```

---

## Setting Up Development Environment

### 1. Install dependencies

```bash
go mod download
go mod tidy
```

### 2. Verify build

```bash
# Build the CLI tool
go build -o mbel ./cmd/mbel

# Should create executable: ./mbel or ./mbel.exe (Windows)
./mbel --help
```

### 3. Run all tests

```bash
# Run all tests in verbose mode
go test -v ./...

# Run with coverage
go test -cover ./...

# Run specific test package
go test -v ./tests
```

### 4. Set up pre-commit checks

```bash
# Format code before committing
go fmt ./...

# Run linter (requires golangci-lint)
golangci-lint run ./...

# Check for common issues
go vet ./...
```

---

## Development Workflow

### 1. Create a feature branch

```bash
git checkout -b feature/my-feature
# or for bug fixes:
git checkout -b fix/issue-description
```

### 2. Make your changes

Edit files in `pkg/mbel/`, `cmd/mbel/`, or `tests/` as needed.

### 3. Test your changes

```bash
# Run all tests
go test -v ./...

# Run specific test
go test -v -run TestName ./tests

# Run with coverage report
go test -cover ./tests
```

### 4. Format and lint

```bash
go fmt ./...
go vet ./...
```

### 5. Commit and push

```bash
git add .
git commit -m "feat: description of changes"
git push origin feature/my-feature
```

### 6. Create pull request

Visit [GitHub MBEL repository](https://github.com/makkiattooo/MBEL) and create a PR.

---

## Building from Source

### Build the CLI tool

```bash
# Development build
go build -o mbel ./cmd/mbel

# Optimized release build
go build -ldflags="-s -w" -o mbel ./cmd/mbel

# For Windows
go build -o mbel.exe ./cmd/mbel
```

### Build a library binary

```bash
# Create a static library
go build -buildmode=c-archive -o libmbel.a ./pkg/mbel

# Or shared library
go build -buildmode=c-shared -o libmbel.so ./pkg/mbel
```

### Cross-compilation

```bash
# Build for Linux on Windows
GOOS=linux GOARCH=amd64 go build -o mbel-linux ./cmd/mbel

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o mbel-mac ./cmd/mbel

# Build for Windows ARM
GOOS=windows GOARCH=arm64 go build -o mbel-arm.exe ./cmd/mbel
```

---

## Running Tests

### Unit tests

```bash
# Run all unit tests
go test -v ./...

# Run specific package
go test -v ./tests

# Run specific test function
go test -v -run TestLexer ./tests

# Stop on first failure
go test -v -failfast ./...
```

### Integration tests

```bash
# Run only integration tests
go test -v -run TestIntegration ./tests

# With timeout
go test -v -timeout 30s ./tests
```

### Benchmarks

```bash
# Run benchmarks
go test -bench=. -benchmem ./tests

# Run specific benchmark
go test -bench=BenchmarkTranslation -benchmem ./tests

# Profile benchmark
go test -bench=. -benchmem -cpuprofile=cpu.prof ./tests
go tool pprof cpu.prof
```

### Code coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage
go tool cover -html=coverage.out

# Show coverage by function
go tool cover -func=coverage.out
```

---

## Debugging

### Using print statements

```go
package main

import "fmt"

func main() {
    fmt.Printf("Variable value: %v\n", myVar)
    fmt.Printf("Type: %T\n", myVar)
}
```

### Using the Go debugger (Delve)

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug a test
dlv test ./tests

# Debug the CLI
dlv debug ./cmd/mbel

# In debugger:
# (dlv) break main.main
# (dlv) continue
# (dlv) next
# (dlv) step
# (dlv) print myVariable
# (dlv) quit
```

### Debug output in code

```go
// Add temporary debug output
fmt.Fprintf(os.Stderr, "[DEBUG] condition=%v\n", condition)

// Later, remove after debugging
```

### Trace errors

```bash
# Run with race detector (detects data races)
go test -race ./...

# Run with memory sanitizer
go test -msan ./...
```

---

## Common Development Tasks

### Adding a new token type

**File: `pkg/mbel/token.go`**

```go
const (
    // ... existing tokens ...
    TOKEN_NEWPLURALFORM TokenType = "NEWPLURALFORM"
)
```

**File: `pkg/mbel/lexer.go`**

```go
func (l *Lexer) NextToken() Token {
    // ... in switch statement ...
    case '@':
        if l.peekChar() == 'n' && /* check for newpluralform */ {
            return Token{Type: TOKEN_NEWPLURALFORM, Literal: "@newpluralform"}
        }
    // ...
}
```

**File: `tests/lexer_test.go`**

```go
func TestNewToken(t *testing.T) {
    l := NewLexer("@newpluralform")
    tok := l.NextToken()
    if tok.Type != TOKEN_NEWPLURALFORM {
        t.Fatalf("wrong token type")
    }
}
```

### Adding a new AST node

**File: `pkg/mbel/ast.go`**

```go
type NewPluralFormNode struct {
    Token       token.Token // The @newpluralform token
    RuleName    string
    Categories  []string
    TokenLiteral string
}

func (n *NewPluralFormNode) expressionNode() {}
func (n *NewPluralFormNode) TokenLiteral() string { return n.Token.Literal }
```

**File: `pkg/mbel/parser.go`**

```go
func (p *Parser) parseNewPluralForm() Expression {
    // Parse the new plural form syntax
    return &NewPluralFormNode{
        Token:    p.curToken,
        RuleName: /* ... */,
    }
}
```

### Adding a new plural rule

**File: `pkg/mbel/plurals.go`**

```go
// For a new language (e.g., Icelandic)
func pluralIcelandic(n int) string {
    // Icelandic plural rules
    // n = 1, 21, 31, ... => "one"
    // else => "other"
    if (n % 10 == 1 && n % 100 != 11) {
        return "one"
    }
    return "other"
}

func init() {
    // Add to PluralRules map
    PluralRules["is"] = pluralIcelandic
}
```

### Adding a compiler optimization

**File: `pkg/mbel/compiler.go`**

```go
// Optimize constant folding
func (c *Compiler) optimizeExpression(expr Expression) Expression {
    // If expression is constant, evaluate at compile time
    if isConstant(expr) {
        return evaluateConstant(expr)
    }
    return expr
}
```

---

## Performance Optimization

### Profiling

```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=. ./tests
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof -bench=. ./tests
go tool pprof mem.prof

# Goroutine profiling
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

### Benchmarking specific operations

```bash
go test -bench=BenchmarkCompile -benchmem -count=5 ./tests
```

### Memory leak detection

```bash
go test -memprofile=mem.prof ./tests
go tool pprof -alloc_space mem.prof  # All allocations
go tool pprof -alloc_objects mem.prof # Number of allocations
```

---

## Documentation

### Update API documentation

Comments in Go files are used to generate documentation:

```go
// Get retrieves a translation for the given language and key.
// Variables are interpolated into the translation string.
// Plural forms are automatically resolved based on the 'n' variable.
//
// Example:
//   manager.Get("en", "items_count", Vars{"n": 5})
//   // Returns: "You have 5 items"
func (m *Manager) Get(lang string, key string, vars Vars) string {
    // ...
}
```

Generate HTML documentation:
```bash
godoc -http=:6060
# Visit http://localhost:6060/pkg/github.com/makkiattooo/MBEL/pkg/mbel/
```

### Update user documentation

Documentation files are in `docs/en/`, `docs/pl/`, etc.

Edit `.md` files directly. For consistency:
- Use heading levels 1-6 (no more than 6)
- Include code examples with language syntax highlighting
- Add cross-references with `[Text](path/to/file.md)`
- Update `docs/*/SUITE.md` to reflect new docs

---

## Release Process

### Version numbering

MBEL uses Semantic Versioning: `v<major>.<minor>.<patch>`

- **Major**: Breaking changes
- **Minor**: New features (backward compatible)
- **Patch**: Bug fixes

### Creating a release

1. **Update version**
```bash
# Edit files with version number
# cmd/mbel/main.go: const VERSION = "X.Y.Z"
# CHANGELOG.md: Add new version section
```

2. **Tag the release**
```bash
git tag -a v1.2.3 -m "Release version 1.2.3"
git push origin v1.2.3
```

3. **Build binaries**
```bash
./scripts/build-release.sh  # Builds for all platforms
```

4. **Create GitHub Release**
- Go to [Releases](https://github.com/makkiattooo/MBEL/releases)
- Create new release from tag
- Upload binaries
- Add release notes from CHANGELOG.md

---

## Continuous Integration

MBEL uses GitHub Actions for CI/CD. See `.github/workflows/`:

### Run locally

```bash
# Install act (GitHub Actions locally)
# https://github.com/nektos/act

# Run all workflows
act

# Run specific workflow
act -j test
```

### CI checks performed

1. **Build**: `go build ./...`
2. **Tests**: `go test -v ./...`
3. **Coverage**: Must maintain >80% coverage
4. **Lint**: `golangci-lint run ./...`
5. **Format**: `go fmt ./...`
6. **Vet**: `go vet ./...`

---

## Common Issues

### Issue: Tests fail with "import not found"

**Solution:**
```bash
go mod tidy
go mod download
go test ./...
```

### Issue: "pkg/mbel" import not recognized

**Solution:** Make sure you're using the correct import path:
```go
import mbel "github.com/makkiattooo/MBEL/pkg/mbel"
```

### Issue: Benchmarks show high memory usage

**Solution:**
```bash
# Profile memory
go test -memprofile=mem.prof -bench=. ./tests

# Analyze
go tool pprof -alloc_space mem.prof
# Check for:
# - Unnecessary allocations
# - Memory leaks
# - Large object creation in hot paths
```

### Issue: Slow tests

**Solution:**
```bash
# Identify slow tests
go test -v -timeout 5s ./tests

# Profile specific test
go test -cpuprofile=cpu.prof -bench=TestName ./tests
go tool pprof cpu.prof
```

---

## Project Structure Philosophy

### Why `pkg/mbel/`?

The `pkg/mbel/` directory structure follows Go conventions:
- **Importable**: Other projects import `github.com/makkiattooo/MBEL/pkg/mbel`
- **Organized**: Separates library code from tools and examples
- **Scalable**: Easy to add more packages (`pkg/utils/`, `pkg/plugins/`, etc.)

### Module boundaries

**`pkg/mbel/`** — Core library (lexer, parser, compiler, runtime, manager)
- Fast, lightweight
- No external dependencies
- Fully testable in isolation

**`cmd/mbel/`** — Command-line tool
- Uses `pkg/mbel`
- Adds file I/O, CLI UX
- Separate build artifact

**`tests/`** — Test suite
- Tests all packages
- Integration tests
- Benchmarks

**`examples/`** — Reference implementations
- Show usage patterns
- Quick start templates
- Real-world scenarios

---

## Future Development

### Planned features

See [ARCHITECTURE.md](ARCHITECTURE.md#future-roadmap) for:
- JS/TypeScript SDK (v1.3)
- Python SDK (v1.4)
- Plugin system (v1.5)
- Custom expression evaluators (v2.0)

### How to contribute

1. Check [open issues](https://github.com/makkiattooo/MBEL/issues)
2. Discuss in issue before starting
3. Fork, branch, implement, test
4. Create descriptive pull request
5. Request review from maintainers

### Code of Conduct

Be respectful, constructive, and inclusive. All contributors are valued.

---

## Getting Help

- **[ARCHITECTURE.md](ARCHITECTURE.md)** — Technical deep dive
- **[QUICKSTART.md](QUICKSTART.md)** — Usage examples
- **[FAQ.md](FAQ.md)** — Common questions
- **[Issues](https://github.com/makkiattooo/MBEL/issues)** — Bug reports, features
- **[Discussions](https://github.com/makkiattooo/MBEL/discussions)** — General questions

---

## Next Steps

1. Run the test suite: `go test -v ./...`
2. Build the CLI: `go build -o mbel ./cmd/mbel`
3. Try the examples: `cd examples && go run server/main.go`
4. Pick an issue and start contributing!
5. Read [ARCHITECTURE.md](ARCHITECTURE.md) for deep technical knowledge

Happy coding! 🎉
