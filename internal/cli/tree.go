package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sembraniteam/esc/internal/config"
	"github.com/spf13/cobra"
)

func TreeCmd(path string) *cobra.Command {
	return &cobra.Command{
		Use:   "tree",
		Short: "Show connections as a tree",
		Long:  "Display SSH connections grouped by dot-separated segments.",
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg, err := config.Load(path)
			if err != nil {
				return err
			}

			tree := map[string]map[string][]string{}

			for _, c := range cfg.Connections {
				parts := strings.Split(c.Name, ".")
				switch len(parts) {
				case 1:
					tree[parts[0]] = nil
				case 2: //nolint:mnd
					tree[parts[0]] = map[string][]string{parts[1]: {}}
				case 3: //nolint:mnd
					if tree[parts[0]] == nil {
						tree[parts[0]] = map[string][]string{}
					}
					tree[parts[0]][parts[1]] = append(
						tree[parts[0]][parts[1]],
						parts[2],
					)
				}
			}

			roots := make([]string, 0, len(tree))
			for r := range tree {
				roots = append(roots, r)
			}

			sort.Strings(roots)

			for _, r := range roots {
				fmt.Println(r)
				for mid, leaves := range tree[r] {
					fmt.Printf(" └─ %s\n", mid)
					sort.Strings(leaves)
					for _, l := range leaves {
						fmt.Printf("    └─ %s\n", l)
					}
				}
			}

			return nil
		},
	}
}
