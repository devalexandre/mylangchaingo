package openAIAssistantRunnable

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/tools"
)

// AgentExecutor is responsible for executing the agent with the provided tools
type AgentExecutor struct {
	Agent *Assistant
	Tools []tools.Tool
}

// NewAgentExecutor creates a new instance of AgentExecutor
func NewAgentExecutor(agent *Assistant, tools []tools.Tool) *AgentExecutor {
	return &AgentExecutor{
		Agent: agent,
		Tools: tools,
	}
}

// Invoke executes the agent with the provided input and returns the response
func (ae *AgentExecutor) Invoke(input map[string]string) (string, error) {
	threadID, err := ae.Agent.CreateThread()
	if err != nil {
		return "", fmt.Errorf("failed to create thread: %w", err)
	}

	_, err = ae.Agent.AddMessage(threadID, "user", input["content"])
	if err != nil {
		return "", fmt.Errorf("failed to add message: %w", err)
	}

	response, err := ae.Agent.RetrieveThreadMessages(threadID, input["content"])
	if err != nil {
		return "", fmt.Errorf("failed to retrieve thread messages: %w", err)
	}

	return response, nil
}

// HandleToolsExecution handles the execution of tools when required
func (ae *AgentExecutor) HandleToolsExecution(threadID, runID string, toolCalls []ToolCall) error {
	for _, toolCall := range toolCalls {
		// Find the tool in the executor's tools
		var t tools.Tool
		for _, to := range ae.Tools {
			if to.Name() == toolCall.Function.Name {
				t = to
				break
			}
		}
		if t == nil {
			return fmt.Errorf("tool %s not found", toolCall.Function.Name)
		}

		// Call the tool
		toolOutput, err := t.Call(context.Background(), toolCall.Function.Arguments)
		if err != nil {
			return fmt.Errorf("failed to execute tool %s: %w", toolCall.Function.Name, err)
		}

		// Submit the tool output
		err = ae.Agent.submitToolOutput(threadID, runID, toolCall.ID, toolOutput)
		if err != nil {
			return fmt.Errorf("failed to submit tool output: %w", err)
		}
	}

	return nil
}
