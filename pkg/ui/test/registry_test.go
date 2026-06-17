package ui_test

import (
	"sync"
	"testing"

	"github.com/Hieun003/ui_tool/pkg/ui"
	"github.com/stretchr/testify/assert"
)

func TestErrorUnknownComponent(t *testing.T) {
	_, err := ui.RenderComponent("unknown", ui.ChatBubbleConfig{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ui: unknown component unknown")
}

func TestErrorConfigMismatch(t *testing.T) {
	_, err := ui.RenderComponent("chat_bubble", ui.CardConfig{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ui: chat_bubble expects")
}

type spoofConfig struct{}

func (spoofConfig) ComponentName() string { return "card" }

func TestDispatchUseNameNotComponentName(t *testing.T) {
	_, err := ui.RenderComponent("chat_bubble", spoofConfig{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "chat_bubble expects")
}

var noopOnce sync.Once

func TestConcurrentRenderAndRegister(t *testing.T) {
	const goroutines = 50
	done := make(chan struct{}, goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer func() { done <- struct{}{} }()
			_, _ = ui.RenderComponent("chat_bubble", ui.ChatBubbleConfig{Message: "ping"})
		}()
	}

	noopOnce.Do(func() {
		err := ui.RegisterComponent("concurrent_test_noop", spoofConfig{}, func(cfg ui.ComponentConfig) (ui.SafeOutput, error) {
			return ui.UnsafeRawOutput("<span>noop</span>"), nil
		})
		assert.NoError(t, err)
	})

	for i := 0; i < goroutines; i++ {
		<-done
	}
}

func TestRegisterComponentEmptyNameRejected(t *testing.T) {
	err := ui.RegisterComponent("", spoofConfig{}, func(cfg ui.ComponentConfig) (ui.SafeOutput, error) {
		return ui.SafeOutput{}, nil
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty name")
}

func TestRegisterComponentNilFnRejected(t *testing.T) {
	err := ui.RegisterComponent("murphy_nil_fn", spoofConfig{}, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nil")
}

func TestRenderComponentJSON(t *testing.T) {
	dataJSON := `{"Message": "hello json", "IsUser": true}`
	got, err := ui.RenderComponentJSON("chat_bubble", []byte(dataJSON))
	assert.NoError(t, err)
	assert.Contains(t, got.String(), "hello json")
	assert.Contains(t, got.String(), "ai-user")
}

