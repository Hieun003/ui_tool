package ui

import (
    "fmt"
    "strings"
    "github.com/Hieun003/ui_tool/internal"
)

// ChatBubbleConfig implements ComponentConfig.
type ChatBubbleConfig struct {
    Message string // raw user message – will be escaped automatically
    Avatar  string // optional URL, empty if not provided
    IsUser  bool   // true = user bubble (right aligned), false = agent bubble (left)
}

func (ChatBubbleConfig) ComponentName() string { return "chat_bubble" }

// chatBubbleTemplate is parsed at init and cached.
const chatBubbleTemplate = `<div class="ai-chat-bubble {{if .IsUser}}ai-user{{else}}ai-agent{{end}}" role="dialog" aria-label="{{if .IsUser}}User{{else}}Agent{{end}} Message">{{if .Avatar}}<img class="ai-avatar" src="{{.Avatar}}" alt="avatar"/>{{end}}<div class="ai-message">{{.Message}}</div></div>`


func init() {
    // Parse and cache the template.
    tmpl, err := internal.MustParse("chat_bubble", chatBubbleTemplate)
    if err != nil {
        panic(err)
    }
    // Register the renderer.
    internal.RegisterBuiltin("chat_bubble", func(cfg interface{}) (string, error) {
        // Type‑assert safely.
        c, ok := cfg.(ChatBubbleConfig)
        if !ok {
            return "", fmt.Errorf("ui: chat_bubble expects %T, got %T", ChatBubbleConfig{}, cfg)
        }
        // Guard: html/template does NOT block javascript: scheme in src attributes
        // in every context. Validate here at the renderer level.
        if c.Avatar != "" && !strings.HasPrefix(c.Avatar, "https://") {
            return "", fmt.Errorf("ui: avatar URL must use https scheme, got %q", c.Avatar)
        }
        var sb strings.Builder
        if err := tmpl.Execute(&sb, c); err != nil {
            return "", fmt.Errorf("ui: rendering chat_bubble failed: %w", err)
        }
        return sb.String(), nil
    })
}
