package models

import (
	"database/sql"
	"fmt"
	"webGo/rand"
)

type Session struct {
	ID     uint
	UserID uint
	// Token is only set when creating a new session. When we look up a session
	// this will be left empty, as we only store the hash of a session
	// token in our DB and we cannot reverse it into a raw token.
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
}

func (ss *SessionService) Create(userID uint) (*Session, error) {
	// TODO Store session in DB.
	token, err := rand.SessionToken()
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	return &Session{
		UserID: userID,
		Token:  *token,
		// TODO set the token's hash.
	}, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	// TODO Implement SessionService.User.
	return nil, nil
}
