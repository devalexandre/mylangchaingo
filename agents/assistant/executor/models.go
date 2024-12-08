package executor

import "github.com/tmc/langchaingo/tools"

type ExecutorOption func(*AgentExecutor)

func WithTools(tools []tools.Tool) ExecutorOption {
	return func(a *AgentExecutor) {
		a.Tools = tools
	}

}
