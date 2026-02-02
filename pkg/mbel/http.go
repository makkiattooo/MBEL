package mbel

import (
	"net/http"
	"strings"
)

// Middleware automatically extracts the locale from the request
// checks Accept-Language header and injects it into the Context
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simple Accept-Language parser (grabs first preference)
		lang := "en" // default fallback if unconfigured

		accept := r.Header.Get("Accept-Language")
		if accept != "" {
			// e.g. "pl-PL,pl;q=0.9,en-US;q=0.8" -> "pl-PL"
			parts := strings.Split(accept, ",")
			if len(parts) > 0 {
				first := strings.TrimSpace(parts[0])
				// remove quality score if present (though conventionally first component doesn't have it)
				semicolon := strings.Index(first, ";")
				if semicolon != -1 {
					first = first[:semicolon]
				}
				lang = first
			}
		}

		// Inject into context
		ctx := WithLocale(r.Context(), lang)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// HandlerFunc wrapper for convenience
func Handler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Middleware(next).ServeHTTP(w, r)
	}
}
