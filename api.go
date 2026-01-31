package mbel

import "context"

// Global instance
var std *Manager

// Vars is a shortcut for map[string]interface{}, useful for template interpolation
type Vars map[string]interface{}

// Init initializes the global MBEL manager
// Call this at the start of your application
func Init(rootPath string, cfg Config) error {
	m, err := NewManager(rootPath, cfg)
	if err != nil {
		return err
	}
	std = m
	return nil
}

// GlobalT translates a key using the global manager and default language.
// Useful for system messages or when context is not available.
func GlobalT(key string, args ...interface{}) string {
	if std == nil {
		return key
	}
	return std.Get(std.defaultLang, key, args...)
}

// T translates a key using the locale found in context
// This is the primary API for localized applications
func T(ctx context.Context, key string, args ...interface{}) string {
	if std == nil {
		return key
	}
	lang := LocaleFromContext(ctx)
	return std.Get(lang, key, args...)
}

// Context handling

type contextKey struct{}

// WithLocale injects the locale code into the context
func WithLocale(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, contextKey{}, lang)
}

// LocaleFromContext retrieves the locale code from the context
// Returns default locale if not found
func LocaleFromContext(ctx context.Context) string {
	if val, ok := ctx.Value(contextKey{}).(string); ok {
		return val
	}
	if std != nil {
		return std.defaultLang
	}
	return "en"
}
