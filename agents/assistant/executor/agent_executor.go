package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/devalexandre/langsmithgo"
	"github.com/devalexandre/mylangchaingo"
	"github.com/devalexandre/mylangchaingo/agents/assistant"
	"github.com/devalexandre/mylangchaingo/agents/assistant/message"
	"github.com/devalexandre/mylangchaingo/agents/assistant/runner"
	"github.com/devalexandre/mylangchaingo/agents/assistant/thread"
	"github.com/tmc/langchaingo/tools"

	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"time"
)

// AgentExecutor is responsible for executing the agent with the provided tools
type AgentExecutor struct {
	Agent           *assistant.Assistant
	Tools           []tools.Tool
	langsmithClient *langsmithgo.Client
}

// NewAgentExecutor creates a new instance of AgentExecutor
func NewAgentExecutor(agent *assistant.Assistant, opts ...ExecutorOption) *AgentExecutor {

	agentExecutor := &AgentExecutor{
		Agent: agent,
	}

	for _, opt := range opts {
		opt(agentExecutor)
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
func (ae *AgentExecutor) Run(input string) (string, error) {
	threads, err := thread.CreateThread()
	if err != nil {
		return "", fmt.Errorf("failed to create thread: %w", err)
	}

	_, err = message.CreateMessage(threads.ID, "user", input)

	if err != nil {
		return "", fmt.Errorf("failed to add message: %w", err)
	}

	run, err := runner.CreateRun(ae.Agent.ID, threads.ID)
	if err != nil {
		return "", fmt.Errorf("failed to create run: %w", err)
	}
	response, err := ae.RetrieveThreadMessages(*run.Id, threads.ID)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve thread messages: %w", err)
	}

	return response, nil
}

// HandleToolsExecution handles the execution of tools when required

func (ae *AgentExecutor) RetrieveThreadMessages(runID, threadID string) (string, error) {
	for {
		status, toolCalls, err := ae.CheckRunStatus(threadID, runID)
		if err != nil {
			return "", err
		}

		if status == "completed" {
			break
		} else if status == "requires_action" {
			// Identifica e executa ferramentas dinamicamente
			err = ae.HandleToolsExecution(threadID, runID, toolCalls)
			if err != nil {
				return "", fmt.Errorf("failed to handle tools execution: %w", err)
			}
		} else if status == "failed" {
			return "", fmt.Errorf("run failed")
		}

		time.Sleep(1 * time.Second)
	}

	// Recupera a resposta final do agente
	messages, err := message.ListMessages(threadID)
	if err != nil {
		return "", err
	}
	for _, msg := range messages.Data {
		if msg.Role == "assistant" {
			return msg.Content[0].Text.Value, nil
		}
	}

	return "", fmt.Errorf("no assistant message found")
}

func (ae *AgentExecutor) CheckRunStatus(threadID, runID string) (string, []assistant.ToolCall, error) {
	url := fmt.Sprintf("%s/threads/%s/runs/%s", assistant.BaseURL, threadID, runID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", nil, err
	}

	respBody, err := assistant.Do(req)
	if err != nil {
		return "", nil, err
	}
	var result struct {
		Status         string `json:"status"`
		RequiredAction *struct {
			SubmitToolOutputs struct {
				ToolCalls []assistant.ToolCall `json:"tool_calls"`
			} `json:"submit_tool_outputs"`
		} `json:"required_action"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", nil, err
	}

	var toolCalls []assistant.ToolCall
	if result.RequiredAction != nil {
		toolCalls = result.RequiredAction.SubmitToolOutputs.ToolCalls
	}

	return result.Status, toolCalls, nil
}

//func (ae *AgentExecutor) HandleToolsExecution(threadID, runID string, toolCalls []assistant.ToolCall) error {
//	for _, toolCall := range toolCalls {
//		// Find the tool in the executor's tools
//
//		for _, to := range ae.Tools {
//			if to.Name() == toolCall.Function.Name {
//
//				payload, err := assistant.ExtractArg1(toolCall.Function.Arguments)
//				if err != nil {
//					return fmt.Errorf("failed to extract payload: %w", err)
//				}
//				if ae.langsmithClient != nil {
//					err := ae.langsmithClient.Run(&langsmithgo.RunPayload{
//						Name:        fmt.Sprintf("%v-%v-%v", langsmithgo.Tool, to.Name(), "AgentExecutor"),
//						SessionName: os.Getenv("LANGCHAIN_PROJECT_NAME"),
//						RunType:     langsmithgo.Tool,
//						RunID:       mylangchaingo.GetRunId(),
//						ParentID:    mylangchaingo.GetParentId(),
//						Inputs: map[string]interface{}{
//							"payload": payload,
//						},
//						Extras: map[string]interface{}{
//							"Metadata": map[string]interface{}{
//								"langsmithgo_version": "v1.0.0",
//								"go_version":          runtime.Version(),
//								"platform":            runtime.GOOS,
//								"arch":                runtime.GOARCH,
//							},
//						},
//					})
//
//					if err != nil {
//						return err
//					}
//				}
//				// Call the tool
//				toolOutput, errCall := to.Call(context.Background(), toolCall.Function.Arguments)
//				if errCall != nil {
//					return fmt.Errorf("failed to execute tool %s: %w", toolCall.Function.Name, err)
//				}
//				output, err := assistant.ExtractArg1(toolOutput)
//				if err != nil {
//					return fmt.Errorf("failed to extract output: %w", err)
//				}
//				if ae.langsmithClient != nil {
//					err := ae.langsmithClient.Run(&langsmithgo.RunPayload{
//						RunID: mylangchaingo.GetRunId(),
//						Outputs: map[string]interface{}{
//							"output": output,
//						},
//					})
//
//					if err != nil {
//						return fmt.Errorf("error running langsmith: %w", err)
//					}
//				}
//
//				// Submit the tool output
//				err = assistant.SubmitToolOutput(threadID, runID, toolCall.ID, toolOutput)
//				if err != nil {
//					return fmt.Errorf("failed to submit tool output: %w", err)
//				}
//			}
//		}
//	}
//
//	return nil
//}

func (ae *AgentExecutor) HandleToolsExecution(threadID, runID string, toolCalls []assistant.ToolCall) error {
	for _, toolCall := range toolCalls {
		tool := ae.findToolByName(toolCall.Function.Name)
		if tool == nil {
			return fmt.Errorf("tool not found: %s", toolCall.Function.Name)
		}

		// extract __arg1 received in arguments
		arg1, err := assistant.ExtractArg1(toolCall.Function.Arguments)
		if err != nil {
			return fmt.Errorf("failed to extract arg1: %w", err)
		}
		// Executa a ferramenta
		toolOutput, err := tool.Call(context.Background(), arg1)
		if err != nil {
			return fmt.Errorf("failed to execute tool %s: %w", tool.Name(), err)
		}

		// Submete a saída da ferramenta de volta ao agente
		err = assistant.SubmitToolOutput(threadID, runID, toolCall.ID, toolOutput)
		if err != nil {
			return fmt.Errorf("failed to submit tool output for tool %s: %w", tool.Name(), err)
		}
	}
	return nil
}

// Busca a ferramenta registrada pelo nome
func (ae *AgentExecutor) findToolByName(name string) tools.Tool {
	for _, tool := range ae.Tools {
		if assistant.FormatString(tool.Name()) == name {
			return tool
		}
	}
	return nil
}
