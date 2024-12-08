package assistant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// NewAssistant inicializa um novo assistente, opcionalmente com um ID de assistente existente.
func NewAssistant(opts ...AssistantOption) (*Assistant, error) {
	assistant := &Assistant{}

	// Aplica quaisquer opções
	for _, opt := range opts {
		opt(assistant)
	}

	// Se um ID de assistente for fornecido, não precisamos criar um novo
	if assistant.ID != "" {
		return assistant, nil
	}

	url := fmt.Sprintf("%s/assistants", BaseURL)

	bodyJSON, err := json.Marshal(assistant)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return nil, err
	}
	respBody, err := Do(req)
	if err != nil {
		return nil, err
	}
	var assistantResponse Assistant
	if err := json.Unmarshal(respBody, &assistantResponse); err != nil {
		return nil, err
	}

	assistant.ID = assistantResponse.ID

	return assistant, nil
}

// Returns a list of assistants.
func (a *Assistant) ListAssistants() ([]Assistant, error) {
	url := fmt.Sprintf("%s/assistants", BaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	respBody, err := Do(req)
	if err != nil {
		return nil, err
	}

	var response AssistantResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

// Retrieve assistan
func (a *Assistant) RetrieveAssistant() (*Assistant, error) {
	url := fmt.Sprintf("%s/assistants/%s", BaseURL, a.ID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	do, err := Do(req)

	if err != nil {
		return nil, err
	}

	var response Assistant
	if err := json.Unmarshal(do, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Modifies an assistant.
func (a *Assistant) UpdateAssistant(assistnt Assistant) (*Assistant, error) {
	url := fmt.Sprintf("%s/assistants/%s", BaseURL, a.ID)

	bodyJSON, err := json.Marshal(assistnt)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return nil, err
	}

	do, err := Do(req)

	if err != nil {
		return nil, err
	}

	var response Assistant
	if err := json.Unmarshal(do, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Delete an assistant.
func (a *Assistant) DeleteAssistant() (*AssistantResponse, error) {
	url := fmt.Sprintf("%s/assistants/%s", BaseURL, a.ID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	do, err := Do(req)

	if err != nil {
		return nil, err
	}

	var response AssistantResponse
	if err := json.Unmarshal(do, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
