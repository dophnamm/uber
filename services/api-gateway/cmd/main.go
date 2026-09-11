package main

import (
	"fmt"
	"log"
	"ride-sharing/services/api-gateway/pkg/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.InitConfig()

	if cfg.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.SetTrustedProxies(nil)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "Ok",
		})
	})

	if err := router.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		log.Fatal(err)
	}
}
