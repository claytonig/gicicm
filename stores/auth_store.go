package stores

import (
	"context"
	"fmt"
	"gicicm/adapters/cache"
	"gicicm/logger"
	"go.uber.org/zap"
	"time"
)

// AuthRepository is a repository layer for all user related operations.
type AuthRepository interface {
	RevokeToken(ctx context.Context, token, email string) error
}

// AuthRepo is responsible for communicating with the data stores via the adapter.
type AuthRepo struct {
	Cache cache.Cache
}

// NewAuthRepository returns a new instance of the user repository.
func NewAuthRepository(cache cache.Cache) AuthRepository {
	return &AuthRepo{
		Cache: cache,
	}
}

func (ar *AuthRepo) RevokeToken(ctx context.Context, token, email string) error {
	key := fmt.Sprintf("token:%s", token)
	_, err := ar.Cache.Set(key, email, time.Hour*24)
	if err != nil {
		logger.Log().Error("error revoking token", zap.String("key", key), zap.String("email", email), zap.Error(err))
	}
	return nil
}
