package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/githamo/stubhub-tc/internal/traffic/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewTrafficResponse(t *testing.T) {
	t.Run("Given valid JSON content, When NewTrafficResponse is called, Then it should return a valid TrafficResponse", func(t *testing.T) {
		validJSON := json.RawMessage(`{"status": "ok"}`)

		response, err := domain.NewTrafficResponse(validJSON)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, validJSON, response.Content)
	})

	t.Run("Given invalid JSON content, When NewTrafficResponse is called, Then it should return ErrInvalidContent", func(t *testing.T) {
		invalidJSON := json.RawMessage(`{invalid-json}`)

		response, err := domain.NewTrafficResponse(invalidJSON)

		assert.Nil(t, response)
		assert.ErrorIs(t, err, domain.ErrInvalidContent)
	})
}
