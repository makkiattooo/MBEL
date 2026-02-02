package mbel

import "strings"

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
	"nb": pluralEnglish,
	"nn": pluralEnglish,

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
	"be": pluralRussian,

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
	"ms": pluralAsian,

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
