package sshx

import (
	"os"

	"golang.org/x/crypto/ssh"
)

func Exec(
	addr string,
	cfg *ssh.ClientConfig,
	cmd string,
	out, errOut *os.File,
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

	sess.Stdout = out
	sess.Stderr = errOut

	return sess.Run(cmd)
}
