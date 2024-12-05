package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/devalexandre/mylangchaingo/agents/assistant"
	"net/http"
)

type Message struct {
	ID          string                  `json:"id,omitempty"`
	Object      string                  `json:"object,omitempty"`
	CreatedAt   int                     `json:"created_at,omitempty,omitempty"`
	ThreadId    string                  `json:"thread_id,omitempty"`
	Role        string                  `json:"role"`
	Content     string                  `json:"content"`
	AssistantId string                  `json:"assistant_id,omitempty"`
	RunId       string                  `json:"run_id,omitempty"`
	Attachments []assistant.Attachments `json:"attachments,omitempty"`
	Metadata    map[string]string       `json:"metadata,omitempty"`
}

type Response struct {
	ID      string           `json:"ID"`
	Object  string           `json:"object"`
	Data    []MessageCreated `json:"data"`
	FirstId string           `json:"first_id,omitempty"`
	LastId  string           `json:"last_id,omitempty"`
	HasMore bool             `json:"has_more,omitempty"`
	Deleted bool             `json:"deleted,omitempty"`
}

type MessageCreated struct {
	Message
	Content []struct {
		Type string `json:"type"`
		Text struct {
			Value       string        `json:"value"`
			Annotations []interface{} `json:"annotations"`
		} `json:"text"`
	} `json:"content"`
}

// NewMessageinicializa um novo assistente, opcionalmente com um ID de assistente existente.
func CreateMessage(threadID, role string, content string, opts ...MessageOption) (*Message, error) {
	message := &Message{
		Role:    role,
		Content: content,
	}

	// Aplica as opções fornecidas
	for _, opt := range opts {
		opt(message)
	}
	if len(content) == 0 {
		return message, nil

	}
	url := fmt.Sprintf("%s/threads/%s/messages", assistant.BaseURL, threadID)

	bodyJSON, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return nil, err
	}

	respBody, err := assistant.Do(req)
	var messageCreated MessageCreated
	if err := json.Unmarshal(respBody, &messageCreated); err != nil {
		return nil, err
	}

	return message, nil
}

// Returns a list of messages for a given thread.
func ListMessages(threadID string) (*Response, error) {
	url := fmt.Sprintf("%s/threads/%s/messages", assistant.BaseURL, threadID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	do, err := assistant.Do(req)

	if err != nil {
		return nil, err
	}

	var response Response
	if err := json.Unmarshal(do, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func RetrieveMessage(threadID, messageId string) (*Message, error) {
	url := fmt.Sprintf("%s/threads/%s/messages/%s", assistant.BaseURL, threadID, messageId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	do, err := assistant.Do(req)

	if err != nil {
		return nil, err
	}

	var response Message
	if err := json.Unmarshal(do, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Modifies an Thread.
func UpdateMessage(message Message) (*Message, error) {
	url := fmt.Sprintf("%s/threads/%s/messages/%s", assistant.BaseURL, message.ThreadId, message.ID)

	bodyJSON, err := json.Marshal(message)
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

	if err := json.Unmarshal(do, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

// Delete an Thread.
func DeleteMessage(threadID, messageId string) (*Response, error) {
	url := fmt.Sprintf("%s/threads/%s/messages/%s", assistant.BaseURL, threadID, messageId)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	do, err := assistant.Do(req)

	if err != nil {
		return nil, err
	}

	var response Response
	if err := json.Unmarshal(do, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
