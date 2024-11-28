package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	switch a := os.Args; a[1] {
	case "hash":
		Must(hash(a[2]))
	case "compare":
		Must(nil, compareHashes(a[2], a[3]))
	default:
		Must(nil, fmt.Errorf("invalid command: %s %s", a[1], a[2]))
	}
}

func hash(password string) (*[]byte, error) {
	saltedHashed, err := bcrypt.GenerateFromPassword([]byte(password),
		bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error when hashing: %w", err)
	}

	fmt.Printf("Salted hash value is %v", hex.EncodeToString(saltedHashed))
	return &saltedHashed, nil
}

func compareHashes(saltedHashed, password string) error {
	err := bcrypt.CompareHashAndPassword(
		[]byte(saltedHashed), []byte(password),
	)
	if err != nil {
		return fmt.Errorf("password invalid: %w", err)
	}

	fmt.Println("Correct password")
	return nil
}

func Must(value any, err error) any {
	if err != nil {
		panic(err)
	}

	return value
}
