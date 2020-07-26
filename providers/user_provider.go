package providers

import (
	"context"

	"gicicm/models"
	"gicicm/stores"
)

// UserProvider is tbe Repository layer for user related operations.
type UserProvider interface {
	Create(ctx context.Context, user *models.User) error
	List(ctx context.Context) ([]models.User, error)
	Delete(ctx context.Context, emailID string) error
}

// userProvider is a struct responsible for communicating with
// all the different stores for the user related operations.
type userProvider struct {
	userStore stores.UserRepository
}

// NewUserProvider returns a new instance of the user repository.
func NewUserProvider(userStore stores.UserRepository) UserProvider {
	return &userProvider{
		userStore: userStore,
	}
}

// List lists all the users.
func (up *userProvider) Create(ctx context.Context, user *models.User) error {
	err := up.userStore.Create(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

// List lists all the users.
func (up *userProvider) List(ctx context.Context) ([]models.User, error) {
	users, err := up.userStore.List(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Delete deletes a user based on the id.
func (up *userProvider) Delete(ctx context.Context, emailID string) error {
	err := up.userStore.Delete(ctx, emailID)
	if err != nil {
		return err
	}

	return nil
}
