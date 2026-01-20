package config

import (
	"errors"
	"time"
)

const (
	FolderName = ".esc"
	FileName   = "config.hcl"
)

const (
	DefaultPort           = 22
	DefaultTimeoutSeconds = 30
)

const (
	AuthAgent    AuthType = "agent"
	AuthKey      AuthType = "key"
	AuthPassword AuthType = "password"
)

type (
	AuthType string

	Schema struct {
		Settings    Settings     `hcl:"settings,block"`
		Connections []Connection `hcl:"connection,block"`
		Aliases     []Alias      `hcl:"alias,block"`
	}

	Settings struct {
		Timeout time.Duration `hcl:"timeout,attr"`
	}

	Connection struct {
		Name string   `hcl:"name,label"`
		Host string   `hcl:"host,attr"`
		Port int64    `hcl:"port,optional"`
		User string   `hcl:"user,attr"`
		Auth Auth     `hcl:"auth,block"`
		Tags []string `hcl:"tags,optional"`
	}

	Auth struct {
		Type          AuthType `hcl:"type,attr"`
		KeyPath       string   `hcl:"key_path,optional"`
		PassphraseEnv string   `hcl:"passphrase_env,optional"`
		PasswordEnv   string   `hcl:"password_env,optional"`
	}

	Alias struct {
		Name   string `hcl:"name,label"`
		Target string `hcl:"target,attr"`
	}
)

func (a AuthType) String() string {
	return string(a)
}

func ToAuthType(s string) (AuthType, error) {
	switch s {
	case "agent":
		return AuthAgent, nil
	case "key":
		return AuthKey, nil
	case "password":
		return AuthPassword, nil
	default:
		return "", errors.New("invalid auth type")
	}
}
