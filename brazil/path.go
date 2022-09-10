package brazil

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/theapemachine/wrkspc/errnie"
)

/*
HomePath does its best to give the caller back the actual home path of the
current user, no matter which OS or environment they are on.
*/
func HomePath() string {
	errnie.Traces()
	home, err := os.UserHomeDir()

	if e := errnie.Handles(err).With(errnie.NOOP); e.Type != errnie.NIL {
		return ""
	}

	return BuildPath(home)
}

/*
CleanPaths removed all paths that were created during `wrkspc` usage.
*/
func CleanPaths() {
	errnie.Traces()
	errnie.Handles(
		os.RemoveAll(BuildPath(HomePath(), ".wrkspc")),
	).With(errnie.NOOP)
}

/*
Workdir returns the current path.
*/
func Workdir() string {
	errnie.Traces()
	wd, err := os.Getwd()
	errnie.Handles(err).With(errnie.NOOP)
	return wd
}

/*
BuildPath joins single strings together into a slash delimited path.
*/
func BuildPath(frags ...string) string {
	errnie.Traces()
	return filepath.FromSlash(strings.Join(frags, "/"))
}

/*
GetFileFromPrefix extracts the filename from a path.
*/
func GetFileFromPrefix(prefix string) string {
	errnie.Traces()
	frags := strings.Split(prefix, "/")
	return frags[len(frags)-1]
}

/*
ReadPath returns everything in path.
*/
func ReadPath(path string) []fs.FileInfo {
	errnie.Traces()
	files, err := ioutil.ReadDir(path)
	errnie.Handles(err).With(errnie.KILL)
	return files
}

/*
MakePath creates a new (nested) path.
*/
func MakePath(path string) {
	errnie.Traces()
	_, err := os.Stat(path)
	if err == nil {
		return
	}

	errnie.Handles(os.MkdirAll(path, os.ModePerm)).With(errnie.NOOP)
}
