package ai

import (
	"context"
	"fmt"
	"log"

	"github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	client  *openai.Client
	model   string
	baseURL string
}

func NewOpenAIClient(apiKey, baseURL string) *OpenAIClient {
	config := openai.DefaultConfig(apiKey)
	// 使用 gpt-4o 模型（copilot-api 支持）
	// model := "claude-sonnet-4.5"
	model := "gpt-4o"
	if baseURL != "" {
		config.BaseURL = baseURL
	}

	return &OpenAIClient{
		client:  openai.NewClientWithConfig(config),
		model:   model,
		baseURL: baseURL,
	}
}

func (c *OpenAIClient) GenerateLaTeX(ctx context.Context, prompt string) (string, error) {
	log.Printf("Generating LaTeX with model: %s, baseURL: %s", c.model, c.baseURL)

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: c.model,
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
		log.Printf("OpenAI API error: %v", err)
		return "", fmt.Errorf("AI API error: %v", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return resp.Choices[0].Message.Content, nil
}

func (c *OpenAIClient) StreamGenerateLaTeX(ctx context.Context, prompt string, streamCh chan<- string) error {
	stream, err := c.client.CreateChatCompletionStream(
		ctx,
		openai.ChatCompletionRequest{
			Model: c.model,
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
