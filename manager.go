package mbel

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Config configures the MBEL manager
type Config struct {
	DefaultLocale string
	Watch         bool // Enable hot-reloading (works only with FileRepository)
}

// Repository defines the interface for loading localization data
type Repository interface {
	// LoadAll returns a map of compiled data broken down by language
	// map[lang] -> map[key]value
	LoadAll() (map[string]map[string]interface{}, error)
}

// Manager manages localization data for multiple languages
type Manager struct {
	mu          sync.RWMutex
	runtimes    map[string]*Runtime // lang -> Runtime
	defaultLang string
	repo        Repository
}

// NewManager creates a standard file-based localization manager
func NewManager(rootPath string, cfg Config) (*Manager, error) {
	repo := &FileRepository{RootPath: rootPath}
	return NewManagerWithRepo(repo, cfg)
}

// NewManagerWithRepo creates a manager with a custom repository (e.g. Database)
func NewManagerWithRepo(repo Repository, cfg Config) (*Manager, error) {
	m := &Manager{
		runtimes:    make(map[string]*Runtime),
		defaultLang: cfg.DefaultLocale,
		repo:        repo,
	}

	if m.defaultLang == "" {
		m.defaultLang = "en"
	}

	if err := m.Load(); err != nil {
		return nil, err
	}

	if cfg.Watch {
		go m.watchLoop()
	}

	return m, nil
}

// Load (re)loads all data from the repository
func (m *Manager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	langData, err := m.repo.LoadAll()
	if err != nil {
		return err
	}

	// Create Runtimes
	newRuntimes := make(map[string]*Runtime)
	for lang, data := range langData {
		newRuntimes[lang] = NewRuntime(data)
	}

	m.runtimes = newRuntimes
	return nil
}

// Get retrieves a localized string
func (m *Manager) Get(lang, key string, args ...interface{}) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Try requested language
	if r, ok := m.runtimes[lang]; ok {
		val := r.Get(key, args...)
		if val != key {
			return val
		}
	}

	// Try partial language match (e.g. en-US -> en)
	if len(lang) > 2 {
		shortLang := lang[:2]
		if r, ok := m.runtimes[shortLang]; ok {
			val := r.Get(key, args...)
			if val != key {
				return val
			}
		}
	}

	// Try default language
	if lang != m.defaultLang {
		if r, ok := m.runtimes[m.defaultLang]; ok {
			return r.Get(key, args...)
		}
	}

	return key // Fallback to key
}

// watchLoop polls for changes
func (m *Manager) watchLoop() {
	// Only support watching if repository is file-based
	fileRepo, ok := m.repo.(*FileRepository)
	if !ok {
		// Watching not supported for non-file repos (yet)
		return
	}

	lastMod := make(map[string]time.Time)
	ticker := time.NewTicker(1 * time.Second)

	for range ticker.C {
		changed := false
		filepath.Walk(fileRepo.RootPath, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || !strings.HasSuffix(path, ".mbel") {
				return nil
			}
			if last, exists := lastMod[path]; !exists || info.ModTime().After(last) {
				lastMod[path] = info.ModTime()
				changed = true
			}
			return nil
		})

		if changed {
			// Reload in background
			m.Load()
		}
	}
}

// ============================================================================
// File Repository Implementation
// ============================================================================

// FileRepository loads MBEL files from the filesystem
type FileRepository struct {
	RootPath string
}

// LoadAll scans the directory and compiles all .mbel files
func (r *FileRepository) LoadAll() (map[string]map[string]interface{}, error) {
	langData := make(map[string]map[string]interface{})

	err := filepath.Walk(r.RootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".mbel") {
			return nil
		}

		// Determine language and namespace
		rel, _ := filepath.Rel(r.RootPath, path)
		parts := strings.Split(rel, string(os.PathSeparator))

		if len(parts) == 0 {
			return nil
		}

		lang := parts[0]
		if strings.HasSuffix(lang, ".mbel") {
			lang = strings.TrimSuffix(lang, ".mbel")
		}

		namespace := ""
		if len(parts) > 1 {
			dir := filepath.Dir(strings.Join(parts[1:], "/"))
			fname := strings.TrimSuffix(filepath.Base(path), ".mbel")

			if dir == "." {
				namespace = fname
			} else {
				namespace = strings.ReplaceAll(dir, "/", ".") + "." + fname
			}
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		l := NewLexer(string(content))
		p := NewParser(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			fmt.Fprintf(os.Stderr, "MBEL Syntax Error in %s: %v\n", path, p.Errors())
		}

		c := NewCompiler()
		res, err := c.Compile(program)
		if err != nil {
			return fmt.Errorf("compilation failed for %s: %w", path, err)
		}

		resMap, ok := res.(map[string]interface{})
		if !ok {
			return nil
		}

		if _, exists := langData[lang]; !exists {
			langData[lang] = make(map[string]interface{})
			langData[lang]["__meta"] = map[string]string{"lang": lang}
		}

		for k, v := range resMap {
			key := k
			if namespace != "" && !strings.HasPrefix(k, "__") {
				key = namespace + "." + k
			}
			langData[lang][key] = v
		}

		return nil
	})

	return langData, err
}
