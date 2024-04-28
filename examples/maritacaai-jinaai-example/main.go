package main

import (
	"context"
	"fmt"
	"github.com/devalexandre/mylangchaingo/embeddings/jina"
	"github.com/devalexandre/mylangchaingo/llms/maritaca"
	"github.com/tmc/langchaingo/llms"
	"os"
)

func main() {

	SetUp() // Set up the environment variables

	token := os.Getenv("MARITACA_KEY")

	opts := append([]maritaca.Option{
		maritaca.WithToken(token),
		maritaca.WithModel("sabia-2-medium"),
	})

	llm, err := maritaca.New(opts...)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "Você é um Historiador, e é apaixonado pela cuntura brasileira. Você está escrevendo um livro sobre a história do Brasil."),
		llms.TextParts(llms.ChatMessageTypeHuman, "Qual a lenda mais conhecida do Brasil?"),
	}

	res, err := llm.GenerateContent(ctx, content)

	if err != nil {
		panic(err)
	}

	//create embbeding with Jinna
	j, err := jina.NewJina(jina.WithModel(jina.BaseModel))
	if err != nil {
		panic(err)
	}

	embs, err := j.EmbedQuery(ctx, res.Choices[0].Content)

	if err != nil {
		panic(err)
	}

	fmt.Println(embs)
}
