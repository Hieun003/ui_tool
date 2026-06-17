package internal

import (
    "fmt"
    "sync"
)

type rendererFunc func(cfg interface{}) (string, error)

// builtinRenderers holds the immutable built‑in component renderers.
var builtinRenderers map[string]rendererFunc

// prototypes holds configuration prototypes for dynamic unmarshaling, protected by RWMutex.
var (
    prototypes        = make(map[string]interface{})
    prototypesMu      sync.RWMutex
    customRenderers   = make(map[string]rendererFunc)
    customRenderersMu sync.RWMutex
)

func init() {
    // This map will be filled by the ui package during its init.
    builtinRenderers = make(map[string]rendererFunc)
}

// RegisterBuiltin registers a built-in renderer and prototype. Called only from ui package init.
// Panics on duplicate name, nil fn, or empty name — all are programming errors
// caught at startup, never at runtime.
func RegisterBuiltin(name string, proto interface{}, fn rendererFunc) {
    if name == "" {
        panic("ui: RegisterBuiltin called with empty name")
    }
    if fn == nil {
        panic(fmt.Sprintf("ui: RegisterBuiltin called with nil fn for %q", name))
    }
    if _, exists := builtinRenderers[name]; exists {
        panic(fmt.Sprintf("builtin renderer %s already registered", name))
    }
    if proto == nil {
        panic(fmt.Sprintf("ui: RegisterBuiltin called with nil prototype for %q", name))
    }
    builtinRenderers[name] = fn
    
    prototypesMu.Lock()
    prototypes[name] = proto
    prototypesMu.Unlock()
}


// RegisterCustom registers a custom renderer and prototype. Returns an error if the name
// collides with a built-in, the fn is nil, or the name is empty.
func RegisterCustom(name string, proto interface{}, fn rendererFunc) error {
    if name == "" {
        return fmt.Errorf("ui: RegisterCustom called with empty name")
    }
    if fn == nil {
        return fmt.Errorf("ui: renderer for %s is nil", name)
    }
    if proto == nil {
        return fmt.Errorf("ui: RegisterCustom called with nil prototype for %s", name)
    }
    if _, exists := builtinRenderers[name]; exists {
        return fmt.Errorf("ui: cannot override built-in component %s", name)
    }
    
    customRenderersMu.Lock()
    defer customRenderersMu.Unlock()
    if _, exists := customRenderers[name]; exists {
        return fmt.Errorf("ui: custom component %s already registered", name)
    }
    customRenderers[name] = fn
    
    prototypesMu.Lock()
    prototypes[name] = proto
    prototypesMu.Unlock()
    return nil
}

// GetPrototype returns a copy/instance of the registered prototype for a component.
func GetPrototype(name string) (interface{}, error) {
    prototypesMu.RLock()
    proto, ok := prototypes[name]
    prototypesMu.RUnlock()
    if !ok {
        return nil, fmt.Errorf("ui: unknown component %s", name)
    }
    return proto, nil
}

// getRenderer returns the appropriate renderer for a component name.
// IMPORTANT: the lock is released *before* the fn pointer is returned so the
// caller can invoke fn without holding any lock. This prevents a deadlock if
// a custom renderer tries to call RegisterCustom (which acquires a write lock).
func getRenderer(name string) (rendererFunc, error) {
    if fn, ok := builtinRenderers[name]; ok {
        return fn, nil // builtinRenderers is written only during init – no lock needed
    }
    customRenderersMu.RLock()
    fn, ok := customRenderers[name]
    customRenderersMu.RUnlock() // explicit release; never hold lock across renderer call
    if ok {
        return fn, nil
    }
    return nil, fmt.Errorf("ui: unknown component %s", name)
}

// GetRenderer is the exported version used by the public API.
func GetRenderer(name string) (rendererFunc, error) {
    return getRenderer(name)
}

