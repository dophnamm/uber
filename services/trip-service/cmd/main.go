package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"ride-sharing/services/trip-service/internal/database"
	tripHandlerV1 "ride-sharing/services/trip-service/internal/infrastructure/http/v1/trip"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/infrastructure/routes"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/services/trip-service/pkg/config"
)

const shutdownTimeout = 10 * time.Second

func main() {
	cfg := config.InitConfig()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := database.InitDB(ctx, cfg)
	if err != nil {
		if ctx.Err() != nil {
			log.Println("Shutdown requested during startup")
			return
		}
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

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}

	go func() {
		log.Printf("Listening and serving HTTP on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	stop()
	log.Println("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown: %v", err)
	}
}
