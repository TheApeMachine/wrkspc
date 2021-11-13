package cmd

import (
	"github.com/spf13/cobra"
	"github.com/theapemachine/wrkspc/amsh"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/fellini"
	"github.com/theapemachine/wrkspc/hefner"
	"github.com/theapemachine/wrkspc/kubrick"
	"github.com/theapemachine/wrkspc/twoface"
)

func init() {
	rootCmd.AddCommand(amshCmd)
}

var amshCmd = &cobra.Command{
	Use:   "amsh",
	Short: "Ape Machine Shell is an interactive console environment that interfaces with wrkspc.",
	Long:  longamshtxt,
	RunE: func(cmd *cobra.Command, args []string) error {
		errnie.Traces()

		// Get a handle on a new Ape Machine Shell Buffer, which will run headless and concurrently
		// putting the host terminal in `raw` mode.
		sh := amsh.NewBuffer(args, twoface.NewDisposer(), hefner.NewPipe(hefner.ProtoPipe{}))

		// Setup the terminal UI.
		screen := kubrick.NewScreen(kubrick.NewLayout(
			kubrick.FullScreenLayout{
				// Pass the shell into the Console template and have Execute return its
				// data channel to the Template, which will direct it to the correct
				// Component.
				Template: fellini.NewTemplate(fellini.Console{}, sh.Execute()),
			},
		))

		// Block the main goroutine until we receive an error.
		return <-screen.Render()
	},
}

var longamshtxt = `
Using the interactive shell allows for more chained and otherwise automatable workflows.
`
