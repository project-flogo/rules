package common

import (
	"crypto/rand"
	"io"
	"time"

	"github.com/oklog/ulid"
)

//unique id (26 chars)
func GetUniqueId() (string, error) {
	entropy, err := generateRandomBytes(10)
	if err != nil {
		return "", err
	}
	uuid, err := ulid.New(ulid.Timestamp(time.Now().UTC()), entropy)
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}

// generateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomBytes(n int) (io.Reader, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return rand.Reader, nil
}
