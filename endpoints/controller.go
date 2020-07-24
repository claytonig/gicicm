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
	router *gin.Engine,
	authProvider providers.AuthProvider,
	userProvider providers.UserProvider) {

	controller := &Controller{
		authProvider: authProvider,
		userProvider: userProvider,
	}

	// root path
	gicicmRoot := router.Group("/gicicm")

	// Unauthenticated endpoints
	gicicmRoot.POST("auth/signup", controller.SignUp)

	// auth middleware
	gicicmRoot.Use(controller.Verify)

	// auth
	gicicmRoot.GET("auth/login", controller.Login)

	// users
	gicicmRoot.GET("/users", controller.ListUsers)
	gicicmRoot.DELETE("/users", controller.DeleteUser)
}
