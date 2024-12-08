package main

import (
	"context"
	"log"
	"os"

	"github.com/devalexandre/mylangchaingo/llms/maritaca"
	"github.com/devalexandre/mylangchaingo/tools/scraper/goquery"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/tools"
)

func main() {

	SetUp() // Set up the environment variables
	ctx := context.Background()

	token := os.Getenv("MARITACA_KEY")

	opts := append([]maritaca.Option{
		maritaca.WithToken(token),
		maritaca.WithModel("sabia-3"),
	})

	llm, err := maritaca.New(opts...)
	if err != nil {
		panic(err)
	}

	scrapr, err := goquery.New()
	if err != nil {
		log.Fatal(err)
	}

	agentTools := []tools.Tool{
		scrapr,
	}

	agent := agents.NewOpenAIFunctionsAgent(llm, agentTools, agents.WithMaxIterations(5))
	executor := agents.NewExecutor(agent)
	if err != nil {
		log.Fatal(err)
	}

	prompt := "Quais os Capítulos, do livro O Saci em https://pt.wikipedia.org/wiki/O_Saci?"
	callOptions := []chains.ChainCallOption{
		chains.WithTemperature(0.6),
	}

	res, err := chains.Run(ctx, executor, prompt, callOptions...)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res)
}
func SetUp() {

}
