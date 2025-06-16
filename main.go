package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
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

	// Prepare request to GPT API
	apiURL := "https://api.openai.com/v1/chat/completions"
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	requestBody := fmt.Sprintf(`{
		"model": "gpt-4.1",
		"messages": [
			{"role": "system", "content": "You are a helpful assistant that writes concise git commit messages."},
			{"role": "user", "content": "Generate a git commit message for the following diff:\n%s"}
		],
		"max_tokens": 60
	}`, string(diffOutput))

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response (simple extraction)
	type Choice struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	type GPTResponse struct {
		Choices []Choice `json:"choices"`
	}
	var gptResp GPTResponse
	if err := json.Unmarshal(body, &gptResp); err != nil {
		return "", fmt.Errorf("failed to parse GPT response: %w", err)
	}
	if len(gptResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in GPT response")
	}
	return strings.TrimSpace(gptResp.Choices[0].Message.Content), nil
}
