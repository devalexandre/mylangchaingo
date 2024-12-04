package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/devalexandre/mylangchaingo/agents/assistant"
	"github.com/devalexandre/mylangchaingo/agents/assistant/thread"
	"net/http"
)

func CreateRun(assistantID, threadId string, opts ...Option) (*Runner, error) {

	runner := &Runner{
		AssistantId: assistantID,
	}

	for _, opt := range opts {
		opt(runner)
	}

	runnerBody, err := json.Marshal(runner)
	if err != nil {
		return nil, err
	}

	//verificar se o threadId é nulo
	thverifica, err := thread.RetrieveThread(threadId)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/threads/%s/runs", assistant.BaseURL, thverifica.ID)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(runnerBody))
	if err != nil {
		return nil, err
	}

	respBody, err := assistant.Do(req)
	if err != nil {
		return nil, err
	}

	var RunnerResponse Runner
	if err := json.Unmarshal(respBody, &RunnerResponse); err != nil {
		return nil, err
	}

	return &RunnerResponse, nil
}

func CreateThreadAndRun(assistantID string, threads thread.Thread, opts ...Option) (*Runner, error) {

	runner := &Runner{
		AssistantId: assistantID,
		Thread:      &threads,
	}

	for _, opt := range opts {
		opt(runner)
	}

	url := fmt.Sprintf("%s/threads/runs", assistant.BaseURL)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	respBody, err := assistant.Do(req)
	if err != nil {
		return nil, err
	}

	var RunnerResponse Runner
	if err := json.Unmarshal(respBody, &RunnerResponse); err != nil {
		return nil, err
	}

	return &RunnerResponse, nil
}
