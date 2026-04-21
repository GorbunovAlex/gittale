package gemini

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

type GeminiClient struct {
	client *genai.Client
	model  string
}

func NewGeminiClient(apiKey string, model string) (*GeminiClient, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, fmt.Errorf("gemini_api_key must be set when model_provider=gemini")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

	return &GeminiClient{client: client, model: model}, nil
}

func (c *GeminiClient) Generate(ctx context.Context, prompt string) (string, error) {
	result, err := c.client.Models.GenerateContent(ctx, c.model, genai.Text(prompt), nil)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result.Text()), nil
}
