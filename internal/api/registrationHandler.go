package api

import (
	"aggron/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegistrationController struct {
	UserRepository *repository.UserRepository
}

/*
Description: Register User to DB (get email)
POST /file

Body (form-data):
- discordID: <string>
- email: <string>
*/
func (c *RegistrationController) RegisterUser(ctx *gin.Context) {
	discordId, exists := ctx.GetPostForm("discordID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, "discordID is required")
		return
	}

	email, exists := ctx.GetPostForm("email")
	if !exists {
		ctx.JSON(http.StatusBadRequest, "email is required")
		return
	}

	_, err := c.UserRepository.CreateUser(ctx, discordId, email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to create User")
		return
	}

	ctx.Status(http.StatusAccepted)
}
