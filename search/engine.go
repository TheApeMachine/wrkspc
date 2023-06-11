package search

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/theapemachine/wrkspc/tweaker"
	"github.com/wrk-grp/errnie"
	customsearch "google.golang.org/api/customsearch/v1"
	"google.golang.org/api/option"
)

/*
Engine is a wrapper around the Google Custom Search Engine API.
*/
type Engine struct {
	client *customsearch.Service
	query  string
	result []*customsearch.Result
	err    error
}

/*
NewEngine creates a new Engine object.
*/
func NewEngine() *Engine {
	apiKey := tweaker.GetString("gcp.searchKey")
	svc, err := customsearch.NewService(context.TODO(), option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}

	return &Engine{
		client: svc,
	}
}

/*
Query returns the query string.
*/
func (e *Engine) Query(target string) string {
	cx := tweaker.GetString("gcp.engineID")
	// Write the incoming bytes to the query string.
	chunks := strings.Split(target, ",")

	// Wrap the chunks in `intext:"<chunk>"` and join them with spaces.
	for i, c := range chunks {
		chunks[i] = `intext:"` + c + `"`
	}
	e.query = strings.Join(chunks, " ")

	// Add fileType:pdf to the query string.
	e.query += ` fileType:pdf`

	// Perform the search.
	var res *customsearch.Search
	if res, e.err = e.client.Cse.List().Cx(cx).Q(e.query).Do(); e.err != nil {
		errnie.Handles(e.err)
		return ""
	}

	// Store the results and return.
	e.result = res.Items
	return ""
}

func (e *Engine) Result() {
	for _, r := range e.result {
		log.Println(r.Title)
		log.Println(r.Link)
		log.Println(r.Snippet)

		// If the link is a PDF, download it, and use the pdfcpu library to
		// extract the text from it.
		if strings.HasSuffix(r.Link, ".pdf") {
			fileURL, err := url.Parse(r.Link)
			if err != nil {
				log.Fatal(err)
			}

			// Download the file.
			resp, err := http.Get(fileURL.String())
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			// Create a temporary file to store the PDF.
			tmpfile, err := ioutil.TempFile("", "wrkspc-*.pdf")
			if err != nil {
				log.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			// Write the PDF to the temporary file.
			_, err = io.Copy(tmpfile, resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			// Error handling omitted for brevity.
			_, r, _ := pdf.Open(tmpfile.Name())
			for no := 1; no < r.NumPage(); no++ { // Loop over each page.
				page := r.Page(no)
				rows, _ := page.GetTextByRow()
				for _, row := range rows { // Loop over each row of text in the page.
					for _, text := range row.Content { // Loop over each piece of text in the row.
						fmt.Println(text)
					}
				}
			}
		}
	}
}
