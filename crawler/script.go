package crawler

import (
	"context"
	"embed"
	"fmt"
	"io"
	"log"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

/*
Script is a wrapper around a javascript file.
*/
type Script struct {
	data   []byte
	module string
}

/*
NewScript creates a new Script.
*/
func NewScript(embedded embed.FS, module string) *Script {
	// Load the script containing the humanSpider function
	jsFile, err := embedded.Open("cfg/scripts/" + module + ".js")
	if err != nil {
		log.Printf("Error opening human-spider.js: %v", err)
		return nil
	}
	defer jsFile.Close()

	script, err := io.ReadAll(jsFile)
	if err != nil {
		log.Printf("Error reading human-spider.js: %v", err)
		return nil
	}

	return &Script{script, module}
}

/*
Run the script.
*/
func (s *Script) Run(ctx context.Context, method, param string, result *string) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		// Evaluate the script to define the humanSpider function in the page
		_, exp, err := runtime.Evaluate(string(s.data)).Do(ctx)
		if err != nil || exp != nil {
			log.Printf("Error evaluating human-spider.js: error, %v, exception: %v", err, exp)
			return nil
		}

		return chromedp.Evaluate(fmt.Sprintf(`
		(() => {
			return %s().%s(%s);
		})()
		`, s.module, method, param,
		), result).Do(ctx)
	})
}
