package goquery

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/devalexandre/langsmithgo"
	"github.com/devalexandre/mylangchaingo"
	"github.com/google/uuid"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/tools"
	"golang.org/x/net/html"

	"github.com/PuerkitoBio/goquery"
)

const (
	DefaultMaxDepth  = 1
	DefaultParallels = 2
	DefaultDelay     = 3
	DefaultAsync     = true
)

var ErrScrapingFailed = errors.New("scraper could not read URL, or scraping is not allowed for provided URL")

type Scraper struct {
	MaxDepth         int
	Parallels        int
	Delay            int64
	Blacklist        []string
	Async            bool
	langsmithClient  *langsmithgo.Client
	CallbacksHandler callbacks.Handler
}

var _ tools.Tool = Scraper{}

func New(options ...Options) (*Scraper, error) {
	scraper := &Scraper{
		MaxDepth:  DefaultMaxDepth,
		Parallels: DefaultParallels,
		Delay:     int64(DefaultDelay),
		Async:     DefaultAsync,
		Blacklist: []string{
			"login", "signup", "signin", "register", "logout", "download", "redirect",
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
	return "Web Scraper will scan a URL and return the content of the web page. Input should be a working URL."
}

func (s Scraper) Call(ctx context.Context, input string) (string, error) {

	if s.CallbacksHandler != nil {
		s.CallbacksHandler.HandleToolStart(ctx, input)
	}

	urlLink, err := ExtractURL(input)
	if err != nil {
		return "", fmt.Errorf("failed to extract URL: %w", err)

	}

	if s.langsmithClient != nil {
		err := s.langsmithClient.Run(&langsmithgo.RunPayload{
			Name:        fmt.Sprintf("%v-%v-%v", langsmithgo.Tool, s.Name(), "GoQuery"),
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
	u, err := url.ParseRequestURI(urlLink)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrScrapingFailed, err)
	}

	res, err := goquery.NewDocument(urlLink)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrScrapingFailed, err)
	}

	var siteData strings.Builder
	siteData.WriteString("\n\nPage URL: " + input)

	title := res.Find("title").Text()
	if title != "" {
		siteData.WriteString("\nPage Title: " + title)
	}

	description, _ := res.Find("meta[name='description']").Attr("content")
	if description != "" {
		siteData.WriteString("\nPage Description: " + description)
	}

	siteData.WriteString("\nHeaders:")
	res.Find("h1, h2, h3, h4, h5, h6").Each(func(i int, sel *goquery.Selection) {
		siteData.WriteString("\n" + sel.Text())
	})

	siteData.WriteString("\nContent:")
	res.Find("p").Each(func(i int, sel *goquery.Selection) {
		siteData.WriteString("\n" + sel.Text())
	})

	res.Find("div#content").Each(func(i int, sel *goquery.Selection) {
		siteData.WriteString("content:" + sel.Text())
	})

	res.Find("main#content").Each(func(i int, sel *goquery.Selection) {
		siteData.WriteString("content:" + sel.Text())
	})

	links := make(map[string]bool)
	res.Find("a[href]").Each(func(i int, sel *goquery.Selection) {
		link, exists := sel.Attr("href")
		if exists {
			absoluteLink, err := u.Parse(link)
			if err != nil {
				log.Println(err)
				return
			}

			if absoluteLink.Hostname() != u.Hostname() {
				return
			}

			for _, item := range s.Blacklist {
				if strings.Contains(absoluteLink.Path, item) {
					return
				}
			}

			if absoluteLink.Path == "/index.html" || absoluteLink.Path == "" {
				absoluteLink.Path = "/"
			}

			if !links[absoluteLink.String()] {
				links[absoluteLink.String()] = true
				siteData.WriteString("\nLink: " + absoluteLink.String())
			}
		}
	})

	siteData.WriteString("\n\nScraped Links:")
	for link := range links {
		siteData.WriteString("\n" + link)
	}

	textExtracted, err := ExtractTextFromHTML(siteData.String())
	if err != nil {
		return "", err

	}
	response := RemoveBlankLines(textExtracted)

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

	if s.CallbacksHandler != nil {
		s.CallbacksHandler.HandleToolEnd(ctx, response)
	}

	return response, nil
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

func ExtractTextFromHTML(htmlContent string) (string, error) {
	// Parse the HTML content
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", err // Retorna um erro se o HTML não puder ser analisado
	}

	var b bytes.Buffer // Um buffer para armazenar o texto extraído
	walkNodes(doc, func(n *html.Node) {
		if n.Type == html.TextNode {
			b.WriteString(n.Data)
		}
	})

	return b.String(), nil // Retorna o texto extraído
}

// walkNodes percorre os nós do documento HTML e aplica a função fn a cada nó
func walkNodes(n *html.Node, fn func(*html.Node)) {
	fn(n) // Aplica a função ao nó atual
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walkNodes(c, fn) // Chama a função recursivamente para cada filho
	}
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
