package rest

import (
	"github.com/githamo/stubhub-tc/pkg/middleware"
	"github.com/gorilla/mux"
)

// SetupRouter configures the API router
func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	// Apply middleware
	router.Use(middleware.Logging)

	return router
}
