package api

import (
	"aggron/internal/cache"
	"aggron/internal/services"
	"aggron/internal/types"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type AuthController struct {
	AuthService  *services.Auth
	RedisService *cache.Redis
}

/*
Description: Callback endpoint after authentication
GET /auth/callback

Body (Raw JSON):
- code: <string>
- state: <{originalURL, fileIdD, senderDiscordID, receiverDiscordID}>
*/
func (c *AuthController) AuthCallback(ctx *gin.Context) {
	state, exists := ctx.GetQuery("state")
	if !exists {
		ctx.JSON(http.StatusBadRequest, "state is required")
		return
	}

	// Verify auth code
	code, exists := ctx.GetQuery("code")
	if !exists {
		ctx.JSON(http.StatusBadRequest, "authorization code is required")
		return
	}

	rawToken, err := c.AuthService.Config.Exchange(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to exchange token, invalid code")
		return
	}

	token, _, err := new(jwt.Parser).ParseUnverified(rawToken.Extra("id_token").(string), jwt.MapClaims{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to parse tokens")
		return
	}

	// Check if the token is valid and extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		ctx.JSON(http.StatusBadRequest, "failed to extract claims, invalid token")
		return
	}

	email := claims["email"].(string)

	// fetch request state info
	stateInfo, err := cache.GetObjTyped[types.StateInfo](c.RedisService, ctx, state)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "state info has not been set")
		return
	}

	var authInfo types.AuthInfo
	var cacheKey string
	var redirectPath string

	// TODO: create different workflow for authenticating sender
	// if stateInfo.SenderEmail == email {
	// 	authInfo = types.AuthInfo{
	// 		DiscordID: stateInfo.SenderDiscordID,
	// 		Email:     email,
	// 	}

	// 	cacheKey = stateInfo.SenderDiscordID

	// 	redirectPath = "/upload"
	// } else {

	// authorization
	if email != stateInfo.ReceiverEmail || email != stateInfo.SenderEmail {
		ctx.JSON(http.StatusUnauthorized, "user is not permitted to received this resource")
		return
	}

	authInfo = types.AuthInfo{
		DiscordID: stateInfo.ReceiverDiscordID,
		Email:     email,
	}

	cacheKey = stateInfo.ReceiverDiscordID
	redirectPath = fmt.Sprintf("/file?fileID=%s&receiverDiscordID=%s", stateInfo.FileID, stateInfo.ReceiverDiscordID)
	// }

	// set auth state for user who just logged in
	err = cache.SetObjTyped(c.RedisService, ctx, cacheKey, authInfo, types.DefaultExpirationTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to set auth info")
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, redirectPath)
}
