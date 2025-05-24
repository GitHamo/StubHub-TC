package domain_test

import (
	"testing"
	"time"

	"github.com/githamo/stubhub-tc/internal/health/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewHealthCheck(t *testing.T) {
	t.Run("Given all components are UP, When NewHealthCheck is called, Then overall status is UP", func(t *testing.T) {
		components := []domain.Component{
			{Name: "db", Status: domain.StatusUp},
			{Name: "cache", Status: domain.StatusUp},
		}

		health := domain.NewHealthCheck(components)

		assert.Equal(t, domain.StatusUp, health.Status)
		assert.Equal(t, components, health.Components)
		assert.WithinDuration(t, time.Now(), health.Timestamp, time.Second)
	})

	t.Run("Given one component is DOWN, When NewHealthCheck is called, Then overall status is DOWN", func(t *testing.T) {
		components := []domain.Component{
			{Name: "db", Status: domain.StatusUp},
			{Name: "api", Status: domain.StatusDown},
		}

		health := domain.NewHealthCheck(components)

		assert.Equal(t, domain.StatusDown, health.Status)
		assert.Equal(t, components, health.Components)
		assert.WithinDuration(t, time.Now(), health.Timestamp, time.Second)
	})

	t.Run("Given no components, When NewHealthCheck is called, Then status is UP and components are empty", func(t *testing.T) {
		health := domain.NewHealthCheck([]domain.Component{})

		assert.Equal(t, domain.StatusUp, health.Status)
		assert.Empty(t, health.Components)
		assert.WithinDuration(t, time.Now(), health.Timestamp, time.Second)
	})
}
