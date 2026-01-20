package cli

import (
	"github.com/sembraniteam/esc/internal/config"
	"github.com/spf13/cobra"
)

func MvCmd(path string) *cobra.Command {
	return &cobra.Command{
		Use:   "mv <old> <new>",
		Short: "Rename an SSH connection",
		Long:  "Rename an existing SSH connection and update all aliases that reference the connection. This operation is atomic and safe.",
		Example: `# Rename a connection
esc mv prod.api.main prod.api.primary
				
# Rename a simple connection
esc mv db database`,
		Args: cobra.ExactArgs(minArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			oldName := args[0]
			newName := args[1]
			lockPath := path + ".lock"

			if err := config.ValidateConnectionName(newName); err != nil {
				return err
			}

			return config.WithLock(lockPath, func() error {
				if err := config.EnsureConfigExists(path); err != nil {
					return err
				}

				file, err := config.Writable(path)
				if err != nil {
					return err
				}

				if err := config.RenameConnection(
					file,
					oldName,
					newName,
				); err != nil {
					return err
				}

				config.SortBlocks(file)
				return config.WriteAtomic(path, file.Bytes())
			})
		},
	}
}
