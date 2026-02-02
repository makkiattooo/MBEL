# Inicio rápido: Su primera aplicación MBEL

Inicie una aplicación multilingüe lista para producción con MBEL en **15 minutos**. Esta guía cubre la configuración del proyecto, archivos de traducción, compilación e implementación.

---

## 1. Inicialice su proyecto

```bash
# Cree el directorio del proyecto
mkdir -p hello-mbel && cd hello-mbel

# Inicialice el módulo Go
go mod init hello-mbel
go get github.com/makkiattooo/MBEL@latest

# Cree la estructura del proyecto
mkdir -p locales cmd dist

# Instale la CLI de MBEL
go install github.com/makkiattooo/MBEL/cmd/mbel@latest
```

---

## 2. Escriba archivos de traducción

### Cree traducciones al español (`locales/es.mbel`)

```mbel
@namespace: hello
@lang: es

app_name = "Hola MBEL"
app_version = "1.0.0"

greeting = "¡Bienvenido, {name}!"
goodbye = "¡Hasta luego, {name}! Tiempo: {time}"

# Español: [one] (1), [other]
items_count(n) {
    [one]   => "Tienes 1 artículo"
    [other] => "Tienes {n} artículos"
}

profile_updated(gender) {
    [male]   => "Actualizó su perfil"
    [female] => "Actualizó su perfil"
    [other]  => "Actualizaron su perfil"
}

ui.menu {
    home = "Inicio"
    about = "Acerca de"
    contact = "Contacto"
    settings = "Configuración"
}

order_total = "Total: {price} (impuestos incluidos)"
```

### Agregue otros idiomas

Cree `locales/en.mbel`, `locales/de.mbel`, etc.

---

## 3. Valide y formatee traducciones

```bash
# Verifique errores de sintaxis
mbel lint locales/

# Formatee automáticamente
mbel fmt locales/

# Muestre estadísticas
mbel stats locales/
```

---

## 4. Compile las traducciones

```bash
# Compile en un archivo JSON
mbel compile locales/ -o dist/translations.json

# Incluya mapa de fuente para depuración
mbel compile locales/ -o dist/translations.json -sourcemap
```

---

## 5. Cree su aplicación Go

### Ejemplo básico (`cmd/main.go`)

```go
package main

import (
	"context"
	"fmt"
	"log"

	mbel "github.com/makkiattooo/MBEL/pkg/mbel"
)

func main() {
	// 1. Inicialice el Manager
	m, err := mbel.NewManager("./dist/translations.json", mbel.Config{
		DefaultLocale: "es",
		FallbackChain: []string{"es", "en"},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 2. Búsqueda simple
	fmt.Println(m.Get("es", "app_name", nil))
	fmt.Println(m.Get("en", "app_name", nil))

	// 3. Interpolación de variables
	vars := mbel.Vars{"name": "María"}
	greeting := m.Get("es", "greeting", vars)
	fmt.Println(greeting)

	// 4. Pluralización
	for _, count := range []int{1, 2, 5} {
		vars := mbel.Vars{"n": count}
		msg := m.Get("es", "items_count", vars)
		fmt.Printf("n=%d: %s\n", count, msg)
	}

	// 5. API global
	mbel.Init(m)
	ctx := context.Background()
	fmt.Println(mbel.T(ctx, "greeting", mbel.Vars{"name": "Juan"}))
}
```

---

## 6. Agregue pruebas

```bash
go test -v ./...
```

---

## 7. Construya y ejecute

```bash
go build -o hello-mbel ./cmd
./hello-mbel
```

---

## 8. Implementación en producción

### Opción A: Incruste JSON

```go
//go:embed dist/translations.json
var translationsJSON []byte
```

### Opción B: Envíe con Docker

```dockerfile
FROM golang:1.21 AS builder
WORKDIR /src
COPY . .
RUN go build -o /tmp/hello-mbel ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /tmp/hello-mbel .
COPY dist/translations.json ./dist/
CMD ["./hello-mbel"]
```

---

## 9. Patrones comunes

### Negociación de idioma

```go
lang := r.Header.Get("Accept-Language")
```

### Almacenamiento en caché y rendimiento

```go
// El Manager almacena en caché automáticamente
for i := 0; i < 1000000; i++ {
	msg := m.Get("es", "app_name", nil)
}
```

---

## 10. Próximos pasos

1. **[Manual](Manual.md)** — Documentación completa
2. **[ARCHITECTURE.md](ARCHITECTURE.md)** — Análisis técnico profundo
3. **[DEVELOPMENT.md](DEVELOPMENT.md)** — Extender MBEL
4. **[Mejores prácticas de seguridad](SECURITY.md)** — Prevención XSS

---

## Solución de problemas

| Problema | Solución |
|----------|----------|
| `no such file or directory: locales/` | `mkdir -p locales` |
| `Undefined: mbel` | Verifique `go.mod` |
| `Syntax error at line 5` | `mbel lint locales/` |
| Traducción no encontrada | Verifique clave, idioma y cadena de retorno |
