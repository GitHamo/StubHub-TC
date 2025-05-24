package encryption_test

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/githamo/stubhub-tc/internal/common/encryption"
	"github.com/stretchr/testify/assert"
)

func givenEnv(t *testing.T, key, value string) func() {
	original := os.Getenv(key)
	os.Setenv(key, value)
	return func() { os.Setenv(key, original) }
}

func TestHelperHash(t *testing.T) {
	t.Run("returns consistent hash for same input", func(t *testing.T) {
		key := "mykey"
		reset := givenEnv(t, "APP_SECRET", base64.StdEncoding.EncodeToString([]byte(key)))
		defer reset()

		helper := encryption.NewHelper()
		hash1 := helper.Hash("value")
		hash2 := helper.Hash("value")

		assert.Equal(t, hash1, hash2)
		assert.Equal(t, "2e8fe004c91ba6f64e3c831f9f30cd7365741f1a934e5a409ba6e372c80ef536", hash1)
	})

	t.Run("returns different hashes for different input", func(t *testing.T) {
		key := "anotherkey"
		reset := givenEnv(t, "APP_SECRET", base64.StdEncoding.EncodeToString([]byte(key)))
		defer reset()

		helper := encryption.NewHelper()
		assert.NotEqual(t, helper.Hash("foo"), helper.Hash("bar"))
	})

	t.Run("returns valid hash with invalid base64", func(t *testing.T) {
		reset := givenEnv(t, "APP_SECRET", "base64:!!!invalid!!!")
		defer reset()

		helper := encryption.NewHelper()
		hash := helper.Hash("value")

		assert.Len(t, hash, 64)
	})

	t.Run("returns valid hash with empty secret", func(t *testing.T) {
		reset := givenEnv(t, "APP_SECRET", "")
		defer reset()

		helper := encryption.NewHelper()
		hash := helper.Hash("value")

		assert.Len(t, hash, 64)
	})
}
