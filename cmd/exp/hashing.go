package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// HMAC hashing.

func HashHMAC(password string) (*[]byte, error) {
	secretKeyForHash := "secret key" // Doesn't have to be a string.

	// Setting up our hashing function.
	h := hmac.New(sha256.New, []byte(secretKeyForHash))
	defer h.Reset()

	// Writing data to our hashing function.
	_, err := h.Write([]byte(password))
	if err != nil {
		return nil, fmt.Errorf("error whriting data to hash function: %w.", err)
	}

	// Get the resulting hash.
	result := h.Sum(nil)

	// The resulting hash is binary so we'll encode it to hex so we can read it
	// as letters and numbers.
	fmt.Println(hex.EncodeToString(result))

	return &result, nil
}
