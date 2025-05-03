package api

import (
	"aggron/internal/cache"
	"aggron/internal/services"
	"aggron/internal/utils"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type FileController struct {
	AuthService      *services.Auth
	RedisService     *cache.Redis
	EncryptorService *services.FileEncryptionService
	S3Service        *services.S3
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

	fileID, err := utils.GenerateId()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to generate fileID")
		return
	}

	// Encrypt file
	openedFile, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to open file")
		return
	}
	defer openedFile.Close()

	plaintext, err := io.ReadAll(openedFile)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to parse file to plaintext")
		return
	}

	cryptText, err := c.EncryptorService.EncryptFile(ctx, fileID, plaintext, senderDiscordID, receiverDiscordID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to encrypt file")
		return
	}

	_, err = c.S3Service.UploadFile(ctx, fileID, cryptText, "application/octet-stream")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to upload crypttext")
		return
	}

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
	fmt.Println(fileId)
	fmt.Println(senderDiscordID)

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
