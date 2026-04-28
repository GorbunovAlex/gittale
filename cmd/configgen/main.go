// configgen reads a gittale YAML config file and prints the -ldflags string
// needed to bake all config values into the main binary at compile time.
// Usage: go run ./cmd/configgen <config-file>
package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const pkg = "gittale/internal/config"

type rawConfig struct {
	Env           string `yaml:"env"`
	ModelProvider string `yaml:"model_provider"`
	GeminiAPIKey  string `yaml:"gemini_api_key"`
	GeminiModel   string `yaml:"gemini_model"`
	OllamaModel   string `yaml:"ollama_model"`
	OllamaURL     string `yaml:"ollama_url"`
	ClaudeAPIKey  string `yaml:"claude_api_key"`
	ClaudeModel   string `yaml:"claude_model"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: configgen <config.yaml>")
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("cannot read config file: %v", err)
	}

	var cfg rawConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("cannot parse config file: %v", err)
	}

	if cfg.Env == "" {
		cfg.Env = "local"
	}
	if cfg.GeminiModel == "" {
		cfg.GeminiModel = "gemini-2.0-flash"
	}
	if cfg.OllamaURL == "" {
		cfg.OllamaURL = "http://localhost:11434"
	}
	if cfg.ClaudeModel == "" {
		cfg.ClaudeModel = "claude-sonnet-4-6"
	}

	fmt.Printf(
		"-X '%s.buildEnv=%s' -X '%s.buildModelProvider=%s' -X '%s.buildGeminiAPIKey=%s' -X '%s.buildGeminiModel=%s' -X '%s.buildOllamaModel=%s' -X '%s.buildOllamaURL=%s' -X '%s.buildClaudeAPIKey=%s' -X '%s.buildClaudeModel=%s'",
		pkg, cfg.Env,
		pkg, cfg.ModelProvider,
		pkg, cfg.GeminiAPIKey,
		pkg, cfg.GeminiModel,
		pkg, cfg.OllamaModel,
		pkg, cfg.OllamaURL,
		pkg, cfg.ClaudeAPIKey,
		pkg, cfg.ClaudeModel,
	)
}
