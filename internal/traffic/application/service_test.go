package application_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/githamo/stubhub-tc/internal/traffic/application"
	"github.com/githamo/stubhub-tc/internal/traffic/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type MockRepository struct {
	content json.RawMessage
	err     error
}

func (m *MockRepository) FindByUUID(uuid string) (json.RawMessage, error) {
	return m.content, m.err
}

func TestTrafficServiceGetContentByUUID(t *testing.T) {
	validUUID := uuid.NewString()
	validContent := json.RawMessage(`{"message": "ok"}`)

	t.Run("Given a valid UUID, When data exists, Then return TrafficResponse", func(t *testing.T) {
		repo := &MockRepository{
			content: validContent,
			err:     nil,
		}
		service := application.NewTrafficService(repo)

		result, err := service.GetResponseByUUID(validUUID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, validContent, result.Content)

	})

	t.Run("Given an invalid UUID, When parsed, Then return ErrInvalidUUID", func(t *testing.T) {

		service := application.NewTrafficService(&MockRepository{})

		result, err := service.GetResponseByUUID("invalid-uuid")

		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrInvalidUUID)
	})

	t.Run("Given a valid UUID, When repository returns error, Then return error", func(t *testing.T) {
		repo := &MockRepository{
			err: errors.New("database error"),
		}
		service := application.NewTrafficService(repo)
		actual, err := service.GetResponseByUUID(validUUID)
		assert.Nil(t, actual)
		assert.EqualError(t, err, "database error")
	})
}
