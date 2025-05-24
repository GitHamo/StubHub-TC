package domain

import "time"

type Status string

const (
	StatusUp   Status = "UP"
	StatusDown Status = "DOWN"
)

type Component struct {
	Name   string
	Status Status
}

type HealthCheck struct {
	Status     Status
	Components []Component
	Timestamp  time.Time
}

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
