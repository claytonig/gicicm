// +build integration

package endpoints

import (
	"bytes"
	"errors"
	"fmt"
	"gicicm/adapters/cache"
	"gicicm/adapters/db"
	"gicicm/config"
	"gicicm/providers"
	"gicicm/stores"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var router *gin.Engine

func TestMain(m *testing.M) {
	// wait for 10 seconds for docker containers to come up.
	//fmt.Println("Waiting for docker containers to start...")
	time.Sleep(time.Second * 10)

	config := config.Config{
		Database: config.DbConfig{
			Host:   "localhost",
			Port:   "5432",
			User:   "goicm",
			Pass:   "pass",
			DBName: "icm",
			DBType: "postgres",
		},
		Cache: config.CacheConfig{
			Host: "localhost:6379",
		},
		SigningKey: "secret",
	}

	// Init adapters
	cache := cache.NewCache(&config)
	database := db.NewDatabaseAdapter(&config)

	// Init stores
	userStore := stores.NewUserRepository(database, cache)
	authStore := stores.NewAuthRepository(cache)

	// Init providers
	authProvider := providers.NewAuthProvider(userStore, authStore, &config)
	userProvider := providers.NewUserProvider(userStore)

	// Init controller
	router = NewController(authProvider, userProvider)

	err := createUserHelper()
	if err != nil {
		log.Fatal(err)
	}

	m.Run()
}

func TestController_Login(t *testing.T) {
	tests := []struct {
		name               string
		reqBody            string
		expectedStatusCode int
		expectedMessage    string
	}{
		{
			name: "successful login",
			reqBody: ` {
						"email":"clayton@gmail.com",
						"password":"Hello@123123"
						}`,
			expectedStatusCode: 200,
			expectedMessage:    "",
		},
		{
			name: "invalid credentials",
			reqBody: ` {
						"email":"clayton@gmail.com",
						"password":"hello123d"
						}`,
			expectedStatusCode: 401,
			expectedMessage:    `{"error":"invalid credentials"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := httptest.NewRecorder()
			req, _ := http.NewRequest(
				"POST",
				"/gicicm/auth/login",
				bytes.NewReader([]byte(tt.reqBody)))
			router.ServeHTTP(res, req)

			got, _ := ioutil.ReadAll(res.Body)

			// Todo: better way would be to parse jwt token and check each claim.
			if tt.expectedMessage == "" {
				assert.NotEqual(t, tt.expectedMessage, string(got))
			} else {
				assert.Equal(t, tt.expectedMessage, string(got))
			}

			assert.Equal(t, tt.expectedStatusCode, res.Code)
		})
	}
}

func TestController_ListUsers(t *testing.T) {
	token := loginHelper("clayton@gmail.com", "Hello@123123")
	expectedResponse := `[{"id":"1","email":"delete@me.com","name":"user to be deleted"},{"id":"2","email":"clayton@test.com","name":"superadmin"},{"id":"3","email":"testtwo@mail.com","name":"test user 2"},{"id":"4","email":"test@mail.com","name":"test user 1"},{"id":"5","email":"clayton@gmail.com","name":"clayton gonsalves"}]`

	res := httptest.NewRecorder()

	req, _ := http.NewRequest(
		"GET",
		"/gicicm/users",
		nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	router.ServeHTTP(res, req)
	got, _ := ioutil.ReadAll(res.Body)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, expectedResponse, string(got))
}

func TestController_CreateUser(t *testing.T) {
	tests := []struct {
		name               string
		reqBody            string
		expectedStatusCode int
		expectedMessage    string
	}{
		{
			name: "Successful creation of account",
			reqBody: ` {
						"email":"foo@bar.com",
						"name":"hello world",
						"password":"Hello@123"
						}`,
			expectedStatusCode: 201,
			expectedMessage:    `{"success":"created"}`,
		},
		{
			name: "Invalid/incomplete request body",
			reqBody: ` {
						"email":"clayton@gmail.com",
						"name":" ",
						"password":"hellO@123"
						}`,
			expectedStatusCode: 400,
			expectedMessage:    `{"error":"invalid input format"}`,
		},
		{
			name: "Invalid password",
			reqBody: ` {
						"email":"clayton@gmail.com",
						"name":" ",
						"password":"h2"
						}`,
			expectedStatusCode: 400,
			expectedMessage:    `{"error":"invalid password, should have more than 8 characters, atleast 1 symbol, 1 uppercase character and a number"}`,
		},
		{
			name: "Invalid email",
			reqBody: ` {
						"email":"claytongmailcom",
						"name":" ",
						"password":"Hello@123123"
						}`,
			expectedStatusCode: 400,
			expectedMessage:    `{"error":"invalid email"}`,
		},
		{
			name: "Account already exists",
			reqBody: ` {
						"email":"clayton@gmail.com",
						"name":"clayton gonsalves",
						"password":"Hello@123123"
						}`,
			expectedStatusCode: 409,
			expectedMessage:    `{"error":"account already exists"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := httptest.NewRecorder()
			req, _ := http.NewRequest(
				"POST",
				"/gicicm/auth/signup",
				bytes.NewReader([]byte(tt.reqBody)))
			router.ServeHTTP(res, req)

			got, _ := ioutil.ReadAll(res.Body)

			assert.Equal(t, tt.expectedStatusCode, res.Code)
			assert.Equal(t, tt.expectedMessage, string(got))
		})
	}
}

func TestController_DeleteUser(t *testing.T) {
	tests := []struct {
		name               string
		email              string
		password           string
		reqParam           string
		expectedStatusCode int
		expectedMessage    string
	}{
		{
			name:               "admin user, resource exists",
			reqParam:           "delete@me.com",
			email:              "clayton@test.com",
			password:           "hello123",
			expectedStatusCode: 200,
			expectedMessage:    `{"result":"Successfully Deleted"}`,
		},
		{
			name:               "admin user, resource does not exist",
			reqParam:           "delete@me.com",
			email:              "clayton@test.com",
			password:           "hello123",
			expectedStatusCode: 404,
			expectedMessage:    `{"error":"account does not exist"}`,
		},
		{
			name:               "admin user, resource exists",
			reqParam:           "delete@me.com",
			email:              "clayton@gmail.com",
			password:           "Hello@123123",
			expectedStatusCode: 403,
			expectedMessage:    `{"error":"not permitted to perform this operation"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := loginHelper(tt.email, tt.password)
			res := httptest.NewRecorder()
			req, _ := http.NewRequest(
				"DELETE",
				fmt.Sprintf("/gicicm/users/%s", tt.reqParam),
				nil)

			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

			router.ServeHTTP(res, req)

			got, _ := ioutil.ReadAll(res.Body)

			assert.Equal(t, tt.expectedStatusCode, res.Code)
			assert.Equal(t, tt.expectedMessage, string(got))
		})
	}
}

func TestController_Logout(t *testing.T) {
	email := "clayton@gmail.com"
	password := "Hello@123123"

	token := loginHelper(email, password)

	res := httptest.NewRecorder()

	req, _ := http.NewRequest(
		"POST",
		"/gicicm/auth/logout",
		nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	getRes := httptest.NewRecorder()

	getReq, _ := http.NewRequest(
		"POST",
		"/gicicm/auth/logout",
		nil)
	getReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(getRes, getReq)

	assert.Equal(t, http.StatusUnauthorized, getRes.Code)
}

func loginHelper(email, password string) string {
	reqBody := fmt.Sprintf(
		`{
			"email":"%s",
			"password":"%s"
		}`,
		email, password)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"POST",
		"/gicicm/auth/login",
		bytes.NewReader([]byte(reqBody)))
	router.ServeHTTP(res, req)

	got, _ := ioutil.ReadAll(res.Body)

	return strings.Trim(string(got), "\"")
}

func createUserHelper() error {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"POST",
		"/gicicm/auth/signup",
		bytes.NewReader([]byte(`{
						"email":"clayton@gmail.com",
						"name":"clayton gonsalves",
						"password":"Hello@123123"
						}`)))
	router.ServeHTTP(res, req)

	if res.Code == http.StatusCreated {
		return nil
	}
	return errors.New("could not create user!")
}
