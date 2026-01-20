package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/sembraniteam/esc/internal/config"
	"github.com/sembraniteam/esc/internal/sshx"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var insecure bool

func ConnectCmd(path string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "connect <name>",
		Short:   "Connect to an SSH server",
		Long:    "Connect to an SSH server using a saved ESC connection. The name can be a connection name or an alias. ESC resolves authentication and opens an interactive SSH session.",
		Example: "esc connect prod.api.main",
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]

			cfg, err := config.Load(path)
			if err != nil {
				return err
			}
			target, err := config.ResolveName(cfg, name)
			if err != nil {
				return err
			}

			conn, err := config.FindConnection(cfg, target)
			if err != nil {
				return err
			}

			auth, err := sshx.Auth(conn.Auth)
			if err != nil {
				return err
			}

			var hk ssh.HostKeyCallback
			if insecure {
				// #nosec G106 -- explicit --insecure flag opt-in.
				hk = ssh.InsecureIgnoreHostKey()
			} else {
				hk, err = sshx.HostKeyCallbackFromKnownHosts()
				if err != nil {
					return err
				}
			}

			clientCfg, err := sshx.ClientConfig(
				conn.User,
				cfg.Settings.Timeout*time.Second,
				auth,
				hk,
			)
			if err != nil {
				return err
			}

			addr := fmt.Sprintf("%s:%d", conn.Host, conn.Port)
			return sshx.ConnectInteractive(
				addr,
				clientCfg,
				os.Stdin,
				os.Stdout,
				os.Stderr,
			)
		},
	}

	cmd.Flags().
		BoolVar(&insecure, "insecure", false, "Disable SSH host key verification (not recommended)")

	return cmd
}
