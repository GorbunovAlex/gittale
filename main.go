package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
)

// main is the entry point of the application.
// It loads environment variables, parses command-line arguments,
// and executes git commands, with special handling for the 'commit' command
// to generate a commit message using an AI model.
func main() {
	err := LoadEnv()
	if err != nil {
		fmt.Println("Error loading environment variables:", err)
	}

	fmt.Println("Args passed:")
	for i, arg := range os.Args {
		fmt.Printf("Arg %d: %s\n", i, arg)
	}

	if len(os.Args) < 2 {
		fmt.Println("Please provide a git command to run, e.g., 'branch' or 'status'")
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "commit" {
		msg, err := generateCommitMessageFromGitChanges()
		if err != nil {
			fmt.Println("Error generating commit message:", err)
			return
		}
		fmt.Println("Generated commit message:", msg)
		runGitCommand("commit", "-m", msg)
		return
	}

	fmt.Printf("\nRunning 'git %s':\n", strings.Join(os.Args[1:], " "))
	runGitCommand(os.Args[1:]...)
}

//go:embed .env
var envFile embed.FS

func LoadEnv() error {
	data, err := envFile.ReadFile(".env")
	if err != nil {
		return fmt.Errorf("failed to read .env file: %w", err)
	}

	envMap, err := godotenv.Unmarshal(string(data))
	if err != nil {
		return fmt.Errorf("failed to unmarshal .env file: %w", err)
	}

	for k, v := range envMap {
		_ = os.Setenv(k, v)
	}

	return nil
}

func runGitCommand(args ...string) {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "git command failed: %v\n", err)
	}
}

func generateCommitMessageFromGitChanges() (string, error) {
	// Get current git diff
	cmd := exec.Command("git", "diff", "--cached")
	diffOutput, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git diff: %w", err)
	}

	// Get current branch name
	branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchOutput, err := branchCmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get branch name: %w", err)
	}
	branchName := strings.TrimSpace(string(branchOutput))

	// Extract prefix from branch name up to "--" if present, else use full branch name
	branchPrefix := branchName
	if idx := strings.Index(branchName, "--"); idx != -1 {
		branchPrefix = branchName[:idx]
	}

	// Prepare prompt for ollama
	prompt := fmt.Sprintf("```\n%s\n```\nGenerate a git commit message for the provided diff, be precise, without overthinking on the purpose. The first line is always a short description without any special symbols, then a description separated by a blank line.", string(diffOutput))

	// Call local ollama API
	modelName := os.Getenv("OLLAMA_MODEL")
	if modelName == "" {
		modelName = "llama3"
	}
	reqBody := map[string]interface{}{
		"model":  modelName,
		"prompt": prompt,
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewReader(reqBytes))
	if err != nil {
		return "", fmt.Errorf("failed to call ollama: %w", err)
	}
	defer resp.Body.Close()

	var respData struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", fmt.Errorf("failed to decode ollama response: %w", err)
	}

	commitMsg := strings.TrimSpace(respData.Response)
	if branchPrefix != "" {
		commitMsg = fmt.Sprintf("%s %s", branchPrefix, commitMsg)
	}

	return commitMsg, nil
}
