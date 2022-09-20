package cmd

import (
	"github.com/spf13/cobra"
	"github.com/theapemachine/wrkspc/eddie"
)

func init() {
	rootCmd.AddCommand(eddieCmd)
}

var eddieCmd = &cobra.Command{
	Use:   "eddie",
	Short: "The Ape Machine editor.",
	Long:  eddietxt,
	RunE: func(_ *cobra.Command, _ []string) error {
		buffer := eddie.NewBuffer()
		buffer.Init()

		return nil
	},
}

var eddietxt = `
A terminal based editor.
`
