package message

import "github.com/devalexandre/mylangchaingo/agents/assistant"

type MessageOption func(*Message)

// WithMessageID configura o ID da mensagem.
func WithMessageID(id string) MessageOption {
	return func(m *Message) {
		m.ID = id
	}
}

// WithObject configura o tipo de objeto da mensagem.
func WithObject(object string) MessageOption {
	return func(m *Message) {
		m.Object = object
	}
}

// WithCreatedAt configura a data de criação da mensagem.
func WithCreatedAt(createdAt int) MessageOption {
	return func(m *Message) {
		m.CreatedAt = createdAt
	}
}

// WithThreadID configura o ID da thread da mensagem.
func WithThreadID(threadID string) MessageOption {
	return func(m *Message) {
		m.ThreadId = threadID
	}
}

// WithRole configura o papel (role) da mensagem.
func WithRole(role string) MessageOption {
	return func(m *Message) {
		m.Role = role
	}
}

// WithContent configura o conteúdo da mensagem.
func WithContent(content string) MessageOption {
	return func(m *Message) {
		m.Content = content
	}
}

// WithAssistantID configura o ID do assistente associado à mensagem.
func WithAssistantID(assistantID string) MessageOption {
	return func(m *Message) {
		m.AssistantId = assistantID
	}
}

// WithRunID configura o ID da execução da mensagem.
func WithRunID(runID string) MessageOption {
	return func(m *Message) {
		m.RunId = runID
	}
}

// WithAttachments configura os anexos da mensagem.
func WithAttachments(attachments []assistant.Attachments) MessageOption {
	return func(m *Message) {
		m.Attachments = attachments
	}
}

// WithMetadata configura os metadados da mensagem.
func WithMetadata(metadata map[string]string) MessageOption {
	return func(m *Message) {
		m.Metadata = metadata
	}
}
