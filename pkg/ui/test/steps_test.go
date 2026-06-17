package ui_test

import (
	"os"
	"strings"
	"testing"

	"github.com/Hieun003/ui_tool/pkg/ui"
	"github.com/stretchr/testify/assert"
)

func TestRenderSteps(t *testing.T) {
	gold, err := os.ReadFile("testdata/steps_default.html")
	if assert.NoError(t, err) {
		cfg := ui.StepsConfig{
			Steps: []ui.StepItem{
				{
					Title:    "Phase 1",
					Duration: "Week 1",
					Items:    []string{"Item 1"},
					Color:    "green",
				},
				{
					Title:    "Phase 2",
					Duration: "Week 2",
					Items:    []string{"Item 2"},
					Color:    "blue",
				},
			},
		}
		got, err := ui.RenderComponent("steps", cfg)
		if assert.NoError(t, err) {
			assert.Equal(t, strings.TrimSpace(string(gold)), strings.TrimSpace(got.String()))
		}
	}
}
