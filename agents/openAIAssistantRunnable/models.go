package openAIAssistantRunnable

import "github.com/tmc/langchaingo/llms"

// ToolConfig is the configuration for a tool that can be used by the assistant.
type ToolConfig struct {
	Type     string                   `json:"type"`
	Function *llms.FunctionDefinition `json:"function,omitempty"`
}

type CreateAssistantRequest struct {
	Instructions string      `json:"instructions"`
	Name         string      `json:"name"`
	Tools        []llms.Tool `json:"tools"`
	Model        string      `json:"model"`
}

type CreateAssistantResponse struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Model        string      `json:"model"`
	Instructions string      `json:"instructions"`
	Tools        []llms.Tool `json:"tools"`
}

type CreateThreadResponse struct {
	ID string `json:"id"`
}

type AddMessageRequest struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AddMessageResponse struct {
	ID string `json:"id"`
}

type CreateRunRequest struct {
	AssistantID  string `json:"assistant_id"`
	Instructions string `json:"instructions"`
}

type CreateRunResponse struct {
	ID             string `json:"id"`
	Status         string `json:"status"`
	RequiredAction struct {
		SubmitToolOutputs struct {
			ToolCalls []ToolCall `json:"tool_calls"`
		} `json:"submit_tool_outputs"`
	} `json:"required_action"`
}

type GetThreadMessagesResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Id          string  `json:"id"`
		Object      string  `json:"object"`
		CreatedAt   int     `json:"created_at"`
		AssistantId *string `json:"assistant_id"`
		ThreadId    string  `json:"thread_id"`
		RunId       *string `json:"run_id"`
		Role        string  `json:"role"`
		Content     []struct {
			Type string `json:"type"`
			Text struct {
				Value       string        `json:"value"`
				Annotations []interface{} `json:"annotations"`
			} `json:"text"`
		} `json:"content"`
		Attachments []interface{} `json:"attachments"`
		Metadata    struct {
		} `json:"metadata"`
	} `json:"data"`
	FirstId string `json:"first_id"`
	LastId  string `json:"last_id"`
	HasMore bool   `json:"has_more"`
}
type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}
