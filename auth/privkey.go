package auth

import (
	"github.com/theapemachine/errnie/v2"

	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/spf13/viper"
)

/*
PrivKey is a wrapper around the user's private ssh key.
*/
type PrivKey struct {
}

/*
NewPrivKey constructs an object holding the user's private ssh key
used to clone private repositories.
*/
func NewPrivKey() PrivKey {
	return PrivKey{}
}

/*
Unlock() returns the users public ssh keys.
*/
func (key PrivKey) Unlock() *ssh.PublicKeys {
	pubKey, err := ssh.NewPublicKeysFromFile(
		viper.GetString("username"),
		viper.GetString("homepath")+"/.ssh/id_rsa", "",
	)

	errnie.Handles(err).With(errnie.KILL)
	return pubKey
}
