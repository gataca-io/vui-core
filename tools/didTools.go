package tools

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/btcsuite/btcutil/base58"
)

/**
Generates a random string with n characters length
base58
*/
func RandSeq(n int) string {
	b, err := GenerateRandomBytes(n)
	if err != nil {
		return ""
	}
	return base58.Encode(b)
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func GenerateRandomStringBase64(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}
