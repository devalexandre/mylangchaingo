package openAIAssistantRunnable

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/tools"
)

// Setup function to set environment variables
func setup() {
	os.Setenv("OPENAI_API_KEY", "")
}

// Test for creating a new AgentExecutor
func TestNewAgentExecutor(t *testing.T) {
	setup()

	tool := tools.Calculator{}
	assistant, err := NewAssistant(
		"Calculator Assistant",
		"You are a personal math tutor.",
		"gpt-3.5-turbo",
		[]tools.Tool{tool},
	)
	assert.NoError(t, err)
	assert.NotNil(t, assistant)

	agentExecutor := NewAgentExecutor(assistant, []tools.Tool{tool})
	assert.NotNil(t, agentExecutor)
	assert.Equal(t, assistant, agentExecutor.Agent)
	assert.Equal(t, []tools.Tool{tool}, agentExecutor.Tools)
}

// Test for invoking the AgentExecutor
func TestAgentExecutor_Invoke(t *testing.T) {
	setup()

	tool := tools.Calculator{}
	assistant, err := NewAssistant(
		"Calculator Assistant",
		"You are a personal math tutor.",
		"gpt-3.5-turbo",
		[]tools.Tool{tool},
	)
	assert.NoError(t, err)

	agentExecutor := NewAgentExecutor(assistant, []tools.Tool{tool})
	assert.NotNil(t, agentExecutor)

	input := map[string]string{"content": "What is 10 + 20?"}
	response, err := agentExecutor.Invoke(input)
	assert.NoError(t, err)
	assert.Contains(t, response, "30")
}
