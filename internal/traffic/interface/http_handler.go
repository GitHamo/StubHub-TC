package interfaces

import (
	"log"
	"net/http"

	"github.com/githamo/stubhub-tc/internal/traffic/application"
	"github.com/githamo/stubhub-tc/internal/traffic/domain"

	"github.com/gorilla/mux"
)

type TrafficHandler struct {
	service *application.TrafficService
}

func NewTrafficHandler(service *application.TrafficService) *TrafficHandler {
	return &TrafficHandler{
		service: service,
	}
}

func (h *TrafficHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/serve/{uuid}", h.GetContentByUUID).Methods("GET")
}

func (h *TrafficHandler) GetContentByUUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	content, err := h.service.GetContentByUUID(uuid)
	if err != nil {
		if err == domain.ErrTrafficDataNotFound {
			http.Error(w, "Traffic data not found", http.StatusNotFound)
			return
		}
		if err == domain.ErrInvalidUUID {
			http.Error(w, "Invalid UUID format", http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Internal server error: %v", err)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(content)
}
