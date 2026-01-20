package cli

import (
	"os"

	"github.com/spf13/cobra"
)

func CompletionCmd(root *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh]",
		Short: "Generate shell completion scripts",
		Long:  "Generate shell completion scripts for ESC.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return root.GenBashCompletion(os.Stdout)
			case "zsh":
				return root.GenZshCompletion(os.Stdout)
			default:
				return cmd.Help()
			}
		},
	}

	return cmd
}
