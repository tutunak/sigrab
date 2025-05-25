package cli

import (
	"fmt"
	jira "github.com/andygrunwald/go-jira"
	"github.com/spf13/cobra"
	"github.com/tutunak/sigrab/internal/config"
)

var (
	url  string
	from string
	to   string
)

var rootCmd = &cobra.Command{
	Use:   "sigrab",
	Short: "Simple issue Grabber for Jira Cloud",
	Long:  "A CLI tool that fetches Jira issues from a Jira Cloud project and writes them to a local JSON file.",
	RunE:  runSigrab,
}

func runSigrab(cmd *cobra.Command, args []string) error {
	_, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	tp := jira.BasicAuthTransport{
		Username: "<username>",
		Password: "<api-token>",
	}

	client := jira.NewClient(url, cfg.APIToken)
	fetcher := jira.NewFetcher(client)
	return nil
}
