package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
func UploadFile(ctx *gin.Context) {
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

Response:
- Downloads file
*/
func RetrieveFile(ctx *gin.Context) {
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

	receiverDiscordID := ctx.Query("receiverDiscordID")

	fmt.Println(fileId)
	fmt.Println(senderDiscordID)
	fmt.Println(receiverDiscordID)

	// TODO: Retrieve Logic (Decrypt file from S3 and check if receiver is authorized to see it)

	ctx.Status(http.StatusCreated)
}

/*
Description: Callback endpoint after authentication
POST /callback/auth

Body (Raw JSON):
- code: <string>
- token: <string>
- state: <{originalURL, fileIdD, senderDiscordID, receiverDiscordID}>
*/
func CallbackAuth(ctx *gin.Context) {
	// TODO: implement authentication verification + set session + redirect to GET /retrieve

	ctx.Redirect(http.StatusPermanentRedirect, "/")
}
