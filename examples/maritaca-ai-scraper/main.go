package main

import (
	"context"
	"github.com/devalexandre/mylangchaingo/llms/maritaca"
	"github.com/devalexandre/mylangchaingo/tools/scraper/goquery"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/tools"
	"log"
	"os"
)

func main() {

	SetUp() // Set up the environment variables
	ctx := context.Background()
	token := os.Getenv("MARITACA_KEY")

	opts := append([]maritaca.Option{
		maritaca.WithToken(token),
		maritaca.WithModel("sabia-2-medium"),
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

	executor, err := agents.Initialize(
		llm,
		agentTools,
		agents.ZeroShotReactDescription,
		agents.WithMaxIterations(5),
	)
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
