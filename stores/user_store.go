package stores

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"gicicm/adapters/cache"
	"gicicm/logger"
	"gicicm/models"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository is a repository layer for all user related operations.
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	List(ctx context.Context) ([]models.User, error)
	Fetch(ctx context.Context, emailID string) (*models.User, error)
	Delete(ctx context.Context, email string) error
}

// userRepo is responsible for communicating with the data stores via the adapter.
type UserRepo struct {
	cache cache.Cache
	db    *sql.DB
}

const (
	listUsersQuery  = "SELECT id,name,email from users"
	fetchUserQuery  = "SELECT id,name,email,password from users where email='%s'"
	createUserQuery = "INSERT INTO users(name,email,password) VALUES('%s','%s','%s')"
	deleteUserQuery = "DELETE FROM users WHERE email='%s'"
)

// NewUserRepository returns a new instance of the user repository.
func NewUserRepository(db *sql.DB, cache cache.Cache) UserRepository {
	return &UserRepo{
		cache: cache,
		db:    db,
	}
}

// Create a new user.
func (ur *UserRepo) Create(ctx context.Context, user *models.User) error {

	tx, err := ur.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log().Error("error while starting transaction", zap.Error(err))
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	query := fmt.Sprintf(createUserQuery, user.Name, user.Email, hashedPassword)
	stmt, err := tx.Prepare(query)

	if err != nil {
		logger.Log().Error("error while preparing query", zap.String("query", query), zap.Error(err))
		rollBackErr := tx.Rollback()
		if rollBackErr != nil {
			logger.Log().Error("Error while rolling back transaction", zap.String("query", query), zap.Error(err))
		}
		return err
	}

	_, err = stmt.Exec()

	if err != nil {
		logger.Log().Error("error while executing query", zap.String("query", query), zap.Error(err))
		rollBackErr := tx.Rollback()
		if rollBackErr != nil {
			logger.Log().Error("Error while rolling back transaction", zap.String("query", query), zap.Error(err))
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		logger.Log().Error("Error while committing transaction", zap.String("query", query), zap.Error(err))
		return err
	}

	return nil
}

// Fetch a user based on id.
func (ur *UserRepo) Fetch(ctx context.Context, emailID string) (*models.User, error) {
	user := new(models.User)

	val, err := ur.cache.Get(fmt.Sprintf("user:%s", emailID))

	if err != nil {
		logger.Log().Error("Error while fetching user from cache", zap.String("key", emailID), zap.Error(err))
	} else {
		bytes := []byte(val)
		err = json.Unmarshal(bytes, &user)
		if err != nil {
			logger.Log().Error("error while unmarshalling user", zap.String("key", emailID), zap.Error(err))
			return nil, err
		}
		return user, nil
	}

	query := fmt.Sprintf(fetchUserQuery, emailID)
	rows, err := ur.db.QueryContext(ctx, query)
	if err != nil {
		logger.Log().Error("error while querying user", zap.String("query", query), zap.Error(err))
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			logger.Log().Error("error while scanning rows", zap.String("query", query), zap.Error(err))
			return nil, err
		}
	}

	bytes, err := json.Marshal(user)
	_, err = ur.cache.Set(emailID, string(bytes), time.Duration(0))
	if err != nil {
		logger.Log().Error("error while setting cache", zap.String("key", emailID), zap.Error(err))
	}

	return user, nil
}

// List users.
func (ur *UserRepo) List(ctx context.Context) ([]models.User, error) {

	var response = []models.User{}

	rows, err := ur.db.QueryContext(ctx, listUsersQuery)

	if err != nil {
		logger.Log().Error("error while querying data", zap.String("query", listUsersQuery), zap.Error(err))
		return nil, err
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			logger.Log().Error("error while closing rows", zap.Error(err))
		}
	}()

	for rows.Next() {
		user := new(models.User)
		err = rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			logger.Log().Error("error while scanning row data into user", zap.String("query", listUsersQuery), zap.Error(err))
			return nil, err
		}

		response = append(response, *user)
	}

	return response, nil
}

// Delete user based on id.
func (ur *UserRepo) Delete(ctx context.Context, email string) error {

	tx, err := ur.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log().Error("error while starting transaction", zap.Error(err))
		return err
	}

	query := fmt.Sprintf(deleteUserQuery, email)

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		logger.Log().Error("error while preparing query", zap.String("query", query), zap.Error(err))
		return err
	}

	result, err := stmt.ExecContext(ctx)
	if err != nil {
		logger.Log().Error("error while executing query", zap.String("query", query), zap.Error(err))
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		logger.Log().Error("error while fetching rows", zap.String("query", query), zap.Error(err))
		return err
	}

	err = tx.Commit()
	if err != nil {
		logger.Log().Error("Error while committing transaction", zap.String("query", query), zap.Error(err))
		return err
	}

	logger.Log().Info("successfully deleted user", zap.String("email", email), zap.Int64("rows affected", rows))

	return nil
}
