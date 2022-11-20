package cmd

import (
	"bytes"
	"embed"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
	"github.com/theapemachine/wrkspc/brazil"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/tweaker"
)

/*
Embed a mini filesystem into the binary with the contents of . Some code below
will check if the user has the configuration file locally, otherwise write a fresh default
from the binary's embedded filesystem. Useful, no management needed.
*/
//go:embed cfg/*
var embedded embed.FS

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "wrkspc",
		Short: "A dynamic workspace environment built on containers.",
		Long:  roottxt,
	}
)

var roottxt = `
wrkspc builds a dynamic working environment on top of containers and
Kubernetes. It requires only a single binary to run and sources all
other tooling dynamically from configures registries and repositories.
`

func Execute() error {
	errnie.Trace()

	rootCmd.PersistentFlags().StringVarP(
		&orchestrator, "orchestrator", "o", "kubernetes",
		"The orchestrator to use <nomad|kubernetes>.",
	)

	rootCmd.AddCommand(runCmd)
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Defines the config file that will be loaded, usually just the name of the service.
	// This should be written to the user's home directory as a hidden file.
	rootCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		".wrkspc.yml",
		"config file (default is $HOME/.wrkspc.yml)",
	)

	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
}

/*
initConfig does the embedded config stuff and sets the entire program up for Viper
based config, which uses the embedded yaml config file a lot.
*/
func initConfig() {
	// Set verbosity level for errnie.
	errnie.Tracing(tweaker.GetBool("errnie.trace"))
	errnie.Debugging(tweaker.GetBool("errnie.debug"))

	errnie.Trace()

	// Get the config file from the user home path.
	chunks := strings.Split(cfgFile, "/")
	fh, err := embedded.Open("cfg/" + chunks[len(chunks)-1])

	errnie.Handles(err)
	defer fh.Close()

	buf, err := io.ReadAll(fh)
	errnie.Handles(err)

	home := brazil.NewPath("~").Location
	brazil.NewFile(home, cfgFile, bytes.NewBuffer(buf))

	viper.AddConfigPath(home)
	viper.SetConfigType("yml")
	viper.SetConfigName(cfgFile)
	viper.AutomaticEnv()
	viper.ReadInConfig()

	// The method errnie is wrapping here writes the markdown documentation for
	// the command line interface, which is automatically generated.
	brazil.NewPath(brazil.NewPath(".").Location, "docs")
	errnie.Handles(
		doc.GenMarkdownTree(rootCmd, "./docs/"),
	)
}
