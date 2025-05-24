package encryption

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"log"
	"os"
	"strings"
)

type Hasher interface {
	Hash(input string) string
}

type Helper struct {
	secretKey []byte
}

func NewHelper() Hasher {
	key := os.Getenv("APP_SECRET")

	key = strings.TrimPrefix(key, "base64:")

	decoded, err := base64.StdEncoding.DecodeString(key)

	if err != nil {
		log.Printf("Failed to decode base64 input: %v", err)

		return &Helper{
			secretKey: []byte(""),
		}
	}

	return &Helper{
		secretKey: []byte(decoded),
	}
}

func (e *Helper) Hash(value string) string {
	h := hmac.New(sha256.New, e.secretKey)
	h.Write([]byte(value))
	return hex.EncodeToString(h.Sum(nil))
}
