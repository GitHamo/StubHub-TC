package application

import (
	"database/sql"

	"github.com/githamo/stubhub-tc/internal/health/domain"
)

// HealthService implements the domain.Service interface
type HealthService struct {
	db *sql.DB
}

// NewHealthService creates a new health service
func NewHealthService(db *sql.DB) *HealthService {
	return &HealthService{
		db: db,
	}
}

// Check performs health checks on all system components
func (s *HealthService) Check() *domain.HealthCheck {
	components := []domain.Component{
		s.checkDatabase(),
	}

	return domain.NewHealthCheck(components)
}

// checkDatabase verifies the database connection
func (s *HealthService) checkDatabase() domain.Component {
	component := domain.Component{
		Name:   "database",
		Status: domain.StatusUp,
	}

	// Ping the database to check connection
	if err := s.db.Ping(); err != nil {
		component.Status = domain.StatusDown
	}

	return component
}
