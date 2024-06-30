package openAIAssistantRunnable

import (
	"context"
)

type EchoTool struct{}

func (e EchoTool) Name() string {
	return "echo"
}

func (e EchoTool) Description() string {
	return "Echoes the input text back to the user."
}

func (e EchoTool) Call(ctx context.Context, input string) (string, error) {
	return input, nil
}
