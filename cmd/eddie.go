package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/theapemachine/wrkspc/eddie"
	"github.com/theapemachine/wrkspc/errnie"
)

func init() {
	rootCmd.AddCommand(eddieCmd)
}

var eddieCmd = &cobra.Command{
	Use:   "eddie",
	Short: "The Ape Machine editor.",
	Long:  eddietxt,
	RunE: func(_ *cobra.Command, args []string) error {
		return errnie.Handles(
			tea.NewProgram(eddie.Buffer{}).Start(),
		).With(errnie.KILL)
	},
}

var eddietxt = `
A terminal based editor.
`
