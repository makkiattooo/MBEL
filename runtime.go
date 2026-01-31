package mbel

import (
	"fmt"
	"regexp"
	"strings"
)

// Runtime provides string resolution with interpolation
type Runtime struct {
	Data     map[string]interface{}
	Terms    map[string]string
	Language string
}

// NewRuntime creates a runtime from compiled data
func NewRuntime(data map[string]interface{}) *Runtime {
	r := &Runtime{
		Data:     data,
		Terms:    make(map[string]string),
		Language: "en",
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

// Get retrieves a string by key with optional arguments for interpolation
func (r *Runtime) Get(key string, args ...interface{}) string {
	val, exists := r.Data[key]
	if !exists {
		return key // Fallback to key itself
	}

	switch v := val.(type) {
	case string:
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
	if s == "" {
		return s
	}

	// Replace term references {-term-name}
	termRe := regexp.MustCompile(`\{-([a-zA-Z_][a-zA-Z0-9_-]*)\}`)
	s = termRe.ReplaceAllStringFunc(s, func(match string) string {
		termName := match[2 : len(match)-1] // Extract "term-name" from "{-term-name}"
		if val, exists := r.Terms[termName]; exists {
			return val
		}
		return match // Keep original if not found
	})

	// Replace argument placeholder {n}, {count}, etc.
	if arg != nil {
		argRe := regexp.MustCompile(`\{([a-zA-Z_][a-zA-Z0-9_]*)\}`)
		s = argRe.ReplaceAllStringFunc(s, func(match string) string {
			key := match[1 : len(match)-1]

			// Case 1: Argument is a map (named parameters)
			if m, ok := arg.(map[string]interface{}); ok {
				if val, exists := m[key]; exists {
					return fmt.Sprintf("%v", val)
				}
				return match // Keep {placeholder} if not found in map
			}

			// Case 2: Argument is scalar (primitive)
			// Replace all placeholders with this value (e.g. for plural count)
			return fmt.Sprintf("%v", arg)
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
	numStr := fmt.Sprintf("%d", numArg)
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

// ============================================================================
// EXTENDED PLURAL RULES (CLDR)
// ============================================================================

// PluralRule represents a language's plural categorization function
type PluralRule func(n int) string

// PluralRules maps language codes to plural rule functions
var PluralRules = map[string]PluralRule{
	// Germanic languages
	"en": pluralEnglish,
	"de": pluralEnglish,
	"nl": pluralEnglish,
	"sv": pluralEnglish,
	"da": pluralEnglish,
	"no": pluralEnglish,
	"nb": pluralEnglish, // Norwegian Bokmål
	"nn": pluralEnglish, // Norwegian Nynorsk

	// Romance languages
	"fr": pluralFrench,
	"es": pluralEnglish,
	"it": pluralEnglish,
	"pt": pluralEnglish,

	// Slavic languages
	"pl": pluralPolish,
	"ru": pluralRussian,
	"uk": pluralRussian,
	"cs": pluralCzech,
	"sk": pluralCzech,
	"hr": pluralRussian,
	"sr": pluralRussian,
	"be": pluralRussian, // Belarusian

	// Other European
	"ro": pluralRomanian,
	"lt": pluralLithuanian,

	// Asian languages (no plural forms)
	"zh": pluralAsian,
	"ja": pluralAsian,
	"ko": pluralAsian,
	"vi": pluralAsian,
	"th": pluralAsian,
	"id": pluralAsian,
	"ms": pluralAsian, // Malay

	// Semitic
	"ar": pluralArabic,
	"he": pluralEnglish,

	// Other
	"tr": pluralEnglish,
	"hu": pluralEnglish,
	"fi": pluralEnglish,
}

// English: one, other
func pluralEnglish(n int) string {
	if n == 1 {
		return "one"
	}
	return "other"
}

// French: one (0, 1), other
func pluralFrench(n int) string {
	if n == 0 || n == 1 {
		return "one"
	}
	return "other"
}

// Polish: one, few, many
func pluralPolish(n int) string {
	if n == 1 {
		return "one"
	}
	mod10 := n % 10
	mod100 := n % 100
	if mod10 >= 2 && mod10 <= 4 && !(mod100 >= 12 && mod100 <= 14) {
		return "few"
	}
	return "many"
}

// Russian/Ukrainian: one, few, many
func pluralRussian(n int) string {
	mod10 := n % 10
	mod100 := n % 100
	if mod10 == 1 && mod100 != 11 {
		return "one"
	}
	if mod10 >= 2 && mod10 <= 4 && !(mod100 >= 12 && mod100 <= 14) {
		return "few"
	}
	return "many"
}

// Czech/Slovak: one, few, other
func pluralCzech(n int) string {
	if n == 1 {
		return "one"
	}
	if n >= 2 && n <= 4 {
		return "few"
	}
	return "other"
}

// Romanian: one, few, other
func pluralRomanian(n int) string {
	if n == 1 {
		return "one"
	}
	if n == 0 || (n%100 >= 1 && n%100 <= 19) {
		return "few"
	}
	return "other"
}

// Lithuanian: one, few, other
func pluralLithuanian(n int) string {
	mod10 := n % 10
	mod100 := n % 100
	if mod10 == 1 && mod100 != 11 {
		return "one"
	}
	if mod10 >= 2 && mod10 <= 9 && !(mod100 >= 11 && mod100 <= 19) {
		return "few"
	}
	return "other"
}

// Arabic: zero, one, two, few, many, other
func pluralArabic(n int) string {
	if n == 0 {
		return "zero"
	}
	if n == 1 {
		return "one"
	}
	if n == 2 {
		return "two"
	}
	mod100 := n % 100
	if mod100 >= 3 && mod100 <= 10 {
		return "few"
	}
	if mod100 >= 11 && mod100 <= 99 {
		return "many"
	}
	return "other"
}

// Asian languages: other only (no plural forms)
func pluralAsian(n int) string {
	return "other"
}

// ResolvePluralCategoryExtended uses the extended plural rules
func ResolvePluralCategoryExtended(lang string, n int) string {
	// Normalize language code (take first 2 chars)
	if len(lang) > 2 {
		lang = strings.ToLower(lang[:2])
	} else {
		lang = strings.ToLower(lang)
	}

	if rule, exists := PluralRules[lang]; exists {
		return rule(n)
	}

	// Default to English rules
	return pluralEnglish(n)
}

// ============================================================================
// SOURCE MAPPING
// ============================================================================

// SourceLocation represents a position in a source file
type SourceLocation struct {
	File   string
	Line   int
	Column int
}

// SourceMap maps keys to their source locations
type SourceMap map[string]SourceLocation

// BuildSourceMap creates a source map from a parsed program
func BuildSourceMap(p *Program, filename string) SourceMap {
	sm := make(SourceMap)

	for _, stmt := range p.Statements {
		switch s := stmt.(type) {
		case *AssignStatement:
			sm[s.Name] = SourceLocation{
				File:   filename,
				Line:   s.Token.Line,
				Column: s.Token.Column,
			}
		case *MetadataStatement:
			sm["@"+s.Key] = SourceLocation{
				File:   filename,
				Line:   s.Token.Line,
				Column: s.Token.Column,
			}
		}
	}

	return sm
}
