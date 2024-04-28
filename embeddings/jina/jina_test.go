package jina

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Setenv("LANGCHAIN_TRACING", "true")
	os.Setenv("LANGSMITH_API_KEY", "")
	os.Setenv("LANGCHAIN_PROJECT_NAME", "jina")
	os.Setenv("OPENAI_API_KEY", "")
	os.Setenv("JINA_API_KEY", "")
	if jinakey := os.Getenv("JINA_API_KEY"); jinakey == "" {
		os.Exit(0)

	}
	os.Exit(m.Run())
}

func TestJinaEmbeddings(t *testing.T) {
	t.Parallel()

	j, err := NewJina()
	require.NoError(t, err)

	_, err = j.EmbedQuery(context.Background(), "Hello world!")
	require.NoError(t, err)

	embeddings, err := j.EmbedDocuments(context.Background(), []string{"Hello world", "The world is ending", "good bye"})
	require.NoError(t, err)
	assert.Len(t, embeddings, 3)
}

// with model option
func TestJinaEmbeddingsWithSamllModel(t *testing.T) {
	t.Parallel()

	j, err := NewJina(WithModel(SmallModel))
	_, err = j.EmbedQuery(context.Background(), "Hello world!")
	require.NoError(t, err)

	embeddings, err := j.EmbedDocuments(context.Background(), []string{"Hello world", "The world is ending", "good bye"})
	require.NoError(t, err)
	assert.Len(t, embeddings, 3)
}

func TestJinaEmbeddingsWithBaseModelModel(t *testing.T) {
	t.Parallel()

	j, err := NewJina(WithModel(BaseModel))
	_, err = j.EmbedQuery(context.Background(), "Hello world!")
	require.NoError(t, err)

	embeddings, err := j.EmbedDocuments(context.Background(), []string{"Hello world", "The world is ending", "good bye"})
	require.NoError(t, err)
	assert.Len(t, embeddings, 3)
}

func TestJinaEmbeddingsWithLargeModelModel(t *testing.T) {
	t.Parallel()

	j, err := NewJina(WithModel(LargeModel))
	_, err = j.EmbedQuery(context.Background(), "Hello world!")
	require.NoError(t, err)

	embeddings, err := j.EmbedDocuments(context.Background(), []string{"Hello world", "The world is ending", "good bye"})
	require.NoError(t, err)
	assert.Len(t, embeddings, 3)
}
