package main

import (
	"chatappv2/config"
	"chatappv2/internal"
	"log"
)

func main() {
	cfg, err := config.LoadApplicationConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	db, _ := config.InitPostgresConnection(cfg)
	defer db.Close()

	r := internal.InitRouter(db)
	r.Run(cfg.Server.Port)
}
