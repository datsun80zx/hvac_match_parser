package main

import (
	"fmt"
	"log"

	"github.com/datsun80zx/hvac_match_parser/config"
	"github.com/datsun80zx/hvac_match_parser/internal/methods"
)

func main() {
	// Load configuration
	cfg := config.ParseFlags()

	// Process matches
	if err := methods.ProcessMatches(cfg); err != nil {
		log.Fatalf("Error processing matches: %v", err)
	}

	fmt.Printf("Successfully created filtered matches in %s\n", cfg.OutputFile)
}
