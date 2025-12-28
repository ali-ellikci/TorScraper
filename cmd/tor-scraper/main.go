package main

import (
	"log"

	"github.com/ali-ellikci/TorScraper/internal/input"
	"github.com/ali-ellikci/TorScraper/internal/scanner"
	"github.com/ali-ellikci/TorScraper/internal/tor"
)

func main() {
	targets, err := input.ReadTargets("configs/targets.yaml")
	if err != nil {
		log.Fatal(err)
	}

	client, err := tor.NewTorClient()
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range targets {
		log.Printf("[INFO] Scanning: %s", t)

		_, err := scanner.ScanTarget(client, t)
		if err != nil {
			log.Printf("[ERR] %s -> %v", t, err)
			continue
		}

		log.Printf("[SUCCESS] %s", t)
	}
}
