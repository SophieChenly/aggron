package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	// Add depndencies as needed here
}

/*
Description: Callback endpoint after authentication
GET /auth/callback

Body (Raw JSON):
- code: <string>
- state: <{originalURL, fileIdD, senderDiscordID, receiverDiscordID}>
*/
func (c *AuthController) AuthCallback(ctx *gin.Context) {
	// TODO: implement authentication verification + set session + redirect to GET /retrieve

	ctx.Redirect(http.StatusPermanentRedirect, "/")
}
