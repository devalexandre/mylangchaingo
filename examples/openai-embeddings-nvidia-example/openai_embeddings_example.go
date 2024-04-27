package main

import (
	"context"
	"fmt"
	"github.com/devalexandre/mylangchaingo/llms/openai"
	"log"
)

func main() {
	opts := []openai.Option{
		openai.WithModel("meta/llama2-70b"),
		openai.WithAPIType(openai.APITypeNvidia),
		openai.WithEmbeddingModel("NV-Embed-QA"),
	}

	llm, err := openai.New(opts...)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	emb, err := llm.CreateEmbedding(ctx, []string{"The first"})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(emb)
}
