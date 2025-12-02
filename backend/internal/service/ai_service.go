package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/qingbingwei/latex_ppt_by_claude/backend/pkg/ai"
)

type AIService struct {
	openaiClient  *ai.OpenAIClient
	claudeClient  *ai.ClaudeClient
	copilotClient *ai.CopilotClient
}

func NewAIService(openaiClient *ai.OpenAIClient, claudeClient *ai.ClaudeClient, copilotClient *ai.CopilotClient) *AIService {
	return &AIService{
		openaiClient:  openaiClient,
		claudeClient:  claudeClient,
		copilotClient: copilotClient,
	}
}

func (s *AIService) GenerateLaTeXPPT(ctx context.Context, prompt string, contextChunks []string, useOpenAI bool) (string, error) {
	// Build enhanced prompt with RAG context
	enhancedPrompt := s.buildPrompt(prompt, contextChunks)

	// 优先使用 Copilot
	if s.copilotClient != nil {
		return s.copilotClient.GenerateLaTeX(ctx, enhancedPrompt)
	}

	if useOpenAI && s.openaiClient != nil {
		return s.openaiClient.GenerateLaTeX(ctx, enhancedPrompt)
	} else if s.claudeClient != nil {
		return s.claudeClient.GenerateLaTeX(ctx, enhancedPrompt)
	}

	return "", fmt.Errorf("no AI client available")
}

func (s *AIService) StreamGenerateLaTeXPPT(ctx context.Context, prompt string, contextChunks []string, streamCh chan<- string) error {
	// Build enhanced prompt with RAG context
	enhancedPrompt := s.buildPrompt(prompt, contextChunks)

	if s.openaiClient != nil {
		return s.openaiClient.StreamGenerateLaTeX(ctx, enhancedPrompt, streamCh)
	}

	return fmt.Errorf("streaming only supported with OpenAI")
}

func (s *AIService) buildPrompt(userPrompt string, contextChunks []string) string {
	var prompt strings.Builder

	prompt.WriteString("You are an expert in creating LaTeX Beamer presentations. ")
	prompt.WriteString("Create a complete, compilable LaTeX Beamer presentation based on the following requirements.\n\n")

	if len(contextChunks) > 0 {
		prompt.WriteString("=== Reference Materials (from knowledge base) ===\n")
		for i, chunk := range contextChunks {
			prompt.WriteString(fmt.Sprintf("\n[Context %d]:\n%s\n", i+1, chunk))
		}
		prompt.WriteString("\n=== End of Reference Materials ===\n\n")
	}

	prompt.WriteString("Requirements:\n")
	prompt.WriteString(userPrompt)
	prompt.WriteString("\n\n")

	prompt.WriteString("Guidelines:\n")
	prompt.WriteString("1. Use \\documentclass[aspectratio=169,11pt]{beamer}\n")
	prompt.WriteString("2. Include Chinese support with \\usepackage[UTF8]{ctex}\n")
	prompt.WriteString("3. Use appropriate beamer theme (e.g., Madrid)\n")
	prompt.WriteString("4. Include title page and table of contents\n")
	prompt.WriteString("5. Organize content into sections and frames\n")
	prompt.WriteString("6. Use itemize/enumerate for lists\n")
	prompt.WriteString("7. Keep each frame concise (3-6 bullet points)\n")
	prompt.WriteString("8. The output must be complete and compilable LaTeX code\n")
	prompt.WriteString("9. Wrap the LaTeX code in ```latex code blocks\n")

	if len(contextChunks) > 0 {
		prompt.WriteString("10. Incorporate relevant information from the reference materials provided\n")
	}

	return prompt.String()
}
