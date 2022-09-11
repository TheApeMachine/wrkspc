package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/theapemachine/wrkspc/datura"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spd"
)

func init() {
	spd.InitCache()
	rootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A command to test functionalities.",
	Long:  testtxt,
	RunE: func(_ *cobra.Command, _ []string) error {
		errnie.Tracing(false)
		errnie.Debugging(true)
		store := datura.NewS3()

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
					errnie.Logs(count, "objs/sec").With(errnie.INFO)
					count = 0
				default:
					// Write a datapoint and increase the count.
					store.Write(spd.NewCached(
						"datapoint", "test", "test.wrkspc.org",
					))

					count++
				}
			}
		}()

		// Run for 10 seconds, then stop.
		time.Sleep(10 * time.Second)
		ticker.Stop()

		done <- struct{}{}
		return nil
	},
}

var testtxt = `
A "playground" command to test functionalities quickly.
`
