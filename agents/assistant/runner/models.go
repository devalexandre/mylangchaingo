package runner

import (
	"github.com/devalexandre/mylangchaingo/agents/assistant"
	"github.com/devalexandre/mylangchaingo/agents/assistant/message"
	"github.com/devalexandre/mylangchaingo/agents/assistant/thread"
)

type TruncationStrategy struct {
	Type         string      `json:"type,omitempty"`
	LastMessages interface{} `json:"last_messages,omitempty"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens,omitempty"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
	TotalTokens      int `json:"total_tokens,omitempty"`
}

type Runner struct {
	Id                     *string             `json:"id,omitempty"`
	Object                 *string             `json:"object,omitempty"`
	CreatedAt              *int                `json:"created_at,omitempty"`
	AssistantId            string              `json:"assistant_id"` // Este campo é obrigatório, sem `omitempty`
	ThreadId               *string             `json:"thread_id,omitempty"`
	Status                 *string             `json:"status,omitempty"`
	StartedAt              *int                `json:"started_at,omitempty"`
	ExpiresAt              *interface{}        `json:"expires_at,omitempty"`
	CancelledAt            *interface{}        `json:"cancelled_at,omitempty"`
	FailedAt               *interface{}        `json:"failed_at,omitempty"`
	CompletedAt            *int                `json:"completed_at,omitempty"`
	RequiredAction         interface{}         `json:"required_action,omitempty"`
	LastError              *interface{}        `json:"last_error,omitempty"`
	Model                  *string             `json:"model,omitempty"`
	Instructions           *interface{}        `json:"instructions,omitempty"`
	AdditionalInstructions *string             `json:"additional_instructions,omitempty"`
	AddicionalMessage      *[]message.Message  `json:"additional_messages,omitempty"`
	Tools                  *[]assistant.Tool   `json:"tools,omitempty"`
	Metadata               *map[string]string  `json:"metadata,omitempty"`
	IncompleteDetails      *interface{}        `json:"incomplete_details,omitempty"`
	Usage                  *Usage              `json:"usage,omitempty"`
	Temperature            *float64            `json:"temperature,omitempty"`
	TopP                   *float64            `json:"top_p,omitempty"`
	MaxPromptTokens        *int                `json:"max_prompt_tokens,omitempty"`
	MaxCompletionTokens    *int                `json:"max_completion_tokens,omitempty"`
	TruncationStrategy     *TruncationStrategy `json:"truncation_strategy,omitempty"`
	ResponseFormat         *string             `json:"response_format,omitempty"`
	ToolChoice             *string             `json:"tool_choice,omitempty"`
	ParallelToolCalls      *bool               `json:"parallel_tool_calls,omitempty"`
	Stream                 *bool               `json:"stream,omitempty"`
	Thread                 *thread.Thread      `json:"thread,omitempty"`
}
