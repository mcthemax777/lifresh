package auth

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type SocialAuth interface {

	// getRedirectUrl 소셜 로그인 페이지 주소 전달
	getRedirectUrl(c *gin.Context)

	// LoginCallBack 로그인 완료시 호출할 callback 함수
	LoginCallBack(c *gin.Context)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(1 * 24 * time.Hour)

	b := make([]byte, 16)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := &http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, cookie)
	return state
}
