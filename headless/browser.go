package headless

import (
	"context"

	"github.com/chromedp/chromedp"
	"github.com/theapemachine/wrkspc/twoface"
)

type Browser struct {
	ctx    *twoface.Context
	remote bool
}

func NewBrowser(ctx *twoface.Context) *Browser {
	return &Browser{
		ctx: ctx,
	}
}

func (browser *Browser) Read(p []byte) (n int, err error) { return }

func (browser *Browser) Write(p []byte) (n int, err error) {
	allocCtx, allocCnl := chromedp.NewRemoteAllocator(
		context.Background(),
		"ws://127.0.0.1:9222/devtools/browser/b720b2f6-3342-4e4a-9b13-c53370ac27e3",
	)
	defer allocCnl()

	cdpCtx, cdpCnl := chromedp.NewContext(allocCtx)
	defer cdpCnl()

	chromedp.Run(cdpCtx, chromedp.ActionFunc(func(ctx context.Context) error {
		return err
	}))

	return
}

func (browser *Browser) Close() error {
	return nil
}
