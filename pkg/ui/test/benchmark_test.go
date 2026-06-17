package ui_test

import (
	"testing"

	"github.com/Hieun003/ui_tool/pkg/ui"
	"github.com/stretchr/testify/require"
)

func BenchmarkRenderChatBubble(b *testing.B) {
	cfg := ui.ChatBubbleConfig{Message: "Hello", Avatar: "", IsUser: true}
	for i := 0; i < b.N; i++ {
		_, err := ui.RenderComponent("chat_bubble", cfg)
		require.NoError(b, err)
	}
}

func BenchmarkRenderCard(b *testing.B) {
	cfg := ui.CardConfig{Title: "Title", Subtitle: "Subtitle", ImageURL: "", Actions: nil}
	for i := 0; i < b.N; i++ {
		_, err := ui.RenderComponent("card", cfg)
		require.NoError(b, err)
	}
}
