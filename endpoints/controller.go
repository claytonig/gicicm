package endpoints

import (
	"gicicm/providers"
	"github.com/gin-gonic/gin"
)

// Controller is responsible for routing the code flow
// to different providers based on the routes.
type Controller struct {
	authProvider providers.AuthProvider
	userProvider providers.UserProvider
}

// NewController returns a new instance of the controller.
func NewController(
	authProvider providers.AuthProvider,
	userProvider providers.UserProvider) *gin.Engine {

	controller := &Controller{
		authProvider: authProvider,
		userProvider: userProvider,
	}

	// new router
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// root path
	gicicmRoot := router.Group("/gicicm")

	// Unauthenticated endpoints
	gicicmRoot.POST("auth/signup", controller.CreateUser)

	// auth
	gicicmRoot.POST("auth/login", controller.Login)

	// auth middleware
	// all endpoint below this are authenticated.
	gicicmRoot.Use(controller.Verify)

	gicicmRoot.POST("auth/logout", controller.Logout)

	// users
	gicicmRoot.GET("/users", controller.ListUsers)
	gicicmRoot.DELETE("/users/:email", controller.DeleteUser)

	return router
}
