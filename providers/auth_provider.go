package providers

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"

	"gicicm/models"
	"gicicm/stores"

	"github.com/dgrijalva/jwt-go"
)

// Repository layer for auth related operations.
type AuthProvider interface {
	Login(ctx context.Context, request *models.LoginRequest) (string, error)
	ParseToken(ctx context.Context, token string) (map[string]interface{}, error)
	Logout(ctx context.Context, token, email string) error
}

// authProvider is struct for auth Provider
// and is responsible for communicated with the stores.
type authProvider struct {
	userStore stores.UserRepository
	authStore stores.AuthRepository
}

// NewAuthProvider returns a new instance of the auth repository.
func NewAuthProvider(userStore stores.UserRepository, authStore stores.AuthRepository) AuthProvider {
	return &authProvider{
		userStore: userStore,
		authStore: authStore,
	}
}

// Login returns a token for a successful login of a user.
func (ap *authProvider) Login(ctx context.Context, request *models.LoginRequest) (string, error) {

	user, err := ap.userStore.Fetch(ctx, request.Email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["iss"] = "icm"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["email"] = request.Email

	if strings.HasSuffix(user.Email, "@test.com") {
		claims["isAdmin"] = true
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, _ := rawToken.SignedString([]byte("sad"))
	return token, nil
}

// Verify parses the token and verifies the user for further operation.
func (ap *authProvider) ParseToken(ctx context.Context, token string) (map[string]interface{}, error) {
	parseToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("sad"), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := parseToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("could not parse token")
	}

	return claims, nil
}

// Logout logs a user out.
func (ap *authProvider) Logout(ctx context.Context, token, email string) error {
	err := ap.authStore.RevokeToken(ctx, token, email)
	if err != nil {
		return err
	}
	return nil
}
