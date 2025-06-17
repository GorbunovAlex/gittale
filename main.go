package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

// Usage Instructions:
//  1. Ensure you have Go installed and git available in your PATH.
//  2. Set your OpenAI API key in the environment variable OPENAI_API_KEY.
//     Example: export OPENAI_API_KEY=your-api-key
//  3. Build the application: go build -o gitgpt
//  4. Run the application with a git command as arguments, e.g.:
//     ./gitgpt status
//     ./gitgpt branch
//     ./gitgpt commit
//     - For 'commit', the app will generate a commit message using GPT based on staged changes.
func main() {
	_ = godotenv.Load()
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

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create GenAI client: %w", err)
	}

	result, err := client.Models.GenerateContent(ctx, "gemini-2.0-flash", genai.Text(fmt.Sprintf("Generate a git commit message for the following diff:\n%s", string(diffOutput))), nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	return strings.TrimSpace(result.Text()), nil
}
