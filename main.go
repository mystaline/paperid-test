package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mystaline/paperid-test/internal/app"
	"github.com/mystaline/paperid-test/internal/config"
)

func main() {
	appConfig := config.LoadConfig()
	a := app.New(appConfig)

	// Listen for interrupt signals to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := a.Run(); err != nil {
			panic(err)
		}
	}()

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.Shutdown(); err != nil {
		panic(err)
	}

	<-ctx.Done()
}
