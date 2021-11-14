package matrix

import (
	"os"

	"github.com/theapemachine/wrkspc/auth"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/git"

	"github.com/spf13/viper"
)

type Resolver struct {
	root string
}

func NewResolver(root string) Resolver {
	errnie.Traces()

	return Resolver{
		root: root,
	}
}

/*
Update the required dependencies into the current
context path, so the final tarball package will
include them on creation.
*/
func (dep Resolver) Update() {
	errnie.Traces()

	// We always want to re-clone because of possible changes, but at the development stage
	// it may not be in a repo yet, so ignore the clone when debugging.
	if !viper.GetBool("debug") {
		// Generate a public version of our private key to compare with the repo.
		key := auth.NewPrivKey()

		// Make sure we have all the dependencies present locally.
		cloner := git.NewCloner(viper.GetString("wrkspc.git.host"), key)
		username := viper.GetString("wrkspc.git.username")
		deps := viper.GetStringSlice("wrkspc.git.dependencies")

		for _, dep := range deps {
			cloner.Get(username + dep)
		}
	}
}

/*
Cleanup removes all the dependencies that were downloaded, since in any kind
of real scenario they will have to be downloaded again on the next run anyway.
*/
func (dep Resolver) Cleanup() {
	errnie.Traces()

	// In debug mode we don't have full access to the updated repo yet, so
	// we can keep things around.
	if !viper.GetBool("debug") {
		errnie.Handles(os.RemoveAll(dep.root)).With(errnie.KILL)
	}
}
