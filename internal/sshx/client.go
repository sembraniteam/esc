package sshx

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func ClientConfig(
	user string,
	timeout time.Duration,
	auth ssh.AuthMethod,
	hostKeyCallback ssh.HostKeyCallback,
) (*ssh.ClientConfig, error) {
	return &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: hostKeyCallback,
		Timeout:         timeout,
	}, nil
}

func ConnectInteractive(
	addr string,
	cfg *ssh.ClientConfig,
	in, out, errOut *os.File,
) error {
	client, err := ssh.Dial("tcp", addr, cfg)
	if err != nil {
		return err
	}

	defer func(client *ssh.Client) {
		if err := client.Close(); err != nil {
			panic(err)
		}
	}(client)

	sess, err := client.NewSession()
	if err != nil {
		return err
	}

	defer func(sess *ssh.Session) {
		if err := sess.Close(); err != nil {
			panic(err)
		}
	}(sess)

	fd := int(in.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return err
	}

	defer func(fd int, oldState *term.State) {
		if err := term.Restore(fd, oldState); err != nil {
			panic(err)
		}
	}(fd, oldState)

	width, height, err := term.GetSize(fd)
	if err != nil {
		width, height = 80, 24
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400, //nolint:mnd
		ssh.TTY_OP_OSPEED: 14400, //nolint:mnd
	}

	if err := sess.RequestPty(
		"xterm-256color",
		height,
		width,
		modes,
	); err != nil {
		return err
	}

	sess.Stdin = in
	sess.Stdout = out
	sess.Stderr = errOut

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGWINCH)
	go func() {
		for range sig {
			if w, h, err := term.GetSize(fd); err == nil {
				_ = sess.WindowChange(h, w)
			}
		}
	}()

	return sess.Shell()
}
