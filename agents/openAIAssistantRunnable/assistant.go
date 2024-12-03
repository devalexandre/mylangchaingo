package openAIAssistantRunnable

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/tools"
)

const baseURL = "https://api.openai.com/v1"

type Assistant struct {
	ID     string
	Name   string
	APIKey string
	Model  string
	Tools  []tools.Tool
}

// Option defines the type for functional options to configure the assistant.
type Option func(*Assistant)

// WithAssistantID configures the assistant with an existing assistant ID.
func WithAssistantID(id string) Option {
	return func(a *Assistant) {
		a.ID = id
	}
}

// NewAssistant inicializa um novo assistente, opcionalmente com um ID de assistente existente.
func NewAssistant(name, instructions, model string, tools []tools.Tool, opts ...Option) (*Assistant, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")

	assistant := &Assistant{
		Name:   name,
		APIKey: apiKey,
		Model:  model,
		Tools:  tools,
	}

	// Aplica quaisquer opções
	for _, opt := range opts {
		opt(assistant)
	}

	// Se um ID de assistente for fornecido, não precisamos criar um novo
	if assistant.ID != "" {
		return assistant, nil
	}

	url := fmt.Sprintf("%s/assistants", baseURL)

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

	requestBody := CreateAssistantRequest{
		Instructions: instructions,
		Name:         name,
		Tools:        toolsPayload,
		Model:        model,
	}

	bodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	var assistantResponse CreateAssistantResponse
	if err := json.Unmarshal(respBody, &assistantResponse); err != nil {
		return nil, err
	}

	assistant.ID = assistantResponse.ID

	return assistant, nil
}

func (a *Assistant) CreateThread() (string, error) {
	url := fmt.Sprintf("%s/threads", baseURL)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.APIKey))
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	var threadResponse CreateThreadResponse
	if err := json.Unmarshal(respBody, &threadResponse); err != nil {
		return "", err
	}

	return threadResponse.ID, nil
}

func (a *Assistant) AddMessage(threadID, role, content string) (string, error) {
	url := fmt.Sprintf("%s/threads/%s/messages", baseURL, threadID)

	requestBody := Message{
		Role:    role,
		Content: content,
	}

	bodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.APIKey))
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	var messageResponse AddMessageResponse
	if err := json.Unmarshal(respBody, &messageResponse); err != nil {
		return "", err
	}

	return messageResponse.ID, nil
}

func (a *Assistant) CreateRun(threadID, instructions string) (string, error) {
	url := fmt.Sprintf("%s/threads/%s/runs", baseURL, threadID)

	requestBody := CreateRunRequest{
		AssistantID:  a.ID,
		Instructions: instructions,
	}

	bodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.APIKey))
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	var runResponse CreateRunResponse
	if err := json.Unmarshal(respBody, &runResponse); err != nil {
		return "", err
	}

	runID := runResponse.ID

	// Wait for the run to complete
	for {
		status, toolCalls, err := a.checkRunStatus(threadID, runID)
		if err != nil {
			return "", err
		}

		if status == "completed" {
			break
		} else if status == "requires_action" {
			err = a.HandleRequiresAction(threadID, runID, toolCalls)
			if err != nil {
				return "", fmt.Errorf("failed to handle requires_action: %w", err)
			}
		} else if status == "failed" {
			return "", fmt.Errorf("run failed")
		}

		time.Sleep(1 * time.Second)
	}

	return runID, nil
}

func (a *Assistant) CreateThreadAndRun(messages []Message) (string, error) {
	url := fmt.Sprintf("%s/threads/runs", baseURL)

	requestBody := CreateThreadAndRunRequest{
		AssistantID: a.ID,
		Thread: Thread{
			Messages: messages,
		},
	}

	bodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.APIKey))
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	var messagesResponse GetThreadMessagesResponse
	if err := json.Unmarshal(respBody, &messagesResponse); err != nil {
		return "", err
	}

	for _, message := range messagesResponse.Data {
		if message.Role == "assistant" {
			return message.Content[0].Text.Value, nil
		}
	}

	return "", fmt.Errorf("no assistant message found")

}

func (a *Assistant) checkRunStatus(threadID, runID string) (string, []ToolCall, error) {
	url := fmt.Sprintf("%s/threads/%s/runs/%s", baseURL, threadID, runID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.APIKey))
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	var result struct {
		Status         string `json:"status"`
		RequiredAction *struct {
			SubmitToolOutputs struct {
				ToolCalls []ToolCall `json:"tool_calls"`
			} `json:"submit_tool_outputs"`
		} `json:"required_action"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", nil, err
	}

	var toolCalls []ToolCall
	if result.RequiredAction != nil {
		toolCalls = result.RequiredAction.SubmitToolOutputs.ToolCalls
	}

	return result.Status, toolCalls, nil
}

func (a *Assistant) RetrieveThreadMessages(threadID string, instructions string) (string, error) {
	runID, err := a.CreateRun(threadID, instructions)
	if err != nil {
		return "", err
	}

	for {
		status, toolCalls, err := a.checkRunStatus(threadID, runID)
		if err != nil {
			return "", err
		}

		if status == "completed" {
			break
		} else if status == "requires_action" {
			// Handle tools if required
			err = a.HandleRequiresAction(threadID, runID, toolCalls)
			if err != nil {
				return "", fmt.Errorf("failed to handle requires_action: %w", err)
			}
		} else if status == "failed" {
			return "", fmt.Errorf("run failed")
		}

		time.Sleep(1 * time.Second)
	}

	url := fmt.Sprintf("%s/threads/%s/messages", baseURL, threadID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.APIKey))
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	var messagesResponse GetThreadMessagesResponse
	if err := json.Unmarshal(respBody, &messagesResponse); err != nil {
		return "", err
	}

	for _, message := range messagesResponse.Data {
		if message.Role == "assistant" {
			return message.Content[0].Text.Value, nil
		}
	}

	return "", fmt.Errorf("no assistant message found")
}

func (a *Assistant) HandleRequiresAction(threadID, runID string, toolCalls []ToolCall) error {
	for _, toolCall := range toolCalls {
		// Find the tool in the assistant's tools
		var t tools.Tool
		for _, to := range a.Tools {
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
		err = a.submitToolOutput(threadID, runID, toolCall.ID, toolOutput)
		if err != nil {
			return fmt.Errorf("failed to submit tool output: %w", err)
		}
	}

	return nil
}

func (a *Assistant) submitToolOutput(threadID, runID, toolCallID, output string) error {
	url := fmt.Sprintf("%s/threads/%s/runs/%s/submit_tool_outputs", baseURL, threadID, runID)

	requestBody := map[string]interface{}{
		"tool_outputs": []map[string]interface{}{
			{
				"tool_call_id": toolCallID,
				"output":       output,
			},
		},
	}

	bodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.APIKey))
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	var result struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return err
	}

	if result.Status != "queued" {
		return fmt.Errorf("unexpected status from submit tool output: %v", result.Status)
	}

	return nil
}
