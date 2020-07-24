package stores

import (
	"context"
	"database/sql"
	"gicicm/adapters/cache"
	"gicicm/models"
)

// UserRepository is a repository layer for all user related operations.
type UserRepository interface {
	Create(ctx context.Context) ([]models.User, error)
	List(ctx context.Context) ([]models.User, error)
	Fetch(ctx context.Context, emailID string) (*models.User, error)
	Delete(ctx context.Context, email string) error
}

// userRepo is responsible for communicating with the data stores via the adapter.
type UserRepo struct {
	cache cache.Cache
	db    *sql.DB
}

// NewUserRepository returns a new instance of the user repository.
func NewUserRepository(db *sql.DB, cache cache.Cache) UserRepository {
	return &UserRepo{
		cache: cache,
		db:    db,
	}
}

// Create a new user.
func (ur *UserRepo) Create(ctx context.Context) ([]models.User, error) {
	return nil, nil
}

// Fetch a user based on id.
func (ur *UserRepo) Fetch(ctx context.Context, emailID string) (*models.User, error) {
	return nil, nil
}

// List users.
func (ur *UserRepo) List(ctx context.Context) ([]models.User, error) {
	return nil, nil
}

// Delete user based on id.
func (ur *UserRepo) Delete(ctx context.Context, email string) error {
	return nil
}
