package main

import (
	"github.com/egregors/rates/internal/server/web"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/egregors/rates/internal/provider"
	"github.com/egregors/rates/internal/server/api"
)

func main() {
	logger := log.Default()

	currencyAPI := provider.NewCurrencyAPI(logger)

	go func() {
		apiSrv := api.New(currencyAPI, logger)

		if err := apiSrv.Run(); err != nil {
			logger.Fatalf("server failed: %v", err)
		}
	}()

	go func() {
		webSrv := web.New(currencyAPI, logger)
		if err := webSrv.Run(); err != nil {
			logger.Fatalf("server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	logger.Println("[INFO] Shutting down the server...")
}
