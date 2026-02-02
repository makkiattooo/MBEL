# Sourcemap & Debugging Guide

## Overview

When compiling MBEL files into JSON, it can be difficult to track which translation key comes from which file and line. **Sourcemaps** solve this problem by maintaining a mapping from compiled keys back to their original source locations.

---

## Generating Sourcemaps

Use the `--sourcemap` flag with the `compile` command:

```bash
mbel compile locales/ -o translations.json -sourcemap
```

This generates two files:

1. **`translations.json`** — Compiled translations (as usual)
2. **`translations.sourcemap.json`** — Map of keys → source locations

---

## Sourcemap Format

The sourcemap is a JSON object where each key maps to its source location:

```json
{
  "greeting": {
    "file": "locales/en.mbel",
    "line": 5,
    "column": 0
  },
  "errors.field_required": {
    "file": "locales/en.mbel",
    "line": 12,
    "column": 4
  },
  "user.profile.name": {
    "file": "locales/sections/user.mbel",
    "line": 3,
    "column": 2
  }
}
```

### Structure:
- **file** — Relative path to the `.mbel` source file
- **line** — 1-indexed line number (where the key is defined)
- **column** — 0-indexed column number (indentation level)

---

## Use Cases

### 1. Debugging: Find Which File Defines a Key

```bash
# Compiled app fails with missing key "cart.checkout_btn"
# Look it up in sourcemap:

cat translations.sourcemap.json | jq '.["cart.checkout_btn"]'

# Output:
# {
#   "file": "locales/features/cart.mbel",
#   "line": 42,
#   "column": 2
# }

# Now edit the file at that location
```

### 2. IDE Integration: Quick Navigation

Tools can use sourcemaps to provide "Go to Source" functionality:

```
❌ Translation key "payment.invalid_card" not found
📍 Click → Opens locales/payment.mbel at line 18
```

### 3. Analytics: Track which files are most-used

```bash
# Count occurrences per file
cat translations.sourcemap.json | jq '[.[].file] | group_by(.) | map({file: .[0], count: length})'

# Output:
# [
#   {"file": "locales/common.mbel", "count": 28},
#   {"file": "locales/features/cart.mbel", "count": 15},
#   {"file": "locales/errors.mbel", "count": 12}
# ]
```

### 4. CI/CD: Validate Consistency

```bash
#!/bin/bash
# Ensure every compiled key has a sourcemap entry

compiled_keys=$(jq 'keys | length' translations.json)
sourcemap_keys=$(jq 'keys | length' translations.sourcemap.json)

if [ "$compiled_keys" != "$sourcemap_keys" ]; then
  echo "ERROR: Sourcemap mismatch!"
  exit 1
fi
```

---

## Integration with Vite / Webpack

```javascript
import translations from './translations.json?raw';
import sourcemap from './translations.sourcemap.json?raw';

const sm = JSON.parse(sourcemap);

function findKeySource(key) {
  if (sm[key]) {
    return sm[key]; // {file, line, column}
  }
  return null;
}

// Usage:
console.log(findKeySource('greeting')); 
// → {file: "locales/en.mbel", line: 5, column: 0}
```

---

## CLI Usage

### Generate with Pretty Printing
```bash
mbel compile locales/ -o dist/translations.json -sourcemap -pretty
```

### Generate Compact JSON (for production)
```bash
mbel compile locales/ -o dist/translations.json -sourcemap -pretty=false
```

### Specify Output Path
Currently, sourcemap is always generated as `<output>.sourcemap.json` (derived from `-o` flag).

Example:
```bash
mbel compile locales/ -o ./dist/i18n/messages.json -sourcemap
# Generates:
# - ./dist/i18n/messages.json
# - ./dist/i18n/messages.sourcemap.json
```

---

## Advanced: Custom Processing

Use sourcemaps to build your own tools:

```go
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type SourceLocation struct {
	File   string `json:"file"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
}

func main() {
	// Load sourcemap
	data, _ := ioutil.ReadFile("translations.sourcemap.json")
	var sm map[string]SourceLocation
	json.Unmarshal(data, &sm)

	// Find all keys from a specific file
	targetFile := "locales/user.mbel"
	for key, loc := range sm {
		if loc.File == targetFile {
			log.Printf("Key %s at line %d", key, loc.Line)
		}
	}
}
```

---

## Performance Notes

- Sourcemap generation adds ~5% overhead to compilation time
- Sourcemap file size is typically 5-10% larger than compiled JSON
- Sourcemaps are **not required** for production; they're primarily development tools

---

## Limitations

- Sourcemap only tracks **definition location**, not usage
- Metadata keys (e.g., `__meta`, `__ai`) have sourcemap entries pointing to file headers
- Range cases in plurals show the location of the plural block start, not each range

