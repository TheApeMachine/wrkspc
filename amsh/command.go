package amsh

import "github.com/theapemachine/wrkspc/errnie"

type Command struct{}

func NewCommand() *Command {
	return &Command{}
}

func (command *Command) Error() errnie.Error {
	return errnie.NewError(nil)
}

func (command *Command) Execute() *Command {
	return command
}
