package myssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
)

type client struct {
	*ssh.Client
	host     string
	user     string
	password string
	port     int
}

func NewClient(host, user, password string, port int) (Client, error) {
	c := &client{
		host:     host,
		user:     user,
		password: password,
		port:     port,
	}

	if err := c.Reset(); err != nil {
		return nil, err
	}
	return c, nil
}

func (self *client) Reset() error {
	if self.Client != nil {
		_ = self.Client.Close()
		self.Client = nil
	}

	config := &ssh.ClientConfig{
		User: self.user,
		Auth: []ssh.AuthMethod{
			ssh.Password(self.password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// 连接SSH服务器
	var err error
	self.Client, err = ssh.Dial("tcp", self.host+":"+fmt.Sprintf("%d", self.port), config)
	if err != nil {
		return err
	}

	return nil
}

func (self *client) Run(cmd string) (output []byte, error error) {
	if self.Client == nil {
		if err := self.Reset(); err != nil {
			return nil, err
		}
	}

	session, err := self.Client.NewSession()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = session.Close()
	}()

	session.Stdout = nil
	session.Stderr = nil
	return session.CombinedOutput(cmd)
}

func (self *client) Runs(fn OutFn, cmd ...string) {
	for _, v := range cmd {
		ret, err := self.Run(v)
		if fn != nil {
			if fn(v, ret, err) {
				continue
			} else {
				return
			}
		}
	}
}

func (self *client) GetClient() *ssh.Client {
	return self.Client
}

func (self *client) Close() error {
	var err error
	if self.Client != nil {
		err = self.Client.Close()
		self.Client = nil
	}

	return err
}
