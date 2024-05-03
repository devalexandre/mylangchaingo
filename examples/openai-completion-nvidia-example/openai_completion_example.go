package main

import (
	"context"
	"fmt"
	"github.com/devalexandre/mylangchaingo/llms/openai"
	"github.com/tmc/langchaingo/llms"
	"log"
)

func main() {
	SetUp()

	opts := []openai.Option{
		openai.WithModel("meta/llama3-70b-instruct"),
		openai.WithAPIType(openai.APITypeNvidia),
	}

	llm, err := openai.New(opts...)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	completion, err := llms.GenerateFromSinglePrompt(ctx,
		llm,
		"The first man to walk on the moon",
		llms.WithTemperature(0.8),
		llms.WithStopWords([]string{"Armstrong"}),
	)
	if err != nil {
		log.Fatal(err)
	}

	embs, err := llm.CreateEmbedding(ctx, []string{completion})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(embs)
}
