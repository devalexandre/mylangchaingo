package thread

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/devalexandre/mylangchaingo/agents/assistant"
	"github.com/devalexandre/mylangchaingo/agents/assistant/message"
	"net/http"
)

type Thread struct {
	ID           string                 `json:"id,omitempty"`
	Object       string                 `json:"object,omitempty"`
	CreatedAt    int                    `json:"created_at,omitempty"`
	Messages     []message.Message      `json:"messages,omitempty"`
	Metadata     map[string]string      `json:"metadata,omitempty"`
	ToolResource assistant.ToolResource `json:"tool_resources,omitempty"`
}

// NewThereadinicializa um novo assistente, opcionalmente com um ID de assistente existente.
func CreateThread() (*Thread, error) {
	url := fmt.Sprintf("%s/threads", assistant.BaseURL)

	// Crie a requisição POST
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute a requisição
	respBody, err := assistant.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	// Parse a resposta
	var threadResponse Thread
	if err := json.Unmarshal(respBody, &threadResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &threadResponse, nil
}

// Retrieve assistan
func RetrieveThread(thredId string) (*Thread, error) {
	url := fmt.Sprintf("%s/threads/%s", assistant.BaseURL, thredId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	do, err := assistant.Do(req)

	if err != nil {
		return nil, err
	}

	var response Thread
	if err := json.Unmarshal(do, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Modifies an Thread.
func UpdateThread(threads Thread) (*Thread, error) {
	url := fmt.Sprintf("%s/threads/%s", assistant.BaseURL, threads.ID)

	bodyJSON, err := json.Marshal(threads)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return nil, err
	}

	do, err := assistant.Do(req)

	if err != nil {
		return nil, err
	}

	var response Thread
	if err := json.Unmarshal(do, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Delete an Thread.
func DeleteThread(thredId string) (*assistant.AssistantResponse, error) {
	url := fmt.Sprintf("%s/threads/%s", assistant.BaseURL, thredId)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	do, err := assistant.Do(req)

	if err != nil {
		return nil, err
	}

	var response assistant.AssistantResponse
	if err := json.Unmarshal(do, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
