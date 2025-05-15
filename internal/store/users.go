package store

import (
	"context"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicateUsername = errors.New("duplicate username")
)

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.text = &password
	p.hash = hash
	return nil
}

type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExpr time.Duration) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		//create user
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}

		// create user invite
		if err := s.createUserInvitation(ctx, tx, token, invitationExpr, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UsersStore) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `SELECT id, username, email, created_at FROM users WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDelay)
	defer cancel()

	user := &User{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return user, nil

}

func (s *UsersStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
	INSERT INTO users (username, password, email)
	VALUES ($1, $2, $3) RETURNING id, created_at;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDelay)
	defer cancel()

	err := tx.QueryRowContext(ctx, query,
		user.Username, user.Password.hash, user.Email).Scan(
		&user.ID,
		&user.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (s *UsersStore) createUserInvitation(ctx context.Context,
	tx *sql.Tx, token string, invitationExpr time.Duration, userID int64) error {
	query := `INSERT INTO user_invetetions (token, user_id, expiry) VALUES ($1, $2, $3);`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDelay)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(invitationExpr))
	if err != nil {
		return err
	}
	return nil
}
