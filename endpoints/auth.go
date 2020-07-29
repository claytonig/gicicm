package endpoints

import (
	"gicicm/common"
	"net/http"
	"strings"

	"gicicm/logger"
	"gicicm/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Login is an endpoint that generates and returns a token
// on successful verification of a user.
func (ctrl *Controller) Login(c *gin.Context) {
	ctx := c.Request.Context()

	response := make(map[string]interface{})
	request := new(models.LoginRequest)

	err := c.BindJSON(request)
	if err != nil {
		logger.Log().Error("error while binding request body to user", zap.Error(err))
		response["error"] = common.BadRequestError
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	token, err := ctrl.authProvider.Login(ctx, request)
	if err != nil {
		if err.Error() == common.InvalidCredentialsError {
			logger.Log().Info("Invalid credentials", zap.Any("request", request), zap.Error(err))
			response["error"] = common.InvalidCredentialsError
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}
		response["error"] = common.InternalServerError
		c.JSON(http.StatusInternalServerError, response)
		c.Abort()
		return
	}

	c.JSON(200, token)
}

// Verify is an endpoint that verifies a user for further operations.
func (ctrl *Controller) Verify(c *gin.Context) {

	ctx := c.Request.Context()

	response := make(map[string]interface{})

	// fetch auth token from headers
	authToken := c.Request.Header.Get("Authorization")
	if authToken == "" {
		logger.Log().Info("Invalid credentials, no token", zap.Any("authToken", authToken))
		response["error"] = "invalid auth token"
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}
	// remove bearer part from header and parse token to get claims
	authToken = strings.Replace(authToken, "Bearer ", "", 1)

	if ctrl.authProvider.IsTokenRevoked(ctx, authToken) {
		logger.Log().Info("Invalid credentials", zap.Any("authToken", authToken))
		response["error"] = "invalid auth token"
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}

	parsedToken, err := ctrl.authProvider.ParseToken(ctx, authToken)
	if err != nil {
		logger.Log().Info("Invalid credentials", zap.Any("authToken", authToken))
		response["error"] = "invalid auth token"
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}

	// set claim in context for later use.
	c.Set("isAdmin", parsedToken["isAdmin"])
	c.Set("email", parsedToken["email"])
	c.Set("token", authToken)
	c.Next()
}

// Logout is an endpoint that logs a user out.
func (ctrl *Controller) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	response := make(map[string]interface{})

	metadata, err := parseContextMetaData(c)
	if err != nil {
		logger.Log().Error("error parsing metadata", zap.Error(err))
		response["error"] = common.InternalServerError
		c.JSON(http.StatusInternalServerError, response)
		c.Abort()
		return
	}

	err = ctrl.authProvider.Logout(ctx, metadata.Token, metadata.Email)
	if err != nil {
		response["error"] = common.InternalServerError
		c.JSON(http.StatusInternalServerError, response)
		c.Abort()
		return
	}
}
