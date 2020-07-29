package endpoints

import (
	"gicicm/common"
	"gicicm/logger"
	"gicicm/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"regexp"
	"strings"
	"unicode"
)

// CreateUser is an endpoint for creating a user.
func (ctrl *Controller) CreateUser(c *gin.Context) {

	ctx := c.Request.Context()

	response := make(map[string]interface{})
	request := new(models.User)
	err := c.BindJSON(request)

	if err != nil {
		logger.Log().Error("error while binding request body to user", zap.Error(err))
		response["error"] = common.BadRequestError
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	// input validation.
	if !isEmailValid(request.Email) {
		logger.Log().Info("invalid email", zap.String("email", request.Email))
		response["error"] = common.EmailValidationError
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	if !isPasswordValid(request.Password) {
		logger.Log().Info("invalid password", zap.String("password", request.Password))
		response["error"] = common.PasswordValidationError
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	if strings.Trim(request.Name, " ") == "" {
		logger.Log().Info("empty name", zap.String("name", request.Name))
		response["error"] = common.BadRequestError
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	err = ctrl.userProvider.Create(ctx, request)
	if err != nil {
		if err.Error() == common.AccountAlreadyExistsError {
			logger.Log().Info("empty name", zap.String("request", request.Email), zap.Error(err))
			response["error"] = common.AccountAlreadyExistsError
			c.JSON(http.StatusConflict, response)
			c.Abort()
			return
		}

		logger.Log().Error("error while creating user", zap.Error(err))
		response["error"] = err.Error()
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	response["success"] = "created"
	c.JSON(http.StatusCreated, response)
}

// ListUsers is an endpoint for listing all users.
func (ctrl *Controller) ListUsers(c *gin.Context) {
	ctx := c.Request.Context()
	response := make(map[string]interface{})

	usersList, err := ctrl.userProvider.List(ctx)

	if err != nil {
		logger.Log().Error("error while getting users", zap.Error(err))
		response["error"] = err.Error()
		c.JSON(http.StatusInternalServerError, response)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, usersList)
}

// DeleteUser deletes a user based on the id.
func (ctrl *Controller) DeleteUser(c *gin.Context) {
	ctx := c.Request.Context()
	response := make(map[string]interface{})

	metadata, err := parseContextMetaData(c)
	if err != nil {
		response["error"] = common.InternalServerError
		c.JSON(http.StatusInternalServerError, response)
		c.Abort()
		return
	}

	// send back a 403 if user is not admin.
	if !metadata.IsAdmin {
		response["error"] = common.UnAuthorizedError
		c.JSON(http.StatusForbidden, response)
		c.Abort()
		return
	}

	email := c.Param("email")

	err = ctrl.userProvider.Delete(ctx, email)
	if err != nil {
		if err.Error() == common.AccountNotFoundError {
			response["error"] = common.AccountNotFoundError
			c.JSON(http.StatusNotFound, response)
			c.Abort()
			return
		}
		logger.Log().Error("error while deleting user", zap.String("email", email), zap.Error(err))
		response["error"] = err.Error()
		c.JSON(http.StatusInternalServerError, response)
		c.Abort()
		return
	}

	response["result"] = "Successfully Deleted"
	c.JSON(http.StatusOK, response)
}

// isPasswordValid validates a password.
// should be more than 8 chars, and should have
// 1 number, 1 uppercase and 1 symbol.
// returns false for an invalid password.
func isPasswordValid(password string) bool {
	var hasNumber, hasUpper, hasSymbol bool
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSymbol = true
		}
	}
	return len(password) > 8 && hasNumber && hasUpper && hasSymbol
}

// isPasswordValid checks whether an email is valid.
// fyi best way to validate an email is to send a mail
// with a verification link.
func isEmailValid(email string) bool {
	re := regexp.MustCompile("^[^@]+@[^@]+[.][^@]+$")
	return re.MatchString(email)
}
