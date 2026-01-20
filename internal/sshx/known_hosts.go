package sshx

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func HostKeyCallbackFromKnownHosts() (ssh.HostKeyCallback, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(home, ".ssh", "known_hosts")

	cb, err := knownhosts.New(path)
	if err != nil {
		return nil, err
	}

	return func(hostname string, remoteAddr net.Addr, key ssh.PublicKey) error {
		if err := cb(hostname, remoteAddr, key); err != nil {
			var ke *knownhosts.KeyError
			if errors.As(err, &ke) {
				if len(ke.Want) == 0 {
					return fmt.Errorf(
						"unknown SSH host key for %s. Add it to known_hosts by connecting once with ssh",
						hostname,
					)
				}

				return fmt.Errorf(
					"SSH host key mismatch for %s. Possible man-in-the-middle attack",
					hostname,
				)
			}

			return err
		}

		return nil
	}, nil
}
