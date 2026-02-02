package mbel

import (
	"fmt"
	"html"
	"regexp"
	"strconv"
)

var (
	termRe = regexp.MustCompile(`\{-([a-zA-Z_][a-zA-Z0-9_-]*)\}`)
	argRe  = regexp.MustCompile(`\{([a-zA-Z_][a-zA-Z0-9_]*)\}`)
)

// Runtime provides string resolution with interpolation
type Runtime struct {
	Data       map[string]interface{}
	Terms      map[string]string
	Language   string
	escapeHTML bool // Enable HTML escaping for interpolated values
}

// NewRuntime creates a runtime from compiled data
func NewRuntime(data map[string]interface{}) *Runtime {
	return NewRuntimeWithOptions(data, false)
}

// NewRuntimeWithOptions creates a runtime with custom options
func NewRuntimeWithOptions(data map[string]interface{}, escapeHTML bool) *Runtime {
	r := &Runtime{
		Data:       data,
		Terms:      make(map[string]string),
		Language:   "en",
		escapeHTML: escapeHTML,
	}

	// Extract terms
	if terms, ok := data["__terms"].(map[string]string); ok {
		r.Terms = terms
	}

	// Extract language from metadata
	if meta, ok := data["__meta"].(map[string]string); ok {
		if lang, exists := meta["lang"]; exists {
			r.Language = lang
		}
	}

	return r
}

// EscapeHTML enables or disables HTML escaping for interpolated values
func (r *Runtime) SetEscapeHTML(escape bool) {
	r.escapeHTML = escape
}

// Get retrieves a string by key with optional arguments for interpolation
func (r *Runtime) Get(key string, args ...interface{}) string {
	recordGetCall()

	val, exists := r.Data[key]
	if !exists {
		return key // Fallback to key itself
	}

	switch v := val.(type) {
	case string:
		if len(args) > 0 {
			return r.interpolate(v, args[0])
		}
		return r.interpolate(v, nil)
	case *RuntimeBlock:
		if len(args) > 0 {
			result := v.ResolveWithLang(args[0], r.Language)
			return r.interpolate(result, args[0])
		}
		return r.interpolate(v.Resolve("other"), nil)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// interpolate replaces {placeholders} and {-term-refs}
func (r *Runtime) interpolate(s string, arg interface{}) string {
	recordInterpolate()

	if s == "" {
		return s
	}

	// Replace term references {-term-name}
	s = termRe.ReplaceAllStringFunc(s, func(match string) string {
		termName := match[2 : len(match)-1] // Extract "term-name" from "{-term-name}"
		if val, exists := r.Terms[termName]; exists {
			return val
		}
		return match // Keep original if not found
	})

	// Replace argument placeholder {n}, {count}, etc.
	if arg != nil {
		s = argRe.ReplaceAllStringFunc(s, func(match string) string {
			key := match[1 : len(match)-1]

			// Accept both named type Vars and raw map[string]interface{}
			switch m := arg.(type) {
			case Vars:
				if val, exists := m[key]; exists {
					valStr := fmt.Sprintf("%v", val)
					if r.escapeHTML {
						valStr = html.EscapeString(valStr)
					}
					return valStr
				}
				return match // Keep {placeholder} if not found in map
			case map[string]interface{}:
				if val, exists := m[key]; exists {
					valStr := fmt.Sprintf("%v", val)
					if r.escapeHTML {
						valStr = html.EscapeString(valStr)
					}
					return valStr
				}
				return match // Keep {placeholder} if not found in map
			default:
				// Scalar (primitive): replace all placeholders with this value
				valStr := fmt.Sprintf("%v", arg)
				if r.escapeHTML {
					valStr = html.EscapeString(valStr)
				}
				return valStr
			}
		})
	}

	return s
}

// ResolveWithLang finds the matching value using language-specific plural rules
func (rb *RuntimeBlock) ResolveWithLang(arg interface{}, lang string) string {
	valToMatch := arg

	// If argument is a map, try to extract the specific argument for this block
	if m, ok := arg.(map[string]interface{}); ok && rb.Argument != "" {
		if v, exists := m[rb.Argument]; exists {
			valToMatch = v
		}
	}

	// Try string match first
	if strArg, ok := valToMatch.(string); ok {
		if val, exists := rb.Cases[strArg]; exists {
			return val
		}
	}

	// Try numeric match
	var numArg int
	switch v := valToMatch.(type) {
	case int:
		numArg = v
	case int64:
		numArg = int(v)
	case float64:
		numArg = int(v)
	default:
		// Fall through to other
	}

	// Check exact number match
	numStr := strconv.Itoa(numArg)
	if val, exists := rb.Cases[numStr]; exists {
		return val
	}

	// Check range matches
	for _, rc := range rb.RangeCases {
		if numArg >= rc.Start && numArg <= rc.End {
			return rc.Value
		}
	}

	// Check plural categories with language
	pluralCat := ResolvePluralCategory(lang, numArg)
	if val, exists := rb.Cases[pluralCat]; exists {
		return val
	}

	// Fallback to "other"
	if val, exists := rb.Cases["other"]; exists {
		return val
	}

	return ""
}
