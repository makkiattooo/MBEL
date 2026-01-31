package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	mbel "github.com/makkiattooo/MBEL"
)

const version = "1.2.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Handle flags provided as the first argument
	arg1 := os.Args[1]
	if arg1 == "-v" || arg1 == "--version" || arg1 == "version" {
		fmt.Printf("mbel version %s (%s/%s)\n", version, runtime.GOOS, runtime.GOARCH)
		return
	}
	if arg1 == "-h" || arg1 == "--help" || arg1 == "help" {
		printUsage()
		return
	}

	// Command Switch
	switch arg1 {
	case "init":
		initCmd()
	case "lint":
		lintCmd(os.Args[2:])
	case "compile":
		compileCmd(os.Args[2:])
	case "watch":
		watchCmd(os.Args[2:])
	case "fmt":
		fmtCmd(os.Args[2:])
	case "stats":
		statsCmd(os.Args[2:])
	case "diff":
		diffCmd(os.Args[2:])
	case "import":
		importCmd(os.Args[2:])
	case "translate":
		translateCmd(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", arg1)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`
/$$      /$$ /$$$$$$$  /$$$$$$$$ /$$      
| $$$    /$$$| $$__  $$| $$_____/| $$      
| $$$$  /$$$$| $$  \ $$| $$      | $$      
| $$ $$/$$ $$| $$$$$$$ | $$$$$   | $$      
| $$  $$$| $$| $$__  $$| $$__/   | $$      
| $$\  $ | $$| $$  \ $$| $$      | $$      
| $$ \/  | $$| $$$$$$$/| $$$$$$$$| $$$$$$$$
|__/     |__/|_______/ |________/|________/

MBEL v` + version + `
"Internationalization for the AI Era"

Usage:
  mbel <command> [arguments]

Core Commands:
  init      ✨ Start here! Interactive project setup
  watch     👁  Watch mode (hot-reload for development)
  compile   📦 Compile .mbel files to JSON (for production)
  lint      🔍 Validate syntax and AI rules

Helpers:
  fmt       🎨 Auto-format .mbel files
  stats     📊 Show project statistics
  diff      ↔  Compare locales (find missing keys)
  import    📥 Import from JSON/YAML
  version   ℹ  Show version info

Flags:
  -v, --version   Show version
  -h, --help      Show this help message

Quick Start:
  mbel init`)
}

// ============================================================================
// INIT COMMAND (Interactive Wizard)
// ============================================================================

func initCmd() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("✨ Welcome to MBEL! Let's set up localization for your project.")
	fmt.Println("-------------------------------------------------------------")

	// 1. Directory
	fmt.Print("Where should we store your locale files? [default: locales]: ")
	dir, _ := reader.ReadString('\n')
	dir = strings.TrimSpace(dir)
	if dir == "" {
		dir = "locales"
	}

	// 2. Default Language
	fmt.Print("What is your default language? [default: en]: ")
	lang, _ := reader.ReadString('\n')
	lang = strings.TrimSpace(lang)
	if lang == "" {
		lang = "en"
	}

	// 3. Create Directory
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("❌ Failed to create directory: %v\n", err)
		os.Exit(1)
	}

	// 4. Create Example File
	examplePath := filepath.Join(dir, lang+".mbel")
	if _, err := os.Stat(examplePath); os.IsNotExist(err) {
		content := fmt.Sprintf(`@AI_Context: "Main application strings"
title = "My App"

# Example logic block
welcome(name) {
    [other] => "Welcome, {name}!"
}

# Example plural
items_count(n) {
    [one]   => "You have one item."
    [other] => "You have {n} items."
}
`)
		if err := ioutil.WriteFile(examplePath, []byte(content), 0644); err != nil {
			fmt.Printf("❌ Failed to create example file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Created %s\n", examplePath)
	} else {
		fmt.Printf("ℹ️  %s already exists, skipping creation.\n", examplePath)
	}

	// 5. Success Message
	fmt.Println("\n🎉 You are all set!")
	fmt.Println("\nNext steps:")
	fmt.Println("1. Install the Go SDK:")
	fmt.Println("   go get github.com/yourusername/mbel")
	fmt.Println("\n2. Initialize in your code:")
	fmt.Printf("   mbel.Init(\"./%s\", mbel.Config{Watch: true})\n", dir)
	fmt.Println("\n3. Run watch mode during development:")
	fmt.Printf("   mbel watch %s\n", dir)
}

// ============================================================================
// FILE DISCOVERY
// ============================================================================

// discoverFiles finds all .mbel files in given paths (files or directories)
// Supports recursive directory scanning up to 10 levels deep
func discoverFiles(paths []string) ([]string, error) {
	var files []string
	seen := make(map[string]bool)

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			// Try as glob pattern
			matches, globErr := filepath.Glob(path)
			if globErr != nil {
				return nil, fmt.Errorf("invalid path or pattern %s: %w", path, err)
			}
			for _, match := range matches {
				if !seen[match] {
					seen[match] = true
					files = append(files, match)
				}
			}
			continue
		}

		if info.IsDir() {
			// Recursive walk
			err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() && strings.HasSuffix(p, ".mbel") {
					if !seen[p] {
						seen[p] = true
						files = append(files, p)
					}
				}
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("error walking directory %s: %w", path, err)
			}
		} else if strings.HasSuffix(path, ".mbel") {
			if !seen[path] {
				seen[path] = true
				files = append(files, path)
			}
		}
	}

	return files, nil
}

// deriveNamespace extracts namespace from file path relative to base
// e.g., locales/en/features/auth/login.mbel -> features.auth
func deriveNamespace(filePath, basePath string) string {
	rel, err := filepath.Rel(basePath, filePath)
	if err != nil {
		return ""
	}

	dir := filepath.Dir(rel)
	if dir == "." {
		return ""
	}

	// Convert path separators to dots
	ns := strings.ReplaceAll(dir, string(filepath.Separator), ".")
	return ns
}

// ============================================================================
// LINT COMMAND
// ============================================================================

func lintCmd(args []string) {
	fs := flag.NewFlagSet("lint", flag.ExitOnError)
	verbose := fs.Bool("v", false, "Verbose output")
	parallel := fs.Int("j", runtime.NumCPU(), "Parallel workers")
	fs.Parse(args)

	paths := fs.Args()
	if len(paths) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No files or directories specified")
		fmt.Fprintln(os.Stderr, "Usage: mbel lint <path> [path2 ...]")
		os.Exit(1)
	}

	files, err := discoverFiles(paths)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "Warning: No .mbel files found")
		os.Exit(0)
	}

	if *verbose {
		fmt.Printf("Found %d files, using %d workers\n", len(files), *parallel)
	}

	// Parallel linting
	type lintResult struct {
		file  string
		err   error
		stats struct {
			statements  int
			annotations int
		}
	}

	results := make(chan lintResult, len(files))
	fileChan := make(chan string, len(files))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < *parallel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range fileChan {
				res := lintResult{file: file}
				content, err := ioutil.ReadFile(file)
				if err != nil {
					res.err = err
					results <- res
					continue
				}

				l := mbel.NewLexer(string(content))
				p := mbel.NewParser(l)
				program := p.ParseProgram()

				if errs := p.Errors(); len(errs) > 0 {
					res.err = fmt.Errorf("syntax errors:\n  %s", strings.Join(errs, "\n  "))
				} else {
					// Validation Rules
					for _, ann := range program.AIAnnotations {
						if ann.Type == "MaxLength" && ann.ForKey != "" {
							for _, stmt := range program.Statements {
								if assign, ok := stmt.(*mbel.AssignStatement); ok && assign.Name == ann.ForKey {
									if sl, ok := assign.Value.(*mbel.StringLiteral); ok {
										if limit, err := strconv.Atoi(ann.Value); err == nil {
											if len(sl.Value) > limit {
												res.err = fmt.Errorf("validation error: %s exceeds max length of %d (got %d)", ann.ForKey, limit, len(sl.Value))
											}
										}
									}
								}
							}
						}
					}

					res.stats.statements = len(program.Statements)
					res.stats.annotations = len(program.AIAnnotations)
				}
				results <- res
			}
		}()
	}

	// Feed files
	for _, f := range files {
		fileChan <- f
	}
	close(fileChan)

	// Wait and collect
	go func() {
		wg.Wait()
		close(results)
	}()

	hasErrors := false
	successCount := 0
	for res := range results {
		if res.err != nil {
			fmt.Fprintf(os.Stderr, "✗ %s: %v\n", res.file, res.err)
			hasErrors = true
		} else {
			successCount++
			if *verbose {
				fmt.Printf("✓ %s (%d statements, %d AI annotations)\n",
					res.file, res.stats.statements, res.stats.annotations)
			}
		}
	}

	if hasErrors {
		os.Exit(1)
	}

	fmt.Printf("✓ %d files valid\n", successCount)
}

// ============================================================================
// COMPILE COMMAND
// ============================================================================

func compileCmd(args []string) {
	fs := flag.NewFlagSet("compile", flag.ExitOnError)
	output := fs.String("o", "", "Output file (default: stdout)")
	pretty := fs.Bool("pretty", true, "Pretty-print JSON")
	parallel := fs.Int("j", runtime.NumCPU(), "Parallel workers")
	withNamespace := fs.Bool("ns", true, "Derive namespace from folder path")
	fs.Parse(args)

	paths := fs.Args()
	if len(paths) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No files or directories specified")
		fmt.Fprintln(os.Stderr, "Usage: mbel compile <path> [-o output.json]")
		os.Exit(1)
	}

	files, err := discoverFiles(paths)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "Warning: No .mbel files found")
		os.Exit(0)
	}

	// Determine base path for namespace derivation
	basePath := ""
	if *withNamespace && len(paths) > 0 {
		info, err := os.Stat(paths[0])
		if err == nil && info.IsDir() {
			basePath = paths[0]
		} else {
			basePath = filepath.Dir(paths[0])
		}
	}

	// Parallel compilation
	type compileResult struct {
		file      string
		namespace string
		data      map[string]interface{}
		err       error
	}

	results := make(chan compileResult, len(files))
	fileChan := make(chan string, len(files))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < *parallel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range fileChan {
				res := compileResult{file: file}

				if *withNamespace && basePath != "" {
					res.namespace = deriveNamespace(file, basePath)
				}

				content, err := ioutil.ReadFile(file)
				if err != nil {
					res.err = err
					results <- res
					continue
				}

				l := mbel.NewLexer(string(content))
				p := mbel.NewParser(l)
				program := p.ParseProgram()

				if errs := p.Errors(); len(errs) > 0 {
					res.err = fmt.Errorf("syntax errors:\n  %s", strings.Join(errs, "\n  "))
					results <- res
					continue
				}

				c := mbel.NewCompiler()
				result, err := c.Compile(program)
				if err != nil {
					res.err = err
					results <- res
					continue
				}

				resultMap, ok := result.(map[string]interface{})
				if !ok {
					res.err = fmt.Errorf("unexpected result type: %T", result)
					results <- res
					continue
				}

				res.data = resultMap
				results <- res
			}
		}()
	}

	// Feed files
	for _, f := range files {
		fileChan <- f
	}
	close(fileChan)

	// Wait and collect
	go func() {
		wg.Wait()
		close(results)
	}()

	// Merge all results
	merged := make(map[string]interface{})
	hasErrors := false

	for res := range results {
		if res.err != nil {
			fmt.Fprintf(os.Stderr, "✗ %s: %v\n", res.file, res.err)
			hasErrors = true
			continue
		}

		// Merge with namespace prefix
		for k, v := range res.data {
			key := k
			if res.namespace != "" && !strings.HasPrefix(k, "__") {
				key = res.namespace + "." + k
			}
			merged[key] = v
		}
	}

	if hasErrors {
		os.Exit(1)
	}

	var jsonData []byte
	if *pretty {
		jsonData, err = json.MarshalIndent(merged, "", "  ")
	} else {
		jsonData, err = json.Marshal(merged)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	if *output != "" {
		if err := ioutil.WriteFile(*output, jsonData, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Compiled %d files to %s\n", len(files), *output)
	} else {
		fmt.Println(string(jsonData))
	}
}

// ============================================================================
// WATCH COMMAND
// ============================================================================

func watchCmd(args []string) {
	fs := flag.NewFlagSet("watch", flag.ExitOnError)
	output := fs.String("o", "", "Output file")
	interval := fs.Int("i", 1000, "Poll interval in milliseconds")
	fs.Parse(args)

	paths := fs.Args()
	if len(paths) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No directory specified")
		fmt.Fprintln(os.Stderr, "Usage: mbel watch <directory> [-o output.json]")
		os.Exit(1)
	}

	fmt.Printf("👁 Watching %s (Ctrl+C to stop)\n", paths[0])

	// Track file modification times
	lastMod := make(map[string]time.Time)

	for {
		files, err := discoverFiles(paths)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			time.Sleep(time.Duration(*interval) * time.Millisecond)
			continue
		}

		changed := false
		for _, file := range files {
			info, err := os.Stat(file)
			if err != nil {
				continue
			}
			if last, exists := lastMod[file]; !exists || info.ModTime().After(last) {
				if exists {
					fmt.Printf("  📝 Changed: %s\n", filepath.Base(file))
					changed = true
				}
				lastMod[file] = info.ModTime()
			}
		}

		if changed && *output != "" {
			// Recompile
			result := make(map[string]interface{})
			hasErrors := false

			for _, file := range files {
				content, err := ioutil.ReadFile(file)
				if err != nil {
					continue
				}

				l := mbel.NewLexer(string(content))
				p := mbel.NewParser(l)
				program := p.ParseProgram()

				if len(p.Errors()) > 0 {
					fmt.Fprintf(os.Stderr, "  ✗ %s: %v\n", filepath.Base(file), p.Errors())
					hasErrors = true
					continue
				}

				c := mbel.NewCompiler()
				compiled, _ := c.Compile(program)
				if compMap, ok := compiled.(map[string]interface{}); ok {
					for k, v := range compMap {
						result[k] = v
					}
				}
			}

			if !hasErrors {
				jsonData, _ := json.MarshalIndent(result, "", "  ")
				ioutil.WriteFile(*output, jsonData, 0644)
				fmt.Printf("  ✓ Compiled to %s\n", *output)
			}
		}

		time.Sleep(time.Duration(*interval) * time.Millisecond)
	}
}

// ============================================================================
// FORMAT COMMAND
// ============================================================================

func fmtCmd(args []string) {
	fs := flag.NewFlagSet("fmt", flag.ExitOnError)
	dryRun := fs.Bool("n", false, "Dry run (show changes without writing)")
	fs.Parse(args)

	paths := fs.Args()
	if len(paths) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No files specified")
		fmt.Fprintln(os.Stderr, "Usage: mbel fmt <files...> [-n]")
		os.Exit(1)
	}

	files, err := discoverFiles(paths)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	formatted := 0
	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", file, err)
			continue
		}

		// Parse and re-emit formatted
		l := mbel.NewLexer(string(content))
		p := mbel.NewParser(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			fmt.Fprintf(os.Stderr, "✗ %s: syntax errors\n", file)
			continue
		}

		newContent := formatProgram(program)

		if string(content) != newContent {
			if *dryRun {
				fmt.Printf("Would format: %s\n", file)
			} else {
				ioutil.WriteFile(file, []byte(newContent), 0644)
				fmt.Printf("Formatted: %s\n", file)
			}
			formatted++
		}
	}

	fmt.Printf("✓ %d files formatted\n", formatted)
}

func formatProgram(p *mbel.Program) string {
	var b strings.Builder

	// Metadata first
	for _, stmt := range p.Statements {
		if ms, ok := stmt.(*mbel.MetadataStatement); ok {
			b.WriteString(fmt.Sprintf("@%s: %s\n", ms.Key, ms.Value))
		}
	}

	// Then sections and assignments
	currentSection := ""
	for _, stmt := range p.Statements {
		switch s := stmt.(type) {
		case *mbel.SectionStatement:
			if b.Len() > 0 {
				b.WriteString("\n")
			}
			b.WriteString(fmt.Sprintf("[%s]\n", s.Name))
			currentSection = s.Name
		case *mbel.AssignStatement:
			if currentSection == "" && b.Len() > 0 {
				b.WriteString("\n")
			}
			if sl, ok := s.Value.(*mbel.StringLiteral); ok {
				if strings.Contains(sl.Value, "\n") {
					b.WriteString(fmt.Sprintf("%s = \"\"\"\n%s\"\"\"\n", s.Name, sl.Value))
				} else {
					b.WriteString(fmt.Sprintf("%s = \"%s\"\n", s.Name, sl.Value))
				}
			} else if _, ok := s.Value.(*mbel.BlockExpression); ok {
				b.WriteString(fmt.Sprintf("%s(...) { ... }\n", s.Name))
			}
		}
	}

	return b.String()
}

// ============================================================================
// STATS COMMAND
// ============================================================================

func statsCmd(args []string) {
	fs := flag.NewFlagSet("stats", flag.ExitOnError)
	fs.Parse(args)

	paths := fs.Args()
	if len(paths) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No path specified")
		os.Exit(1)
	}

	files, err := discoverFiles(paths)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	totalKeys := 0
	totalStrings := 0
	totalBlocks := 0
	totalAnnotations := 0
	keyCount := make(map[string]int)

	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}

		l := mbel.NewLexer(string(content))
		p := mbel.NewParser(l)
		program := p.ParseProgram()

		totalAnnotations += len(program.AIAnnotations)

		for _, stmt := range program.Statements {
			if as, ok := stmt.(*mbel.AssignStatement); ok {
				totalKeys++
				keyCount[as.Name]++
				if _, ok := as.Value.(*mbel.StringLiteral); ok {
					totalStrings++
				} else if _, ok := as.Value.(*mbel.BlockExpression); ok {
					totalBlocks++
				}
			}
		}
	}

	// Find duplicates
	duplicates := []string{}
	for key, count := range keyCount {
		if count > 1 {
			duplicates = append(duplicates, fmt.Sprintf("%s (%d)", key, count))
		}
	}
	sort.Strings(duplicates)

	fmt.Println("📊 MBEL Statistics")
	fmt.Println("──────────────────")
	fmt.Printf("Files:          %d\n", len(files))
	fmt.Printf("Total keys:     %d\n", totalKeys)
	fmt.Printf("  Strings:      %d\n", totalStrings)
	fmt.Printf("  Logic blocks: %d\n", totalBlocks)
	fmt.Printf("AI annotations: %d\n", totalAnnotations)

	if len(duplicates) > 0 {
		fmt.Printf("\n⚠️  Duplicate keys (%d):\n", len(duplicates))
		for _, d := range duplicates {
			fmt.Printf("  - %s\n", d)
		}
	}
}

// ============================================================================
// DIFF COMMAND
// ============================================================================

func diffCmd(args []string) {
	fs := flag.NewFlagSet("diff", flag.ExitOnError)
	fs.Parse(args)

	paths := fs.Args()
	if len(paths) < 2 {
		fmt.Fprintln(os.Stderr, "Error: Need two paths to compare")
		fmt.Fprintln(os.Stderr, "Usage: mbel diff <path1> <path2>")
		os.Exit(1)
	}

	keys1 := collectKeys(paths[0])
	keys2 := collectKeys(paths[1])

	// Find missing in path2
	missing := []string{}
	for key := range keys1 {
		if _, exists := keys2[key]; !exists {
			missing = append(missing, key)
		}
	}
	sort.Strings(missing)

	// Find extra in path2
	extra := []string{}
	for key := range keys2 {
		if _, exists := keys1[key]; !exists {
			extra = append(extra, key)
		}
	}
	sort.Strings(extra)

	fmt.Printf("🔍 Comparing %s ↔ %s\n", paths[0], paths[1])
	fmt.Println("──────────────────────────")

	if len(missing) == 0 && len(extra) == 0 {
		fmt.Println("✓ All keys match!")
		return
	}

	if len(missing) > 0 {
		fmt.Printf("\n❌ Missing in %s (%d):\n", paths[1], len(missing))
		for _, k := range missing {
			fmt.Printf("  - %s\n", k)
		}
	}

	if len(extra) > 0 {
		fmt.Printf("\n➕ Extra in %s (%d):\n", paths[1], len(extra))
		for _, k := range extra {
			fmt.Printf("  + %s\n", k)
		}
	}
}

func collectKeys(path string) map[string]bool {
	keys := make(map[string]bool)

	files, err := discoverFiles([]string{path})
	if err != nil {
		return keys
	}

	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}

		l := mbel.NewLexer(string(content))
		p := mbel.NewParser(l)
		program := p.ParseProgram()

		for _, stmt := range program.Statements {
			if as, ok := stmt.(*mbel.AssignStatement); ok {
				keys[as.Name] = true
			}
		}
	}

	return keys
}

// ============================================================================
// IMPORT COMMAND
// ============================================================================

func importCmd(args []string) {
	fs := flag.NewFlagSet("import", flag.ExitOnError)
	output := fs.String("o", "", "Output .mbel file")
	namespace := fs.String("ns", "", "Namespace for imported keys")
	fs.Parse(args)

	files := fs.Args()
	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No JSON file specified")
		fmt.Fprintln(os.Stderr, "Usage: mbel import <file.json> [-o output.mbel]")
		os.Exit(1)
	}

	content, err := ioutil.ReadFile(files[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	var b strings.Builder

	if *namespace != "" {
		b.WriteString(fmt.Sprintf("@namespace: %s\n\n", *namespace))
	}

	// Sort keys for consistent output
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		val := data[key]
		switch v := val.(type) {
		case string:
			if strings.Contains(v, "\n") {
				b.WriteString(fmt.Sprintf("%s = \"\"\"\n%s\"\"\"\n", key, v))
			} else {
				b.WriteString(fmt.Sprintf("%s = \"%s\"\n", key, v))
			}
		default:
			// Skip non-string values
		}
	}

	result := b.String()

	if *output != "" {
		if err := ioutil.WriteFile(*output, []byte(result), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Converted %d keys to %s\n", len(keys), *output)
	} else {
		fmt.Print(result)
	}
}

// ============================================================================
// TRANSLATE COMMAND (SCAFFOLD)
// ============================================================================

func translateCmd(args []string) {
	fs := flag.NewFlagSet("translate", flag.ExitOnError)
	toLang := fs.String("to", "", "Target language code (e.g. pl, de)")
	model := fs.String("model", "gpt-4", "AI model to use")
	output := fs.String("o", "", "Output file")
	fs.Parse(args)

	if *toLang == "" {
		fmt.Fprintln(os.Stderr, "Error: --to language required")
		os.Exit(1)
	}

	files := fs.Args()
	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No input files specified")
		os.Exit(1)
	}

	fmt.Printf("🤖 Translating %d files to %s using %s...\n", len(files), *toLang, *model)

	// Simulation
	for _, file := range files {
		fmt.Printf("  Processing %s...\n", file)
		time.Sleep(500 * time.Millisecond) // Simulate work
	}

	if *output != "" {
		ioutil.WriteFile(*output, []byte("# Translated content would go here\n"), 0644)
		fmt.Printf("✓ Check %s for results (Placeholder)\n", *output)
	} else {
		fmt.Println("✓ Done (Placeholder mode - no API key configured)")
		fmt.Println("  To enable real translation, configure MBEL_OPENAI_KEY")
	}
}
