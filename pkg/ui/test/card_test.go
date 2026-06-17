package ui_test

import (
	"os"
	"strings"
	"testing"

	"github.com/Hieun003/ui_tool/pkg/ui"
	"github.com/stretchr/testify/assert"
)

func TestRenderCard(t *testing.T) {
	gold, err := os.ReadFile("testdata/card_default.html")
	if assert.NoError(t, err) {
		cfg := ui.CardConfig{Title: "Title", Subtitle: "Subtitle", ImageURL: "", Actions: nil}
		got, err := ui.RenderComponent("card", cfg)
		if assert.NoError(t, err) {
			assert.Equal(t, strings.TrimSpace(string(gold)), strings.TrimSpace(got.String()))
		}
	}
}

func TestCardImageURLSchemeRejected(t *testing.T) {
	_, err := ui.RenderComponent("card", ui.CardConfig{
		Title:    "Test",
		ImageURL: "javascript:alert(1)",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "https scheme")
}

func TestCardActionURLInjectionRejected(t *testing.T) {
	cfg := ui.CardConfig{
		Title: "Test",
		Actions: []struct {
			Label string
			URL   string
		}{
			{Label: "Click", URL: "javascript:alert(document.cookie)"},
		},
	}
	_, err := ui.RenderComponent("card", cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "action[0] URL must use https scheme")
}
