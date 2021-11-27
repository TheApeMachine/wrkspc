package brazil

import (
	"os"
)

/*
Environment is a set of ENV variables shielded from the user's underlying evironment.
*/
type Environment struct {
	path string
}

/*
NewEnvironment gives us a handler to manipulating the user's environment.
*/
func NewEnvironment() *Environment {
	// We want to override te executable paths of the user for a while so we contain them to
	// only the embdded tooling, such that we affect the system in the least possible way.
	return &Environment{path: os.Getenv("PATH")}
}

/*
Initialize the shielded environment.
*/
func (env *Environment) Initialize() *Environment {
	os.Setenv("PATH", HomePath()+"/wrkspc")
	return env
}

/*
Restore the environment of the user so we make no impact.
*/
func (env *Environment) Restore() {
	os.Setenv("PATH", env.path)
}
