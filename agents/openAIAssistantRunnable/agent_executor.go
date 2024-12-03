package openAIAssistantRunnable

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/devalexandre/langsmithgo"
	"github.com/devalexandre/mylangchaingo"
	"github.com/google/uuid"
	"github.com/tmc/langchaingo/tools"
	"log"
	"os"
	"runtime"
)

// AgentExecutor is responsible for executing the agent with the provided tools
type AgentExecutor struct {
	Agent           *Assistant
	Tools           []tools.Tool
	langsmithClient *langsmithgo.Client
}

// NewAgentExecutor creates a new instance of AgentExecutor
func NewAgentExecutor(agent *Assistant, tools []tools.Tool) *AgentExecutor {

	agentExecutor := &AgentExecutor{
		Agent: agent,
		Tools: tools,
	}
	if os.Getenv("LANGCHAIN_TRACING") != "" && os.Getenv("LANGCHAIN_TRACING") != "false" {
		client, err := langsmithgo.NewClient()
		if err != nil {
			log.Fatal(err)
			return nil
		}
		agentExecutor.langsmithClient = client
		root := uuid.New().String()
		mylangchaingo.SetRunId(root)

	}

	return agentExecutor
}

// Run executes the agent with the provided input and returns the response
func (ae *AgentExecutor) Run(input map[string]string) (string, error) {
	threadID, err := ae.Agent.CreateThread()
	if err != nil {
		return "", fmt.Errorf("failed to create thread: %w", err)
	}

	_, err = ae.Agent.AddMessage(threadID, "user", input["content"])
	if err != nil {
		return "", fmt.Errorf("failed to add message: %w", err)
	}

	runId, err := ae.Agent.CreateRun(threadID)
	response, err := ae.Agent.RetrieveThreadMessages(runId, threadID)
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
		payload, err := extractArg1(toolCall.Function.Arguments)
		if err != nil {
			return fmt.Errorf("failed to extract payload: %w", err)
		}
		if ae.langsmithClient != nil {
			err := ae.langsmithClient.Run(&langsmithgo.RunPayload{
				Name:        fmt.Sprintf("%v-%v-%v", langsmithgo.Tool, t.Name(), "AgentExecutor"),
				SessionName: os.Getenv("LANGCHAIN_PROJECT_NAME"),
				RunType:     langsmithgo.Tool,
				RunID:       mylangchaingo.GetRunId(),
				ParentID:    mylangchaingo.GetParentId(),
				Inputs: map[string]interface{}{
					"payload": payload,
				},
				Extras: map[string]interface{}{
					"Metadata": map[string]interface{}{
						"langsmithgo_version": "v1.0.0",
						"go_version":          runtime.Version(),
						"platform":            runtime.GOOS,
						"arch":                runtime.GOARCH,
					},
				},
			})

			if err != nil {
				return err
			}
		}
		// Call the tool
		toolOutput, err := t.Call(context.Background(), toolCall.Function.Arguments)
		if err != nil {
			return fmt.Errorf("failed to execute tool %s: %w", toolCall.Function.Name, err)
		}
		output, err := extractArg1(toolOutput)
		if err != nil {
			return fmt.Errorf("failed to extract output: %w", err)
		}
		if ae.langsmithClient != nil {
			err := ae.langsmithClient.Run(&langsmithgo.RunPayload{
				RunID: mylangchaingo.GetRunId(),
				Outputs: map[string]interface{}{
					"output": output,
				},
			})

			if err != nil {
				return fmt.Errorf("error running langsmith: %w", err)
			}
		}

		// Submit the tool output
		err = ae.Agent.submitToolOutput(threadID, runID, toolCall.ID, toolOutput)
		if err != nil {
			return fmt.Errorf("failed to submit tool output: %w", err)
		}
	}

	return nil
}

func extractArg1(jsonStr string) (string, error) {
	// Cria um mapa para armazenar os dados JSON.
	var data map[string]string

	// Desserializa a string JSON no mapa.
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return "", err
	}

	// Retorna o valor de __arg1.
	val, ok := data["__arg1"]
	if !ok {
		return "", fmt.Errorf("__arg1 key not found")
	}

	return val, nil
}
