package brazil

import (
	"os"
	"strings"

	"github.com/theapemachine/wrkspc/errnie"
)

/*
Path wraps the structure of a filesystem to provide wrkspc with
some specific functionality, and improving developer ergonomics.
*/
type Path struct {
	Location string
	err      *errnie.Error
}

/*
NewPath takes a variatic input of path segments, which it will
join with a slash (/) to make it behave inline as with most
file system path.
It will also create the path if it does not already exist.
*/
func NewPath(segments ...string) *Path {
	path := &Path{}
	path.Location = path.toPrefix(segments...)
	return path
}

/*
toPrefix takes the segments that were passed in and joins
them into a traditional file system path shape.
It is also cabable of handling file system location wildcards.
*/
func (path *Path) toPrefix(segments ...string) string {
	var out string
	var err error

	switch segments[0] {
	case "~":
		out, err = os.UserHomeDir()
	case ".":
		out, err = os.Getwd()
	default:
		out = strings.Join(segments, "/")
		path.makePath(out)
	}

	if errnie.Handles(err) != nil {
		path.err = errnie.NewError(err)
		return ""
	}

	return out
}

/*
makePath writes a new path on the file system if it was not
already found.
*/
func (path *Path) makePath(prefix string) {
	_, err := os.Stat(prefix)

	if errnie.Handles(err) != nil {
		path.err = errnie.NewError(err)
		errnie.Handles(os.MkdirAll(prefix, os.ModePerm))
		errnie.Informs("created new path at", prefix)
	}

}
