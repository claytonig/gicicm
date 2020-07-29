// +build !integration

package stores

import (
	"context"
	"encoding/json"
	"fmt"
	cacheMock "gicicm/adapters/cache/mocks"
	"gicicm/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestUserStore_Fetch_CacheMiss(t *testing.T) {
	emailID := "test@test.com"
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected while setting up the mock db", err)
	}
	mockCache := new(cacheMock.Cache)
	key := fmt.Sprintf("user:%s", emailID)
	mockCache.On("Get", key).Return("", nil)
	mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything).Return("", nil)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "email", "password"}).
		AddRow(1, "testuser", "test@test.com", "asda9sdu9as8hda9sgca86sdtfas68")

	query := fmt.Sprintf(fetchUserQuery, emailID)

	mockSQL.ExpectQuery(query).WillReturnRows(rows)

	c := NewUserRepository(db, mockCache)

	user, err := c.Fetch(context.TODO(), emailID)

	assert.NoError(t, err)
	assert.Equal(t, emailID, user.Email)
	mockCache.AssertExpectations(t)
	_ = mockSQL.ExpectationsWereMet()

}

func TestUserStore_Fetch_CacheHit(t *testing.T) {
	mockCache := new(cacheMock.Cache)
	mockUser := &models.User{
		ID:    "1",
		Email: "test@test.com",
		Name:  "testUser",
	}

	b, err := json.Marshal(mockUser)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when Marshalling", err)
	}

	key := fmt.Sprintf("user:%s", mockUser.Email)
	mockCache.On("Get", key).Return(string(b), nil)

	userRepo := NewUserRepository(nil, mockCache)
	user, err := userRepo.Fetch(context.TODO(), mockUser.Email)

	assert.NoError(t, err)
	assert.ObjectsAreEqualValues(mockUser, user)
	mockCache.AssertExpectations(t)
}

func TestUserStore_CreateUser(t *testing.T) {
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected while setting up the mock db", err)
	}

	mockUser := &models.User{
		ID:       "1",
		Email:    "test@test.com",
		Name:     "testUser",
		Password: "asdasd",
	}
	funcGenerate = func(pass string) ([]byte, error) {
		return []byte("asdsad"), nil
	}

	defer db.Close()
	defer func() {
		funcGenerate = generateHash
	}()

	query := fmt.Sprintf("['INSERT INTO users(name,email,password) VALUES('%s','%s','%s')']", mockUser.Name, mockUser.Email, "asdasd")

	mockSQL.ExpectBegin()
	mockSQL.ExpectPrepare(query).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(nil)
	mockSQL.ExpectCommit()

	userRepo := NewUserRepository(db, nil)

	err = userRepo.Create(context.TODO(), mockUser)
	assert.NoError(t, err)
}

func TestUserStore_ListUsers(t *testing.T) {
	db, mockSQL, _ := sqlmock.New()

	mockUsers := []*models.User{
		{
			ID:       "1",
			Email:    "test@test.com",
			Name:     "testUser",
			Password: "asdasd",
		},
		{
			ID:       "2",
			Email:    "test1@test.com",
			Name:     "test1User",
			Password: "asdasd",
		},
		{
			ID:       "3",
			Email:    "test3@test.com",
			Name:     "test3User",
			Password: "asdasd",
		},
	}

	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow("1", "testUser", "test@test.com").
		AddRow("2", "test1User", "test1@test.com").
		AddRow("3", "test3User", "test3@test.com")

	query := listUsersQuery

	mockSQL.ExpectQuery(query).WillReturnRows(rows)

	userRepo := NewUserRepository(db, nil)
	users, err := userRepo.List(context.TODO())

	assert.NoError(t, err)
	assert.ObjectsAreEqualValues(len(mockUsers), len(users))
	assert.ObjectsAreEqualValues(mockUsers, users)
}

func TestUserStore_DeleteUser(t *testing.T) {
	db, mockSQL, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected while setting up the mock db", err)
	}
	mockCache := new(cacheMock.Cache)

	mockUser := &models.User{
		ID:       "1",
		Email:    "test@test.com",
		Name:     "testUser",
		Password: "asdasd",
	}

	key := fmt.Sprintf("user:%s", mockUser.Email)

	mockCache.On("Del", key).Return(nil)

	defer db.Close()

	query := fmt.Sprintf(deleteUserQuery, mockUser.Email)

	mockSQL.ExpectBegin()
	mockSQL.ExpectPrepare(query).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(nil)
	mockSQL.ExpectCommit()

	userRepo := NewUserRepository(db, mockCache)
	err = userRepo.Delete(context.TODO(), mockUser.Email)

	assert.NoError(t, err)
	mockCache.AssertExpectations(t)
}
