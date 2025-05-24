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

// @desc start application and listen to shutdown signal
func main() {
	log.Println("Starting TrafficController service...")

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}

	server := setupServer()
	startGracefulShutdown(server)
}

// @desc build and return the HTTP server with all dependencies
func setupServer() *http.Server {
	port := getPort()

	db, err := commonInfra.NewDatabaseConnection()
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	crypto := encryption.NewHelper()
	trafficRepo := trafficInfra.NewMySQLRepository(db, crypto)

	trafficService := trafficApp.NewTrafficService(trafficRepo)
	healthService := healthApp.NewHealthService(db)

	trafficHandler := trafficInterfaces.NewTrafficHandler(trafficService)
	healthHandler := healthInterfaces.NewHealthHandler(healthService)

	router := rest.SetupRouter()
	trafficHandler.RegisterRoutes(router)
	healthHandler.RegisterRoutes(router)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// start server in a goroutine so it doesn't block
	go func() {
		log.Printf("🚀 Server is running on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server failed: %v", err)
		}
	}()

	return server
}

// @desc load port from env vars or use default
func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	if port := os.Getenv("APP_PORT"); port != "" {
		return port
	}
	return "8080"
}

// @desc wait for termination signal and gracefully shutdown
func startGracefulShutdown(server *http.Server) {
	// wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("🔻 Shutting down server...")

	// create a deadline to wait for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Forced to shutdown server: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}
