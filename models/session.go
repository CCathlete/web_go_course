package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
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
	fmt.Printf("models.Create: session is %v\n", session)

	return &session, nil
}

func (ss *SessionService) User(us *UserService, token string) (*User, error) {
	var user User
	// Hashing the sessions token.
	tokenHash := ss.hash(token)
	// Querying for the session with that hash.
	row := ss.DB.QueryRow(`
		select
			users.id,
			users.email,
			users.password_hash
		from
			sessions
			join users on users.id = sessions.user_id
		where
			sessions.token_hash = $1
		;`, tokenHash)
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session.User: invalid session token")
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
	// If user_id is not in sessions we create an entry.
	// If it is in sessions we update the token_hash to the new one.
	row := ss.DB.QueryRow(`
	insert into
		sessions (user_id, token_hash)
	values
		($1, $2) on conflict (user_id) do
	update
	set
		token_hash = $2 returning id
	;`, s.UserID, s.TokenHash)
	if err := row.Scan(&s.ID); err != nil {
		return fmt.Errorf("store token hash: %w", err)
	}

	return nil
}

func (ss *SessionService) Delete(token string) error {
	tokenHash := ss.hash(token)
	_, err := ss.DB.Exec(`
		delete from sessions
		where token_hash = $1;
		`, tokenHash)
	if err != nil {
		return fmt.Errorf("models.Delete (session.go): %w", err)
	}
	return nil
}

// // Older version of storeTokenHash.
// func (ss *SessionService) storeTokenHash(s Session) error {
// 	isNewUser, err := ss.isNewUser(s.UserID, &s.ID)
// 	if err == nil {
// 		if isNewUser { // A new user_id so we create a new entry.
// 			row := ss.DB.QueryRow(`
// 		insert into sessions (user_id, token_hash)
// 		values ($1, $2)
// 		returning id;
// 		`, s.UserID, s.TokenHash)
// 			if err := row.Scan(&s.ID); err != nil {
// 				return fmt.Errorf("store token hash: %w", err)
// 			}
// 		} else { // user_id exists so we update the entry.
// 			_, err := ss.DB.Exec(`
// 		update sessions set token_hash = $1
// 		where id = $2 and user_id = $3;
// 		`, s.TokenHash, s.ID, s.UserID)
// 			if err != nil {
// 				return fmt.Errorf("store token hash: %w", err)
// 			}
// 		}
// 	} else { // There has been some error with querying the DB.
// 		return fmt.Errorf("store token hash: %w", err)
// 	}

// 	return nil
// }

// func (ss *SessionService) isNewUser(userID uint, pSessionID *uint) (bool, error) {
// 	row := ss.DB.QueryRow(`
// 		  select id from sessions
// 			where user_id=$1;
// 		`, userID)
// 	switch err := row.Scan(pSessionID); err {
// 	case sql.ErrNoRows:
// 		log.Printf("User is new and does not have a session.")
// 		return true, nil
// 	case nil:
// 		log.Printf("User already exists and has a session.")
// 		return false, nil
// 	default:
// 		return false, fmt.Errorf("query session: %w", err)
// 	}
// }
