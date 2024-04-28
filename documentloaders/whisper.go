package documentloaders

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/devalexandre/langsmithgo"
	"github.com/devalexandre/mylangchaingo"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	lgdl "github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
)

// WhisperOpenAILoader is a struct for loading and transcribing audio files using Whisper OpenAI model.
type WhisperOpenAILoader struct {
	model               string  // the model to use for transcription
	audioFilePath       string  // path to the audio file
	language            string  // language of the audio
	temperature         float64 // transcription temperature
	token               string  // authentication token for OpenAI API
	langsmithClient     *langsmithgo.Client
	langsmithgoParentId string
}

// Ensure WhisperOpenAILoader implements the Loader interface.
var _ lgdl.Loader = &WhisperOpenAILoader{}

// TranscribeAudioResponse represents the JSON response from the transcription API.
type TranscribeAudioResponse struct {
	Text string `json:"text"`
}

// WhisperOpenAIOption defines a function type for configuring WhisperOpenAILoader.
type WhisperOpenAIOption func(loader *WhisperOpenAILoader)

// NewWhisperOpenAI creates a new WhisperOpenAILoader with given API key and options.
func NewWhisperOpenAI(apiKey string, opts ...WhisperOpenAIOption) *WhisperOpenAILoader {

	loader := &WhisperOpenAILoader{
		model:       "whisper-1",
		temperature: 0.7,
		language:    "en",
		token:       apiKey,
	}
	// Apply options to configure the loader.
	for _, opt := range opts {
		opt(loader)
	}

	if os.Getenv("LANGCHAIN_TRACING") != "" && os.Getenv("LANGCHAIN_TRACING") != "false" {
		client := langsmithgo.NewClient(os.Getenv("LANGSMITH_API_KEY"))
		loader.langsmithClient = client

		if mylangchaingo.GetRunId() == "" {
			mylangchaingo.SetRunId(uuid.New().String())
		}
	}

	return loader
}

// WithModel sets the model for the WhisperOpenAILoader.
func WithModel(model string) WhisperOpenAIOption {
	return func(w *WhisperOpenAILoader) {
		w.model = model
	}
}

// WithAudioPath sets the audio file path for the WhisperOpenAILoader.
func WithAudioPath(path string) WhisperOpenAIOption {
	return func(w *WhisperOpenAILoader) {
		w.audioFilePath = path
	}
}

// WithLanguage allows setting a custom language.
// doc for language: https://platform.openai.com/docs/guides/speech-to-text/supported-languages
func WithLanguage(language string) WhisperOpenAIOption {
	return func(opts *WhisperOpenAILoader) {
		opts.language = language
	}
}

// WithTemperature sets the transcription temperature for the WhisperOpenAILoader.
func WithTemperature(temperature float64) WhisperOpenAIOption {
	return func(w *WhisperOpenAILoader) {
		w.temperature = temperature
	}
}

// WithLangsmithParentId sets the parent id for langsmith.
func WithLangsmithParentId(parentId string) WhisperOpenAIOption {
	return func(w *WhisperOpenAILoader) {
		w.langsmithgoParentId = parentId
	}
}

func (c *WhisperOpenAILoader) Load(ctx context.Context) ([]schema.Document, error) {

	if strings.Contains(c.audioFilePath, "http") {
		audioFilePath, err := downloadFileFromURL(c.audioFilePath)
		if err != nil {
			return nil, err
		}

		c.audioFilePath = audioFilePath
	}
	if c.langsmithClient != nil {
		err := c.langsmithClient.Run(&langsmithgo.RunPayload{
			Name:        "whisper - Load",
			SessionName: os.Getenv("LANGCHAIN_PROJECT_NAME"),
			RunType:     langsmithgo.Parser,
			RunID:       mylangchaingo.GetRunId(),
			ParentID:    mylangchaingo.GetParentId(),
			Inputs: map[string]interface{}{
				"prompt":      c.audioFilePath,
				"model":       c.model,
				"temperature": c.temperature,
				"language":    c.language,
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

	transcribe, err := c.transcribe(ctx, c.audioFilePath)
	if err != nil {
		return nil, err
	}

	if c.langsmithClient != nil {
		err := c.langsmithClient.Run(&langsmithgo.RunPayload{
			RunID: mylangchaingo.GetRunId(),
			Outputs: map[string]interface{}{
				"output": string(transcribe),
			},
		})

		if err != nil {
			return nil, fmt.Errorf("error running langsmith: %w", err)
		}
	}

	//update valies runId and ParentId
	mylangchaingo.SetParentId(mylangchaingo.GetRunId())
	mylangchaingo.SetRunId(uuid.New().String())

	// create a virtual file
	tmpOutputFile, err := os.CreateTemp("", "*.txt")
	if err != nil {
		return nil, fmt.Errorf("erro ao criar arquivo temporário de texto: %w", err)
	}

	defer os.Remove(tmpOutputFile.Name())

	// Write in virtual file
	if _, err := tmpOutputFile.Write(transcribe); err != nil {
		return nil, fmt.Errorf("erro ao escrever no arquivo temporário de texto: %w", err)
	}

	// read file
	file, err := os.Open(tmpOutputFile.Name())
	if err != nil {
		return nil, fmt.Errorf("erro ao ler o arquivo de texto gerado: %w", err)
	}
	txtLoader := lgdl.NewText(file)

	if c.langsmithClient != nil {
		err := c.langsmithClient.Run(&langsmithgo.RunPayload{
			Name:        fmt.Sprintf("%s-%s", os.Getenv("LANGCHAIN_PROJECT_NAME"), mylangchaingo.GetRunId()),
			SessionName: os.Getenv("LANGCHAIN_PROJECT_NAME"),
			RunType:     langsmithgo.Retriever,
			RunID:       mylangchaingo.GetRunId(),
			ParentID:    mylangchaingo.GetParentId(),
			Inputs: map[string]interface{}{
				"prompt": file.Name(),
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

	retriaver, errRetriaver := txtLoader.Load(ctx)

	if c.langsmithClient != nil {
		err := c.langsmithClient.Run(&langsmithgo.RunPayload{
			RunID: mylangchaingo.GetRunId(),
			Outputs: map[string]interface{}{
				"output": retriaver,
			},
		})

		if err != nil {
			return nil, fmt.Errorf("error running langsmith: %w", err)
		}
	}

	//update valies runId and ParentId
	mylangchaingo.SetParentId(mylangchaingo.GetRunId())

	return retriaver, errRetriaver
}

func (c *WhisperOpenAILoader) LoadAndSplit(ctx context.Context, splitter textsplitter.TextSplitter) ([]schema.Document, error) {
	docs, err := c.Load(ctx)
	if err != nil {
		return nil, err
	}

	return textsplitter.SplitDocuments(splitter, docs)
}

// downloadFileFromURL downloads a file from the provided URL and saves it to a temporary file.
// It returns the path to the temporary file and any error encountered.
//
// nolint
func downloadFileFromURL(fileURL string) (string, error) {
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	// Additional schema verification can be performed here if necessary

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", fmt.Errorf("URL scheme must be HTTP or HTTPS")
	}

	// Configuring an http.Client with timeout
	netClient := &http.Client{
		Timeout: time.Second * 10, // Set the timeout as needed
	}

	resp, err := netClient.Get(fileURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Rest of the code for file manipulation...
	tmpFile, err := os.CreateTemp("", "downloaded_file_*") // Adjust the default according to the file type
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

// transcribe performs the audio file transcription using the Whisper OpenAI model.
func (c *WhisperOpenAILoader) transcribe(ctx context.Context, audioFilePath string) ([]byte, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, err := os.Open(audioFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a form file part in the multipart writer.
	part1, err := writer.CreateFormFile("file", filepath.Base(audioFilePath))
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(part1, file); err != nil {
		return nil, err
	}

	// Add other fields to the multipart form.
	_ = writer.WriteField("model", c.model)
	_ = writer.WriteField("response_format", "json")
	_ = writer.WriteField("temperature", fmt.Sprintf("%f", c.temperature))
	_ = writer.WriteField("language", c.language)
	if err = writer.Close(); err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/audio/transcriptions", payload)
	if err != nil {
		return nil, err
	}

	// Set request headers.
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", writer.FormDataContentType()) // Correctly set the Content-Type for multipart form data.

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var transcriptionResponse TranscribeAudioResponse
	if err = json.Unmarshal(body, &transcriptionResponse); err != nil {
		return nil, err
	}
	return []byte(transcriptionResponse.Text), nil
}
