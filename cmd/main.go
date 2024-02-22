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
	c := rates.New(
		backends.NewCurrencyAPI(),
		rates.WithLogger(logger),
		rates.WithCache(cache.NewInMem[map[string]float64](10*time.Second)),
	)

	go func() {
		if err := api.New(c, logger).Run(); err != nil {
			logger.Fatalf("server failed: %v", err)
		}
	}()

	go func() {
		if err := web.New(c, logger).Run(); err != nil {
			logger.Fatalf("server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	logger.Println("[INFO] Shutting down the server...")
}
