package output

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type ScanRecord struct {
	URL            string    `json:"url"`
	Status         string    `json:"status"`
	StatusCode     int64     `json:"status_code,omitempty"`
	IPAddress      string    `json:"ip_address,omitempty"`
	Timestamp      time.Time `json:"timestamp"`
	ScreenshotPath string    `json:"screenshot_path,omitempty"`
	HTMLPath       string    `json:"html_path,omitempty"`
	Error          string    `json:"error,omitempty"`
}

type ScanReport struct {
	StartTime    time.Time    `json:"start_time"`
	EndTime      time.Time    `json:"end_time,omitempty"`
	TotalTargets int          `json:"total_targets"`
	SuccessCount int          `json:"success_count"`
	FailCount    int          `json:"fail_count"`
	Records      []ScanRecord `json:"records"`
}

type ReportWriter struct {
	report     *ScanReport
	reportFile string
}

func NewReportWriter(totalTargets int) *ReportWriter {
	reportFile := filepath.Join("output", fmt.Sprintf("scan_report_%s.json", time.Now().Format("20060102_150405")))

	report := &ScanReport{
		StartTime:    time.Now(),
		TotalTargets: totalTargets,
		Records:      make([]ScanRecord, 0),
	}

	return &ReportWriter{
		report:     report,
		reportFile: reportFile,
	}
}

func (rw *ReportWriter) AddSuccess(url string, statusCode int64, ipAddr, screenshotPath, htmlPath string) {
	record := ScanRecord{
		URL:            url,
		Status:         "SUCCESS",
		StatusCode:     statusCode,
		IPAddress:      ipAddr,
		Timestamp:      time.Now(),
		ScreenshotPath: screenshotPath,
		HTMLPath:       htmlPath,
	}

	rw.report.Records = append(rw.report.Records, record)
	rw.report.SuccessCount++
}

func (rw *ReportWriter) AddError(url string, errMsg string) {
	record := ScanRecord{
		URL:       url,
		Status:    "FAILED",
		Timestamp: time.Now(),
		Error:     errMsg,
	}

	rw.report.Records = append(rw.report.Records, record)
	rw.report.FailCount++
}

func (rw *ReportWriter) Save() error {
	rw.report.EndTime = time.Now()

	outputDir := filepath.Dir(rw.reportFile)
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	data, err := json.MarshalIndent(rw.report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %v", err)
	}

	err = os.WriteFile(rw.reportFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save report: %v", err)
	}

	return nil
}

func (rw *ReportWriter) GetReportPath() string {
	return rw.reportFile
}

func (rw *ReportWriter) GetStats() (int, int) {
	return rw.report.SuccessCount, rw.report.FailCount
}
