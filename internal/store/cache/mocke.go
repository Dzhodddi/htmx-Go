package cache

import (
	"context"
	"project/internal/store"
)

func NewMockCacheStorage() Storage {
	return Storage{
		Users: MockUserStorage{},
	}
}

type MockUserStorage struct{}

func (m MockUserStorage) Get(ctx context.Context, id int64) (*store.User, error) {
	return &store.User{}, nil
}

func (m MockUserStorage) Set(ctx context.Context, user *store.User) error {
	return nil
}
