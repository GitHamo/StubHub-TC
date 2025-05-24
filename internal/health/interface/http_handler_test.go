package interfaces_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/githamo/stubhub-tc/internal/health/domain"
	interfaces "github.com/githamo/stubhub-tc/internal/health/interface"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type MockHealthService struct {
	mockResponse *domain.HealthCheck
}

func (m *MockHealthService) Check() *domain.HealthCheck {
	return m.mockResponse
}

func TestHealthHandlerHealthCheck(t *testing.T) {
	t.Run("returns 200 OK with UP status and JSON body", func(t *testing.T) {
		// Given
		timestamp := time.Now()
		service := &MockHealthService{
			mockResponse: &domain.HealthCheck{
				Status: domain.StatusUp,
				Components: []domain.Component{
					{Name: "db", Status: domain.StatusUp},
				},
				Timestamp: timestamp,
			},
		}
		handler := interfaces.NewHealthHandler(service)
		router := mux.NewRouter()
		handler.RegisterRoutes(router)

		req := httptest.NewRequest("GET", "/health", nil)
		resp := httptest.NewRecorder()

		// When
		router.ServeHTTP(resp, req)

		// Then
		assert.Equal(t, http.StatusOK, resp.Code)

		var body map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &body)
		assert.NoError(t, err)
		assert.Equal(t, "UP", body["status"])
		assert.NotEmpty(t, body["components"])
		assert.NotEmpty(t, body["timestamp"])
	})

	t.Run("returns 503 Service Unavailable with DOWN status", func(t *testing.T) {
		service := &MockHealthService{
			mockResponse: &domain.HealthCheck{
				Status:     domain.StatusDown,
				Components: []domain.Component{},
				Timestamp:  time.Now(),
			},
		}
		handler := interfaces.NewHealthHandler(service)
		router := mux.NewRouter()
		handler.RegisterRoutes(router)

		req := httptest.NewRequest("GET", "/health", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusServiceUnavailable, resp.Code)

		var body map[string]interface{}
		_ = json.Unmarshal(resp.Body.Bytes(), &body)
		assert.Equal(t, "DOWN", body["status"])
	})
}
