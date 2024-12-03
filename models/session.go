package models

import "database/sql"

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
	// TODO Implement SessionService.Create.
	return nil, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	// TODO Implement SessionService.User.
	return nil, nil
}
