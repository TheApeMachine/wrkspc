package crawler

import (
	"context"
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/theapemachine/wrkspc/luthor"
)

/*
Browser is a wrapper around chromedp.
*/
type Browser struct {
	embedded embed.FS
	result   string
}

/*
NewBrowser creates a new Browser.
*/
func NewBrowser(embedded embed.FS) *Browser {
	return &Browser{embedded, ""}
}

/*
Run a new browser session.
*/
func (b *Browser) Run(target string) error {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	if err := chromedp.Run(ctx,
		chromedp.Tasks{
			chromedp.Navigate("https://www.nu.nl/spanningen-oekraine/6262108/onrust-en-onduidelijkheid-na-vermeende-droneaanval-op-kremlin.html"),
			NewScript(b.embedded, "humanize").Run(ctx, "getpage", "", &b.result),
			chromedp.ActionFunc(func(ctx context.Context) error {
				return luthor.NewParser(b.result).GetArticle(&b.result)
			}),
		},
	); err != nil {
		log.Fatal(err)
	}

	fmt.Println(b.result)

	return nil
}
