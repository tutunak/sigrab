package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tutunak/sigrab/internal/config"
	"github.com/tutunak/sigrab/internal/jira"
	"github.com/tutunak/sigrab/internal/utils"
)

var (
	url  string
	from string
	to   string
	path string
)

var rootCmd = &cobra.Command{
	Use:   "sigrab",
	Short: "Simple issue Grabber for Jira Cloud",
	Long:  "A CLI tool that fetches Jira issues from a Jira Cloud project and writes them to a local JSON file.",
	RunE:  runSigrab,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVar(&url, "url", "", "Jira Cloud URL (required)")
	rootCmd.Flags().StringVar(&from, "from", "", "Starting Jira issue key (e.g., DEV-123)")
	rootCmd.Flags().StringVar(&to, "to", "", "Ending Jira issue key (e.g., DEV-140)")
	rootCmd.Flags().StringVar(&path, "path", "", "Path to save the JSON file (default: current directory)")

	rootCmd.MarkFlagRequired("url")

}

func runSigrab(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fullPath := utils.InitDir(path)

	fetcher := jira.NewFetcher(jira.NewClient(cfg.UserEmail, cfg.APIToken, url))

	err = fetcher.FetchBackward(to, fullPath)

	if err != nil {
		return fmt.Errorf("failed to fetch issues: %w", err)
	}
	return nil
}
