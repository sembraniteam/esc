package cli

import (
	"fmt"

	"github.com/sembraniteam/esc/internal/config"
	"github.com/spf13/cobra"
)

func RmCmd(path string) *cobra.Command {
	return &cobra.Command{
		Use:   "rm <name>",
		Short: "Remove an existing SSH connection",
		Long: `Remove an SSH connection from the ESC configuration.
				The connection name must match exactly.
				This operation permanently deletes the connection entry
				from the config file.`,
		Example: `# Remove a specific connection
				  esc rm prod.api.main
				
				  # Remove an ungrouped connection
				  esc rm db`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]
			lockPath := path + ".lock"

			return config.WithLock(lockPath, func() error {
				file, err := config.Writable(path)
				if err != nil {
					return err
				}

				if !config.RemoveConnection(file, name) {
					return fmt.Errorf("connection '%s' not found", name)
				}

				return config.WriteAtomic(path, file.Bytes())
			})
		},
	}
}

func AliasRmCmd(path string) *cobra.Command {
	return &cobra.Command{
		Use:     "rm <name>",
		Short:   "Remove an alias",
		Long:    "Remove an existing alias from the ESC configuration.",
		Example: `esc alias rm p`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]
			lockPath := path + ".lock"

			return config.WithLock(lockPath, func() error {
				if err := config.EnsureConfigExists(path); err != nil {
					return err
				}

				file, err := config.Writable(path)
				if err != nil {
					return err
				}

				if !config.RemoveAlias(file, name) {
					return fmt.Errorf("alias '%s' not found", name)
				}

				config.SortBlocks(file)
				return config.WriteAtomic(path, file.Bytes())
			})
		},
	}
}
