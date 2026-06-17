package ui

import (
	"fmt"
	"strings"

	"github.com/Hieun003/ui_tool/internal"
)

type StepItem struct {
	Title    string   `json:"title"`
	Duration string   `json:"duration,omitempty"` // e.g., "Tuần 1-2"
	Items    []string `json:"items,omitempty"`    // Bullet points
	Color    string   `json:"color,omitempty"`    // e.g., "green", "blue", "purple"
}

type StepsConfig struct {
	Steps []StepItem `json:"steps"`
}

func (StepsConfig) ComponentName() string { return "steps" }

// stepsTemplate is flat (no indentation or newlines) to prevent markdown code blocks.
const stepsTemplate = `<div class="ai-steps-container">{{range $index, $step := .Steps}}{{if $index}}<div class="ai-step-arrow"><svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 5v14M19 12l-7 7-7-7"/></svg></div>{{end}}<div class="ai-step-box ai-step-{{$step.Color}}" role="listitem"><div class="ai-step-header"><h3 class="ai-step-title">{{$step.Title}}</h3>{{if $step.Duration}}<span class="ai-step-duration">{{$step.Duration}}</span>{{end}}</div>{{if $step.Items}}<ul class="ai-step-items">{{range $step.Items}}<li class="ai-step-item">{{.}}</li>{{end}}</ul>{{end}}</div>{{end}}</div>`

func init() {
	tmpl, err := internal.MustParse("steps", stepsTemplate)
	if err != nil {
		panic(err)
	}
	internal.RegisterBuiltin("steps", StepsConfig{}, func(cfg interface{}) (string, error) {
		c, ok := cfg.(StepsConfig)
		if !ok {
			return "", fmt.Errorf("ui: steps expects %T, got %T", StepsConfig{}, cfg)
		}
		var sb strings.Builder
		if err := tmpl.Execute(&sb, c); err != nil {
			return "", fmt.Errorf("ui: rendering steps failed: %w", err)
		}
		return sb.String(), nil
	})
}
