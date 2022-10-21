package nomad

import (
	"os"

	// These packages have init() funcs which check os.Args and drop directly
	// into their command logic. This is because they are run as separate
	// processes along side of a task. By early importing them we can avoid
	// additional code being imported and thus reserving memory.
	_ "github.com/hashicorp/nomad/client/logmon"
	_ "github.com/hashicorp/nomad/drivers/docker/docklog"
	_ "github.com/hashicorp/nomad/drivers/shared/executor"
	"github.com/theapemachine/wrkspc/errnie"

	// Don't move any other code imports above the import block above!
	"github.com/hashicorp/nomad/command"
	"github.com/hashicorp/nomad/version"
	"github.com/mitchellh/cli"
	"github.com/sean-/seed"
)

func init() {
	seed.Init()
}

type Infra struct {
}

func NewInfra() *Infra {
	return &Infra{}
}

func (infra *Infra) Run(args []string) int {
	// Create the meta object
	metaPtr := new(command.Meta)
	metaPtr.SetupUi(args)

	// The Nomad agent never outputs color
	agentUi := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	return infra.exec(metaPtr, agentUi, []string{})
}

func (infra *Infra) exec(
	metaPtr *command.Meta,
	agentUi *cli.BasicUi,
	args []string,
) int {
	commands := command.Commands(metaPtr, agentUi)

	cli := &cli.CLI{
		Name:     "nomad",
		Version:  version.GetVersion().FullVersionNumber(true),
		Args:     args,
		Commands: commands,
	}

	exitCode, err := cli.Run()
	errnie.Handles(err)
	return exitCode
}
