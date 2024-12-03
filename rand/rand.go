package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// Reads the content of a random num generator and writes it into the underlying array of the byteslice.
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	nRead, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("error generating numbers: %w", err)
	}
	if nRead < n {
		return nil, fmt.Errorf(
			"tried to generate %d numbers, generatued %d: %w", n, nRead, err)
	}

	return b, nil
}

// Returns a random string using crypto/rand (random byes need to be converted correctly to a string otherwise they won't have ascii characters that will describe them and we'll have problems using this in a cookie).
// n is the number of bytes being used to ogenerate a random string.
func String(n int) (*string, error) {
	b, err := Bytes(n)
	if err != nil {
		return nil, err
	}

	sob := base64.URLEncoding.EncodeToString(b)
	return &sob, nil
}
