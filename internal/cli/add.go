package cli

import (
	"fmt"

	"github.com/sembraniteam/esc/internal/config"
	"github.com/spf13/cobra"
)

func AddCmd(path string) *cobra.Command {
	var (
		name     string
		host     string
		user     string
		authType string
		port     int
	)

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add new SSH connection",
		Long:  "Add a new SSH connection to ESC configuration. This command creates a new connection entry in the HCL config file. The connection name may optionally use dot-separated grouping (e.g. db or prod.api.main). All connection names must be unique.",
		Example: `# Add a production API server
esc add --name prod.api.main --host 10.10.10.1 --user ubuntu

# Add a simple ungrouped connection
esc add --name db --host 127.0.0.1 --user postgres`,
		RunE: func(_ *cobra.Command, _ []string) error {
			lockPath := path + ".lock"

			if err := config.ValidateConnectionName(name); err != nil {
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

				if config.HasConnection(file, name) {
					return fmt.Errorf("connection '%s' already exists", name)
				}

				authT, err := config.ToAuthType(authType)
				if err != nil {
					return err
				}

				if err := config.AppendConnection(file, config.Connection{
					Name: name,
					Host: host,
					User: user,
					Port: int64(port),
					Auth: config.Auth{Type: authT},
				}); err != nil {
					return err
				}

				config.SortBlocks(file)
				return config.WriteAtomic(path, file.Bytes())
			})
		},
	}

	cmd.Flags().
		StringVar(&name, "name", "", "connection name (e.g. db or prod.api.main)")
	cmd.Flags().StringVar(&host, "host", "", "host address")
	cmd.Flags().StringVar(&user, "user", "", "ssh user")
	cmd.Flags().IntVar(&port, "port", config.DefaultPort, "ssh port")
	cmd.Flags().
		StringVar(&authType, "auth-type", config.AuthAgent.String(), "ssh auth type. Available ssh auth type 'agent', 'key', 'password'")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("host")
	_ = cmd.MarkFlagRequired("user")

	return cmd
}

func AliasAddCmd(path string) *cobra.Command {
	return &cobra.Command{
		Use:   "add <name> <target>",
		Short: "Add a new alias",
		Long:  `Create a short alias for an existing SSH connection. Aliases allow you to connect using a shorter name.`,
		Example: `# Create alias 'p' for prod.api.main
esc alias add p prod.api.main`,
		Args: cobra.ExactArgs(minArgs),
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]
			target := args[1]

			if err := config.ValidateConnectionName(name); err != nil {
				return err
			}

			lockPath := path + ".lock"
			return config.WithLock(lockPath, func() error {
				if err := config.EnsureConfigExists(path); err != nil {
					return err
				}

				cfg, err := config.Load(path)
				if err != nil {
					return err
				}

				if err := config.EnsureTargetExists(cfg, target); err != nil {
					return err
				}

				file, err := config.Writable(path)
				if err != nil {
					return err
				}

				if config.HasAlias(file, name) {
					return fmt.Errorf("alias '%s' already exists", name)
				}

				config.SortBlocks(file)
				config.AppendAlias(file, name, target)
				return config.WriteAtomic(path, file.Bytes())
			})
		},
	}
}
