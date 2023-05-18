package main

import (
	"github.com/theapemachine/wrkspc/cmd"
	"github.com/wrk-grp/errnie"
)

func main() {
	errnie.Handles(cmd.Execute())
}
