package scanner

import (
	"fmt"
	"io"
	"net/http"
)

func ScanTarget(client *http.Client, targetURL string) ([]byte, error) {
	resp, err := client.Get(targetURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
