package providers

import (
	"context"
	"gicicm/models"
	"gicicm/stores"
)

// Repository layer for auth related operations.
type AuthProvider interface {
	SignUp(ctx context.Context, newUserRequest *models.SignUpRequest) error
	Login(ctx context.Context, request *models.LoginRequest) error
	Verify(ctx context.Context, request *models.LoginRequest) error
}

// authProvider is struct for auth Provider
// and is responsible for communicated with the stores.
type authProvider struct {
	userStore stores.UserRepository
}

// NewAuthRepository returns a new instance of the auth repository.
func NewAuthRepository(userStore stores.UserRepository) AuthProvider {
	return &authProvider{
		userStore: userStore,
	}
}

// SignUp creates a new user in the system.
func (ap *authProvider) SignUp(ctx context.Context, newUserRequest *models.SignUpRequest) error {
	return nil
}

// Login returns a token for a successful login of a user.
func (ap *authProvider) Login(ctx context.Context, request *models.LoginRequest) error {
	return nil
}

// Verify parses the token and verifies the user for further operation.
func (ap *authProvider) Verify(ctx context.Context, request *models.LoginRequest) error {
	return nil
}
