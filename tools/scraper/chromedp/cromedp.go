package chromedp

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/devalexandre/langsmithgo"
	"github.com/devalexandre/mylangchaingo"
	"github.com/google/uuid"
	"github.com/tmc/langchaingo/tools"
	"golang.org/x/net/html"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

const (
	DefualtMaxDept   = 1
	DefualtParallels = 2
	DefualtDelay     = 3
	DefualtAsync     = true
)

var ErrScrapingFailed = errors.New("scraper could not read URL, or scraping is not allowed for provided URL")
var _ tools.Tool = Scraper{}

type Scraper struct {
	MaxDepth        int
	Parallels       int
	Delay           int64
	Blacklist       []string
	Async           bool
	Timeout         time.Duration
	Await           time.Duration
	langsmithClient *langsmithgo.Client
}

func New(options ...Options) (*Scraper, error) {
	scraper := &Scraper{
		MaxDepth:  DefualtMaxDept,
		Parallels: DefualtParallels,
		Delay:     int64(DefualtDelay),
		Async:     DefualtAsync,
		Blacklist: []string{
			"login",
			"signup",
			"signin",
			"register",
			"logout",
			"download",
			"redirect",
		},
	}

	for _, opt := range options {
		opt(scraper)
	}

	if os.Getenv("LANGCHAIN_TRACING") != "" && os.Getenv("LANGCHAIN_TRACING") != "false" {
		client, err := langsmithgo.NewClient()
		if err != nil {
			return nil, err
		}
		scraper.langsmithClient = client
		root := uuid.New().String()
		mylangchaingo.SetRunId(root)

	}

	return scraper, nil
}

func (s Scraper) Name() string {
	return "Web Scraper"
}

func (s Scraper) Description() string {
	return "Web Scraper will scan a url and return the content of the web page. Input should be a working url."
}

func (s Scraper) Call(ctx context.Context, input string) (string, error) {
	if s.langsmithClient != nil {
		err := s.langsmithClient.Run(&langsmithgo.RunPayload{
			Name:        fmt.Sprintf("%v-%v-%v", langsmithgo.Tool, s.Name(), "CromeDP"),
			SessionName: os.Getenv("LANGCHAIN_PROJECT_NAME"),
			RunType:     langsmithgo.Tool,
			RunID:       mylangchaingo.GetRunId(),
			ParentID:    mylangchaingo.GetParentId(),
			Inputs: map[string]interface{}{
				"payload": input,
			},
			Extras: map[string]interface{}{
				"Metadata": map[string]interface{}{
					"langsmithgo_version": "v1.0.0",
					"go_version":          runtime.Version(),
					"platform":            runtime.GOOS,
					"arch":                runtime.GOARCH,
				},
			},
		})

		if err != nil {
			return "", err
		}
	}

	url, err := ExtractURL(input)
	if err != nil {
		return "", fmt.Errorf("failed to extract URL: %w", err)

	}
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"),
		chromedp.Flag("headless", true),      // ou false para ver o navegador
		chromedp.Flag("disable-http2", true), // adicione esta linha para desativar HTTP/2
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	// Create a new browser context
	ctx, cancel = chromedp.NewContext(allocCtx)
	defer cancel()

	// Create a timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	var headers, paragraphs, contentMain, contentDiv string
	err = chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(10*time.Second), // Optional: wait for specific elements to ensure the page has loaded.
		chromedp.Text(`h1, h2, h3, h4, h5, h6`, &headers, chromedp.ByQueryAll),
		chromedp.Text(`p`, &paragraphs, chromedp.ByQueryAll),
		chromedp.Text(`div#content`, &contentDiv, chromedp.ByQueryAll),
		chromedp.Text(`main#content`, &contentMain, chromedp.ByQueryAll),
	)
	if err != nil {
		return "", fmt.Errorf("failed to scrape the website: %w", err)
	}
	combinedText := headers + "\n" + paragraphs + "\n" + contentMain + "\n" + contentDiv
	response := RemoveBlankLines(combinedText)

	if s.langsmithClient != nil {
		err := s.langsmithClient.Run(&langsmithgo.RunPayload{
			RunID: mylangchaingo.GetRunId(),
			Outputs: map[string]interface{}{
				"output": response,
			},
		})

		if err != nil {
			return "", fmt.Errorf("error running langsmith: %w", err)
		}
	}

	return response, nil
}

// walkNodes percorre os nós do documento HTML e aplica a função fn a cada nó
func walkNodes(n *html.Node, fn func(*html.Node)) {
	fn(n) // Aplica a função ao nó atual
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walkNodes(c, fn) // Chama a função recursivamente para cada filho
	}
}

func ExtractURL(str string) (string, error) {
	// Este padrão de expressão regular é projetado para corresponder à maioria das URLs.
	// Modifique conforme necessário para atender a casos de uso mais específicos.
	// O padrão atual busca por protocolos comuns como http, https, ftp, etc.
	pattern := `https?://[^\s]+`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err // Retorna um erro se a expressão regular for inválida
	}

	// Encontra a primeira correspondência na string fornecida.
	match := re.FindString(str)
	if match == "" {
		return "", fmt.Errorf("no URL found in string")
	}

	//remove > from the end of the URL
	if strings.HasSuffix(match, ">") {
		match = match[:len(match)-1]
	}
	return match, nil
}

func RemoveBlankLines(input string) string {
	var builder strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(input))

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" {
			builder.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		// Tratamento de erro, se necessário
		return ""
	}

	return strings.TrimRight(builder.String(), "\n")
}
