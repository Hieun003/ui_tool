package ui

import (
	"fmt"
	"strings"

	"github.com/Hieun003/ui_tool/internal"
)

// ComponentConfig is a marker interface for component configuration structs.
// The method ComponentName is **ignored at runtime** – it exists only so that
// external packages can implement the interface and get compile‑time type safety.
type ComponentConfig interface {
	ComponentName() string
}

// SafeOutput wraps rendered HTML that has either been auto-escaped or explicitly sanitized.
type SafeOutput struct {
	html string
}

// String returns the underlying safe HTML content.
func (s SafeOutput) String() string {
	return s.html
}

// NewSafeOutput constructs a SafeOutput by running a basic sanitization pass over the input HTML.
// This is the default constructor for custom renderers returning dynamic HTML.
func NewSafeOutput(s string) SafeOutput {
	return SafeOutput{html: sanitize(s)}
}

// UnsafeRawOutput constructs a SafeOutput bypass check.
// Using this function is explicit and easily auditable during code review.
func UnsafeRawOutput(s string) SafeOutput {
	return SafeOutput{html: s}
}

// basic sanitization for custom HTML elements (blocks script/iframe)
func sanitize(s string) string {
	r := strings.NewReplacer(
		"<script", "&lt;script",
		"</script", "&lt;/script",
		"<iframe", "&lt;iframe",
		"</iframe", "&lt;/iframe",
	)
	return r.Replace(s)
}

// RenderComponent renders a known component to a SafeOutput.
// name – the component name (e.g. "chat_bubble", "card").
// cfg – a struct that implements ComponentConfig.
func RenderComponent(name string, cfg ComponentConfig) (SafeOutput, error) {
	// Dispatcher must use the *name* argument, never cfg.ComponentName().
	renderer, err := internal.GetRenderer(name)
	if err != nil {
		return SafeOutput{}, err
	}
	// Underlying renderer expects the concrete config type; it will perform its
	// own comma‑ok check and return a descriptive error if mismatched.
	out, err := renderer(cfg)
	if err != nil {
		return SafeOutput{}, err
	}
	// Built-in renderers use html/template and are pre-escaped safely.
	return SafeOutput{html: out}, nil
}

// RegisterComponent registers a custom renderer under the given name.
// name may not clash with a built‑in component.
func RegisterComponent(name string, fn func(cfg ComponentConfig) (SafeOutput, error)) error {
	if fn == nil {
		return fmt.Errorf("ui: renderer for %s is nil", name)
	}
	// Wrap the user function to the internal renderer signature.
	wrapped := func(cfg interface{}) (string, error) {
		if cc, ok := cfg.(ComponentConfig); ok {
			so, err := fn(cc)
			if err != nil {
				return "", err
			}
			return so.html, nil
		}
		return "", fmt.Errorf("ui: %s expects ComponentConfig, got %T", name, cfg)
	}
	return internal.RegisterCustom(name, wrapped)
}

