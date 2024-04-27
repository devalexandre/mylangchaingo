package jina

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) {
	os.Setenv("JINA_API_KEY", "jina_eda6a00a90ac48daabac72b6d3ba5e3d7Dl_rdQSZwyZ04aRdkcVYIzOjtd7")
	if jinakey := os.Getenv("JINA_API_KEY"); jinakey == "" {
		t.Skip("JINA_API_KEY not set")

	}
}
func TestJinaEmbeddings(t *testing.T) {
	t.Parallel()

	setup(t)

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

	setup(t)

	j, err := NewJina(WithModel(SmallModel))
	_, err = j.EmbedQuery(context.Background(), "Hello world!")
	require.NoError(t, err)

	embeddings, err := j.EmbedDocuments(context.Background(), []string{"Hello world", "The world is ending", "good bye"})
	require.NoError(t, err)
	assert.Len(t, embeddings, 3)
}

func TestJinaEmbeddingsWithBaseModelModel(t *testing.T) {
	t.Parallel()

	setup(t)

	j, err := NewJina(WithModel(BaseModel))
	_, err = j.EmbedQuery(context.Background(), "Hello world!")
	require.NoError(t, err)

	embeddings, err := j.EmbedDocuments(context.Background(), []string{"Hello world", "The world is ending", "good bye"})
	require.NoError(t, err)
	assert.Len(t, embeddings, 3)
}

func TestJinaEmbeddingsWithLargeModelModel(t *testing.T) {
	t.Parallel()

	setup(t)

	j, err := NewJina(WithModel(LargeModel))
	_, err = j.EmbedQuery(context.Background(), "Hello world!")
	require.NoError(t, err)

	embeddings, err := j.EmbedDocuments(context.Background(), []string{"Hello world", "The world is ending", "good bye"})
	require.NoError(t, err)
	assert.Len(t, embeddings, 3)
}
