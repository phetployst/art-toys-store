package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/phetployst/art-toys-store/config"
	"github.com/phetployst/art-toys-store/server"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Error: .env path is required")
	}

	configProvider := config.ConfigProvider{
		Getter: &config.OsEnvGetter{},
		Loader: &config.GodotenvLoader{},
	}

	if err := configProvider.LoadEnvFile(os.Args[1]); err != nil {
		log.Fatalf("Failed to load .env file from path %s: %v", os.Args[1], err)
	}

	config, err := configProvider.GetConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	server.StartHTTPServer(ctx, &config)
}
