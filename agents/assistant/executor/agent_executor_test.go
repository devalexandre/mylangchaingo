package executor

import (
	assistant2 "github.com/devalexandre/mylangchaingo/agents/assistant"
	c "github.com/devalexandre/mylangchaingo/tools"
	"github.com/tmc/langchaingo/tools"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Setup function to set environment variables
func setup() {

}

// Test for creating a new AgentExecutor
func TestNewAgentExecutor(t *testing.T) {
	setup()

	tool := c.Calculator{}
	assistant, err := assistant2.NewAssistant(
		assistant2.WithName("Calculator Assistant"),
		assistant2.WithDescription("You are a personal math tutor."),
		assistant2.WithModel("gpt-3.5-turbo"),
		assistant2.WithTools([]tools.Tool{tool}),
	)
	assert.NoError(t, err)
	assert.NotNil(t, assistant)

	agentExecutor := NewAgentExecutor(assistant, WithTools([]tools.Tool{tool}))
	assert.NotNil(t, agentExecutor)
	assert.Equal(t, assistant, agentExecutor.Agent)
	assert.Equal(t, []tools.Tool{tool}, agentExecutor.Tools)
}

// Test for invoking the AgentExecutor
func TestAgentExecutor_Invoke(t *testing.T) {
	setup()

	tool := c.Calculator{}
	assistant, err := assistant2.NewAssistant(
		assistant2.WithName("Calculator Assistant"),
		assistant2.WithDescription("You are a personal math tutor."),
		assistant2.WithModel("gpt-3.5-turbo"),
	)
	assert.NoError(t, err)

	agentExecutor := NewAgentExecutor(assistant, WithTools([]tools.Tool{tool}))
	assert.NotNil(t, agentExecutor)

	input := "What is 10 + 20?"

	response, err := agentExecutor.Run(input)
	assert.NoError(t, err)
	assert.Contains(t, response, "30")
}
