package main

import (
	"context"
	"os"
	"strings"

	"gittale/internal/config"
	gitservice "gittale/internal/services/git"
	llmservice "gittale/internal/services/llm"
	"gittale/pkg/sl"
)

func main() {
	cfg := config.MustLoad()
	log := sl.SetupLogger(string(cfg.Env))

	gitService := gitservice.New()
	llmService, err := llmservice.NewFromConfig(cfg)
	if err != nil {
		log.Error("failed to initialize llm service", sl.Error(err))
		return
	}

	if len(os.Args) < 2 {
		log.Error("please provide a git command to run, e.g., 'branch' or 'status'")
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "commit" {
		ctx := context.Background()

		diff, err := gitService.StagedDiff(ctx)
		if err != nil {
			log.Error("failed to get staged diff", sl.Error(err))
			return
		}

		if strings.TrimSpace(diff) == "" {
			log.Error("no staged changes found. stage files before committing")
			return
		}

		branchName, err := gitService.CurrentBranch(ctx)
		if err != nil {
			log.Error("failed to get current branch", sl.Error(err))
			return
		}

		msg, err := llmService.GenerateCommitMessage(ctx, diff, branchName)
		if err != nil {
			log.Error("failed to generate commit message", sl.Error(err))
			return
		}

		log.Info("generated commit message", "message", msg)

		if err := gitService.Run("commit", "-m", msg); err != nil {
			log.Error("git commit failed", sl.Error(err))
		}
		return
	}

	log.Info("running git command", "args", strings.Join(os.Args[1:], " "))
	if err := gitService.Run(os.Args[1:]...); err != nil {
		log.Error("git command failed", sl.Error(err))
	}
}
