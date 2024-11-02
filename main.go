package main

import (
	"fmt"
	"log"
	"os"

	"github.com/phetployst/art-toys-store/config"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Error: .env path is required")
	}

	config.LoadEnvFile(os.Args[1])

	osGetter := &config.OsEnvGetter{}

	configProvider := config.ConfigProvider{Getter: osGetter}
	config := configProvider.GetConfig()

	fmt.Println(config)
}
