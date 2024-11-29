package models

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uint
	Email        string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}

// Inserts a new entry to the DB (postgres)
func (us UserService) Create(email, password string) (*User, error) {
	email = strings.ToLower(email)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password),
		bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}
	hashString := string(hashedBytes)

	row := us.DB.QueryRow(`
	insert into users (email, password_hash)
	values ($1, $2)
	returning id;
	`, email, hashString)

	var id uint
	err = row.Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("error when inserting into DB: %w", err)
	}

	return &User{
		ID:           id,
		Email:        email,
		PasswordHash: hashString,
	}, nil
}

func (us *UserService) Authenticate(email, password string) (*User, error) {
	email = strings.ToLower(email)
	row := us.DB.QueryRow(`
	select id, password_hash 
	from users 
	where email = $1;
	`, email)

	var id uint
	var passHash string
	err := row.Scan(&id, &passHash)
	if err != nil {
		log.Printf("Authentication: %v", err)
		return nil, fmt.Errorf("error when pulling from DB for authentication")
	}

	// Hashing the password and comparing it to our stored hash.
	err = bcrypt.CompareHashAndPassword([]byte(passHash), []byte(password))
	if err != nil {
		log.Printf("Authentication: %v", err)
		return nil, fmt.Errorf("authentication error: wrong password")
	}
	log.Println("Authentication successful.")

	return &User{
		ID:           id,
		Email:        email,
		PasswordHash: passHash,
	}, nil
}
