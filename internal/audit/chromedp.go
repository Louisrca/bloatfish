package audit

import (
	"context"
	"fmt"
	"math"
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
	HTMLSizeKB    int
	ExternalReqs  int
	Score         int
	Grade         string
	Deep          bool
}

func QDomNodes(domNodes int) int {

	if domNodes <= 300 {
		return 10
	} else if domNodes >= 301 && domNodes <= 800 {
		return 30
	} else if domNodes >= 801 && domNodes <= 1500 {
		return 60
	} else if domNodes > 1500 {
		return 4
	} else {
		return 0
	}
}

func QHTMLSize(htmlSizeKB int) int {

	if htmlSizeKB <= 100 {
		return 10
	} else if htmlSizeKB >= 101 && htmlSizeKB <= 300 {
		return 30
	} else if htmlSizeKB >= 301 && htmlSizeKB <= 600 {
		return 60
	} else if htmlSizeKB > 600 {
		return 4
	} else {
		return 0
	}
}

func QRequests(requests int) int {
	if requests <= 15 {
		return 10
	} else if requests >= 16 && requests <= 40 {
		return 30
	} else if requests >= 41 && requests <= 80 {
		return 60
	} else if requests > 80 {
		return 4
	} else {
		return 0
	}
}

func ScoreWebAudit(requests int, HTMLSizeKB int, DomNodes int) (int, string) {
	score := 100 - (QRequests(requests)+QHTMLSize(HTMLSizeKB)+QDomNodes(DomNodes))/3
	if score < 0 {
		return 0, "F"
	}

	if score > 100 {
		return 100, "A"
	} else if score >= 80 {
		return score, "B"
	} else if score >= 60 {
		return score, "C"
	} else if score >= 40 {
		return score, "D"
	} else if score >= 20 {
		return score, "E"
	} else if score >= 0 {
		return score, "F"
	}

	return score, "N/A"
}

func chromeDPAudit(url string) (*utils.ChromeDPReport, error) {
	var domNodes int
	var HTMLSizeKB string

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var totalBytes int64
	var requestCount int

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

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
		chromedp.OuterHTML("html", &HTMLSizeKB),
		chromedp.Evaluate(`
			document.getElementsByTagName("*").length
		`, &domNodes),
		chromedp.Sleep(10*time.Second),
	)
	if err != nil {
		fmt.Println("Error during chromedp run:", err)
		return nil, err
	}

	score, grade := ScoreWebAudit(QRequests(requestCount), QHTMLSize(int(math.Round(float64(len(HTMLSizeKB))/1024))), QDomNodes(domNodes))

	fmt.Println("grade:", grade, "score:", score)

	return &utils.ChromeDPReport{
		URL:           url,
		Requests:      requestCount,
		TransferredKB: float64(totalBytes) / 1024,
		DomNodes:      domNodes,
		HTMLSizeKB:    int(math.Round(float64(len(HTMLSizeKB)) / 1024)),
		Score:         score,
		Grade:         grade,
		Deep:          true,
	}, nil
}

func DeepAudit(urls []string) []*utils.ChromeDPReport {
	var deepAuditReports []*utils.ChromeDPReport
	for _, url := range urls {
		fmt.Printf("Auditing %s...\n", url)
		result, err := chromeDPAudit(url)
		if err != nil {
			fmt.Printf("❌ Error auditing %s: %v\n", url, err)
			continue
		}
		fmt.Printf("✅ Audit completed for %s: %+v\n", url, result)
		deepAuditReports = append(deepAuditReports, result)
	}

	return deepAuditReports

}
