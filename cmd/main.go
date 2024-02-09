package main

import (
	"log"

	"github.com/egregors/rates/internal/provider"
	"github.com/egregors/rates/internal/server/api"
)

func main() {
	logger := log.Default()

	currencyAPI := provider.NewCurrencyAPI(logger)

	srv := api.New(currencyAPI, logger)
	if err := srv.Run(); err != nil {
		logger.Fatalf("server failed: %v", err)
	}
}
