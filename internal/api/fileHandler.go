package api

import (
	"aggron/internal/cache"
	"aggron/internal/services"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type FileController struct {
	AuthService  *services.Auth
	RedisService *cache.Redis
}

/*
Description: Upload File
POST /file

Body (form-data):
- file: <file>
- senderDiscordID: <string>
- receiverDiscordID: <string>

Response:
- fileID: <string>
*/
func (c *FileController) UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "file is required")
		return
	}

	senderDiscordID, exists := ctx.GetPostForm("senderDiscordID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, "senderDiscordID is required")
		return
	}

	receiverDiscordID := ctx.PostForm("receiverDiscordID")

	fmt.Println(file.Filename)
	fmt.Println(senderDiscordID)
	fmt.Println(receiverDiscordID)

	// TODO: Upload Logic (Encrypt file and upload to S3 and store key + authorized userId)

	// temporary -- remove later
	tempFileId := "file.pdf"

	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", tempFileId))
	ctx.Header("Content-Type", "application/pdf")
	ctx.Status(http.StatusCreated)
}

/*
Description: Retrieve File
GET /file

Query params:
- fileID: <string>
- senderDiscordID: <string>
- receiverDiscordID: <string>
- state: <string>

Response:
- <file>
*/
func (c *FileController) RetrieveFile(ctx *gin.Context) {
	fileId, exists := ctx.GetQuery("fileID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, "fileID is required")
		return
	}

	senderDiscordID, exists := ctx.GetQuery("senderDiscordID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, "discordID is required")
		return
	}

	state, exists := ctx.GetQuery("state")
	if !exists {
		ctx.JSON(http.StatusBadRequest, "state is required")
		return
	}

	receiverDiscordID := ctx.Query("receiverDiscordID")

	// check if authenticated, if not then redirect to auth url
	isAuthenticated, err := c.RedisService.Exists(context.TODO(), receiverDiscordID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to verify existence of cached object")
		return
	}

	if !isAuthenticated {
		authorizeURL := c.AuthService.Authorize(state, oauth2.AccessTypeOnline)
		ctx.Redirect(http.StatusFound, authorizeURL)
	}

	// TODO: Retrieve Logic (Decrypt file from S3 and check if receiver is authorized to see it)

	ctx.Status(http.StatusOK)
}
