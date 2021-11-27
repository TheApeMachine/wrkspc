package matrix

import (
	"bufio"
	"os"

	"github.com/containerd/console"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/moby/term"
	"github.com/spf13/viper"
	"github.com/theapemachine/wrkspc/errnie"
)

/*
Scanner is a wrapper around bufio.Scanner to output Docker containter build log streams.
*/
type Scanner struct {
	handle *bufio.Scanner
	res    types.ImageBuildResponse
}

/*
NewScanner returns a configured Scanner.
*/
func NewScanner(res types.ImageBuildResponse) Scanner {
	return Scanner{
		handle: bufio.NewScanner(res.Body),
		res:    res,
	}
}

/*
Scan the builder output log.
*/
func (scanner Scanner) Scan() {
	// It's an io reader, so it wants to be closed.
	defer scanner.res.Body.Close()

	if !viper.GetBool("debug") {
		// Set the terminal to raw mode.
		current := console.Current()
		errnie.Handles(current.SetRaw()).With(errnie.KILL)
		// Back to cooked at return time.
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
print unmarshals the json string and prints it to stdout.
*/
// func (scanner Scanner) print(str string, res map[string]interface{}) {
// 	// Move the cursor to the left.
// 	fmt.Print("\033[G\033[K")

// 	json.Unmarshal([]byte(str), &res)
// 	fmt.Print(res["stream"])

// 	// Move the cursor up.
// 	fmt.Print("\033[A")
// }

/*
printalt is an alternative method to print the log stream.
*/
func (scanner Scanner) printalt() {
	termFd, isTerm := term.GetFdInfo(os.Stderr)
	jsonmessage.DisplayJSONMessagesStream(scanner.res.Body, os.Stderr, termFd, isTerm, nil)
}
