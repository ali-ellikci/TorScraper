package scanner

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type ScanResult struct {
	Screenshot []byte
	HTML       []byte
	StatusCode int64
	IPAddress  string
}

func ScanTarget(ctx context.Context, client *http.Client, targetURL string) (*ScanResult, error) {

	ipAddr, err := verifyTORConnection(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("[FAILED] TOR connection verification failed: %v", err)
	}
	fmt.Printf("[INFO] Using TOR IP: %s\n", ipAddr)

	htmlContent, statusCode, err := fetchHTMLWithClient(ctx, client, targetURL)
	if err != nil {
		return nil, fmt.Errorf("[FAILED] failed to fetch HTML: %v", err)
	}

	screenshot, err := takeScreenshot(ctx, targetURL)
	if err != nil {
		return nil, fmt.Errorf("[FAILED] failed to capture screenshot: %v", err)
	}

	result := &ScanResult{
		Screenshot: screenshot,
		HTML:       htmlContent,
		StatusCode: statusCode,
		IPAddress:  ipAddr,
	}

	return result, nil
}

func verifyTORConnection(ctx context.Context, client *http.Client) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://check.torproject.org/api/ip", nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to check TOR IP: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func fetchHTMLWithClient(ctx context.Context, client *http.Client, targetURL string) ([]byte, int64, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	htmlContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return htmlContent, int64(resp.StatusCode), nil
}

func takeScreenshot(ctx context.Context, targetURL string) ([]byte, error) {
	var screenshot []byte
	var statusCode int64

	ctxTimeout, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
	defer cancelTimeout()

	opts := []chromedp.ExecAllocatorOption{
		chromedp.ProxyServer("socks5://127.0.0.1:9050"),
		chromedp.DisableGPU,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-plugins", true),
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(ctxTimeout, opts...)
	defer allocCancel()

	ctxChrome, cancelChrome := chromedp.NewContext(allocCtx)
	defer cancelChrome()

	chromedp.ListenTarget(ctxChrome, func(ev interface{}) {
		if ev, ok := ev.(*network.EventResponseReceived); ok {
			if ev.Type == network.ResourceTypeDocument {
				statusCode = int64(ev.Response.Status)
			}
		}
	})

	err := chromedp.Run(
		ctxChrome,
		network.Enable(),
		chromedp.Navigate(targetURL),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			for {
				if statusCode != 0 {
					return nil
				}

				select {
				case <-timeoutCtx.Done():
					return fmt.Errorf("HTTP status beklenirken timeout")
				default:
					time.Sleep(100 * time.Millisecond)
				}
			}
		}),

		chromedp.FullScreenshot(&screenshot, 90),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to capture screenshot: %v", err)
	}

	return screenshot, nil
}
