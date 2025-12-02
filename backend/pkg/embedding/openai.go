package embedding

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type OpenAIEmbedding struct {
	client *openai.Client
}

func NewOpenAIEmbedding(apiKey, baseURL string) *OpenAIEmbedding {
	config := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		config.BaseURL = baseURL
	}
	return &OpenAIEmbedding{
		client: openai.NewClientWithConfig(config),
	}
}

func (e *OpenAIEmbedding) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	resp, err := e.client.CreateEmbeddings(
		ctx,
		openai.EmbeddingRequest{
			Input: []string{text},
			Model: openai.AdaEmbeddingV2,
		},
	)

	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, nil
	}

	return resp.Data[0].Embedding, nil
}

func (e *OpenAIEmbedding) GenerateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	resp, err := e.client.CreateEmbeddings(
		ctx,
		openai.EmbeddingRequest{
			Input: texts,
			Model: openai.AdaEmbeddingV2,
		},
	)

	if err != nil {
		return nil, err
	}

	embeddings := make([][]float32, len(resp.Data))
	for i, data := range resp.Data {
		embeddings[i] = data.Embedding
	}

	return embeddings, nil
}
