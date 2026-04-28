// Package config provides configuration loading and management for the application.
// All values are read from a YAML config file at build time by the Makefile and
// stamped into the binary via -ldflags. No file is read at runtime.
package config

import "log"

type AppEnv string

const (
	EnvLocal AppEnv = "local"
	EnvDev   AppEnv = "dev"
	EnvProd  AppEnv = "prod"
)

type ModelProvider string

const (
	LLMProviderOllama ModelProvider = "ollama"
	GeminiProvider    ModelProvider = "gemini"
	ClaudeProvider    ModelProvider = "claude"
)

type Config struct {
	Env           AppEnv
	ModelProvider ModelProvider
	GeminiAPIKey  string
	GeminiModel   string
	OllamaModel   string
	OllamaURL     string
	ClaudeAPIKey  string
	ClaudeModel   string
}

// Stamped at build time by the Makefile via -ldflags -X.
var (
	buildEnv           = "local"
	buildModelProvider string
	buildGeminiAPIKey  string
	buildGeminiModel   = "gemini-2.0-flash"
	buildOllamaModel   string
	buildOllamaURL     = "http://localhost:11434"
	buildClaudeAPIKey  string
	buildClaudeModel   = "claude-sonnet-4-6"
)

func MustLoad() *Config {
	if buildModelProvider == "" {
		log.Fatal("no provider compiled in: fill in config/config.yaml and run 'make build' or 'make install'")
	}
	return &Config{
		Env:           AppEnv(buildEnv),
		ModelProvider: ModelProvider(buildModelProvider),
		GeminiAPIKey:  buildGeminiAPIKey,
		GeminiModel:   buildGeminiModel,
		OllamaModel:   buildOllamaModel,
		OllamaURL:     buildOllamaURL,
		ClaudeAPIKey:  buildClaudeAPIKey,
		ClaudeModel:   buildClaudeModel,
	}
}
