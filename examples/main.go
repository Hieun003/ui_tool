package main

import (
	"fmt"
	"github.com/Hieun003/ui_tool/pkg/ui"
)

func main() {
	// Chat bubble user
	chatHTML, err := ui.RenderComponent("chat_bubble", ui.ChatBubbleConfig{Message: "Hello", Avatar: "", IsUser: true})
	if err != nil {
		panic(err)
	}
	fmt.Println("--- CHAT BUBBLE SNIPPET ---")
	fmt.Println(chatHTML.String())

	// Card default
	cardHTML, err := ui.RenderComponent("card", ui.CardConfig{Title: "Title", Subtitle: "Subtitle", ImageURL: "", Actions: nil})
	if err != nil {
		panic(err)
	}
	fmt.Println("\n--- CARD SNIPPET ---")
	fmt.Println(cardHTML.String())
}
