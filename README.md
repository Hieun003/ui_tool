# UI Tool

A lightweight, high-performance **Go module** that renders safe HTML UI component snippets for AI agents. It provides a single public package `pkg/ui` with:

- **Component rendering** (`RenderComponent`) that returns a type-safe `SafeOutput` HTML fragment.
- **Style sheet** (`StyleSheet`) – a single `<style>` tag with design-token CSS, embedded via `//go:embed`.
- **Built-in components**: `ChatBubble` and `Card`.
- **Custom component registration** via `RegisterComponent`.
- **Thread-safety & Deadlock-free design** – lock-free invocation patterns prevent deadlock.
- **Strict URL Sanitization** – rejects non-HTTPS schemes (e.g. `javascript:`) in sensitive link attributes.

---

## Installation
```bash
go get github.com/Hieun003/ui_tool/pkg/ui
```

---

## 1. Quick Start
```go
package main

import (
	"fmt"
	"github.com/Hieun003/ui_tool/pkg/ui"
)

func main() {
	// 1. Print the stylesheet (include once per page)
	fmt.Println(ui.StyleSheet())

	// 2. Render a chat bubble (user message)
	html, err := ui.RenderComponent("chat_bubble", ui.ChatBubbleConfig{
		Message: "Hello AI!",
		IsUser:  true,
	})
	if err != nil {
		panic(err)
	}
	
	// Print HTML snippet
	fmt.Println(html.String())
}
```

---

## 2. Built-in Components

### 2.1 Chat Bubble (`chat_bubble`)
Renders conversational bubbles for users or agents with optional avatars.
```go
type ChatBubbleConfig struct {
	Message string // Raw message – automatically HTML-escaped
	Avatar  string // Optional URL (must start with "https://")
	IsUser  bool   // True for user bubble (right-aligned), false for agent (left-aligned)
}
```

### 2.2 Card (`card`)
Renders rich media and action cards.
```go
type CardConfig struct {
	Title    string
	Subtitle string
	ImageURL string // Optional image URL (must start with "https://")
	Actions  []struct {
		Label string
		URL   string // Action links (must start with "https://")
	}
}
```

---

## 3. Custom Component Registration

You can extend the library by registering custom component renderers. Custom renderers must return `SafeOutput` to ensure explicit XSS safety reviews.

```go
package main

import (
	"fmt"
	"github.com/Hieun003/ui_tool/pkg/ui"
)

// Define your custom config struct. It must implement ui.ComponentConfig.
type AlertConfig struct {
	Text string
}

func (AlertConfig) ComponentName() string { return "alert" }

func main() {
	// Register the alert component
	err := ui.RegisterComponent("alert", func(cfg ui.ComponentConfig) (ui.SafeOutput, error) {
		c, ok := cfg.(AlertConfig)
		if !ok {
			return ui.SafeOutput{}, fmt.Errorf("expected AlertConfig, got %T", cfg)
		}
		
		// Use ui.NewSafeOutput for automatic HTML sanitization,
		// or ui.UnsafeRawOutput for audited raw HTML blocks.
		html := fmt.Sprintf("<div class='ai-alert'>%s</div>", c.Text)
		return ui.NewSafeOutput(html), nil
	})
	if err != nil {
		panic(err)
	}

	// Render the custom component
	out, _ := ui.RenderComponent("alert", AlertConfig{Text: "Danger Alert!"})
	fmt.Println(out.String())
}
```

---

## Golden-File Testing
The package tests output against manual golden files under `pkg/ui/testdata/` to verify layout consistency:

*   `chat_bubble_user.html` — User bubble expected output.
*   `chat_bubble_xss.html` — Expected escaped output for XSS attempts.
*   `card_default.html` — Card expected output layout.

The test suite compares rendering outputs using a canonical DOM tree parser defined in `testhelper_test.go` (ignored by production builds).

---

## Verification & Benchmarks
Ensure all changes pass concurrency race checks and strict performance thresholds ($\le 50\mu s$ per render):
```bash
# Run test suite twice with concurrency race detection
go test ./... -race -count=2

# Run render benchmarks
go test -bench=. ./...
```
