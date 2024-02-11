package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/egregors/rates/conv"
	"github.com/egregors/rates/internal/backends"
	"github.com/egregors/rates/internal/server/api"
	"github.com/egregors/rates/internal/server/web"
	"github.com/egregors/rates/lib/cache"
)

func main() {
	logger := log.Default()
	c := conv.New(
		backends.NewCurrencyAPI(),
		conv.WithLogger(logger),
		conv.WithCache(cache.NewInMem[map[string]float64]()),
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
