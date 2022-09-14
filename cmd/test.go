package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/theapemachine/wrkspc/datura"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spd"
)

func init() {
	rootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A command to test functionalities.",
	Long:  testtxt,
	RunE: func(_ *cobra.Command, _ []string) error {
		errnie.Tracing(false)
		errnie.Debugging(false)
		// store := datura.NewS3()
		store := datura.NewRadix()

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
					errnie.Logs(count, "objs/sec", store.PoolSize()).With(errnie.INFO)
					count = 0
				default:
					// Write a datapoint and increase the count.
					dg := spd.NewCached(
						"datapoint", "test", "test.wrkspc.org",
						"test",
					)

					store.Write(dg)
					count++
				}
			}
		}()

		// Run for 10 seconds, then stop.
		time.Sleep(2 * time.Second)
		done <- struct{}{}

		go func() {
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					// Output the count every second, then reset for
					// the next sample.
					errnie.Logs(count, "objs/sec", store.PoolSize()).With(errnie.INFO)
					count = 0
				default:
					// Write a datapoint and increase the count.
					dg := spd.NewCached(
						"datapoint", "test", "test.wrkspc.org",
						"v4.0.0/datapoint/test/test.wrkspc.org",
					)

					store.Read(dg)
					count++
				}
			}
		}()

		// Run for 10 seconds, then stop.
		time.Sleep(2 * time.Second)
		ticker.Stop()
		done <- struct{}{}

		return nil
	},
}

var testtxt = `
A "playground" command to test functionalities quickly.
`
