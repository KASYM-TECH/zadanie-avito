package main

import (
	"avito/assembly"
	"avito/config"
	"avito/log"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func run() error {
	cfg := config.LoadFromEnv("dev.env")
	ctx := context.Background()
	if _, err := cfg.Validate(ctx); err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger := log.NewLogger(cfg.AppMode)
	assembler := assembly.NewAssembler(logger)
	go func() {
		<-ctx.Done()
		if err := assembler.Close(ctx); err != nil {
			logger.Fatal(context.Background(), "failed to close assembler")
		}
		os.Exit(0)
	}()

	router, err := assembler.Assemble(ctx, cfg)
	if err != nil {
		logger.Fatal(ctx, err.Error())
		return err
	}

	if err = http.ListenAndServe(cfg.ServerAddress, router); err != nil {
		return err
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		os.Exit(1)
	}
}
