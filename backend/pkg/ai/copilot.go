package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	copilotTokenURL      = "https://api.github.com/copilot_internal/v2/token"
	copilotChatURL       = "https://api.githubcopilot.com/chat/completions"
	copilotEditorVersion = "vscode/1.95.0"
	copilotEditorPlugin  = "copilot-chat/0.22.0"
)

type CopilotClient struct {
	githubToken string
	accessToken string
	tokenExpiry time.Time
	httpClient  *http.Client
}

type copilotTokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

type copilotChatRequest struct {
	Model       string               `json:"model"`
	Messages    []copilotChatMessage `json:"messages"`
	Temperature float64              `json:"temperature,omitempty"`
	MaxTokens   int                  `json:"max_tokens,omitempty"`
	Stream      bool                 `json:"stream,omitempty"`
}

type copilotChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type copilotChatResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func NewCopilotClient(githubToken string) *CopilotClient {
	return &CopilotClient{
		githubToken: githubToken,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// refreshToken 获取或刷新 Copilot 访问令牌
func (c *CopilotClient) refreshToken() error {
	if c.accessToken != "" && time.Now().Before(c.tokenExpiry) {
		return nil // Token still valid
	}

	req, err := http.NewRequest("GET", copilotTokenURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create token request: %v", err)
	}

	req.Header.Set("Authorization", "token "+c.githubToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Editor-Version", copilotEditorVersion)
	req.Header.Set("Editor-Plugin-Version", copilotEditorPlugin)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get copilot token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to get copilot token: status %d, body: %s", resp.StatusCode, string(body))
	}

	var tokenResp copilotTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode token response: %v", err)
	}

	c.accessToken = tokenResp.Token
	c.tokenExpiry = time.Unix(tokenResp.ExpiresAt, 0).Add(-5 * time.Minute) // 提前5分钟过期

	log.Printf("Copilot token refreshed, expires at: %v", c.tokenExpiry)
	return nil
}

func (c *CopilotClient) GenerateLaTeX(ctx context.Context, prompt string) (string, error) {
	// 刷新令牌
	if err := c.refreshToken(); err != nil {
		return "", fmt.Errorf("failed to refresh token: %v", err)
	}

	log.Printf("Generating LaTeX with GitHub Copilot")

	chatReq := copilotChatRequest{
		Model: "gpt-4o", // Copilot 支持的模型
		Messages: []copilotChatMessage{
			{
				Role:    "system",
				Content: "You are an expert in creating LaTeX Beamer presentations. Generate complete, compilable LaTeX code for presentations.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.7,
		MaxTokens:   4096,
		Stream:      false,
	}

	reqBody, err := json.Marshal(chatReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", copilotChatURL, bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Editor-Version", copilotEditorVersion)
	req.Header.Set("Editor-Plugin-Version", copilotEditorPlugin)
	req.Header.Set("Copilot-Integration-Id", "vscode-chat")
	req.Header.Set("Openai-Intent", "conversation-panel")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Copilot API error: status %d, body: %s", resp.StatusCode, string(body))
		return "", fmt.Errorf("copilot API error: status %d", resp.StatusCode)
	}

	var chatResp copilotChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from Copilot")
	}

	log.Printf("Copilot response received, tokens used: %d", chatResp.Usage.TotalTokens)
	return chatResp.Choices[0].Message.Content, nil
}
