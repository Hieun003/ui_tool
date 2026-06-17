package ui_test

import (
	"os"
	"strings"
	"testing"

	"github.com/Hieun003/ui_tool/pkg/ui"
	"github.com/stretchr/testify/assert"
)

func TestRenderChatBubble(t *testing.T) {
	gold, err := os.ReadFile("testdata/chat_bubble_user.html")
	if assert.NoError(t, err) {
		got, err := ui.RenderComponent("chat_bubble", ui.ChatBubbleConfig{Message: "Hello", Avatar: "", IsUser: true})
		if assert.NoError(t, err) {
			assert.Equal(t, strings.TrimSpace(string(gold)), strings.TrimSpace(got.String()))
		}
	}
}

func TestRenderChatBubbleXSS(t *testing.T) {
	gold, err := os.ReadFile("testdata/chat_bubble_xss.html")
	if assert.NoError(t, err) {
		got, err := ui.RenderComponent("chat_bubble", ui.ChatBubbleConfig{Message: "<script>alert(1)</script>", Avatar: "", IsUser: true})
		if assert.NoError(t, err) {
			assert.Equal(t, strings.TrimSpace(string(gold)), strings.TrimSpace(got.String()))
		}
	}
}

func TestAvatarURLSchemeRejected(t *testing.T) {
	badURLs := []string{
		"javascript:alert(1)",
		"data:text/html,<h1>xss</h1>",
		"http://example.com/avatar.png",
		"//example.com/avatar.png",
	}
	for _, url := range badURLs {
		_, err := ui.RenderComponent("chat_bubble", ui.ChatBubbleConfig{Message: "hi", Avatar: url})
		assert.Errorf(t, err, "expected error for avatar URL %q", url)
		assert.Contains(t, err.Error(), "https scheme", "wrong error for %q", url)
	}
}

func TestAvatarURLSchemeAccepted(t *testing.T) {
	_, err := ui.RenderComponent("chat_bubble", ui.ChatBubbleConfig{
		Message: "hi",
		Avatar:  "https://cdn.example.com/avatar.png",
		IsUser:  true,
	})
	assert.NoError(t, err)
}

func TestNewSafeOutputSanitizes(t *testing.T) {
	input := "<div><script>alert(1)</script><iframe src='javascript:void(0)'></iframe></div>"
	got := ui.NewSafeOutput(input)
	assert.NotContains(t, got.String(), "<script")
	assert.NotContains(t, got.String(), "<iframe")
	assert.Contains(t, got.String(), "&lt;script")
	assert.Contains(t, got.String(), "&lt;iframe")
}

func TestUnsafeRawOutputPreserves(t *testing.T) {
	input := "<div><script>alert(1)</script></div>"
	got := ui.UnsafeRawOutput(input)
	assert.Equal(t, input, got.String())
}
