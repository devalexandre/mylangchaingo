package openAIAssistantRunnable

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/tools"
)

const instructions = `You are an assistant equipped with a powerful calculator tool. You can solve mathematical expressions provided by the user. When the user asks a math question, use the calculator tool to evaluate the expression and return the result.

Examples:
- User: What is 5 + 7?
- Assistant: 5 + 7 = 12

- User: Can you calculate (3 * (2 + 4)) / 3?
- Assistant: (3 * (2 + 4)) / 3 = 6`

func SetUp() {
	os.Setenv("OPENAI_API_KEY", "")
}

func TestMain(m *testing.M) {
	SetUp()
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestNewAssistant(t *testing.T) {
	tool := tools.Calculator{}
	tools := []tools.Tool{tool}

	assistant, err := NewAssistant(
		"Test Assistant",
		"This is a test assistant.",
		"gpt-3.5-turbo-0125",
		tools,
	)
	assert.NoError(t, err)
	assert.NotNil(t, assistant)
	assert.Equal(t, "Test Assistant", assistant.Name)
	assert.Equal(t, "gpt-3.5-turbo-0125", assistant.Model)
}

func TestAssistant_CreateThread(t *testing.T) {
	tool := tools.Calculator{}
	tools := []tools.Tool{tool}

	assistant, err := NewAssistant(
		"Test Assistant",
		"This is a test assistant.",
		"gpt-3.5-turbo-0125",
		tools,
	)
	assert.NoError(t, err)

	threadID, err := assistant.CreateThread()
	assert.NoError(t, err)
	assert.NotEmpty(t, threadID)
}

func TestAssistant_AddMessage(t *testing.T) {
	tool := tools.Calculator{}
	tools := []tools.Tool{tool}

	assistant, err := NewAssistant(
		"Test Assistant",
		"This is a test assistant.",
		"gpt-3.5-turbo-0125",
		tools,
	)
	assert.NoError(t, err)

	threadID, err := assistant.CreateThread()
	assert.NoError(t, err)

	messageID, err := assistant.AddMessage(threadID, "user", "Hello, assistant!")
	assert.NoError(t, err)
	assert.NotEmpty(t, messageID)
}

func TestAssistant_GenerateResponse(t *testing.T) {
	tool := tools.Calculator{}
	tools := []tools.Tool{tool}

	assistant, err := NewAssistant(
		"Test Assistant",
		"This is a test assistant.",
		"gpt-3.5-turbo-0125",
		tools,
	)
	assert.NoError(t, err)

	threadID, err := assistant.CreateThread()
	assert.NoError(t, err)

	messageID, err := assistant.AddMessage(threadID, "user", "2 + 2")
	assert.NoError(t, err)
	assert.NotEmpty(t, messageID)

	response, err := assistant.RetrieveThreadMessages(threadID, "")
	assert.NoError(t, err)
	assert.NotEmpty(t, response)
	assert.Contains(t, response, "4")
}

func TestAssistant_HandleRequiresAction(t *testing.T) {
	tool := tools.Calculator{}
	tools := []tools.Tool{tool}

	assistant, err := NewAssistant(
		"Test Assistant",
		"This is a test assistant.",
		"gpt-3.5-turbo-0125",
		tools,
	)
	assert.NoError(t, err)

	threadID, err := assistant.CreateThread()
	assert.NoError(t, err)

	messageID, err := assistant.AddMessage(threadID, "user", "2 + 2")
	assert.NoError(t, err)
	assert.NotEmpty(t, messageID)

	runID, err := assistant.CreateRun(threadID)
	assert.NoError(t, err)
	assert.NotEmpty(t, runID)
}

func TestAssistant_RunFailure(t *testing.T) {
	tool := tools.Calculator{}
	tools := []tools.Tool{tool}

	assistant, err := NewAssistant(
		"Test Assistant",
		"This is a test assistant.",
		"gpt-3.5-turbo-0125",
		tools,
	)
	assert.NoError(t, err)

	threadID, err := assistant.CreateThread()
	assert.NoError(t, err)

	// Simulando falha na execução
	_, err = assistant.CreateRun(threadID)
	assert.Error(t, err)
}

func TestAssistant_CheckRunStatus(t *testing.T) {
	tool := tools.Calculator{}
	tools := []tools.Tool{tool}

	assistant, err := NewAssistant(
		"Test Assistant",
		"This is a test assistant.",
		"gpt-3.5-turbo-0125",
		tools,
	)
	assert.NoError(t, err)

	threadID, err := assistant.CreateThread()
	assert.NoError(t, err)

	runID, err := assistant.CreateRun(threadID)
	assert.NoError(t, err)
	assert.NotEmpty(t, runID)

	status, toolCalls, err := assistant.checkRunStatus(threadID, runID)
	assert.NoError(t, err)
	assert.NotEmpty(t, status)
	assert.NotNil(t, toolCalls)
}

func TestAssistant_SubmitToolOutput(t *testing.T) {
	tool := tools.Calculator{}
	tools := []tools.Tool{tool}

	assistant, err := NewAssistant(
		"Calculator",
		instructions,
		"gpt-3.5-turbo-0125",
		tools,
	)
	assert.NoError(t, err)

	threadID, err := assistant.CreateThread()
	assert.NoError(t, err)

	messageID, err := assistant.AddMessage(threadID, "user", "quanto 25 + 54 ?")
	assert.NoError(t, err)
	assert.NotEmpty(t, messageID)

	response, err := assistant.RetrieveThreadMessages(threadID, "")
	assert.NoError(t, err)

	assert.NotEmpty(t, response)
}

func TestAssistant_CreateThreadAndRun(t *testing.T) {

	assistant, err := NewAssistant(
		"Burguer Beer",
		instructions,
		"gpt-3.5-turbo-0125",
		nil,
		WithAssistantID("asst_qSusDWoOKM3lJFEabFqO7j4w"),
	)

	assert.NoError(t, err)

	messages := []Message{
		{
			Role:    "user",
			Content: "Quais os lanches mais baratos?",
		},
	}

	response, err := assistant.CreateThreadAndRun(messages)
	assert.NoError(t, err)

	assert.NotEmpty(t, response)

}
>>>>>>> 58e74fd (chore: update assistente create run)
