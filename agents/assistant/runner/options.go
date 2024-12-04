package runner

import (
	"github.com/devalexandre/mylangchaingo/agents/assistant/message"
	"github.com/tmc/langchaingo/tools"
)

// Option define o tipo para opções funcionais para configurar o runner.
type Option func(*Runner)

// Funções para configurar cada campo da struct Runner
func WithObject(object string) Option {
	return func(r *Runner) {
		r.Object = &object
	}
}

func WithThreadID(threadID string) Option {
	return func(r *Runner) {
		r.ThreadId = &threadID
	}
}

func WithStatus(status string) Option {
	return func(r *Runner) {
		r.Status = &status
	}
}

func WithModel(model string) Option {
	return func(r *Runner) {
		r.Model = &model
	}
}

func WithInstructions(instructions interface{}) Option {
	return func(r *Runner) {
		r.Instructions = &instructions
	}
}

func WithAdditionalInstructions(instructions string) Option {
	return func(r *Runner) {
		r.AdditionalInstructions = &instructions
	}
}

func WithAddicionalMessage(messages []message.Message) Option {
	return func(r *Runner) {
		r.AddicionalMessage = &messages
	}
}

func WithTools(tools []tools.Tool) Option {

	return func(r *Runner) {
		//r.Tools = &tools
	}
}

func WithMetadata(metadata map[string]string) Option {
	return func(r *Runner) {
		r.Metadata = &metadata
	}
}

func WithTemperature(temperature float64) Option {
	return func(r *Runner) {
		r.Temperature = &temperature
	}
}

func WithTopP(topP float64) Option {
	return func(r *Runner) {
		r.TopP = &topP
	}
}

func WithMaxPromptTokens(tokens int) Option {
	return func(r *Runner) {
		r.MaxPromptTokens = &tokens
	}
}

func WithMaxCompletionTokens(tokens int) Option {
	return func(r *Runner) {
		r.MaxCompletionTokens = &tokens
	}
}

func WithTruncationStrategy(strategy TruncationStrategy) Option {
	return func(r *Runner) {
		r.TruncationStrategy = &strategy
	}
}

func WithResponseFormat(format string) Option {
	return func(r *Runner) {
		r.ResponseFormat = &format
	}
}

func WithToolChoice(toolChoice string) Option {
	return func(r *Runner) {
		r.ToolChoice = &toolChoice
	}
}

func WithParallelToolCalls(parallel bool) Option {
	return func(r *Runner) {
		r.ParallelToolCalls = &parallel
	}
}

func WithStream(stream bool) Option {
	return func(r *Runner) {
		r.Stream = &stream
	}
}
