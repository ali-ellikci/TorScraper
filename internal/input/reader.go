package input

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func ReadTargets(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Could not open targets file: %v", err)
		return make([]string, 0), err
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	targets := make([]string, 0)

	for scanner.Scan() {

		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		targets = append(targets, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return targets, nil

}
