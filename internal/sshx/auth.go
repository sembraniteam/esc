package sshx

import (
	"fmt"
	"net"
	"os"

	"github.com/sembraniteam/esc/internal/config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func Auth(a config.Auth) (ssh.AuthMethod, error) {
	switch a.Type {
	case config.AuthAgent:
		return agentAuth()
	case config.AuthKey:
		//if a.KeyPath == nil {
		//	return nil, fmt.Errorf("key auth requires key_path")
		//}

		return keyAuth(a.KeyPath, a.PassphraseEnv)
	case config.AuthPassword:
		//if a.PasswordEnv == nil {
		//	return nil, fmt.Errorf("password auth requires password_env")
		//}

		return passwordAuth(a.PasswordEnv)
	default:
		return nil, fmt.Errorf("unknown auth type: %s", a.Type)
	}
}

func agentAuth() (ssh.AuthMethod, error) {
	sock := os.Getenv("SSH_AUTH_SOCK")
	if sock == "" {
		return nil, fmt.Errorf("SSH_AUTH_SOCK is not set")
	}

	conn, err := net.Dial("unix", sock)
	if err != nil {
		return nil, err
	}

	ag := agent.NewClient(conn)
	return ssh.PublicKeysCallback(ag.Signers), nil
}

func keyAuth(path, passphraseEnv string) (ssh.AuthMethod, error) {
	// #nosec G304 -- path is user config location, not arbitrary input.
	key, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var signer ssh.Signer
	if passphraseEnv == "" {
		pass := os.Getenv(passphraseEnv)
		if pass == "" {
			return nil, fmt.Errorf(
				"passphrase env '%s' is not set",
				passphraseEnv,
			)
		}

		signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(pass))
	} else {
		signer, err = ssh.ParsePrivateKey(key)
	}

	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(signer), nil
}

func passwordAuth(env string) (ssh.AuthMethod, error) {
	pass := os.Getenv(env)
	if pass == "" {
		return nil, fmt.Errorf("password env '%s' is not set", env)
	}

	return ssh.Password(pass), nil
}
