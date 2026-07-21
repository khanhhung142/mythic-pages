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
	var extractLLM string
	var extractModel string
	var judgeLLM string
	var judgeModel string
	var llmProvider string // legacy: overrides both
	var llmModel string
	var searchProvider string

	cmd := &cobra.Command{
		Use:   "audit <entry|entry.md>",
		Short: "Run full audit pipeline on an entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rt, err := NewRuntime(AIConfig{
				ExtractProvider: extractLLM,
				ExtractModel:    extractModel,
				JudgeProvider:   judgeLLM,
				JudgeModel:      judgeModel,
				LLMProvider:     llmProvider,
				LLMModel:        llmModel,
				SearchProvider:  searchProvider,
			})
			if err != nil {
				return err
			}
			return runAudit(args[0], outputPath, verbose, rt)
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output filename under audit/ (default: <slug>-audit.md)")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print progress to stderr")
	cmd.Flags().StringVar(&extractLLM, "extract-llm", envOr("AUDIT_EXTRACT_LLM", "deepseek"), "Extract LLM: deepseek, gemini, claude, openai")
	cmd.Flags().StringVar(&extractModel, "extract-model", os.Getenv("AUDIT_EXTRACT_MODEL"), "Override extract model")
	cmd.Flags().StringVar(&judgeLLM, "judge-llm", envOr("AUDIT_JUDGE_LLM", "claude"), "Judge LLM (recommend claude): claude, openai, deepseek, gemini")
	cmd.Flags().StringVar(&judgeModel, "judge-model", os.Getenv("AUDIT_JUDGE_MODEL"), "Override judge model")
	cmd.Flags().StringVar(&llmProvider, "llm", os.Getenv("AUDIT_LLM"), "Legacy: set both extract and judge to one provider (batch experiment)")
	cmd.Flags().StringVar(&llmModel, "llm-model", os.Getenv("AUDIT_LLM_MODEL"), "Legacy: model when using --llm")
	cmd.Flags().StringVar(&searchProvider, "search", envOr("AUDIT_SEARCH", "perplexity"), "Search provider: perplexity")

	return cmd
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
