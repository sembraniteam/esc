package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/sembraniteam/esc/internal/config"
	"github.com/sembraniteam/esc/internal/sshx"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var (
	parallel bool
	limit    int
)

const defaultLimit = 5

func ExecCmd(path string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec <name> <command>",
		Short: "Execute a command on SSH server(s)",
		Long:  "Execute a command on one or more SSH servers. If the target matches multiple connections (e.g. prod.*), the command can be executed sequentially or in parallel.",
		Example: `# Execute on a single server
esc exec prod.api.main "uptime"
				
# Execute on multiple servers sequentially
esc exec prod.* "df -h"
				
# Execute on multiple servers in parallel (limit 5)
esc exec prod.* "uptime" --parallel --limit 5`,
		Args: cobra.MinimumNArgs(minArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]
			command := args[1]

			cfg, err := config.Load(path)
			if err != nil {
				return err
			}

			names := config.ResolveTargets(cfg, target)
			if len(names) == 0 {
				return fmt.Errorf("no connections matched '%s'", target)
			}

			if !parallel || len(names) == 1 {
				for _, name := range names {
					if err := execOne(
						cfg,
						name,
						command,
						os.Stdout,
						os.Stderr,
					); err != nil {
						return err
					}
				}
				return nil
			}

			if limit <= 0 {
				limit = 5
			}

			return execParallel(cfg, names, command, limit)
		},
	}

	cmd.Flags().
		BoolVar(&parallel, "parallel", false, "Execute command in parallel")
	cmd.Flags().
		IntVar(&limit, "limit", defaultLimit, "Maximum parallel executions")

	return cmd
}

func execOne(
	cfg *config.Schema,
	name, command string,
	out, errOut *os.File,
) error {
	conn, err := config.FindConnection(cfg, name)
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

	return sshx.Exec(addr, clientCfg, command, out, errOut)
}

func execParallel(
	cfg *config.Schema,
	names []string,
	command string,
	limit int,
) error {
	type result struct {
		name   string
		output string
		err    error
	}

	if limit < 1 {
		limit = 1
	}

	sem := make(chan struct{}, limit)
	results := make(chan result, len(names))

	var wg sync.WaitGroup

	for _, name := range names {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			stdoutR, stdoutW, err := os.Pipe()
			if err != nil {
				results <- result{name: name, err: err}
				return
			}
			stderrR, stderrW, err := os.Pipe()
			if err != nil {
				_ = stdoutR.Close()
				_ = stdoutW.Close()
				results <- result{name: name, err: err}
				return
			}

			var stdoutBuf bytes.Buffer
			var stderrBuf bytes.Buffer

			stdoutDone := make(chan struct{})
			stderrDone := make(chan struct{})

			go func() {
				_, _ = io.Copy(&stdoutBuf, stdoutR)
				close(stdoutDone)
			}()
			go func() {
				_, _ = io.Copy(&stderrBuf, stderrR)
				close(stderrDone)
			}()

			err = execOne(cfg, name, command, stdoutW, stderrW)
			_ = stdoutW.Close()
			_ = stderrW.Close()
			<-stdoutDone
			<-stderrDone
			_ = stdoutR.Close()
			_ = stderrR.Close()

			var buf bytes.Buffer
			_, _ = buf.Write(stdoutBuf.Bytes())
			_, _ = buf.Write(stderrBuf.Bytes())
			results <- result{
				name:   name,
				output: buf.String(),
				err:    err,
			}
		}(name)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var hadError bool
	for r := range results {
		fmt.Printf("[%s]\n", r.name)
		if r.output != "" {
			fmt.Print(r.output)
		}
		if r.err != nil {
			hadError = true
			_, err := fmt.Fprintf(os.Stderr, "error: %v\n", r.err)
			if err != nil {
				return err
			}
		}

		fmt.Println()
	}

	if hadError {
		return fmt.Errorf("one or more executions failed")
	}

	return nil
}
