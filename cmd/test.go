package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/theapemachine/wrkspc/datura"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spd"
	"github.com/theapemachine/wrkspc/twoface"
)

func init() {
	rootCmd.AddCommand(testCmd)
}

func do(manager twoface.Employer, dg []byte) {
	ticker := time.NewTicker(1000 * time.Millisecond)
	done := make(chan struct{})
	count := 0

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				// Output the count every second, then reset for
				// the next sample.
				errnie.Logs(count, "objs/sec", manager.PoolSize()).With(errnie.INFO)
				count = 0
			default:

				manager.Write(dg)
				count++
			}
		}
	}()

	time.Sleep(10 * time.Second)
	done <- struct{}{}
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A command to test functionalities.",
	Long:  testtxt,
	RunE: func(_ *cobra.Command, _ []string) error {
		errnie.Tracing(false)
		errnie.Debugging(false)
		manager := datura.NewManager()
		// store := datura.NewRadix()

		// Write a datapoint and increase the count.
		dg := spd.NewCached(
			"datapoint", "test", "test.wrkspc.org",
			"test",
		)

		errnie.Debugs("writing...")
		do(manager, dg)

		dg = spd.NewCached(
			"datapoint", "test", "test.wrkspc.org",
			"v4.0.0/datapoint/test/test.wrkspc.org",
		)

		errnie.Debugs("reading...")
		do(manager, dg)

		return nil
	},
}

var testtxt = `
A "playground" command to test functionalities quickly.
`
