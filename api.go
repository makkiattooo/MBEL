package mbel

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

// Global instance
var std *Manager

// Vars is a shortcut for map[string]interface{}, useful for template interpolation
type Vars map[string]interface{}

// Metrics holds basic telemetry for runtime
type Metrics struct {
	mu             sync.RWMutex
	GetCalls       int64 // Total calls to Get
	InterpolateOps int64 // Total interpolation operations
	CacheHits      int64 // Cache hits (future)
	CacheMisses    int64 // Cache misses (future)
}

// Global metrics instance
var metrics = &Metrics{}

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

// ============================================================================
// CONVENIENCE HELPERS
// ============================================================================

// TDefault translates using an explicit language and default locale fallback
func TDefault(key, lang string, args ...interface{}) string {
	if std == nil {
		return key
	}
	return std.Get(lang, key, args...)
}

// MustT translates and panics if key is not found (for critical strings)
func MustT(ctx context.Context, key string, args ...interface{}) string {
	result := T(ctx, key, args...)
	if result == key {
		panic(fmt.Sprintf("translation not found for key: %s", key))
	}
	return result
}

// TWithLocale is a convenience wrapper for T with explicit context and locale override
func TWithLocale(ctx context.Context, lang, key string, args ...interface{}) string {
	if std == nil {
		return key
	}
	return std.Get(lang, key, args...)
}

// ============================================================================
// METRICS
// ============================================================================

// RecordGetCall increments the Get call counter (called internally by Runtime.Get)
func recordGetCall() {
	atomic.AddInt64(&metrics.GetCalls, 1)
}

// RecordInterpolate increments interpolation counter
func recordInterpolate() {
	atomic.AddInt64(&metrics.InterpolateOps, 1)
}

// GetMetrics returns a copy of current metrics
func GetMetrics() map[string]int64 {
	return map[string]int64{
		"get_calls":       atomic.LoadInt64(&metrics.GetCalls),
		"interpolate_ops": atomic.LoadInt64(&metrics.InterpolateOps),
		"cache_hits":      atomic.LoadInt64(&metrics.CacheHits),
		"cache_misses":    atomic.LoadInt64(&metrics.CacheMisses),
	}
}

// ResetMetrics clears all metrics counters
func ResetMetrics() {
	atomic.StoreInt64(&metrics.GetCalls, 0)
	atomic.StoreInt64(&metrics.InterpolateOps, 0)
	atomic.StoreInt64(&metrics.CacheHits, 0)
	atomic.StoreInt64(&metrics.CacheMisses, 0)
}
