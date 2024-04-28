package maritaca

import (
	"context"
	"errors"
	"fmt"
	"github.com/devalexandre/langsmithgo"
	"github.com/devalexandre/mylangchaingo"
	"github.com/devalexandre/mylangchaingo/llms/maritaca/internal/maritacaclient"
	"github.com/google/uuid"
	"net/http"
	"os"
	"runtime"

	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/llms"
)

var (
	ErrEmptyResponse       = errors.New("no response")
	ErrIncompleteEmbedding = errors.New("no all input got emmbedded")
)

// LLM is a maritaca LLM implementation.
type LLM struct {
	CallbacksHandler callbacks.Handler
	client           *maritacaclient.Client
	options          options
	langsmithClient  *langsmithgo.Client
	langsmithRunId   string
}

var _ llms.Model = (*LLM)(nil)

// New creates a new maritaca LLM implementation.
func New(opts ...Option) (*LLM, error) {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	if o.httpClient == nil {
		o.httpClient = http.DefaultClient
	}

	client, err := maritacaclient.NewClient(o.httpClient)
	if err != nil {
		return nil, err
	}
	llms := &LLM{client: client, options: o}
	if os.Getenv("LANGCHAIN_TRACING") != "" && os.Getenv("LANGCHAIN_TRACING") != "false" {
		client := langsmithgo.NewClient(os.Getenv("LANGSMITH_API_KEY"))
		llms.langsmithClient = client
	}

	err = llms.setUpLangsmithClient()
	if err != nil {
		return nil, err

	}

	return llms, nil
}

// Call Implement the call interface for LLM.
func (o *LLM) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return llms.GenerateFromSinglePrompt(ctx, o, prompt, options...)
}

// GenerateContent implements the Model interface.
// nolint: goerr113
func (o *LLM) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) { // nolint: lll, cyclop, funlen
	if o.CallbacksHandler != nil {
		o.CallbacksHandler.HandleLLMGenerateContentStart(ctx, messages)
	}

	opts := llms.CallOptions{}
	for _, opt := range options {
		opt(&opts)
	}

	// Override LLM model if set as llms.CallOption
	model := o.options.model
	if opts.Model != "" {
		model = opts.Model
	}

	// Our input is a sequence of MessageContent, each of which potentially has
	// a sequence of Part that could be text, images etc.
	// We have to convert it to a format maritaca undestands: ChatRequest, which
	// has a sequence of Message, each of which has a role and content - single
	// text + potential images.
	chatMsgs := make([]*maritacaclient.Message, 0, len(messages))
	for _, mc := range messages {
		msg := &maritacaclient.Message{Role: typeToRole(mc.Role)}

		// Look at all the parts in mc; expect to find a single Text part and
		// any number of binary parts.
		var text string
		foundText := false

		for _, p := range mc.Parts {
			switch pt := p.(type) {
			case llms.TextContent:
				if foundText {
					return nil, errors.New("expecting a single Text content")
				}
				foundText = true
				text = pt.Text

			default:
				return nil, errors.New("only support Text and BinaryContent parts right now")
			}
		}

		msg.Content = text

		chatMsgs = append(chatMsgs, msg)
	}

	format := o.options.format
	if opts.JSONMode {
		format = "json"
	}

	// Get our maritacaOptions from llms.CallOptions
	maritacaOptions := makemaritacaOptionsFromOptions(o.options.maritacaOptions, opts)
	req := &maritacaclient.ChatRequest{
		Model:    model,
		Format:   format,
		Messages: chatMsgs,
		Options:  maritacaOptions,
		Stream:   func(b bool) *bool { return &b }(opts.StreamingFunc != nil),
	}

	var fn maritacaclient.ChatResponseFunc
	streamedResponse := ""
	var resp maritacaclient.ChatResponse

	fn = func(response maritacaclient.ChatResponse) error {
		if opts.StreamingFunc != nil && response.Text != "" {
			if err := opts.StreamingFunc(ctx, []byte(response.Text)); err != nil {
				return err
			}
		}
		switch response.Event {
		case "message":
			streamedResponse += response.Text
		case "end":
			resp.Answer = streamedResponse
		case "nostream":
			resp = response
		}

		return nil
	}
	o.client.Token = o.options.maritacaOptions.Token

	o.options.langsmithgoRunId = mylangchaingo.GetRunId()

	if o.options.langsmithgoParentId == "" {
		o.options.langsmithgoParentId = o.langsmithRunId
	}

	if o.langsmithClient != nil {

		err := o.langsmithClient.Run(&langsmithgo.RunPayload{
			Name:        "MariatacaAI - GenerateContent",
			SessionName: os.Getenv("LANGCHAIN_PROJECT_NAME"),
			RunType:     langsmithgo.LLM,
			RunID:       o.options.langsmithgoRunId,
			ParentID:    o.options.langsmithgoParentId,
			Inputs: map[string]interface{}{
				"payload": req,
			},
			Metadata: map[string]interface{}{
				"go_version": runtime.Version(),
				"platform":   runtime.GOOS,
				"arch":       runtime.GOARCH,
			},
		})
		o.options.langsmithgoParentId = o.options.langsmithgoRunId

		if err != nil {
			return nil, fmt.Errorf("error running langsmith: %w", err)

		}
	}
	err := o.client.Generate(ctx, req, fn)
	if err != nil {
		if o.CallbacksHandler != nil {
			o.CallbacksHandler.HandleLLMError(ctx, err)
		}
		return nil, err
	}

	choices := createChoice(resp)

	response := &llms.ContentResponse{Choices: choices}

	if o.langsmithClient != nil {
		err := o.langsmithClient.Run(&langsmithgo.RunPayload{
			RunID: o.options.langsmithgoRunId,
			Outputs: map[string]interface{}{
				"output": response,
			},
		})

		if err != nil {
			return nil, fmt.Errorf("error running langsmith: %w", err)
		}
	}

	if o.langsmithClient != nil {
		//update valies runId and ParentId
		mylangchaingo.SetParentId(mylangchaingo.GetRunId())
		mylangchaingo.SetRunId(uuid.New().String())
	}

	if o.CallbacksHandler != nil {
		o.CallbacksHandler.HandleLLMGenerateContentEnd(ctx, response)
	}

	return response, nil
}

func typeToRole(typ llms.ChatMessageType) string {
	switch typ {
	case llms.ChatMessageTypeSystem:
		return "system"
	case llms.ChatMessageTypeAI:
		return "assistant"
	case llms.ChatMessageTypeHuman:
		fallthrough
	case llms.ChatMessageTypeGeneric:
		return "user"
	case llms.ChatMessageTypeFunction:
		return "function"
	case llms.ChatMessageTypeTool:
		return "tool"
	}
	return ""
}

func makemaritacaOptionsFromOptions(maritacaOptions maritacaclient.Options, opts llms.CallOptions) maritacaclient.Options {
	// Load back CallOptions as maritacaOptions
	maritacaOptions.MaxTokens = opts.MaxTokens
	maritacaOptions.Model = opts.Model
	maritacaOptions.TopP = opts.TopP
	maritacaOptions.RepetitionPenalty = opts.RepetitionPenalty
	maritacaOptions.StoppingTokens = opts.StopWords
	maritacaOptions.Stream = opts.StreamingFunc != nil

	return maritacaOptions
}

func createChoice(resp maritacaclient.ChatResponse) []*llms.ContentChoice {
	return []*llms.ContentChoice{
		{
			Content: resp.Answer,
			GenerationInfo: map[string]any{
				"CompletionTokens": resp.Metrics.Usage.CompletionTokens,
				"PromptTokens":     resp.Metrics.Usage.PromptTokens,
				"TotalTokens":      resp.Metrics.Usage.TotalTokens,
			},
		},
	}
}

func (o *LLM) setUpLangsmithClient() error {
	if o.langsmithClient != nil {

		mylangchaingo.SetRunId(uuid.New().String())

		if o.langsmithRunId == "" {
			o.langsmithRunId = mylangchaingo.GetRunId()
		}

		if o.options.langsmithgoParentId == "" {
			o.options.langsmithgoParentId = mylangchaingo.GetParentId()
		}

		if o.options.langsmithgoRunId == "" {
			o.options.langsmithgoRunId = o.langsmithRunId
		}

		err := o.langsmithClient.Run(&langsmithgo.RunPayload{
			Name:        "MaritacaAI",
			SessionName: os.Getenv("LANGCHAIN_PROJECT_NAME"),
			RunType:     langsmithgo.LLM,
			RunID:       o.langsmithRunId,
			ParentID:    o.options.langsmithgoParentId,
			Inputs: map[string]interface{}{
				"payload": nil,
			},
			Metadata: map[string]interface{}{
				"go_version": runtime.Version(),
				"platform":   runtime.GOOS,
				"arch":       runtime.GOARCH,
			},
		})
		o.options.langsmithgoParentId = o.langsmithRunId

		if err != nil {
			return fmt.Errorf("error running langsmith: %w", err)

		}

		err = o.langsmithClient.Run(&langsmithgo.RunPayload{
			RunID: o.options.langsmithgoRunId,
			Outputs: map[string]interface{}{
				"output": "",
			},
		})

		if err != nil {
			return fmt.Errorf("error running langsmith: %w", err)
		}

		//update valies runId and ParentId
		mylangchaingo.SetParentId(mylangchaingo.GetRunId())

		if o.options.langsmithgoParentId != "" {
			mylangchaingo.SetParentId(o.options.langsmithgoParentId)
		}

		mylangchaingo.SetRunId(uuid.New().String()) //every call to LLM will have a new runId

	}

	return nil
}
