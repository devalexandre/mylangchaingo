package assistant

import (
	"github.com/tmc/langchaingo/llms"
)

const BaseURL = "https://api.openai.com/v1"

type ToolType string

const (
	ToolTypeFunction    ToolType = "function"
	ToolTypeFileSearch  ToolType = "file_search"
	CodeInterpreterType ToolType = "code_interpreter"
)

type CodeInterpreter struct {
	Type string `json:"type"` //The type of tool being defined: code_interpreter
}

type RankingOption struct {
	Ranker         string  `json:"ranker"`
	ScoreThreshold float64 `json:"score_threshold"`
}

type FileSearch struct {
	MaxNumResults int           `json:"max_num_results"`
	RankingOption RankingOption `json:"ranking_option"`
}

type StaticChunkingStrategy struct {
	MaxChunkSizeTokens int `json:"max_chunk_size_tokens"`
	ChunkOverlapTokens int `json:"chunk_overlap_tokens"`
}
type ChunkingStrategy struct {
	Type   string                 `json:"type"`
	Static StaticChunkingStrategy `json:"static"`
}

type VectorStore struct {
	FileIDs          []string         `json:"file_ids,omitempty"`
	ChunkingStrategy ChunkingStrategy `json:"chunking_strategy,omitempty"`
}
type CodeInterpreterToolResource struct {
	FileIDs []string `json:"file_ids"`
}

type FileSearchToolResource struct {
	VectorStoreIDs []string      `json:"vector_store_ids,omitempty"`
	VectorStores   []VectorStore `json:"vector_stores,omitempty"`
}

type ToolResource struct {
	CodeInterpreter *CodeInterpreterToolResource `json:"code_interpreter,omitempty"`
	FileSearch      *FileSearchToolResource      `json:"file_search,omitempty"`
}

type ToolFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ToolCall struct {
	ID       string       `json:"id,omitempty"`
	Type     ToolType     `json:"type"`
	Function ToolFunction `json:"function,omitempty"`
}

type Tool struct {
	ID       string    `json:"id,omitempty"`
	Type     ToolType  `json:"type"`
	Function llms.Tool `json:"function,omitempty"`
}
type ContentTextValue struct {
	Value       string       `json:"value"`
	Annotations []Annotation `json:"annotations"`
}
type ContentText struct {
	Type string           `json:"type"`
	Text ContentTextValue `json:"text"`
}

type URL struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}
type File struct {
	FileID string `json:"file_id"`
	Detail string `json:"detail,omitempty"`
}
type ImageURL struct {
	Type     string `json:"type"`
	ImageURL URL    `json:"image_url"`
}

type ImageFile struct {
	Type      string `json:"type"`
	ImageFile File   `json:"image_file"`
}

type Annotation struct {
	FileCitation string `json:"file_citation"`
	FilePath     string `json:"file_path"`
}
type Content struct {
	Text      ContentText `json:"text,omitempty"`
	ImageURL  ImageURL    `json:"image_url,omitempty"`
	ImageFile ImageFile   `json:"image_file,omitempty"`
}

type Attachments struct {
	FileID string `json:"file_id"`
	Tools  []Tool `json:"tools"`
}

type Assistant struct {
	ID           string            `json:"id,omitempty"`
	Model        string            `json:"model"`
	Name         string            `json:"name,omitempty"`
	Description  string            `json:"description,omitempty"`
	Instructions string            `json:"instructions,omitempty"`
	Tools        *[]llms.Tool      `json:"tools,omitempty"`
	ToolResource *ToolResource     `json:"tool_resources,omitempty"`
	Temperature  *float64          `json:"temperature,omitempty"`
	TopP         *float64          `json:"top_p,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

type AssistantResponse struct {
	ID      string      `json:"id"`
	Object  string      `json:"object"`
	Data    []Assistant `json:"data,omitempty"`
	Deleted bool        `json:"deleted,omitempty"`
}
