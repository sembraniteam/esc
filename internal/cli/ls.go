package cli

import (
	"fmt"
	"sort"

	"github.com/sembraniteam/esc/internal/config"
	"github.com/spf13/cobra"
)

func ListCmd(path string) *cobra.Command {
	return &cobra.Command{
		Use:   "ls [pattern]",
		Short: "List SSH connections",
		Long:  "List all SSH connections defined in the ESC configuration. This command shows connection names only. Aliases are not included.",
		Example: `# List all connections
esc ls
# List connections by group
esc ls prod.*`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			cfg, err := config.Load(path)
			if err != nil {
				return err
			}

			pattern := "*"
			if len(args) == 1 {
				pattern = args[0]
			}

			var names []string
			for _, c := range cfg.Connections {
				if config.Match(pattern, c.Name) {
					names = append(names, c.Name)
				}
			}

			sort.Strings(names)

			for _, name := range names {
				fmt.Println(name)
			}

			return nil
		},
	}
}

func AliasListCmd(path string) *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List aliases",
		Long:  "List all defined aliases and their target connections.",
		RunE: func(_ *cobra.Command, args []string) error {
			cfg, err := config.Load(path)
			if err != nil {
				return err
			}

			if len(cfg.Aliases) == 0 {
				fmt.Println("No aliases defined.")
				return nil
			}

			for _, a := range cfg.Aliases {
				fmt.Printf("%s -> %s\n", a.Name, a.Target)
			}

			return nil
		},
	}
}
