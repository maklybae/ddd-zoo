package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maklybae/ddd-zoo/internal/application/services"
	"github.com/maklybae/ddd-zoo/internal/infrastructure/persistence/inmemory"
	httpserver "github.com/maklybae/ddd-zoo/internal/presentation/http"
	v1 "github.com/maklybae/ddd-zoo/internal/types/openapi/v1"
	"github.com/maklybae/ddd-zoo/pkg/events"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Initialize repositories
	animalRepo := inmemory.NewAnimalRepository()
	enclosureRepo := inmemory.NewEnclosureRepository()
	feedingScheduleRepo := inmemory.NewFeedingScheduleRepository()

	// Initialize events dispatcher
	eventsDispatcher := events.NewEventDispatcher()

	// Initialize services
	timeProvider := services.NewRealTimeProvider()
	animalTransferSvc := services.NewAnimalTransfer(animalRepo, enclosureRepo, eventsDispatcher, timeProvider)
	feedingOrganizationSvc := services.NewFeedingOrganization(animalRepo, feedingScheduleRepo, eventsDispatcher, timeProvider)
	statisticsSvc := services.NewZooStatistics(animalRepo, enclosureRepo, feedingScheduleRepo)

	// Initialize HTTP server
	server := httpserver.NewServer(
		animalRepo,
		enclosureRepo,
		feedingScheduleRepo,
		animalTransferSvc,
		feedingOrganizationSvc,
		statisticsSvc,
		timeProvider,
	)

	// Initialize Gin router
	router := gin.Default()

	// Serve OpenAPI specification file statically
	router.StaticFile("/api/openapi.yaml", "./api/openapi/v1/ddd_zoo.yaml")

	// Register OpenAPI handlers
	v1.RegisterHandlers(router, server)

	// Setup Swagger UI using our OpenAPI specification
	url := ginSwagger.URL("/api/openapi.yaml") // The URL pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Configure and start HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Println("Starting HTTP server on :8080")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
		return
	}

	log.Println("Server exited")
}
