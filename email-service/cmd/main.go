package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"email-service/internal/container"

	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env != "docker" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

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
