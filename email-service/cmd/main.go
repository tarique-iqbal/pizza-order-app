package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"email-service/internal/container"
)

func main() {
	app, err := container.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go app.Consumer.Run(ctx, app.Dispatcher)

	<-sigs
	log.Println("Exiting...")
	cancel()
	app.Consumer.Close()
}
