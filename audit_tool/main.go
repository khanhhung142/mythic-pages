package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "vietmyth-auditor",
		Short: "Audit vietmyth.vn entries for factual accuracy",
	}

	root.AddCommand(auditCmd())

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func auditCmd() *cobra.Command {
	var outputPath string
	var verbose bool
	var llmProvider string
	var llmModel string
	var searchProvider string

	cmd := &cobra.Command{
		Use:   "audit <entry|entry.md>",
		Short: "Run full audit pipeline on an entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rt, err := NewRuntime(AIConfig{
				LLMProvider:    llmProvider,
				LLMModel:       llmModel,
				SearchProvider: searchProvider,
			})
			if err != nil {
				return err
			}
			return runAudit(args[0], outputPath, verbose, rt)
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output filename under audit/ (default: <slug>-audit.md)")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print progress to stderr")
	cmd.Flags().StringVar(&llmProvider, "llm", envOr("AUDIT_LLM", "claude"), "LLM provider: claude, openai, deepseek, gemini")
	cmd.Flags().StringVar(&llmModel, "llm-model", os.Getenv("AUDIT_LLM_MODEL"), "Override default model for the selected LLM provider")
	cmd.Flags().StringVar(&searchProvider, "search", envOr("AUDIT_SEARCH", "perplexity"), "Search provider: perplexity")

	return cmd
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
