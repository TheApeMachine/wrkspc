package matrix

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/containerd/console"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/moby/term"
	"github.com/spf13/viper"
	"github.com/theapemachine/wrkspc/errnie"
)

/*
Scanner is a wrapper around bufio.Scanner specifically geared
towards reading out streams from the builder daemon.
*/
type Scanner struct {
	handle *bufio.Scanner
	res    types.ImageBuildResponse
}

/*
NewScanner returns an instance of Scanner, used to retrieve lines
from an input stream.
*/
func NewScanner(res types.ImageBuildResponse) Scanner {
	return Scanner{
		handle: bufio.NewScanner(res.Body),
		res:    res,
	}
}

/*
Scan the log stream from the builder daemon so we can get
the lines to output.
*/
func (scanner Scanner) Scan() {
	// It's an io reader, so it wants to be closed.
	defer scanner.res.Body.Close()

	if !viper.GetBool("debug") {
		// Set the terminal to raw mode so we can better print
		// the streaming output and make sure it is reset to cooked when done.
		current := console.Current()
		errnie.Handles(current.SetRaw()).With(errnie.KILL)
		defer current.Reset()
	}

	// var res map[string]interface{}

	// // Loop over the scanner to print out the log stream.
	// for scanner.handle.Scan() {
	// 	scanner.print(scanner.handle.Text(), res)
	// }

	scanner.printalt()

	// Handle any errors that may have occureed during scanning.
	errnie.Handles(scanner.handle.Err()).With(errnie.KILL)
}

/*
print unmarshals the json string we are getting from the builder
daemon and prints it to stdout.
*/
func (scanner Scanner) print(str string, res map[string]interface{}) {
	if !viper.GetBool("debug") {
		// Move the cursor to the left.
		fmt.Print("\033[G\033[K")
	}

	json.Unmarshal([]byte(str), &res)
	fmt.Print(res["stream"])

	if !viper.GetBool("debug") {
		// Move the cursor up.
		fmt.Print("\033[A")
	}
}

func (scanner Scanner) printalt() {
	termFd, isTerm := term.GetFdInfo(os.Stderr)
	jsonmessage.DisplayJSONMessagesStream(scanner.res.Body, os.Stderr, termFd, isTerm, nil)
}
