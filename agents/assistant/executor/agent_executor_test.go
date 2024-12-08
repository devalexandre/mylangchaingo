package executor

import (
	assistant2 "github.com/devalexandre/mylangchaingo/agents/assistant"
	"github.com/tmc/langchaingo/tools"

	"testing"

	"github.com/devalexandre/mylangchaingo/tools/scraper/goquery"

	"github.com/stretchr/testify/assert"
)

// Setup function to set environment variables
func setup() {

}

// Test for creating a new AgentExecutor
func TestNewAgentExecutor(t *testing.T) {
	setup()

	assistant, err := assistant2.NewAssistant(
		nil,
		assistant2.WithName("Web Scraper Assistant"),
		assistant2.WithDescription("You are a personal math tutor."),
		assistant2.WithModel("gpt-3.5-turbo"),
	)
	assert.NoError(t, err)
	assert.NotNil(t, assistant)

	agentExecutor := NewAgentExecutor(assistant)
	assert.NotNil(t, agentExecutor)
	assert.Equal(t, assistant, agentExecutor.Agent)
	//assert.Equal(t, []tools.Tool{tool}, agentExecutor.Tools)
}

// Test for invoking the AgentExecutor
func TestAgentExecutor_Invoke(t *testing.T) {
	setup()

	tool, _ := goquery.New()
	assistant, err := assistant2.NewAssistant(
		assistant2.WithName("Web Scraper Assistant"),
		assistant2.WithDescription("Busca produtos em uma loja especifica, e retorna um json com nome, descrição, valor"),
		assistant2.WithModel("gpt-3.5-turbo"),
		assistant2.WithTools([]tools.Tool{tool}),
	)
	assert.NoError(t, err)

	agentExecutor := NewAgentExecutor(assistant, WithTools([]tools.Tool{tool}))
	assert.NotNil(t, agentExecutor)

	input := "Quais produtos tem em https://www.vilanova.com.br"

	response, err := agentExecutor.Run(input)
	if err != nil {
		t.Errorf("error invoking agent: %v", err)
	}

	assert.NotEmpty(t, response)
}
