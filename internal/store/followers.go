package store

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
)

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowerStore struct {
	db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, followerID int64, userId int64) error {
	if followerID == userId {
		return ErrConflict
	}
	query := `INSERT INTO followers (user_id, follower_id) VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDelay)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userId, followerID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && (pqErr.Code == "23505" || pqErr.Code == "23503") {
			return ErrConflict
		}
	}
	return err
}

func (s *FollowerStore) Unfollow(ctx context.Context, followerID int64, userId int64) error {
	if followerID == userId {
		return ErrConflict
	}

	query := `DELETE FROM followers WHERE user_id = $1 AND follower_id = $2`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDelay)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userId, followerID)
	return err
}
