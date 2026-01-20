package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sembraniteam/esc/internal/config"
	"github.com/spf13/cobra"
)

func EditCmd(path string) *cobra.Command {
	return &cobra.Command{
		Use:   "edit <name>",
		Short: "Edit an existing SSH connection",
		Long:  "Interactively edit fields of an existing SSH connection.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
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

				block := config.FindConnectionBlock(file, name)
				if block == nil {
					return fmt.Errorf("connection '%s' not found", name)
				}

				in := bufio.NewReader(os.Stdin)
				body := block.Body()

				fmt.Print("Host: ")
				host, _ := in.ReadString('\n')
				body.SetAttributeValue("host", config.StringValTrim(host))

				fmt.Print("User: ")
				user, _ := in.ReadString('\n')
				body.SetAttributeValue("user", config.StringValTrim(user))

				fmt.Print("Port: ")
				port, _ := in.ReadString('\n')
				body.SetAttributeValue("port", config.StringValTrim(port))

				config.SortBlocks(file)
				return config.WriteAtomic(path, file.Bytes())
			})
		},
	}
}
