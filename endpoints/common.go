package endpoints

import (
	"errors"
	"gicicm/models"
	"github.com/gin-gonic/gin"
)

func parseContextMetaData(c *gin.Context) (*models.RequestMetaData, error) {
	metadata := new(models.RequestMetaData)

	isAdmin, ok := c.Keys["isAdmin"].(bool)
	if !ok {
		return nil, errors.New("cannot get metadata from request")
	} else {
		metadata.IsAdmin = isAdmin
	}

	email, ok := c.Keys["email"].(string)
	if !ok {
		return nil, errors.New("cannot get metadata from request")
	} else {
		metadata.Email = email
	}

	token, ok := c.Keys["token"].(string)
	if !ok {
		return nil, errors.New("cannot get metadata from request")
	} else {
		metadata.Token = token
	}

	return metadata, nil
}
