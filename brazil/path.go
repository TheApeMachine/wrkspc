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
	home, err := os.UserHomeDir()
	errnie.Handles(err).With(errnie.NOOP)
	return BuildPath(home)
}

/*
CleanPaths removed all paths that were created during `wrkspc` usage.
*/
func CleanPaths() {
	errnie.Handles(
		os.RemoveAll(BuildPath(HomePath(), ".wrkspc")),
	).With(errnie.NOOP)
}

/*
Workdir returns the current path.
*/
func Workdir() string {
	wd, err := os.Getwd()
	errnie.Handles(err).With(errnie.NOOP)
	return wd
}

/*
BuildPath joins single strings together into a slash delimited path.
*/
func BuildPath(frags ...string) string {
	return filepath.FromSlash(strings.Join(frags, "/"))
}

/*
GetFileFromPrefix extracts the filename from a path.
*/
func GetFileFromPrefix(prefix string) string {
	frags := strings.Split(prefix, "/")
	return frags[len(frags)-1]
}

/*
ReadPath returns everything in path.
*/
func ReadPath(path string) []fs.FileInfo {
	files, err := ioutil.ReadDir(path)
	errnie.Handles(err).With(errnie.KILL)
	return files
}
