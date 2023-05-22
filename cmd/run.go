package cmd

import (
	"context"
	"io/ioutil"
	"log"
	"fmt"
	language "cloud.google.com/go/language/apiv1"
	"cloud.google.com/go/language/apiv1/languagepb"
	"github.com/spf13/cobra"
)

/*
runCmd is a proxy for running any terminal command using a container
which is dynamically built from an image in a configured registry.
*/
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the service with the ~/.wrkspc.yml config values.",
	Long:  runtxt,
	RunE: func(_ *cobra.Command, _ []string) (err error) {
		// Load the data.txt file from the current directory and call analyzeEntities
		content, err := ioutil.ReadFile("data.txt")
		if err != nil {
			log.Fatal(err)
		}

		ctx := context.Background()
		client, err := language.NewClient(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		st, err := analyzeEntities(ctx, client, string(content))

		if err != nil {
			log.Fatal(err)
		}

		for _, e := range st.GetEntities() {
			fmt.Println(e.GetName())
			fmt.Println(e.GetType())
		}

		return nil
	},
}

/*
runtxt lives here to keep the command definition section cleaner.
*/
var runtxt = `
Use this sub command to proxy any terminal command through and it will
look for an existing image in the configured registry which has the command
included, build that image into a container and deploy it onto the
Kubernetes cluster that will be created first.
`
func analyzeEntities(ctx context.Context, client *language.Client, text string) (*languagepb.AnalyzeEntitiesResponse, error) {
	return client.AnalyzeEntities(ctx, &languagepb.AnalyzeEntitiesRequest{
			Document: &languagepb.Document{
					Source: &languagepb.Document_Content{
							Content: text,
					},
					Type: languagepb.Document_PLAIN_TEXT,
			},
			EncodingType: languagepb.EncodingType_UTF8,
	})
}