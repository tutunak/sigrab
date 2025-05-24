package cli

import (
	"github.com/spf13/cobra"
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
