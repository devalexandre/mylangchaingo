package assistant

import (
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/tools"
)

// Option define o tipo para opções funcionais para configurar o assistente.
type AssistantOption func(*Assistant)

// WithAssistantID configura o assistente com um ID existente.
func WithAssistantID(id string) AssistantOption {
	return func(a *Assistant) {
		a.ID = id
	}
}

// WithName configura o nome do assistente.
func WithName(name string) AssistantOption {
	return func(a *Assistant) {
		a.Name = name
	}
}

// WithDescription configura a descrição do assistente.
func WithDescription(description string) AssistantOption {
	return func(a *Assistant) {
		a.Description = description
	}
}

// WithInstructions configura as instruções do assistente.
func WithInstructions(instructions string) AssistantOption {
	return func(a *Assistant) {
		a.Instructions = instructions
	}
}

// WithModel configura o modelo do assistente.
func WithModel(model string) AssistantOption {
	return func(a *Assistant) {
		a.Model = model
	}
}

// WithTools configura as ferramentas do assistente.
func WithTools(tools []tools.Tool) AssistantOption {
	assistantTools := make([]llms.Tool, len(tools))
	for i, tool := range tools {
		toolsPayload := llms.FunctionDefinition{
			Name:        FormatString(tool.Name()),
			Description: tool.Description(),
			Parameters: map[string]any{
				"properties": map[string]any{
					"__arg1": map[string]string{"title": "__arg1", "type": "string"},
				},
				"required": []string{"__arg1"},
				"type":     "object",
			},
		}
		assistantTools[i].Type = "function"
		assistantTools[i].Function = &toolsPayload
	}

	return func(a *Assistant) {
		a.Tools = &assistantTools
	}

}

// WithToolResource configura o recurso da ferramenta do assistente.
func WithToolResource(toolResource ToolResource) AssistantOption {
	return func(a *Assistant) {
		a.ToolResource = &toolResource
	}
}

// WithTemperature configura a temperatura do assistente.
func WithTemperature(temperature float64) AssistantOption {
	return func(a *Assistant) {
		a.Temperature = &temperature
	}
}

// WithTopP configura o valor TopP do assistente.
func WithTopP(topP float64) AssistantOption {
	return func(a *Assistant) {
		a.TopP = &topP
	}
}

// WithMetadata configura os metadados do assistente.
func WithMetadata(metadata map[string]string) AssistantOption {
	return func(a *Assistant) {
		a.Metadata = metadata
	}
}
