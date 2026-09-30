package exporter

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// exportWithChromedp renders a presentation HTML to 1080p landscape PDF via Chrome DevTools Protocol
func exportWithChromedp(chromeBin, htmlPath, pdfOut string, timeout time.Duration) error {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(chromeBin),
		chromedp.NoSandbox,
		chromedp.DisableGPU,
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("headless", "new"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	url := fmt.Sprintf("file://%s", htmlPath)
	var buf []byte

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body", chromedp.ByQuery),
		// Wait for web fonts (Poppins, Plus Jakarta Sans, JetBrains Mono) to finish downloading
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, _, err := runtime.Evaluate(`document.fonts.ready`).WithAwaitPromise(true).Do(ctx)
			return err
		}),
		// Debounce for final CSS layout calculations
		chromedp.Sleep(200*time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			buf, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPreferCSSPageSize(true).
				WithDisplayHeaderFooter(false).
				WithMarginTop(0).
				WithMarginBottom(0).
				WithMarginLeft(0).
				WithMarginRight(0).
				Do(ctx)
			return err
		}),
	)
	if err != nil {
		return fmt.Errorf("chromedp export failed: %w", err)
	}

	if len(buf) == 0 {
		return fmt.Errorf("chromedp generated 0 bytes for PDF")
	}

	return os.WriteFile(pdfOut, buf, 0644)
}
