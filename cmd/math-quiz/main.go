package main

import (
	"context"
	"log"
	"math-quiz/internal/app"
	"os"
	"os/signal"
)

func main() {
	if err := realMain(); err != nil {
		log.Fatal(err)
	}
}

func realMain() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := app.Run(ctx); err != nil {
		return err
	}

	return nil
}
