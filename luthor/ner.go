package luthor

import (
	"context"
	"fmt"
	"log"

	language "cloud.google.com/go/language/apiv1"
	"cloud.google.com/go/language/apiv1/languagepb"
)

/*
NER is the top level object which directs the parsing and extraction of
named entities from a string of text.
*/
type NER struct {
	text string
}

/*
NewNER creates a new NER.
*/
func NewNER(text string) *NER {
	return &NER{text}
}

/*
Read implements the io.Reader interface, and every call should return
the named entities found in the text.
*/
func (ner *NER) Read(p []byte) (n int, err error) {
	ctx := context.Background()
	client, err := language.NewClient(ctx)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()
	st, err := analyzeEntities(ctx, client, string(""))

	if err != nil {
		log.Fatal(err)
	}

	for _, e := range st.GetEntities() {
		fmt.Println(e.GetName())
		fmt.Println(e.GetType())
	}

	return 0, nil
}

/*
analyzeEntities is a helper function which calls the Google Cloud Natural Language API.
*/
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
