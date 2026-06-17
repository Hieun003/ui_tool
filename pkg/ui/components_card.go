package ui

import (
    "fmt"
    "strings"
    "github.com/Hieun003/ui_tool/internal"
)

// CardConfig implements ComponentConfig.
type CardConfig struct {
    Title    string
    Subtitle string
    ImageURL string // optional, empty if none
    Actions  []struct {
        Label string
        URL   string
    }
}

func (CardConfig) ComponentName() string { return "card" }

const cardTemplate = `<div class="ai-card" role="region" aria-label="Card">{{if .ImageURL}}<img class="ai-card-image" src="{{.ImageURL}}" alt="card image"/>{{end}}<div class="ai-card-content"><h2 class="ai-card-title">{{.Title}}</h2>{{if .Subtitle}}<h3 class="ai-card-subtitle">{{.Subtitle}}</h3>{{end}}{{if .Actions}}<div class="ai-card-actions">{{range .Actions}}<a class="ai-card-action" href="{{.URL}}">{{.Label}}</a>{{end}}</div>{{end}}</div></div>`


func init() {
    tmpl, err := internal.MustParse("card", cardTemplate)
    if err != nil {
        panic(err)
    }
    internal.RegisterBuiltin("card", func(cfg interface{}) (string, error) {
        c, ok := cfg.(CardConfig)
        if !ok {
            return "", fmt.Errorf("ui: card expects %T, got %T", CardConfig{}, cfg)
        }
        // Guard against javascript: injection in img src — html/template does not
        // block this scheme in all attribute contexts.
        if c.ImageURL != "" && !strings.HasPrefix(c.ImageURL, "https://") {
            return "", fmt.Errorf("ui: card image URL must use https scheme, got %q", c.ImageURL)
        }
        // Guard action hrefs — same javascript: risk applies to href attributes.
        for i, a := range c.Actions {
            if a.URL != "" && !strings.HasPrefix(a.URL, "https://") {
                return "", fmt.Errorf("ui: card action[%d] URL must use https scheme, got %q", i, a.URL)
            }
        }
        var sb strings.Builder
        if err := tmpl.Execute(&sb, c); err != nil {
            return "", fmt.Errorf("ui: rendering card failed: %w", err)
        }
        return sb.String(), nil
    })
}
