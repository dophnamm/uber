package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"ride-sharing/services/trip-service/internal/database"
	tripHandlerV1 "ride-sharing/services/trip-service/internal/infrastructure/http/v1/trip"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/infrastructure/routes"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/services/trip-service/pkg/config"
)

func main() {
	cfg := config.InitConfig()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := database.InitDB(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.SetTrustedProxies(nil)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Repositories
	tripRepo := repository.NewTripRepository(db)
	rideFareRepo := repository.NewRideFareRepository(db)

	// Services
	tripService := service.NewTripService(tripRepo, rideFareRepo)

	// Handlers
	tripHandler := tripHandlerV1.NewTripHandlerV1(tripService)

	// Router
	routes.NewTripRoutesV1(router, tripHandler)

	if err := router.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		log.Fatal(err)
	}
}
