package cli

import (
	"github.com/sembraniteam/esc/internal/config"
	"github.com/spf13/cobra"
)

func NewShowCmd(configPath string) *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Show details of an SSH connection",
		Long: `Show detailed information about a single SSH connection.

The name may be a connection name or an alias.`,
		Example: `
  # Show a connection
  esc show prod.api.main

  # Show using alias
  esc show p
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			input := args[0]

			cfg, err := config.Load(configPath)
			if err != nil {
				return err
			}

			name, err := config.ResolveName(cfg, input)
			if err != nil {
				return err
			}

			conn, err := config.FindConnection(cfg, name)
			if err != nil {
				return err
			}

			printConnection(conn)
			return nil
		},
	}
}
