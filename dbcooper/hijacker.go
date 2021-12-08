package dbcooper

import (
	"os"
	"strings"

	"github.com/theapemachine/wrkspc/brazil"
)

/*
Hijacker is responsible for (re)configuring the user's environment and leaving it clean afterwards.
*/
type Hijacker struct {
	shell    string
	realPath string
}

/*
NewHijacker constructs a Hijacker and returns a pointer reference to the new instance.
*/
func NewHijacker() *Hijacker {
	return &Hijacker{
		realPath: os.Getenv("PATH"),
	}
}

/*
TakeOver the environment.
*/
func (hijacker *Hijacker) TakeOver() *Hijacker {
	// Let's play: guess that shell.
	hijacker.shell = hijacker.guessShell()

	// And use it to backup the current `.rc` file
	brazil.Copy(
		brazil.BuildPath(brazil.HomePath(), "."+hijacker.shell+"rc"),
		brazil.BuildPath(brazil.HomePath(), "."+hijacker.shell+"rc.bak"),
	)

	// Make sure we only look in the `wrkspc` path for any executables.
	os.Setenv("PATH", brazil.HomePath()+"/.wrkspc")
	return hijacker
}

/*
Release the contained environment and go back to the original state.
*/
func (hijacker *Hijacker) Release() *Hijacker {
	// Delete our temporary `.rc` file and restore the previous one.
	brazil.DeleteFile(brazil.BuildPath(brazil.HomePath(), "."+hijacker.shell+"rc"))

	brazil.Copy(
		brazil.BuildPath(brazil.HomePath(), "."+hijacker.shell+"rc.bak"),
		brazil.BuildPath(brazil.HomePath(), "."+hijacker.shell+"rc"),
	)

	os.Setenv("PATH", hijacker.realPath)
}

/*
guessShell uses a pretty dumb way to figure out which shell the user is using underneath.
*/
func (hijacker *Hijacker) guessShell() string {
	shell := strings.Split(os.Getenv("SHELL"), "/")
	return shell[len(shell)-1]
}
