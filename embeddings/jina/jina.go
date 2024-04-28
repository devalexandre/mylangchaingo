package jina

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/devalexandre/langsmithgo"
	"github.com/devalexandre/mylangchaingo"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/tmc/langchaingo/embeddings"
)

type Jina struct {
	Model               string
	InputText           []string
	StripNewLines       bool
	BatchSize           int
	APIBaseURL          string
	APIKey              string
	langsmithClient     *langsmithgo.Client
	langsmithgoParentId string
}

type EmbeddingRequest struct {
	Input []string `json:"input"`
	Model string   `json:"model"`
}

type EmbeddingResponse struct {
	Model  string `json:"model"`
	Object string `json:"object"`
	Usage  struct {
		TotalTokens  int `json:"total_tokens"`
		PromptTokens int `json:"prompt_tokens"`
	} `json:"usage"`
	Data []struct {
		Object    string    `json:"object"`
		Index     int       `json:"index"`
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

var _ embeddings.Embedder = &Jina{}

func NewJina(opts ...Option) (*Jina, error) {
	v := applyOptions(opts...)

	if os.Getenv("LANGCHAIN_TRACING") != "" && os.Getenv("LANGCHAIN_TRACING") != "false" {
		client := langsmithgo.NewClient(os.Getenv("LANGSMITH_API_KEY"))
		v.langsmithClient = client
		mylangchaingo.SetRunId(uuid.New().String())
	}

	return v, nil
}

// EmbedDocuments returns a vector for each text.
func (j *Jina) EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error) {
	batchedTexts := embeddings.BatchTexts(
		embeddings.MaybeRemoveNewLines(texts, j.StripNewLines),
		j.BatchSize,
	)

	emb := make([][]float32, 0, len(texts))
	for _, batch := range batchedTexts {
		curBatchEmbeddings, err := j.CreateEmbedding(ctx, batch)
		if err != nil {
			return nil, err
		}
		emb = append(emb, curBatchEmbeddings...)
	}

	return emb, nil
}

// EmbedQuery returns a vector for a single text.
func (j *Jina) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	if j.StripNewLines {
		text = strings.ReplaceAll(text, "\n", " ")
	}

	emb, err := j.CreateEmbedding(ctx, []string{text})
	if err != nil {
		return nil, err
	}

	return emb[0], nil
}

// CreateEmbedding sends texts to the Jina API and retrieves their embeddings.
func (j *Jina) CreateEmbedding(ctx context.Context, texts []string) ([][]float32, error) {

	requestBody := EmbeddingRequest{
		Input: texts,
		Model: j.Model,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, j.APIBaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+j.APIKey)

	if j.langsmithClient != nil {
		err := j.langsmithClient.Run(&langsmithgo.RunPayload{
			Name:        "Jina - Create Embedding",
			SessionName: os.Getenv("LANGCHAIN_PROJECT_NAME"),
			RunType:     langsmithgo.Embedding,
			RunID:       mylangchaingo.GetRunId(),
			ParentID:    mylangchaingo.GetParentId(),
			Inputs: map[string]interface{}{
				"Input": texts,
				"Model": j.Model,
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

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("API request failed with status: " + resp.Status)
	}

	var embeddingResponse EmbeddingResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &embeddingResponse)
	if err != nil {
		return nil, err
	}

	embs := make([][]float32, 0, len(embeddingResponse.Data))
	for _, data := range embeddingResponse.Data {
		embs = append(embs, data.Embedding)
	}

	if j.langsmithClient != nil {
		err := j.langsmithClient.Run(&langsmithgo.RunPayload{
			RunID: mylangchaingo.GetRunId(),
			Outputs: map[string]interface{}{
				"output": embs,
			},
		})

		if err != nil {
			return nil, fmt.Errorf("error running langsmith: %w", err)
		}
	}

	//update valies runId and ParentId
	mylangchaingo.SetParentId(mylangchaingo.GetRunId())
	mylangchaingo.SetRunId(uuid.New().String())

	return embs, nil
}
