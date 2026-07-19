package renderer

import (
	"context"
	"os"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// ChromedpRenderer implementa generatecv.PDFRenderer usando Chrome headless.
type ChromedpRenderer struct{}

func NewChromedpRenderer() *ChromedpRenderer {
	return &ChromedpRenderer{}
}

func (r *ChromedpRenderer) RenderPDF(ctx context.Context, html string) ([]byte, error) {
	tmpFile, err := os.CreateTemp("", "cv-*.html")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(html); err != nil {
		tmpFile.Close()
		return nil, err
	}
	tmpFile.Close()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, opts...)
	defer cancelAlloc()

	taskCtx, cancelTask := chromedp.NewContext(allocCtx)
	defer cancelTask()

	var pdfBuf []byte
	err = chromedp.Run(taskCtx,
		chromedp.Navigate("file://"+tmpFile.Name()),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			if err != nil {
				return err
			}
			pdfBuf = buf
			return nil
		}),
	)
	if err != nil {
		return nil, err
	}

	return pdfBuf, nil
}
