package api

import (
	"aggron/internal/cache"
	"aggron/internal/repository"
	"aggron/internal/services"
	"aggron/internal/types"
	"aggron/internal/utils"
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
	UserRepository   *repository.UserRepository
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

	sender, err := c.UserRepository.FindByDiscordID(ctx, senderDiscordID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to get sender email")
		return
	}

	receiver, err := c.UserRepository.FindByDiscordID(ctx, receiverDiscordID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to get receiver email")
		return
	}

	state := types.StateInfo{
		SenderDiscordID:   senderDiscordID,
		SenderEmail:       sender.Email,
		ReceiverDiscordID: receiverDiscordID,
		ReceiverEmail:     receiver.Email,
		FileID:            fileID,
	}

	err = cache.SetObjTyped(c.RedisService, ctx, fileID, state, types.DefaultExpirationTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to set state")
		return
	}

	ctx.String(http.StatusCreated, fileID)
}

/*
Description: Retrieve File
GET /file

Query params:
- fileID: <string>
- receiverDiscordID: <string>

Response:
- <file>
*/
func (c *FileController) RetrieveFile(ctx *gin.Context) {
	fileId, exists := ctx.GetQuery("fileID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, "fileID is required")
		return
	}

	receiverDiscordID := ctx.Query("receiverDiscordID")

	stateInfo, err := cache.GetObjTyped[types.StateInfo](c.RedisService, ctx, fileId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to fetch state")
		return
	}

	if stateInfo.ReceiverDiscordID != receiverDiscordID && receiverDiscordID != "" {
		ctx.JSON(http.StatusUnauthorized, "user is not authorized to receive this resource")
		return
	}

	// check if authenticated, if not then redirect to auth url
	isAuthenticated, err := c.RedisService.Exists(ctx, receiverDiscordID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to verify existence of cached object")
		return
	}

	if !isAuthenticated {
		authorizeURL := c.AuthService.Authorize(fileId, oauth2.AccessTypeOnline)
		ctx.Redirect(http.StatusFound, authorizeURL)
	}

	encryptedData, err := c.S3Service.DownloadFile(ctx, fileId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to download file from storage")
		return
	}

	decryptedData, err := c.EncryptorService.DecryptFile(ctx, fileId, encryptedData, receiverDiscordID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, "failed to decrypt file")
		return
	}

	ftype := utils.GetFileType(decryptedData)
	file := fmt.Sprintf("%s.%s", fileId, ftype)

	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file))
	ctx.Header("Content-Type", "application/octet-stream")

	ctx.Data(http.StatusOK, "application/octet-stream", decryptedData)
}
