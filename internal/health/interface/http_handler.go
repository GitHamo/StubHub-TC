package interfaces

import (
	"encoding/json"
	"net/http"

	"github.com/githamo/stubhub-tc/internal/health/domain"
	"github.com/gorilla/mux"
)

type HealthHandler struct {
	service domain.Service
}

func NewHealthHandler(service domain.Service) *HealthHandler {
	return &HealthHandler{
		service: service,
	}
}

func (h *HealthHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/health", h.HealthCheck).Methods("GET")
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	health := h.service.Check()

	// Return health check results as JSON
	w.Header().Set("Content-Type", "application/json")

	// Set HTTP status code based on health status
	if health.Status == domain.StatusDown {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     health.Status,
		"components": health.Components,
		"timestamp":  health.Timestamp,
	})
}
