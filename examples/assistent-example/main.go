package main

import (
	"fmt"
	"github.com/devalexandre/mylangchaingo/agents/openAIAssistantRunnable"
	"github.com/tmc/langchaingo/tools"
)

func main() {

	tool := tools.Calculator{}
	assistant, err := openAIAssistantRunnable.NewAssistant(
		"Calculator Assistant",
		"You are a personal math tutor.",
		"gpt-3.5-turbo",
		[]tools.Tool{tool},
	)
	if err != nil {
		fmt.Println("Error creating assistant:", err)
		return
	}

	agentExecutor := openAIAssistantRunnable.NewAgentExecutor(assistant, []tools.Tool{tool})

	input := map[string]string{"content": "What is 10 + 20?"}
	response, err := agentExecutor.Run(input)
	if err != nil {
		fmt.Println("Error invoking agent:", err)
		return
	}

	fmt.Println("Response:", response)
}