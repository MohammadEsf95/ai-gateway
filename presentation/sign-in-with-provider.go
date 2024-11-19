package presentation

import (
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func SignInWithProvider(c *gin.Context) {

	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(c.Writer, c.Request)
}
