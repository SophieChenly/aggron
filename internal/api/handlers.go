package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
POST /upload

Body (form-data)
- file: <file>
- discordID: <string>
*/
func UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "File is required")
		return
	}

	discordID, exists := ctx.GetPostForm("discordID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, "discordID is required")
		return
	}

	// just testing here, can remove later
	fmt.Println(file.Filename)
	fmt.Println(discordID)

	// TODO: Upload Logic (Encrypt file and upload to S3)

	ctx.Status(http.StatusCreated)
}
