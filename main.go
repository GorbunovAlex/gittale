package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
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

	gemini_api_key := os.Getenv("API_KEY")
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  gemini_api_key,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create GenAI client: %w", err)
	}

	result, err := client.Models.GenerateContent(ctx, "gemini-2.0-flash", genai.Text(fmt.Sprintf("Generate a git commit message for the following diff, be precise, without overthinking on the purpose:\n%s", string(diffOutput))), nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	commitMsg := strings.TrimSpace(result.Text())
	if branchPrefix != "" {
		commitMsg = fmt.Sprintf("%s %s", branchPrefix, commitMsg)
	}

	return commitMsg, nil
}
