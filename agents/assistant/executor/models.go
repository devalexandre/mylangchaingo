package executor

import (
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/tools"
)

type ExecutorOption func(*AgentExecutor)

func WithTools(tools []tools.Tool) ExecutorOption {
	toolsPayload := make([]llms.Tool, len(tools))
	for i, tool := range tools {
		toolsPayload[i] = llms.Tool{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        tool.Name(),
				Description: tool.Description(),
				Parameters: map[string]any{
					"properties": map[string]any{
						"__arg1": map[string]string{"title": "__arg1", "type": "string"},
					},
					"required": []string{"__arg1"},
					"type":     "object",
				},
			},
		}
	}
	return func(a *AgentExecutor) {
		a.Tools = toolsPayload
	}

}
