package claude

import (
	"context"
	"fmt"
	"strings"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

const defaultMaxTokens = 2048

type ClaudeClient struct {
	client anthropic.Client
	model  string
}

func NewClaudeClient(apiKey, model string) (*ClaudeClient, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, fmt.Errorf("claude_api_key must be set when model_provider=claude")
	}
	client := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &ClaudeClient{client: client, model: model}, nil
}

func (c *ClaudeClient) Generate(ctx context.Context, prompt string) (string, error) {
	msg, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     c.model,
		MaxTokens: defaultMaxTokens,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		return "", fmt.Errorf("claude api error: %w", err)
	}

	var sb strings.Builder
	for _, block := range msg.Content {
		if block.Type == "text" {
			sb.WriteString(block.Text)
		}
	}
	return strings.TrimSpace(sb.String()), nil
}
