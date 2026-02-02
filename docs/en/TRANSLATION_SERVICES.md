# Translation Services Integration

## Overview

MBEL can integrate with professional translation services and platforms for automated or professional human translation of your content.

---

## Integration Patterns

### 1. Export for Professional Translation

Export your base language translations to standard formats for CAT (Computer-Assisted Translation) tools:

```bash
# Export to XLIFF (future feature)
mbel export --format=xliff locales/en.mbel -o translations.xliff

# For professional translators who use:
# - memoQ
# - Trados
# - Phrase (formerly Memsource)
# - Crowdin
```

### 2. API-Based Translation

Call LLM APIs with MBEL context:

```go
package main

import (
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/makkiattooo/MBEL"
)

type AITranslator struct {
	client *openai.Client
	m      *mbel.Manager
}

func (t *AITranslator) Translate(sourceKey, targetLang string) (string, error) {
	// 1. Get source translation
	source, err := t.m.Get("en", sourceKey, nil)
	if err != nil {
		return "", err
	}

	// 2. Get AI context
	annotations := t.extractAnnotations(sourceKey)

	// 3. Build prompt with context
	prompt := fmt.Sprintf(`
Translate this string to %s:

Source: "%s"
Context: %s
Tone: %s
Constraints: %s

Provide ONLY the translated string, no explanations.
	`, targetLang, source, annotations["Context"], annotations["Tone"], annotations["Constraints"])

	// 4. Call OpenAI
	resp, err := t.client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})

	return resp.Choices[0].Message.Content, nil
}

func (t *AITranslator) extractAnnotations(key string) map[string]string {
	// Load annotations from compiled data
	// Implementation: parse __ai field
	return map[string]string{
		"Context": "...",
		"Tone": "...",
		"Constraints": "...",
	}
}
```

### 3. Crowdin Integration

Sync translations with Crowdin:

```bash
#!/bin/bash
# sync-crowdin.sh

# 1. Export from MBEL
mbel compile locales/en/ -o crowdin-source.json

# 2. Upload to Crowdin API
curl -X POST https://api.crowdin.com/api/v2/projects/123/files/upload \
  -H "Authorization: Bearer $CROWDIN_TOKEN" \
  -F "file=@crowdin-source.json"

# 3. Download translated files (after translation)
curl https://api.crowdin.com/api/v2/projects/123/translations/exports/123 \
  -H "Authorization: Bearer $CROWDIN_TOKEN" \
  -o translations.zip

# 4. Extract and integrate back
unzip translations.zip -d locales/
```

### 4. Google Translate API

Batch translate untranslated strings:

```go
package main

import (
	"cloud.google.com/go/translate"
)

func batchTranslate(m *mbel.Manager, targetLangs []string) error {
	client, _ := translate.NewClient(context.Background())
	defer client.Close()

	// Get source keys
	sourceData, _ := m.Get("en", "", nil)

	for lang := range targetLangs {
		for key, value := range sourceData {
			// Check if translation exists
			translated, _ := m.Get(lang, key, nil)
			if translated == "" {
				// Translate using Google
				result, _ := client.Translate(context.Background(),
					[]string{value.(string)},
					&translate.Options{TargetLanguage: lang},
				)
				fmt.Printf("Translated %s to %s: %s\n", key, lang, result[0].TranslatedText)
			}
		}
	}

	return nil
}
```

---

## Workflow: Translation Lifecycle

```
1. Development
   └─ Add/modify strings in en.mbel
   └─ Add AI_Context annotations
   
2. Export
   └─ mbel compile locales/ -o dist/translations.json
   └─ mbel export --format=xliff -o for-translation.xliff
   
3. Translation (Professional or AI)
   └─ Send to translators or LLM
   └─ Use AI annotations for context
   
4. Quality Assurance
   └─ Validate translated strings
   └─ Check tone & constraints
   └─ Test in staging environment
   
5. Integration
   └─ Import translated files
   └─ Run mbel lint to validate
   └─ Merge into locales/
   
6. Deployment
   └─ Compile all languages
   └─ Generate sourcemap
   └─ Deploy with CI/CD
```

---

## Validation & Quality

### Pre-Import Validation

```bash
#!/bin/bash
# validate-translations.sh

for file in translations/*.json; do
    # 1. JSON syntax
    jq empty "$file" || { echo "Invalid JSON: $file"; exit 1; }
    
    # 2. Required keys (compare with source)
    source_keys=$(jq 'keys | length' en.json)
    target_keys=$(jq 'keys | length' "$file")
    
    if [ "$source_keys" != "$target_keys" ]; then
        echo "Missing keys in $file"
        jq 'keys' "$file" | diff - <(jq 'keys' en.json)
    fi
    
    # 3. Placeholder consistency
    mbel lint "$file" || { echo "Lint failed: $file"; exit 1; }
done
```

### Constraint Checking

```go
func validateTranslation(key, sourceValue, translatedValue string, constraints map[string]interface{}) error {
    // Check max length
    if maxLen, ok := constraints["MaxLength"].(int); ok {
        if len(translatedValue) > maxLen {
            return fmt.Errorf("translation exceeds max length %d", maxLen)
        }
    }
    
    // Check required words
    if required, ok := constraints["MustContain"].([]string); ok {
        for _, word := range required {
            if !strings.Contains(strings.ToLower(translatedValue), strings.ToLower(word)) {
                return fmt.Errorf("translation missing required word: %s", word)
            }
        }
    }
    
    // Check forbidden words
    if forbidden, ok := constraints["MustNotContain"].([]string); ok {
        for _, word := range forbidden {
            if strings.Contains(strings.ToLower(translatedValue), strings.ToLower(word)) {
                return fmt.Errorf("translation contains forbidden word: %s", word)
            }
        }
    }
    
    return nil
}
```

---

## Real-World Examples

### Example 1: E-commerce Platform

```
1. Create base translations:
   locales/en.mbel
   
2. Add product-specific context:
   @AI_Context: "Product price visible in search results"
   @AI_MaxLength: 20
   price_format = "${price}"
   
3. Export to Crowdin:
   mbel compile locales/en -o temp.json
   curl ... upload to Crowdin
   
4. Translators see context annotations
   and are constrained by MaxLength
   
5. Download Spanish translation:
   {"price_format": "${precio}"}
   
6. Integrate back:
   cp locales/es/from_crowdin.json locales/es/prices.mbel
   
7. Run tests:
   go test ./...
   mbel lint locales/
```

### Example 2: Saas App with Multiple Regions

```
Target languages: en (base), de, fr, es, ja, zh
Workflow:

1. Week 1: Add new feature strings
   └─ Create feature.mbel with full annotations
   
2. Week 2: Export to Crowdin
   └─ Review context with translators
   
3. Week 3: Receive translations
   └─ Validate with AI constraints
   
4. Week 4: QA in staging
   └─ Test all 6 language variants
   └─ Check tone consistency
   
5. Week 5: Deploy
   └─ mbel compile with --sourcemap
   └─ Push to production
```

---

## Machine Translation Quality

### Pre-Translate for Speed

```go
// First pass: machine translation (fast, cost-effective)
func preTranslate(text, targetLang string) (string, error) {
	client, _ := translate.NewClient(ctx)
	result, _ := client.Translate(ctx, []string{text},
		&translate.Options{TargetLanguage: targetLang})
	return result[0].TranslatedText, nil
}

// Manual review flag for human translators
type TranslationStatus struct {
	Key string
	Source string
	MachineTranslation string
	HumanTranslation *string // nil means needs review
	Quality float64 // 0.0-1.0 confidence
	Flags []string // ["tone_mismatch", "truncated", "placeholder_error"]
}
```

### Confidence Scoring

```go
func scoreTranslation(source, translated string, constraints map[string]interface{}) float64 {
	score := 1.0
	
	// Deduct for length violations
	if maxLen, ok := constraints["MaxLength"].(int); ok {
		if len(translated) > maxLen {
			score -= 0.1
		}
	}
	
	// Deduct for placeholder mismatch
	sourcePhCount := strings.Count(source, "{")
	transPhCount := strings.Count(translated, "{")
	if sourcePhCount != transPhCount {
		score -= 0.2
	}
	
	// Deduct for very different lengths
	lengthRatio := float64(len(translated)) / float64(len(source))
	if lengthRatio < 0.5 || lengthRatio > 2.5 {
		score -= 0.15
	}
	
	return score
}
```

---

## Troubleshooting

### Missing Translations

```bash
# Find strings without translations
jq 'keys as $en | input | keys as $target | 
    $en - $target' en.json de.json | jq -r '.[]'
```

### Placeholder Mismatches

```bash
# Check placeholders
jq '.[] | select(test("{")) | . as $source | 
    input[.] | . as $target |
    if (($source | [scan("{[^}]+}")]) != ($target | [scan("{[^}]+}")])) 
    then "MISMATCH: " + . else empty end' en.json de.json
```

### Encoding Issues

```bash
# Validate UTF-8
file -i locales/*/translations.json | grep -v utf-8 || echo "All files UTF-8 OK"
```

