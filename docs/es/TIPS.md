# MBEL: Consejos y Mejores Prácticas

Sigue estos patrones para sacar el máximo provecho a MBEL.

## 1. Convenciones de Nomenclatura
Usa claves jerárquicas en lugar de nombres planos como `label1`.
*   **Patrón**: `[feature].[screen].[component].[element]`

## 2. Organización de Archivos
Divide tus localizaciones en unidades lógicas. La estructura de carpetas define automáticamente los namespaces con la opción `--ns`.

## 3. Uso del Contexto de IA
No escatimes en `@AI_Context`. Es el seguro de calidad para tus traducciones automatizadas.

## 4. Bloques Lógicos vs Interpolación en Código
Evita construir frases en tu código Go.
*   **Mal (Go)**: `fmt.Sprintf(mbel.T(ctx, "hola") + " " + name)`
*   **Bien (MBEL)**:
    ```mbel
    welcome(name) {
        [other] => "¡Bienvenido/a, {name}!"
    }
    ```
