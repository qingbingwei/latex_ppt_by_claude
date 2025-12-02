package ai

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	client *openai.Client
}

func NewOpenAIClient(apiKey, baseURL string) *OpenAIClient {
	config := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		config.BaseURL = baseURL
	}
	return &OpenAIClient{
		client: openai.NewClientWithConfig(config),
	}
}

func (c *OpenAIClient) GenerateLaTeX(ctx context.Context, prompt string) (string, error) {
	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are an expert in creating LaTeX Beamer presentations. Generate complete, compilable LaTeX code for presentations.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7,
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (c *OpenAIClient) StreamGenerateLaTeX(ctx context.Context, prompt string, streamCh chan<- string) error {
	stream, err := c.client.CreateChatCompletionStream(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are an expert in creating LaTeX Beamer presentations. Generate complete, compilable LaTeX code for presentations.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7,
			Stream:      true,
		},
	)

	if err != nil {
		return err
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if err != nil {
			close(streamCh)
			return err
		}

		if len(response.Choices) > 0 {
			streamCh <- response.Choices[0].Delta.Content
		}
	}
}
