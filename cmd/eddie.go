package cmd

import (
	"io/fs"
	"os"

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
		var fh *os.File
		var err error

		if len(args) > 0 {
			fh, err = os.OpenFile(
				args[0], os.O_APPEND, fs.FileMode(os.O_CREATE),
			)

			errnie.Handles(err)
		}

		buffer := eddie.NewBuffer(fh)
		buffer.Init().Focus()

		return nil
	},
}

var eddietxt = `
A terminal based editor.
`
