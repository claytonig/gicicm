package endpoints

import (
	"gicicm/logger"
	"gicicm/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// CreateUser is an endpoint for creating a user.
func (ctrl *Controller) CreateUser(c *gin.Context) {

	ctx := c.Request.Context()

	response := make(map[string]interface{})
	request := new(models.User)
	err := c.BindJSON(request)

	if err != nil {
		logger.Log().Error("error while binding request body to user", zap.Error(err))
		response["error"] = err.Error()
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err = ctrl.userProvider.Create(ctx, request)
	if err != nil {
		logger.Log().Error("error while creating user", zap.Error(err))
		response["error"] = err.Error()
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response["success"] = "created"
	c.JSON(http.StatusOK, response)
	return
}

// ListUsers is an endpoint for listing all users.
func (ctrl *Controller) ListUsers(c *gin.Context) {
	ctx := c.Request.Context()
	response := make(map[string]interface{})

	usersList := make([]models.User, 0)
	var err error

	usersList, err = ctrl.userProvider.List(ctx)

	if err != nil {
		logger.Log().Error("error while getting users", zap.Error(err))
		response["error"] = err.Error()
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusOK, usersList)
	return
}

// DeleteUser deletes a user based on the id.
func (ctrl *Controller) DeleteUser(c *gin.Context) {
	ctx := c.Request.Context()
	response := make(map[string]interface{})

	metadata, err := parseContextMetaData(c)
	if err != nil {
		response["error"] = "something went wrong"
		c.JSON(http.StatusInternalServerError, response)
		c.Abort()
		return
	}

	if !metadata.IsAdmin {
		response["error"] = "not permitted to perform this operation"
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
		return
	}

	email := c.Param("email")

	err = ctrl.userProvider.Delete(ctx, email)

	if err != nil {
		logger.Log().Error("error while deleting user", zap.String("email", email), zap.Error(err))
		response["error"] = err.Error()
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response["result"] = "Successfully Deleted"
	c.JSON(http.StatusOK, response)
	return
}
