package cmd

import (
	"github.com/spf13/cobra"
	"github.com/theapemachine/wrkspc/apeterm"
	"github.com/theapemachine/wrkspc/eddie"
)

func init() {
	rootCmd.AddCommand(eddieCmd)
}

var eddieCmd = &cobra.Command{
	Use:   "eddie",
	Short: "The Ape Machine editor.",
	Long:  eddietxt,
	RunE: func(_ *cobra.Command, args []string) error {
		buffer := eddie.NewBuffer(
			apeterm.NewUI(), &args[0],
		)
		buffer.SetFocus()
		return nil
	},
}

var eddietxt = `
A terminal based editor.
`
