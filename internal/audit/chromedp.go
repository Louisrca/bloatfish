package audit

import (
	"context"
	"time"

	"github.com/Louisrca/bloatfish/internal/utils"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type WebAuditResult struct {
	URL           string
	Requests      int
	TransferredKB float64
	DomNodes      int
	ExternalReqs  int
	Deep          bool
}

func DeepAudit(url string) (*WebAuditResult, error) {
	// Create a new Chrome context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var totalBytes int64
	var requestCount int

	// Create a new context with the allocator
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// Set up the listener
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch e := ev.(type) {
		case *network.EventLoadingFinished:
			totalBytes += int64(e.EncodedDataLength)
			requestCount++
		}
	})

	err := chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.Sleep(10*time.Second),
	)
	if err != nil {
		return nil, err
	}

	utils.WriteJSONReport(WebAuditResult{
		URL:           url,
		Requests:      requestCount,
		TransferredKB: float64(totalBytes) / 1024,
		Deep:          true,
	}, "web_audit_report.json")

	return &WebAuditResult{
		URL:           url,
		Requests:      requestCount,
		TransferredKB: float64(totalBytes) / 1024,
		Deep:          true,
	}, nil
}
