package cli

import (
	"os"
	"path/filepath"

	"github.com/sembraniteam/esc/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

const minArgs = 2

func Execute() error {
	path := defaultConfigPath()
	rootCmd = &cobra.Command{
		Use:   "esc",
		Short: "Easy SSH Connection",
		Long:  `ESC (Easy SSH Connection) is a CLI tool to manage and connect to SSH servers using a simple, structured configuration.`,
	}

	rootCmd.AddCommand(
		AddCmd(path),
		RmCmd(path),
		MvCmd(path),
		EditCmd(path),
		ListCmd(path),
		TreeCmd(path),
		NewShowCmd(path),
		ConnectCmd(path),
		ExecCmd(path),
		AliasCmd(path),
	)

	rootCmd.AddCommand(
		CompletionCmd(rootCmd),
	)

	return rootCmd.Execute()
}

func defaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return config.FileName
	}

	return filepath.Join(home, config.FolderName, config.FileName)
}
