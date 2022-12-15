package cmd

import (
	"bytes"
	"embed"
	"io"
	"strings"

	"github.com/elewis787/boa"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
	"github.com/theapemachine/wrkspc/brazil"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/tweaker"
)

/*
Embed a mini filesystem into the binary to hold the config file,
and some front end templates. This will be compiled into the
binary, so it is easier to manage.
*/
//go:embed cfg/* tmpl/*
var embedded embed.FS

var (
	cfgFile      string
	orchestrator string
	provision    string

	rootCmd = &cobra.Command{
		Use:   "wrkspc",
		Short: "A dynamic workspace environment built on containers.",
		Long:  roottxt,
	}
)

var roottxt = `
wrkspc allows you to easily build and deploy platforms, and contains
everything you need from development to infrastructure.
`

/*
Execute configures the CLI and executes the program with the
values that were passed in from the command line.
*/
func Execute() error {
	errnie.Trace()

	rootCmd.SetUsageFunc(boa.UsageFunc)
	rootCmd.SetHelpFunc(boa.HelpFunc)

	// Set the orchestrator we will use to build the infrastructure
	// of our platform.
	rootCmd.PersistentFlags().StringVarP(
		&orchestrator, "orchestrator", "o", "kubernetes",
		"The orchestrator to use <nomad|kubernetes>.",
	)

	// Set the provisioning mode.
	rootCmd.PersistentFlags().StringVarP(
		&provision, "provision", "p", "auto",
		"The provisioning mode to run with.",
	)

	// Add the `run` command to the CLI.
	rootCmd.AddCommand(runCmd)

	// Run the program and return any error that may happen.
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
initConfig unpacks the embedded file system and writes the config
file to the home directory of the user if it is not present.
It also automatically generates the CLI documentation for wrkspc.
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
