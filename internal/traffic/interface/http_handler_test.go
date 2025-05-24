package interfaces_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/githamo/stubhub-tc/internal/traffic/domain"
	interfaces "github.com/githamo/stubhub-tc/internal/traffic/interface"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const servePath = "/serve/"

type MockService struct {
	Response *domain.TrafficResponse
	Err      error
}

func (m *MockService) GetResponseByUUID(uuid string) (*domain.TrafficResponse, error) {
	return m.Response, m.Err
}

func TestTrafficHandlerRequestTraffic(t *testing.T) {
	validUUID := "123e4567-e89b-12d3-a456-426614174000"
	validContent := []byte(`{"status":"ok"}`)

	t.Run("Given valid UUID, When service returns content, Then return 200 and content", func(t *testing.T) {
		service := &MockService{
			Response: &domain.TrafficResponse{Content: validContent},
			Err:      nil,
		}
		handler := interfaces.NewTrafficHandler(service)

		req := httptest.NewRequest(http.MethodGet, servePath+validUUID, nil)
		rec := httptest.NewRecorder()

		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		assert.JSONEq(t, string(validContent), rec.Body.String())
	})

	t.Run("Given invalid UUID, When service returns ErrInvalidUUID, Then return 400", func(t *testing.T) {
		service := &MockService{Err: domain.ErrInvalidUUID}
		handler := interfaces.NewTrafficHandler(service)

		req := httptest.NewRequest(http.MethodGet, "/serve/invalid-uuid", nil)
		rec := httptest.NewRecorder()

		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid UUID format")
	})

	t.Run("Given valid UUID, When service returns ErrTrafficDataNotFound, Then return 404", func(t *testing.T) {
		service := &MockService{Err: domain.ErrTrafficDataNotFound}
		handler := interfaces.NewTrafficHandler(service)

		req := httptest.NewRequest(http.MethodGet, servePath+validUUID, nil)
		rec := httptest.NewRecorder()

		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "Traffic data not found")
	})

	t.Run("Given valid UUID, When service returns unexpected error, Then return 500", func(t *testing.T) {
		service := &MockService{Err: errors.New("db error")}
		handler := interfaces.NewTrafficHandler(service)

		req := httptest.NewRequest(http.MethodGet, servePath+validUUID, nil)
		rec := httptest.NewRecorder()

		router := mux.NewRouter()
		handler.RegisterRoutes(router)
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Internal server error")
	})
}
