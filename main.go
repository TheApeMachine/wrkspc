package main

import (
	"github.com/theapemachine/wrkspc/cmd"
	"github.com/theapemachine/wrkspc/errnie"
)

func main() {
	errnie.Handles(cmd.Execute())
}
