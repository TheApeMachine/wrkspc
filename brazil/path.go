package brazil

import (
	"io/fs"
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

	if e := errnie.Handles(err); e.Type != errnie.NIL {
		return ""
	}

	return BuildPath(home)
}

/*
CleanPaths removed all paths that were created during `wrkspc` usage.
*/
func CleanPaths() {
	errnie.Handles(
		os.RemoveAll(BuildPath(HomePath(), ".wrkspc")),
	)
}

/*
Workdir returns the current path.
*/
func Workdir() string {
	wd, err := os.Getwd()
	errnie.Handles(err)
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

func GeneratePath(prefix string) chan string {
	out := make(chan string)

	go func() {
		defer close(out)

		filepath.Walk(
			prefix, func(p string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					out <- p
				}

				return err
			},
		)
	}()

	return out
}

/*
ReadPath returns everything in path.
*/
func ReadPath(path string) []fs.DirEntry {
	files, err := os.ReadDir(path)
	errnie.Handles(err)
	return files
}

/*
MakePath creates a new (nested) path.
*/
func MakePath(path string) {
	_, err := os.Stat(path)
	if err == nil {
		return
	}

	errnie.Handles(os.MkdirAll(path, os.ModePerm))
}
