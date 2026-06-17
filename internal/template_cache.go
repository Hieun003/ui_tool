package internal

import (
    "fmt"
    "html/template"
    "sync"
)

// tmplCache stores parsed templates for built‑in components.
var (
    tmplCache   = make(map[string]*template.Template)
    tmplCacheMu sync.RWMutex
)

// MustParse parses a template string and returns the compiled template.
// It does **not** panic; instead it returns an error that callers can test.
func MustParse(name, src string) (*template.Template, error) {
    tmpl, err := template.New(name).Parse(src)
    if err != nil {
        return nil, fmt.Errorf("ui: built‑in template %s parse error: %w", name, err)
    }
    tmplCacheMu.Lock()
    tmplCache[name] = tmpl
    tmplCacheMu.Unlock()
    return tmpl, nil
}

// getTemplate retrieves a cached template by name. Returns nil if not found.
func getTemplate(name string) *template.Template {
    tmplCacheMu.RLock()
    defer tmplCacheMu.RUnlock()
    return tmplCache[name]
}
