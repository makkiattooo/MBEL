# MBEL Go SDK Reference

The `github.com/makkiattooo/mbel` package provides correct runtime resolution for MBEL files.

## 1. Manager & Initialization

### `mbel.Init(path string, cfg mbel.Config)`
Initializes the global singleton.
*   `path`: Directory containing `.mbel` files.
*   `Config.Watch`: If true, enables hot-reload (polling).

### `mbel.NewManagerWithRepo(repo Repository, cfg Config)`
Create a manager with a custom data source (e.g. Database).

```go
type Repository interface {
    LoadAll() (map[string]map[string]interface{}, error)
}
```

## 2. Translation

### `mbel.T(ctx context.Context, key string, args ...interface{})`
The main function to get a string.
*   Uses `LocaleFromContext` to determine language.
*   Falls back to `DefaultLocale` if missing.
*   Falls back to `key` string if translation missing.

### Arguments
*   **Simple value**: `T(ctx, "key", "value")` -> Replaces `{n}` or checks conditions against this value.
*   **Named variables**: `T(ctx, "key", mbel.Vars{"name": "X", "gender": "Y"})` -> Supports complex interpolation and logic.

## 3. Middleware

### `mbel.Middleware(next http.Handler)`
Automatically parses `Accept-Language` header from HTTP requests and injects the best matching locale into `r.Context()`.

## 4. Other Helpers

*   `mbel.WithLocale(ctx, lang)`: Manually set locale in context.
*   `mbel.LocaleFromContext(ctx)`: Get current locale.
*   `mbel.GlobalT(key)`: Translate using default locale (no context).
