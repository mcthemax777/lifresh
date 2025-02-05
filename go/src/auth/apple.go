package auth

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type AppleAuth struct {
	googleOauthConfig oauth2.Config
	oauthUrl          string
}

func (g AppleAuth) init() {
}

func (g AppleAuth) getRedirectUrl(c *gin.Context) string {
	return ""
}

func (g AppleAuth) LoginCallBack(c *gin.Context) {

}

func (g AppleAuth) loginHandler(c *gin.Context) {

}