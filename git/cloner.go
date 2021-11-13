package git

import (
	"os"

	"github.com/theapemachine/wrkspc/auth"

	"github.com/theapemachine/errnie/v2"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/viper"
)

type Cloner struct {
	provider string
	key      auth.PrivKey
	repo     *git.Repository
}

func NewCloner(provider string, key auth.PrivKey) Cloner {
	return Cloner{
		provider: provider,
		key:      key,
	}
}

func (cloner Cloner) Get(slug string) {
	var err error

	cloner.repo, err = git.PlainClone(viper.GetString("homepath")+"/"+slug, false, &git.CloneOptions{
		URL:      "https://" + cloner.provider + "/" + slug,
		Auth:     cloner.key.Unlock(),
		Progress: os.Stdout,
	})

	_ = cloner.repo
	errnie.Handles(err).With(errnie.KILL)
}
