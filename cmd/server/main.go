package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/githamo/stubhub-tc/api/rest"
	encryption "github.com/githamo/stubhub-tc/internal/common/encryption"
	commonInfra "github.com/githamo/stubhub-tc/internal/common/infrastructure"
	healthApp "github.com/githamo/stubhub-tc/internal/health/application"
	healthInterfaces "github.com/githamo/stubhub-tc/internal/health/interface"
	trafficApp "github.com/githamo/stubhub-tc/internal/traffic/application"
	trafficInfra "github.com/githamo/stubhub-tc/internal/traffic/infrastructure"
	trafficInterfaces "github.com/githamo/stubhub-tc/internal/traffic/interface"
)

func main() {
	log.Println("Starting TrafficController service...")

	// load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}

	// set default port if not provided
	port := os.Getenv("PORT") // for heroku deployment
	if port == "" {
		port = os.Getenv("APP_PORT")
		if port == "" {
			port = "8080"
		}
	}

	// init database connection
	db, err := commonInfra.NewDatabaseConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// repositories
	crypto := encryption.NewHelper()
	trafficRepo := trafficInfra.NewMySQLRepository(db, crypto)

	// application services
	trafficService := trafficApp.NewTrafficService(trafficRepo)
	healthService := healthApp.NewHealthService(db)

	// HTTP handlers
	trafficHandler := trafficInterfaces.NewTrafficHandler(trafficService)
	healthHandler := healthInterfaces.NewHealthHandler(healthService)

	// router
	router := rest.SetupRouter()

	// routes
	trafficHandler.RegisterRoutes(router)
	healthHandler.RegisterRoutes(router)

	// appplication server
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// start server in a goroutine so it doesn't block
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// create a deadline to wait for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
