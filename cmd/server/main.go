package main

import (
	"log"

	"github.com/lvestera/slot-machine/internal/server"
	"github.com/lvestera/slot-machine/internal/server/config"
)

func main() {
	cfg := config.NewConfig()
	server := server.NewServer(cfg)

	if err := server.Run(); err != nil {
		log.Fatal("Can't start server: " + err.Error())
	}
}
