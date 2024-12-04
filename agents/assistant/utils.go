package assistant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"io/ioutil"
	"net/http"
	"os"
)

func Do(req *http.Request) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_API_KEY")))
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func SubmitToolOutput(threadID, runID, toolCallID, output string) error {
	url := fmt.Sprintf("%s/threads/%s/runs/%s/submit_tool_outputs", BaseURL, threadID, runID)

	requestBody := map[string]interface{}{
		"tool_outputs": []map[string]interface{}{
			{
				"tool_call_id": toolCallID,
				"output":       output,
			},
		},
	}

	bodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return err
	}

	respBody, err := Do(req)
	if err != nil {
		return err
	}

	var result struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return err
	}

	if result.Status != "queued" {
		return fmt.Errorf("unexpected status from submit tool output: %v", result.Status)
	}

	return nil
}

// toolFromTool converts an llms.Tool to a Tool.
func ToolFromTool(t Tool) (Tool, error) {
	tool := Tool{
		Type: ToolType(t.Type),
	}
	switch t.Type {
	case ToolTypeFunction:
		tool.Function = FunctionDefinition{
			Name:        t.Function.Name,
			Description: t.Function.Description,
			Parameters:  t.Function.Parameters,
			Strict:      t.Function.Strict,
		}
	default:
		return Tool{}, fmt.Errorf("tool type %v not supported", t.Type)
	}
	return tool, nil
}

// toolCallsFromToolCalls converts a slice of llms.ToolCall to a slice of ToolCall.
func ToolCallsFromToolCalls(tcs []llms.ToolCall) []ToolCall {
	toolCalls := make([]ToolCall, len(tcs))
	for i, tc := range tcs {
		toolCalls[i] = toolCallFromToolCall(tc)
	}
	return toolCalls
}

// toolCallFromToolCall converts an llms.ToolCall to a ToolCall.
func toolCallFromToolCall(tc llms.ToolCall) ToolCall {
	return ToolCall{
		ID:   tc.ID,
		Type: ToolType(tc.Type),
		Function: ToolFunction{
			Name:      tc.FunctionCall.Name,
			Arguments: tc.FunctionCall.Arguments,
		},
	}
}

func ExtractArg1(jsonStr string) (string, error) {
	// Cria um mapa para armazenar os dados JSON.
	var data map[string]string

	// Desserializa a string JSON no mapa.
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return "", err
	}

	// Retorna o valor de __arg1.
	val, ok := data["__arg1"]
	if !ok {
		return "", fmt.Errorf("__arg1 key not found")
	}

	return val, nil
}
