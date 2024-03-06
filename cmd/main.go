package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/egregors/rates"
	"github.com/egregors/rates/backends"
	"github.com/egregors/rates/internal/server/api"
	"github.com/egregors/rates/internal/server/web"
	"github.com/egregors/rates/lib/cache"
)

func main() {
	logger := log.Default()
	cNPM := rates.New(
		backends.NewCurrencyAPI("https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies/"),
		rates.WithLogger(logger),
		rates.WithCache(cache.NewInMem[map[string]float64](10*time.Second)),
	)

	// TODO: add fallback backends: https://currency-api.pages.dev/v1/currencies/

	go func() {
		if err := api.New(cNPM, logger).Run(); err != nil {
			logger.Fatalf("server failed: %v", err)
		}
	}()

	go func() {
		if err := web.New(cNPM, logger).Run(); err != nil {
			logger.Fatalf("server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	logger.Println("[INFO] Shutting down the server...")
}
