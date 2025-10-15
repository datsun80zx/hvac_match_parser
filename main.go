package main

import (
	"fmt"
	"log"

	"github.com/datsun80zx/hvac_match_parser/config"
	"github.com/datsun80zx/hvac_match_parser/internal/methods"
)

func main() {
	cfg := config.ParseFlags()

	if err := methods.ProcessMatches(cfg); err != nil {
		log.Fatalf("Error processing matches: %v", err)
	}

	fmt.Println("Successfully created filtered matches!")
}
