package domain

import "time"

// Status represents the health status of a component
type Status string

const (
	StatusUp   Status = "UP"
	StatusDown Status = "DOWN"
)

// Component represents a system component to be health checked
type Component struct {
	Name   string
	Status Status
}

// HealthCheck represents the health check result
type HealthCheck struct {
	Status     Status
	Components []Component
	Timestamp  time.Time
}

// NewHealthCheck creates a new health check result
func NewHealthCheck(components []Component) *HealthCheck {
	status := StatusUp

	for _, component := range components {
		// If any component is down, the overall status is down
		if component.Status == StatusDown {
			status = StatusDown
			break
		}
	}

	return &HealthCheck{
		Status:     status,
		Components: components,
		Timestamp:  time.Now(),
	}
}

type Service interface {
	Check() *HealthCheck
}
