# MBEL: Manual de Referencia Completo

**Versión:** 1.2.0
**Fecha:** Enero 2026

---

## 📖 Tabla de Contenidos

1.  [Introducción](#1-introducción)
2.  [El Lenguaje MBEL](#2-el-lenguaje-mbel)
    *   [Estructura del Archivo](#21-estructura-del-archivo)
    *   [Tipos de Datos](#22-tipos-de-datos)
    *   [Interpolación y Variables](#23-interpolación-y-variables)
    *   [Lógica y Control](#24-lógica-y-control)
    *   [Reglas de Pluralización](#25-reglas-de-pluralización)
    *   [Metadatos de IA](#26-metadatos-de-ia)
3.  [Herramientas CLI](#3-herramientas-cli)
4.  [Integración Go SDK](#4-integración-go-sdk)

---

## 1. Introducción

Modern Billed-English Language (MBEL) es un formato de localización diseñado para la era de la IA.

### Filosofía
1.  **El Contexto es el Rey**.
2.  **La Lógica pertenece a los Datos**.
3.  **Primero la IA**.

---

## 2. El Lenguaje MBEL

### 2.1 Estructura del Archivo

```mbel
@namespace: "features.auth"  # Metadatos

[Pantalla de Login]          # Sección

title = "Iniciar Sesión"     # Asignación
```

### 2.2 Tipos de Datos

*   **Cadena Literal**: Entre comillas dobles.
*   **Cadena Multilínea**: Entre triples comillas `"""`.

### 2.3 Interpolación y Variables

Variables entre llaves `{}`.

*   **Sintaxis**: `Hola, {user_name}!`

### 2.4 Lógica y Control


@AI_Tone: "Formal"
submit_btn = "Enviar"
```

---


### 3.1 Instalación

```bash
go install github.com/yourusername/mbel/cmd/mbel@latest
```

### 3.2 Comandos

*   **`init`**: Asistente de configuración.
*   **`lint`**: Validador de sintaxis.
*   **`compile`**: Compilar a JSON.
*   **`watch`**: Recarga en caliente.
*   **`stats`**: Estadísticas.
*   **`fmt`**: Formateador.

---

## 4. Integración Go SDK

### 4.1 Arquitectura

*   **Manager**: Punto de entrada central.
*   **Runtime**: Instancia para un idioma específico.
*   **Repository**: Interfaz de fuente de datos.

### 4.2 Inicialización

```go
import "github.com/yourusername/mbel"

func init() {
    mbel.Init("./locales", mbel.Config{
        DefaultLocale: "es",
        Watch:         true,
    })
}
```

### 4.3 Uso (Función T)

`T` (Translate) resuelve la cadena según el contexto.

```go
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // 1. Clave Simple
    title := mbel.T(ctx, "page_title")

    // 2. Con Variables
    msg := mbel.T(ctx, "welcome", mbel.Vars{"name": "Ana"})

    // 3. Con Plural
    items := mbel.T(ctx, "cart_items", 5) 
}
```

### 4.4 Middleware HTTP

MBEL analiza automáticamente el encabezado `Accept-Language`.

```go
router.Use(mbel.Middleware)
```
