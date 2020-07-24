package providers

import (
	"context"
	"gicicm/models"
	"gicicm/stores"
)

// UserProvider is tbe Repository layer for user related operations.
type UserProvider interface {
	List(ctx context.Context) ([]models.User, error)
	Delete(ctx context.Context, emailID string) error
}

// userProvider is a struct responsible for communicating with
// all the different stores for the user related operations.
type userProvider struct {
	userStore stores.UserRepository
}

// NewUserRepository returns a new instance of the user repository.
func NewUserRepository(userStore stores.UserRepository) UserProvider {
	return &userProvider{
		userStore: userStore,
	}
}

// List lists all the users.
func (up *userProvider) List(ctx context.Context) ([]models.User, error) {
	return nil, nil
}

// Delete deletes a user based on the id.
func (up *userProvider) Delete(ctx context.Context, emailID string) error {
	return nil
}
