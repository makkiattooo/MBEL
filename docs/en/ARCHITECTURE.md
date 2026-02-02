# Architecture & Extension Points

## Core Architecture

```
Lexer (lexer.go)
  ↓ Tokens
Parser (parser.go)
  ↓ AST (Abstract Syntax Tree)
Compiler (compiler.go)
  ↓ Runtime (compiled bytecode structure)
Runtime (runtime.go)
  ↓ Translated strings
```

### Each Component's Role

**Lexer**: Converts raw MBEL text into tokens (strings, variables, keywords, symbols).
**Parser**: Groups tokens into logical blocks (assignments, metadata, patterns).
**Compiler**: Transforms AST into executable Runtime structures.
**Runtime**: Resolves keys, interpolates variables, handles pluralization.

---

## Extension Points (Current & Future)

### 1. Custom Repository (Current)

Implement the `Repository` interface to load translations from anywhere:

```go
type Repository interface {
    Get(lang string) (Data, error)
}

type MyDatabaseRepository struct {
    db *sql.DB
}

func (r *MyDatabaseRepository) Get(lang string) (mbel.Data, error) {
    // Load from database
    return compiledData, nil
}

// Use it
m, _ := mbel.NewManagerWithRepository(myDBRepo)
```

---

### 2. Custom Expression Evaluators (Planned)

Future support for custom expression languages in place of simple `{var}` substitution:

```go
// Concept (not yet implemented)
mbel.RegisterExpressionHandler("lua", luaEvaluator)
mbel.RegisterExpressionHandler("jinja", jinjaEvaluator)

// Then in .mbel files:
msg = "Hello {lua:string.upper(name)}"
```

---

### 3. Export Formats (Planned)

Convert compiled MBEL to other formats for external tools:

```bash
mbel export --format=po locales/ -o translations.po    # GNU gettext
mbel export --format=xliff locales/ -o translations.xliff
mbel export --format=json locales/ -o messages/
```

---

### 4. Plugin System (Future Vision)

Extensible plugin architecture for custom:
- Transliteration engines
- Formatting rules
- Validation hooks
- Custom plural categories

```go
// Future API concept
type Plugin interface {
    Name() string
    Init(mbel.Manager) error
    OnCompile(lang string, data mbel.Data) error
    OnGet(key string, vars mbel.Vars) (string, error)
}

mbel.RegisterPlugin(&MyTransliterationPlugin{})
```

---

## Data Flow Diagram

```
User Code
   ↓
mbel.Get("user.message", vars) ← Entry point
   ↓
Manager.Get(lang, key)
   ↓
Repository.Get(lang) → Cache check
   ↓
Compile on miss
   ↓
Runtime.Get(key, vars)
   ↓
Resolve default lang → Fallback chain
   ↓
Interpolate values
   ↓
Pluralize (if needed)
   ↓
Return string
```

---

## Compiler Optimization Opportunities

### Current Optimizations:
- **Regex Compilation Caching** (termRe, argRe compile once)
- **File Caching** (mtime-based invalidation in FileRepository)
- **Lazy Runtime Loading** (Manager.LazyLoad for memory efficiency)

### Future Optimization Ideas:
1. **JIT Compilation**: Pre-compile hot keys during first month of operation
2. **String Interning**: Reuse identical string objects across runtimes
3. **Pattern Specialization**: Detect simple keys and use fast path (no interpolation)
4. **Parallel Compilation**: Use goroutines to compile multiple languages simultaneously

---

## Sourcemap & Debugging

The `sourcemap.go` module tracks original file locations for compiled keys:

```go
sm := mbel.BuildSourceMap(data) // Build index
loc := sm["user.message"]
// loc.File: "en.mbel"
// loc.Line: 42
// loc.Column: 8
```

### Debugging Workflow:
```bash
1. Compile with --sourcemap flag
2. Get translation key that failed
3. Look up location in sourcemap.json
4. Edit source file at that line
5. Recompile and test
```

---

## Performance Characteristics

| Operation | Time | Memory |
|-----------|------|--------|
| Load 1 language | ~5-10ms | 50-200KB |
| Lazy-load startup | ~1ms | <10KB |
| Runtime.Get (cached) | ~870ns | 461B (per call) |
| Compile large file | ~300µs | ~200KB (temp) |
| Pluralization resolve | ~50ns | inline |

---

## Internal State Management

### Manager
- **Cached Runtimes**: `data map[string]*Runtime` (eager) or `allData map[string]Data` (lazy)
- **Hot Reload**: Watches file system, reloads on change
- **Metrics**: Global atomic counters for Get calls and interpolations

### Runtime
- **Key Index**: Map of keys → values (resolved at compile time)
- **Plurals Map**: Language-specific plural rule functions
- **Escape Mode**: Boolean flag for HTML escaping

### Compiler
- **Pattern Extraction**: Finds all `{var}` and `{name|pattern}` placeholders
- **Plural Compilation**: Groups plural blocks, validates rules
- **Metadata Preservation**: Stores AI annotations for external tools

---

## Future: SDK for Other Languages

### JavaScript/TypeScript
```typescript
import { Manager } from "@mbel/js";

const m = new Manager("./locales");
const msg = await m.get("en", "user.message", { name: "Alice" });
```

### Python
```python
from mbel import Manager

m = Manager("./locales")
msg = m.get("en", "user.message", {"name": "Alice"})
```

### Rust
```rust
use mbel::{Manager, Vars};

let m = Manager::new("./locales")?;
let msg = m.get("en", "user.message", vars!{"name" => "Alice"})?;
```

These would share the compiled `.mbel.json` format for consistency.

