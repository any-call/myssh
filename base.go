package myssh

import "golang.org/x/crypto/ssh"

type (
	OutFn func(cmd string, output []byte, err error) bool

	Client interface {
		Reset() error
		Run(cmd string) (output []byte, error error)
		Runs(fn OutFn, cmd ...string)
		GetClient() *ssh.Client
		Close() error
	}
)
