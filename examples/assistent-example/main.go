package main

import (
	"fmt"

	"github.com/devalexandre/mylangchaingo/agents/assistant"
	"github.com/devalexandre/mylangchaingo/agents/assistant/executor"
	"github.com/tmc/langchaingo/tools"
)

func main() {

	tool := tools.Calculator{}
	assistant, err := assistant.NewAssistant(
		assistant.WithName("Calculator Assistant"),
		assistant.WithDescription("You are a personal math tutor."),
		assistant.WithModel("gpt-3.5-turbo"),
		assistant.WithTools([]tools.Tool{tool}),
	)
	if err != nil {
		fmt.Println("Error creating assistant:", err)
		return
	}

	agentExecutor := executor.NewAgentExecutor(assistant, executor.WithTools([]tools.Tool{tool}))

	input := "What is 10 + 20?"
	response, err := agentExecutor.Run(input)
	if err != nil {
		fmt.Println("Error invoking agent:", err)
		return
	}

	fmt.Println("Response:", response)
}
