package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"webGo/rand"
)

const (
	// Min munber of byter to be used for each session token.
	MinBytesPerToken = 32
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
	// BytesPerToken is used to determine how many bytes to use when generating
	// each session token. If this value is not set or is less than
	// MinBytesPerToken const it will be ignored and MinBytesPerToken
	// will be used.
	BytesPerToken int
}

func (ss *SessionService) Create(userID uint) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}
	session := Session{
		UserID:    userID,
		Token:     *token,
		TokenHash: ss.hash(*token),
	}
	ss.storeTokenHash(session)

	return &session, nil
}

func (ss *SessionService) User(us *UserService, token string) (*User, error) {
	var user User
	// Hashing the sessions token.
	tokenHash := ss.hash(token)
	// Querying for the session with that hash.
	row := ss.DB.QueryRow(`
		  select user_id from sessions
			where token_hash=$1;
		`, tokenHash)
	if err := row.Scan(&user.ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session.User: user doesn't have an opem session")
		}
		return nil, fmt.Errorf("session.User: problem with db query -  %w", err)
	}
	// userID was found in the sessions table so we query the user db (in this case, same db different tables)
	row = us.DB.QueryRow(`
		  select email, password_hash from users
			where user_id=$1;
		`, user.ID)
	if err := row.Scan(&user.Email, &user.PasswordHash); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session.User: user doesn't exist")
		}
		return nil, fmt.Errorf("session.User: problem with db query -  %w", err)
	}

	return &user, nil
}

// Hashing the session token using SHA256.
// We're not using bcrypt since the token is long and random not like a password
// so there's no need for salt which complicates things.
// HMAC is also not needed since we don't really need the use of a secret key
// like we do in digital signatures.
func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}

// Receives an initialised user session, creates an entry in
// the sessions DB, stores its token hash and user_id
// and assigns the session id into the session object.
func (ss *SessionService) storeTokenHash(s Session) error {
	isNewUser, err := ss.isNewUser(s.UserID)
	if err == nil {
		if isNewUser { // A new user_id so we create a new entry.
			row := ss.DB.QueryRow(`
		insert into sessions (user_id, token_hash)
		values ($1, $2)
		returning id;
		`, s.UserID, s.TokenHash)
			if err := row.Scan(&s.ID); err != nil {
				return fmt.Errorf("store token hash: %w", err)
			}
		} else { // user_id exists so we update the entry.
			_, err := ss.DB.Exec(`
		update sessions set token_hash = $1 
		where id = $2 and user_id = $3;
		`, s.TokenHash, s.ID, s.UserID)
			if err != nil {
				return fmt.Errorf("store token hash: %w", err)
			}
		}
	} else { // There has been some error with querying the DB.
		return fmt.Errorf("store token hash: %w", err)
	}

	return nil
}

func (ss *SessionService) isNewUser(userID uint) (bool, error) {
	var sessionID uint
	row := ss.DB.QueryRow(`
		  select id from sessions
			where user_id=$1;
		`, userID)
	switch err := row.Scan(&sessionID); err {
	case sql.ErrNoRows:
		log.Printf("User is new and does not have a session.")
		return true, nil
	case nil:
		log.Printf("User already exists and has a session.")
		return false, nil
	default:
		return false, fmt.Errorf("query session: %w", err)
	}
}
