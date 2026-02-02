# Translation Services Integration Guide

This guide covers integrating MBEL with automated translation services like Google Translate, DeepL, and AWS Translate for efficient multi-language support.

---

## Overview

MBEL complements automated translation services by providing:
1. **Structure** — Define translation formats and validate content
2. **Caching** — Avoid repeated API calls for the same keys
3. **Fallbacks** — Handle missing translations gracefully
4. **Post-processing** — Fix AI translation issues and maintain consistency

---

## 1. Google Cloud Translation

### Setup

```bash
# Install Google Cloud client
go get cloud.google.com/go/translate

# Configure credentials
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/key.json"
```

### Basic Integration

```go
package main

import (
	"context"
	"cloud.google.com/go/translate"
	"github.com/makkiattooo/MBEL/pkg/mbel"
)

func translateWithGoogle(text, targetLang string) (string, error) {
	ctx := context.Background()
	
	client, err := translate.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{text}, &translate.Options{
		TargetLanguage: targetLang,
	})
	if err != nil {
		return "", err
	}

	return resp[0].TranslatedText, nil
}
```

### Cache Integration

```go
type TranslationCache struct {
	manager   *mbel.Manager
	cache     map[string]string
}

func (tc *TranslationCache) Get(lang, key string) (string, error) {
	// First, check if already in compiled MBEL
	existing := tc.manager.Get(lang, key, nil)
	if existing != "" {
		return existing, nil
	}

	// Cache key: lang + key
	cacheKey := lang + ":" + key
	if cached, ok := tc.cache[cacheKey]; ok {
		return cached, nil
	}

	// Get English source and translate
	source := tc.manager.Get("en", key, nil)
	if source == "" {
		return "", fmt.Errorf("key not found: %s", key)
	}

	translated, err := translateWithGoogle(source, lang)
	if err != nil {
		return "", err
	}

	// Cache result
	tc.cache[cacheKey] = translated
	return translated, nil
}
```

---

## 2. DeepL Integration

DeepL provides high-quality translations for European languages.

### Setup

```bash
go get github.com/DeepLcom/deepl-go
```

### Implementation

```go
import (
	"github.com/DeepLcom/deepl-go"
	"os"
)

func translateWithDeepL(text, targetLang string) (string, error) {
	auth := deepl.NewAuthenticator(os.Getenv("DEEPL_AUTH_KEY"))
	translator := deepl.NewTranslator(auth)

	result, err := translator.TranslateText(
		context.Background(),
		text,
		"EN", // source language
		"DE", // target language
	)
	if err != nil {
		return "", err
	}

	return result.Text, nil
}
```

### Cost Optimization

```go
// Batch API calls to reduce costs
func batchTranslate(texts []string, targetLang string) ([]string, error) {
	results := make([]string, len(texts))
	
	for i, text := range texts {
		// Rate limit: 50 requests/minute for free tier
		if i > 0 && i%10 == 0 {
			time.Sleep(1 * time.Second)
		}
		
		result, err := translateWithDeepL(text, targetLang)
		if err != nil {
			return nil, err
		}
		results[i] = result
	}
	
	return results, nil
}
```

---

## 3. AWS Translate

### Setup

```bash
go get github.com/aws/aws-sdk-go-v2/service/translate
```

### Implementation

```go
import (
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/translate"
	"github.com/aws/aws-sdk-go-v2/aws"
)

func translateWithAWS(text, targetLang string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return "", err
	}

	client := translate.NewFromConfig(cfg)
	
	resp, err := client.TranslateText(context.Background(), &translate.TranslateTextInput{
		Text:               aws.String(text),
		SourceLanguageCode: aws.String("en"),
		TargetLanguageCode: aws.String(targetLang),
	})

	if err != nil {
		return "", err
	}

	return *resp.TranslatedText, nil
}
```

---

## 4. Crowdin Integration

Sync translations with Crowdin platform:

```bash
#!/bin/bash
# sync-crowdin.sh

# 1. Upload source (English)
curl -X POST https://api.crowdin.com/api/v2/projects/PROJECT_ID/files \
  -H "Authorization: Bearer $CROWDIN_TOKEN" \
  -F "file=@locales/en.mbel"

# 2. Download translated files (after translation)
curl https://api.crowdin.com/api/v2/projects/PROJECT_ID/translations/exports \
  -H "Authorization: Bearer $CROWDIN_TOKEN" \
  -o translations.zip

# 3. Extract and validate
unzip translations.zip -d temp_translations/
for file in temp_translations/*.mbel; do
	mbel lint "$file" || echo "Validation failed: $file"
done

# 4. Integrate
cp temp_translations/*.mbel locales/
```

---

## 5. Automated Workflow

### Generate MBEL Files from API Response

```go
func generateMBELFromTranslations(translations map[string]string, lang string) error {
	content := fmt.Sprintf("@namespace: app\n@lang: %s\n\n", lang)
	
	for key, value := range translations {
		// Escape special characters
		value = strings.ReplaceAll(value, `"`, `\"`)
		value = strings.ReplaceAll(value, "\n", `\n`)
		content += fmt.Sprintf("%s = %q\n", key, value)
	}

	return ioutil.WriteFile(fmt.Sprintf("locales/%s.mbel", lang), []byte(content), 0644)
}

// Usage: Translate all languages
func main() {
	// 1. Get English base
	baseTranslations := map[string]string{
		"greeting": "Welcome, {name}!",
		"farewell": "Goodbye!",
	}

	// 2. Translate to all target languages
	targetLangs := []string{"de", "fr", "es", "it", "ja"}
	for _, lang := range targetLangs {
		translated := make(map[string]string)
		for key, val := range baseTranslations {
			result, _ := translateWithGoogle(val, lang)
			translated[key] = result
		}
		generateMBELFromTranslations(translated, lang)
	}

	// 3. Compile and test
	exec.Command("mbel", "compile", "locales/", "-o", "dist/translations.json").Run()
	exec.Command("go", "test", "./...").Run()
}
```

---

## 6. Quality Assurance

### Post-Translation Validation

```go
func validateTranslations(filePath string) error {
	// 1. Check syntax
	output, err := exec.Command("mbel", "lint", filePath).Output()
	if err != nil {
		return fmt.Errorf("lint failed: %s", output)
	}

	// 2. Check variable consistency
	content, _ := ioutil.ReadFile(filePath)
	fileContent := string(content)
	
	// 3. Check no hardcoded secrets
	if strings.Contains(fileContent, "sk-") || 
	   strings.Contains(fileContent, "password") ||
	   strings.Contains(fileContent, "secret") {
		return fmt.Errorf("suspected secrets in translations")
	}

	// 4. Check placeholder consistency with English
	enContent, _ := ioutil.ReadFile("locales/en.mbel")
	enPlaceholders := countPlaceholders(string(enContent))
	filePlaceholders := countPlaceholders(fileContent)
	
	if enPlaceholders != filePlaceholders {
		return fmt.Errorf("placeholder mismatch")
	}

	return nil
}

func countPlaceholders(content string) map[string]int {
	counts := make(map[string]int)
	re := regexp.MustCompile(`\{([^}]+)\}`)
	for _, match := range re.FindAllStringSubmatch(content, -1) {
		counts[match[1]]++
	}
	return counts
}
```

### A/B Testing

```go
// Compare translations from different services
type TranslationComparison struct {
	Key     string
	English string
	Google  string
	DeepL   string
	AWS     string
}

func compareTranslations(key, english string) TranslationComparison {
	google, _ := translateWithGoogle(english, "de")
	deepl, _ := translateWithDeepL(english, "de")
	aws, _ := translateWithAWS(english, "de")

	return TranslationComparison{
		Key:     key,
		English: english,
		Google:  google,
		DeepL:   deepl,
		AWS:     aws,
	}
}
```

---

## 7. Cost Analysis

### Typical Costs (Per 1 Million Characters)

| Service | Cost | Pros | Cons |
|---------|------|------|------|
| **Google Translate** | $20 | High quality, many languages | Expensive for scale |
| **DeepL API** | $15 | Excellent quality (EU) | Limited language pairs |
| **AWS Translate** | $15 | Good quality, AWS integration | Variable by language |
| **Human Translators** | $100-500 | Best quality, cultural aware | Slow, expensive |

**Optimization strategy:**
1. Use automated services for initial translations
2. Have humans review and correct key content
3. Cache results to avoid re-translation
4. Use cheaper services for internal-only strings

---

## 8. Workflow Examples

### Development Workflow

```bash
# 1. Create English base
echo 'greeting = "Hello, {name}!"' > locales/en.mbel

# 2. Auto-translate
go run translate.go --lang de,fr,es

# 3. Review in editor  
code locales/

# 4. Compile
mbel compile locales/ -o dist/translations.json

# 5. Test
go test ./...

# 6. Commit
git add locales/ dist/
git commit -m "feat: add German, French, Spanish translations"
```

### Production Workflow (CI/CD)

```yaml
# .github/workflows/translate.yml
name: Auto-Translate
on:
  push:
    paths:
      - 'locales/en.mbel'
jobs:
  translate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Run translation
        env:
          GOOGLE_APPLICATION_CREDENTIALS: ${{ secrets.GCP_CREDENTIALS }}
        run: |
          go run cmd/translate.go --langs de,fr,es,it,ja
      
      - name: Validate
        run: |
          mbel lint locales/
          mbel compile locales/ -o dist/translations.json
      
      - name: Create PR
        uses: peter-evans/create-pull-request@v4
        with:
          commit-message: "chore: auto-translated strings"
          title: "Auto-translated new strings"
          labels: "translation"
```

---

## 9. Troubleshooting

### Issue: "Over quota"

**Solution:**
- Use local caching
- Batch requests together
- Switch to cheaper service tier
- Reduce update frequency

### Issue: "Quality is poor"

**Solution:**
- Use human proofreaders
- Post-process with regex rules
- Use DeepL for better quality on EU languages
- Add AI context annotations

### Issue: "Inconsistent terminology"

**Solution:**
- Create glossary/terminology database
- Use service glossary features
- Implement post-processing normalization
- Use MBEL AI annotations for guidance

---

## 10. Best Practices

1. **Always use caching** — Avoid repeated API calls
2. **Validate output** — Use MBEL linting after translation
3. **Version translations** — Track service/version used
4. **Monitor costs** — Set budget alerts on API accounts
5. **Have fallbacks** — Always fall back to English if translation fails
6. **Human review** — Especially for user-facing text
7. **Test thoroughly** — Include translated strings in test coverage
8. **Document context** — Use MBEL AI annotations

---

See also:
- [SECURITY.md](SECURITY.md) — Protect your API keys
- [QUICKSTART.md](QUICKSTART.md) — Basic usage
- [ARCHITECTURE.md](ARCHITECTURE.md) — How Manager handles multiple languages

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

