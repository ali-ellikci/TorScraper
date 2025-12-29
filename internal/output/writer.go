package output

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func SaveScreenshot(data []byte, targetURL string) (string, error) {
	screenshotDir := "output/screenshots"

	err := os.MkdirAll(screenshotDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create screenshots directory: %v", err)
	}

	filename := generateSafeFilename(targetURL, ".png")
	filePath := filepath.Join(screenshotDir, filename)

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to save screenshot: %v", err)
	}

	return filePath, nil
}

func SaveHTML(htmlContent []byte, targetURL string) (string, error) {
	htmlDir := "output/html"

	err := os.MkdirAll(htmlDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create html directory: %v", err)
	}

	filename := generateSafeFilename(targetURL, ".html")
	filePath := filepath.Join(htmlDir, filename)

	err = os.WriteFile(filePath, htmlContent, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to save html: %v", err)
	}

	return filePath, nil
}

func generateSafeFilename(targetURL string, extension string) string {

	parsedURL, err := url.Parse(targetURL)
	var hostname string

	if err == nil && parsedURL.Hostname() != "" {
		hostname = parsedURL.Hostname()
	} else {

		hostname = strings.TrimPrefix(targetURL, "http://")
		hostname = strings.TrimPrefix(hostname, "https://")
		hostname = strings.TrimSuffix(hostname, "/")
	}

	// Remove invalid characters and limit length
	hostname = strings.NewReplacer(
		":", "_",
		"/", "_",
		"?", "_",
		"#", "_",
		"&", "_",
		"=", "_",
	).Replace(hostname)

	// Add timestamp to avoid collisions
	timestamp := time.Now().Format("20060102_150405")

	// Combine and add extension
	filename := fmt.Sprintf("%s_%s%s", hostname, timestamp, extension)

	// Limit filename length (max 255 chars on most filesystems)
	if len(filename) > 250 {
		filename = filename[:250] + extension
	}

	return filename
}
