package controller

import (
	"auth/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

type OAuthController struct {
	userService service.UserService
}

func NewOAuthController(userService service.UserService) *OAuthController {
	return &OAuthController{userService: userService}
}

func (c *OAuthController) BeginAuth(ctx *gin.Context) {
	provider := ctx.Param("provider")
	q := ctx.Request.URL.Query()
	q.Add("provider", provider)
	ctx.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

func (c *OAuthController) CallbackHandler(ctx *gin.Context) {
	provider := ctx.Param("provider")
	q := ctx.Request.URL.Query()
	q.Add("provider", provider)
	ctx.Request.URL.RawQuery = q.Encode()

	gothUser, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, err := c.userService.UpsertGoogleUser(gothUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
