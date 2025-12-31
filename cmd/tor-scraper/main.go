package main

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/ali-ellikci/TorScraper/internal/input"
	"github.com/ali-ellikci/TorScraper/internal/logger"
	"github.com/ali-ellikci/TorScraper/internal/output"
	"github.com/ali-ellikci/TorScraper/internal/scanner"
	"github.com/ali-ellikci/TorScraper/internal/tor"
)

func main() {
	appLogger, err := logger.New()
	if err != nil {
		log.Fatal(err)
	}
	defer appLogger.Close()

	targets, err := input.ReadTargets("configs/targets.yaml")
	if err != nil {
		appLogger.Error("Failed to read targets: %v", err)
		log.Fatal(err)
	}

	client, err := tor.NewTorClient()
	if err != nil {
		appLogger.Error("Failed to create TOR client: %v", err)
		log.Fatal(err)
	}

	appLogger.Info("Starting TOR Scraper with %d targets", len(targets))

	reportWriter := output.NewReportWriter(len(targets))

	scanTargets(targets, client, appLogger, reportWriter)

	// Save report
	err = reportWriter.Save()
	if err != nil {
		appLogger.Error("Failed to save report: %v", err)
	} else {
		appLogger.Info("Report saved: %s", reportWriter.GetReportPath())
	}

	appLogger.Info("========================================")
	successCount, failCount := reportWriter.GetStats()
	appLogger.Info("Total: %d, Success: %d, Failed: %d", len(targets), successCount, failCount)
	appLogger.Info("Screenshots: output/screenshots/")
	appLogger.Info("HTML files: output/html/")
	appLogger.Info("Log file: output/scan_report_*.log")
	appLogger.Info("JSON Report: %s", reportWriter.GetReportPath())
	appLogger.Info("========================================")
}

func scanTargets(targets []string, client interface{}, appLogger *logger.Logger, reportWriter *output.ReportWriter) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 3)

	for i, t := range targets {
		wg.Add(1)
		go func(index int, targetURL string) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			appLogger.Info("[%d/%d] Scanning: %s", index+1, len(targets), targetURL)

			result, err := scanner.ScanTarget(context.Background(), client.(*http.Client), targetURL)
			if err != nil {
				appLogger.Error("%s -> %v", targetURL, err)
				reportWriter.AddError(targetURL, err.Error())
				return
			}

			screenshotPath, err := output.SaveScreenshot(result.Screenshot, targetURL)
			if err != nil {
				appLogger.Warn("Screenshot save failed for %s: %v", targetURL, err)
				screenshotPath = ""
			} else {
				appLogger.Success("Screenshot saved: %s", screenshotPath)
			}

			htmlPath, err := output.SaveHTML(result.HTML, targetURL)
			if err != nil {
				appLogger.Warn("HTML save failed for %s: %v", targetURL, err)
				htmlPath = ""
			} else {
				appLogger.Success("HTML saved: %s", htmlPath)
			}

			appLogger.Success("%s (Status: %d, IP: %s)", targetURL, result.StatusCode, result.IPAddress)
			reportWriter.AddSuccess(targetURL, result.StatusCode, result.IPAddress, screenshotPath, htmlPath)
		}(i+1, t)
	}

	wg.Wait()
}
