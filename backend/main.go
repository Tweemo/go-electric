package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tweemo/go-electric/api"
	"github.com/tweemo/go-electric/config"
	"github.com/tweemo/go-electric/rates"
)

func main() {
	cfg := config.Load()
	if cfg.GinMode != "" {
		gin.SetMode(cfg.GinMode)
	}

	r, err := rates.Load("data/rates.json")
	if err != nil {
		log.Fatalf("load rates: %v", err)
	}

	router := api.NewRouter(cfg, r)
	log.Fatal(router.Run(cfg.Addr()))
}
