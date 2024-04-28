package openaiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/devalexandre/langsmithgo"
	"github.com/devalexandre/mylangchaingo"
	"github.com/google/uuid"
	"net/http"
	"os"
	"runtime"
)

const (
	defaultEmbeddingModel       = "text-embedding-ada-002"
	defaultEmbeddingModelNvidia = "NV-Embed-QA"
)

type embeddingPayload struct {
	Model     string   `json:"model"`
	Input     []string `json:"input"`
	InputType string   `json:"input_type,omitempty"`
}

type embeddingResponsePayload struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

// nolint:lll
func (c *Client) createEmbedding(ctx context.Context, payload *embeddingPayload) (*embeddingResponsePayload, error) {

	if c.baseURL == "" {
		c.baseURL = defaultBaseURL
	}

	if c.apiType == APITypeOpenAI {
		payload.Model = c.embeddingsModel
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	if c.apiType == APITypeNvidia {
		payloadNvidia := c.nvidiaEmbedding(payload)
		payloadBytes, err = json.Marshal(payloadNvidia)
		if err != nil {
			return nil, fmt.Errorf("marshal payload: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.buildURL("/embeddings", c.embeddingsModel), bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	c.setHeaders(req)

	if c.langsmithClient != nil {
		if c.langsmithgoParentId == "" {
			c.langsmithgoParentId = mylangchaingo.GetParentId()
		}
		err := c.langsmithClient.Run(&langsmithgo.RunPayload{
			Name:        "OpenAI - Create Embedding",
			SessionName: os.Getenv("LANGCHAIN_PROJECT_NAME"),
			RunType:     langsmithgo.Embedding,
			RunID:       mylangchaingo.GetRunId(),
			ParentID:    c.langsmithgoParentId,
			Inputs: map[string]interface{}{
				"Input":     payload.Input,
				"Model":     payload.Model,
				"InputType": payload.InputType,
			},
			Metadata: map[string]interface{}{
				"go_version": runtime.Version(),
				"platform":   runtime.GOOS,
				"arch":       runtime.GOARCH,
			},
		})

		if err != nil {
			return nil, fmt.Errorf("error running langsmith: %w", err)

		}
	}

	r, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("API returned unexpected status code: %d", r.StatusCode)

		// No need to check the error here: if it fails, we'll just return the
		// status code.
		var errResp errorMessage
		if err := json.NewDecoder(r.Body).Decode(&errResp); err != nil {
			return nil, errors.New(msg) // nolint:goerr113
		}

		return nil, fmt.Errorf("%s: %s", msg, errResp.Error.Message) // nolint:goerr113
	}

	var response embeddingResponsePayload

	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if c.langsmithClient != nil {
		err := c.langsmithClient.Run(&langsmithgo.RunPayload{
			RunID: mylangchaingo.GetRunId(),
			Outputs: map[string]interface{}{
				"output": response,
			},
		})

		if err != nil {
			return nil, fmt.Errorf("error running langsmith: %w", err)
		}
	}

	if c.langsmithClient != nil {
		//update valies runId and ParentId
		mylangchaingo.SetParentId(mylangchaingo.GetRunId())
		mylangchaingo.SetRunId(uuid.New().String())
	}

	return &response, nil
}

// nvidiaEmbedding is a helper function to set the correct parameters for the Nvidia API.
// It sets the correct base URL, model, and input type.
func (c *Client) nvidiaEmbedding(payload *embeddingPayload) embeddingPayload {
	payload.InputType = "query"

	if c.apiType == APITypeNvidia {
		c.baseURL = defaultEmbeddingURLNvidia
		if c.embeddingsModel == "" {
			payload.Model = defaultEmbeddingModelNvidia
		}
		if c.embeddingsModel != "" {
			payload.Model = c.embeddingsModel
		}
	}

	return *payload
}
