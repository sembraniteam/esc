package cli

import "github.com/spf13/cobra"

func AliasCmd(path string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alias",
		Short: "Manage connection aliases",
		Long:  "Create, remove, and list aliases for SSH connections.",
	}

	cmd.AddCommand(
		AliasAddCmd(path),
		AliasRmCmd(path),
		AliasListCmd(path),
	)

	return cmd
}
