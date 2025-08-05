package main

import (
	"log"

	"github.com/Babushkin05/wb-orders-service/internal/config"
)

func main() {
	cfg := config.MustLoad()
	log.Printf("Loaded config: %+v\n", cfg)
}
