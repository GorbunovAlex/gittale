package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/go-resty/resty/v2"
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

	apiEndpoint := os.Getenv("OLLAMA_ADDRESS")
	if apiEndpoint == "" {
		apiEndpoint = "http://localhost:11434/api/chat"
	}
	modelName := os.Getenv("OLLAMA_MODEL")
	if modelName == "" {
		modelName = "gemma3:4b-it-qat"
	}
	client := resty.New()

	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model": modelName,
			"messages": []interface{}{
				map[string]interface{}{"role": "system", "content": "You are a helpful assistant that writes concise git commit messages."},
				map[string]interface{}{"role": "user", "content": fmt.Sprintf("Generate a git commit message for the following diff:\n%s", string(diffOutput))},
			},
			"stream": false,
		}).
		Post(apiEndpoint)
	if err != nil {
		return "", fmt.Errorf("error while sending the request: %w", err)
	}

	body := response.Body()
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", fmt.Errorf("error while decoding JSON response: %w", err)
	}

	// For Ollama, the response content is typically in data["message"]["content"]
	message, ok := data["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected response format from Ollama")
	}
	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("no content in Ollama response")
	}

	return strings.TrimSpace(content), nil
}
