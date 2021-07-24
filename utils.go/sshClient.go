package utils

import (
	"fmt"
	"time"

	gossh "golang.org/x/crypto/ssh"
)

func GetSSHClient(host string, port int, user string, password string, key string) (*gossh.Client, error) {
	// 先尝试用password验证
	config := &gossh.ClientConfig{
		Timeout:         5 * time.Second,
		User:            user,
		HostKeyCallback: gossh.InsecureIgnoreHostKey(),
		Auth:            []gossh.AuthMethod{gossh.Password(password)},
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	client, err := gossh.Dial("tcp", addr, config)
	// 如果验证不成功，再用key验证
	if err != nil {
		// Create the Signer for this private key.
		signer, err := gossh.ParsePrivateKey([]byte(key))
		if err != nil {
			return nil, fmt.Errorf("unable to parse private key: %v", err)
		}
		keyConfig := &gossh.ClientConfig{
			Timeout:         5 * time.Second,
			User:            user,
			HostKeyCallback: gossh.InsecureIgnoreHostKey(),
			Auth:            []gossh.AuthMethod{gossh.PublicKeys(signer)},
		}
		client, err = gossh.Dial("tcp", addr, keyConfig)
		if err != nil {
			return nil, fmt.Errorf("unable to get client: %v", err)
		}
	}
	return client, nil
}
