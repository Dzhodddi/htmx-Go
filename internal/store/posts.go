package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserId    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	Version   int64     `json:"version"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Comments  []Comment `json:"comment"`
	User      User      `json:"user"`
}
type PostsStore struct {
	db *sql.DB
}

type PostWithMetadata struct {
	Post
	CommentsCount int64 `json:"comments_count"`
}

func (s *PostsStore) GetUserFeed(ctx context.Context, userId int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error) {
	query :=
		`
		SELECT
			p.id, p.user_id, p.title, p.content, p.created_at, p.tags,
			u.username,
			COUNT(c.id) AS comments_count
		FROM posts p
		LEFT JOIN COMMENTS c ON c.post_id = p.id
		LEFT JOIN users u ON p.user_id = u.id
		JOIN followers f ON f.follower_id = p.user_id OR p.user_id = $1
		WHERE 
		    f.user_id = $1 AND
			(p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%') AND
			(p.tags @> $5 OR $5 ='{}') AND
			p.created_at >= $6 AND p.created_at <= $7
		GROUP BY p.id, u.username
		ORDER BY p.created_at ` + fq.SortBy + `
		LIMIT $2 OFFSET $3;
`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDelay)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userId, fq.Limit, fq.Offset, fq.Search, pq.Array(fq.Tags), fq.Since, fq.Until)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var feeds []PostWithMetadata
	for rows.Next() {
		var p PostWithMetadata
		err := rows.Scan(
			&p.ID,
			&p.UserId,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			pq.Array(&p.Tags),
			&p.User.Username,
			&p.CommentsCount)
		if err != nil {
			return nil, err
		}
		feeds = append(feeds, p)
	}
	return feeds, nil
}

func (s *PostsStore) Create(ctx context.Context, post *Post) error {
	query := `
	INSERT INTO posts (content, title, user_id, tags)
	VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDelay)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query,
		post.Content, post.Title, post.UserId, pq.Array(post.Tags)).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostsStore) GetByID(ctx context.Context, postId int64) (*Post, error) {
	query := `
	SELECT id, user_id, title, content, tags, created_at, updated_at, version
	FROM posts 
	WHERE ID =  $1;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDelay)
	defer cancel()

	var post Post
	err := s.db.QueryRowContext(ctx, query, postId).Scan(
		&post.ID,
		&post.UserId,
		&post.Title,
		&post.Content,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}

	}
	return &post, nil
}

func (s *PostsStore) Delete(ctx context.Context, postID int64) error {
	query := `DELETE FROM posts WHERE ID =  $1;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDelay)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostsStore) Edit(ctx context.Context, post *Post) error {
	query := `UPDATE posts 
SET title = $1, content = $2, version = version + 1
WHERE ID = $3
WHERE VERSION = $4
RETURNING version;`

	err := s.db.QueryRowContext(ctx,
		query, post.Title, post.Content, post.ID, post.Version).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return nil
		}
	}
	return nil
}
