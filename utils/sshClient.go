package utils

import (
	"errors"
	"fmt"
	"time"

	gossh "golang.org/x/crypto/ssh"
)

// 创建ssh client，认证方式支持password和key
func GetSSHClient(host string, port int, user string, authType int, password string, key string) (*gossh.Client, error) {
	if authType != 1 && authType != 2 {
		return nil, errors.New("auth type is only 1 or 2")
	}
	var client *gossh.Client
	var errs error
	addr := fmt.Sprintf("%s:%d", host, port)
	if authType == 1 {
		config := &gossh.ClientConfig{
			Timeout:         5 * time.Second,
			User:            user,
			HostKeyCallback: gossh.InsecureIgnoreHostKey(),
			Auth:            []gossh.AuthMethod{gossh.Password(password)},
		}
		client, errs = gossh.Dial("tcp", addr, config)
	} else {
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
		client, errs = gossh.Dial("tcp", addr, keyConfig)
	}
	if errs != nil {
		return nil, errs
	}
	return client, nil
}
