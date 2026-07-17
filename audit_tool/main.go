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

	cmd := &cobra.Command{
		Use:   "audit <entry|entry.md>",
		Short: "Run full audit pipeline on an entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAudit(args[0], outputPath, verbose)
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output filename under audit/ (default: <slug>-audit.md)")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print progress to stderr")

	return cmd
}
