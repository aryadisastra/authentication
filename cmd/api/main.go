package main

import (
	"log"

	"github.com/aryadisastra/authentication/internal/config"
	"github.com/aryadisastra/authentication/internal/db"
	"github.com/aryadisastra/authentication/internal/router"

	_ "github.com/aryadisastra/authentication/internal/swagger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg := config.Load()
	gdb := db.Open(cfg.DBDsn)

	r := router.New(gdb, cfg.JWTSecret, cfg.JWTExpiresMin)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("auth service listening on :%s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
