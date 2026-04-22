package llm

import (
	"context"
	"fmt"
	"strings"

	"gittale/internal/config"
	"gittale/internal/services/llm/gemini"
	"gittale/internal/services/llm/ollama"
)

const (
	defaultMaxBatchChars = 12000
	defaultGeminiModel   = "gemini-2.0-flash"
	defaultOllamaURL     = "http://localhost:11434"
)

type TextGenerator interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

type Service struct {
	client        TextGenerator
	maxBatchChars int
}

func splitDiffIntoBatches(diff string, maxChars int) []string {
	if maxChars <= 0 {
		maxChars = defaultMaxBatchChars
	}

	lines := strings.Split(diff, "\n")
	batches := make([]string, 0)
	var current strings.Builder

	flush := func() {
		if current.Len() > 0 {
			batches = append(batches, current.String())
			current.Reset()
		}
	}

	for _, line := range lines {
		lineWithNL := line + "\n"
		if len(lineWithNL) > maxChars {
			flush()
			remaining := lineWithNL
			for len(remaining) > maxChars {
				batches = append(batches, remaining[:maxChars])
				remaining = remaining[maxChars:]
			}
			if remaining != "" {
				current.WriteString(remaining)
			}
			continue
		}

		if current.Len()+len(lineWithNL) > maxChars {
			flush()
		}
		current.WriteString(lineWithNL)
	}

	flush()
	return batches
}

func buildBatchSummaryPrompt(batch string, idx int, total int) string {
	return fmt.Sprintf(
		"You are summarizing a git diff chunk (%d/%d).\n"+
			"List the concrete file-level and behavior-level changes in short bullets.\n"+
			"Do not invent details. Do not write a commit message yet.\n"+
			"Do not use markdown code fences or backticks in your response.\n\n"+
			"Diff chunk:\n%s",
		idx,
		total,
		batch,
	)
}

func buildCommitMessagePrompt(summaries []string, branchPrefix string) string {
	titleRule := "short imperative title (max 72 chars)"
	if branchPrefix != "" {
		titleRule = fmt.Sprintf("short imperative title (max 72 chars), must start with \"%s \"", branchPrefix)
	}

	return fmt.Sprintf(
		"Generate a git commit message based on the diff summaries below.\n\n"+
			"Rules:\n"+
			"- Output ONLY the commit message. No explanations, no markdown, no backticks, no code fences.\n"+
			"- First line: %s.\n"+
			"- Second line: blank.\n"+
			"- Third line onwards: optional body with specific details, plain text only.\n\n"+
			"Summaries:\n%s",
		titleRule,
		strings.Join(summaries, "\n\n"),
	)
}

func extractBranchPrefix(branchName string) string {
	name := strings.TrimSpace(branchName)
	if name == "" {
		return ""
	}
	if idx := strings.Index(name, "--"); idx != -1 {
		return strings.TrimSpace(name[:idx])
	}
	return name
}

func (s *Service) GenerateCommitMessage(ctx context.Context, diff string, branchName string) (string, error) {
	if strings.TrimSpace(diff) == "" {
		return "", fmt.Errorf("diff is empty")
	}

	batches := splitDiffIntoBatches(diff, s.maxBatchChars)
	summaries := make([]string, 0, len(batches))

	for i, batch := range batches {
		summary, err := s.client.Generate(ctx, buildBatchSummaryPrompt(batch, i+1, len(batches)))
		if err != nil {
			return "", fmt.Errorf("failed to summarize diff batch %d: %w", i+1, err)
		}
		if strings.TrimSpace(summary) != "" {
			summaries = append(summaries, summary)
		}
	}

	if len(summaries) == 0 {
		return "", fmt.Errorf("unable to produce commit summary from diff")
	}

	branchPrefix := extractBranchPrefix(branchName)
	message, err := s.client.Generate(ctx, buildCommitMessagePrompt(summaries, branchPrefix))
	if err != nil {
		return "", fmt.Errorf("failed to generate commit message: %w", err)
	}

	message = strings.TrimSpace(message)
	if message == "" {
		return "", fmt.Errorf("llm returned empty commit message")
	}

	return message, nil
}

func NewFromConfig(cfg *config.Config) (*Service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}

	var client TextGenerator
	var err error

	switch cfg.ModelProvider {
	case config.GeminiProvider:
		model := strings.TrimSpace(cfg.GeminiModel)
		if model == "" {
			model = defaultGeminiModel
		}
		client, err = gemini.NewGeminiClient(strings.TrimSpace(cfg.GeminiAPIKey), model)
		if err != nil {
			return nil, err
		}
	case config.LLMProviderOllama:
		model := strings.TrimSpace(cfg.OllamaModel)
		if model == "" {
			return nil, fmt.Errorf("ollama_model must be set when model_provider=ollama")
		}
		url := strings.TrimSpace(cfg.OllamaURL)
		if url == "" {
			url = defaultOllamaURL
		}
		client = ollama.NewOllamaClient(url, model)
	default:
		return nil, fmt.Errorf("unsupported llm provider: %s", cfg.ModelProvider)
	}

	return &Service{
		client:        client,
		maxBatchChars: defaultMaxBatchChars,
	}, nil
}
