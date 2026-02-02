# AI Annotations & Translation Context

## Overview

MBEL treats translation management as a partnership between developers and AI systems. **AI Annotations** are structured metadata attached to translation keys that provide context, tone guidance, and constraints to AI translators.

Unlike traditional i18n tools that provide zero context to LLMs, MBEL creates a deterministic pipeline for high-quality automated translations.

---

## Basic Syntax

### Single-Line Annotations

Use comment syntax with `AI_` prefix:

```mbel
# AI_Context: User greeting message displayed on homepage
# AI_Tone: Friendly, warm, welcoming
# AI_MaxLength: 50
greeting = "Hello, {name}! Welcome back."
```

### Supported Annotation Types

| Type | Purpose | Example |
|------|---------|---------|
| `AI_Context` | Where and why this string is used | "Button in header for navigating to settings" |
| `AI_Tone` | Emotional tone or style | "Professional, formal", "Playful, casual" |
| `AI_Audience` | Who sees this string | "Non-technical users", "Developers" |
| `AI_MaxLength` | Character limit | 80 |
| `AI_Constraints` | Hard rules | "No exclamation marks", "Must start with verb" |
| `AI_Examples` | Reference translations | "Spanish: \"Hola\"", "French: \"Bonjour\"" |

---

## Advanced: Multi-Line Annotations

For complex context, use curly braces:

```mbel
# AI_Context: {
#   Error message shown when payment fails
#   Suggest next steps (retry, contact support)
#   Should be concise but helpful
#   No jargon - clear to non-technical users
# }
# AI_Constraints: {
#   - Do NOT mention "transaction failed"
#   - DO suggest retry option
#   - Max 80 characters
#   - No all-caps words
# }
# AI_Examples: {
#   German: "Zahlungsvorgang konnte nicht abgeschlossen werden."
#   Spanish: "No pudimos procesar tu pago."
#   French: "Votre paiement n'a pas abouti."
# }
payment_error = "We couldn't process your payment. Please try again or contact support."
```

### Multi-Line Format Rules

- Start with `# AI_Type: {`
- Each line is a comment line (`#`)
- Close with `# }`
- Content inside is preserved with newlines

---

## How Annotations Flow to AI

When you export translations for AI processing:

```bash
mbel compile locales/ -o dist/translations.json
# Also generates dist/translations.sourcemap.json and AI context data
```

The `__ai` metadata field contains all annotations:

```json
{
  "__ai": {
    "greeting": [
      {"type": "Context", "value": "User greeting message..."},
      {"type": "Tone", "value": "Friendly, warm, welcoming"},
      {"type": "MaxLength", "value": "50"}
    ],
    "payment_error": [...]
  },
  "greeting": "Hello, {name}! Welcome back.",
  "payment_error": "..."
}
```

---

## SDK Usage

### Access Annotations in Code

```go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	
	mbel "github.com/makkiattooo/MBEL"
)

func main() {
	data, _ := ioutil.ReadFile("dist/translations.json")
	var compiled map[string]interface{}
	json.Unmarshal(data, &compiled)

	// Access AI annotations
	if ai, ok := compiled["__ai"].(map[string]interface{}); ok {
		if greetingAnnotations, ok := ai["greeting"].([]interface{}); ok {
			fmt.Printf("Greeting has %d annotations\n", len(greetingAnnotations))
			for _, ann := range greetingAnnotations {
				fmt.Printf("  %v\n", ann)
			}
		}
	}
}
```

### Custom Translation Pipeline

```go
func translateWithContext(key string, targetLang string) (string, error) {
	// 1. Load compiled data with annotations
	data, _ := loadCompiledData()
	
	// 2. Extract annotations for this key
	annotations := extractAnnotations(data, key)
	
	// 3. Call LLM with context
	prompt := buildTranslationPrompt(key, targetLang, annotations)
	translation := callLLM(prompt)
	
	// 4. Validate against constraints
	if valid := validateTranslation(translation, annotations); !valid {
		return "", fmt.Errorf("translation violates constraints")
	}
	
	return translation, nil
}
```

---

## Best Practices

### ✅ Do:
- **Be specific** — "Button that opens user profile menu in header" not "Button"
- **Include examples** — Show reference translations from other languages
- **Set constraints early** — If max length matters, say so
- **Use clear tone names** — "Professional/formal", "Casual/friendly", "Technical"
- **Document context** — Why does this string exist? Where is it shown?

### ❌ Don't:
- Use vague context — "UI text" is not helpful
- Set conflicting constraints — "Be brief" and "Explain everything"
- Forget about tone — Neutral is not always right
- Assume the translator knows your product
- Use code jargon without explanation

---

## Example: Complete Annotation Pattern

```mbel
# E-commerce product listing page

# AI_Context: Price label shown under product image
# AI_Tone: Clear, factual, neutral
# AI_Audience: Customers shopping online
# AI_Constraints: Must include currency symbol
# AI_Examples: German: "Preis: €49,99" | Spanish: "Precio: $49,99"
price_label = "Price: ${price}"

# AI_Context: Button to add item to shopping cart
# AI_Tone: Action-oriented, encouraging
# AI_MaxLength: 15
# AI_Constraints: No punctuation, verb-first
add_to_cart = "Add to Cart"

# AI_Context: {
#   Warning shown when product is out of stock
#   Should NOT use alarm-inducing words
#   Provide guidance on alternatives
# }
# AI_Constraints: {
#   - Suggest "notify me" option availability
#   - Use "unavailable" or "out of stock", not "sold out"
#   - Keep to 50 characters
# }
out_of_stock = "This item is currently unavailable. Notify me when it returns."
```

---

## Integration with Translation Tools

### Export for CAT Tools

```bash
# Future: Export annotations for professional translators
mbel export --format=xliff --include-annotations locales/
```

### Import from Translation Memory

```bash
# Load previous translations with context
mbel import --source translations.xliff --merge-context
```

---

## Parsing Details

The parser handles:
- Single-line: `# AI_Type: value` 
- Multi-line: `# AI_Type: { ... }`
- Nested braces (counted for proper termination)
- Escape sequences in values
- Line continuation in curly-brace blocks

Annotations are **associated with the next key** they precede:

```mbel
# AI_Context: Greeting
greeting = "Hello"

# AI_Context: Farewell  ← This is for the "farewell" key
# AI_Tone: Warm
farewell = "Goodbye"
```

