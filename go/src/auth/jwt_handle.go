package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"lifresh/define"
	"net/http"
	"os"
	"strings"
	"time"
)

type AccessClaims struct {
	AccountID define.SnowflakeID `json:"account_id"`
	UID       string             `json:"uid"`
	jwt.RegisteredClaims
}
type RefreshClaims struct {
	AccountID define.SnowflakeID `json:"account_id"`
	UID       string             `json:"uid"`
	JTI       string             `json:"jti"`
	jwt.RegisteredClaims
}

func SignAccess(accountID define.SnowflakeID, uid string, ttl time.Duration) (string, error) {
	secret := os.Getenv("JWT_ACCESS_SECRET")
	if secret == "" {
		secret = "dev-access-secret"
	}
	now := time.Now()
	cl := AccessClaims{
		AccountID: accountID,
		UID:       uid,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   uid,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, cl)
	return t.SignedString([]byte(secret))
}

func SignRefresh(accountID define.SnowflakeID, uid string, ttl time.Duration) (string, string, error) {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret == "" {
		secret = "dev-refresh-secret"
	}
	now := time.Now()
	jti := strings.ReplaceAll(time.Now().Format("20060102150405.000000000"), ".", "")
	cl := RefreshClaims{
		AccountID: accountID,
		UID:       uid,
		JTI:       jti,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   uid,
			ID:        jti,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, cl)
	signed, err := t.SignedString([]byte(secret))
	return signed, jti, err
}

func parseAccess(tokenStr string) (*AccessClaims, error) {
	secret := os.Getenv("JWT_ACCESS_SECRET")
	if secret == "" {
		secret = "dev-access-secret"
	}
	tok, err := jwt.ParseWithClaims(tokenStr, &AccessClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !tok.Valid {
		return nil, errors.New("invalid access token")
	}
	return tok.Claims.(*AccessClaims), nil
}

func ParseRefresh(tokenStr string) (*RefreshClaims, error) {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret == "" {
		secret = "dev-refresh-secret"
	}
	tok, err := jwt.ParseWithClaims(tokenStr, &RefreshClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !tok.Valid {
		return nil, errors.New("invalid refresh token")
	}
	return tok.Claims.(*RefreshClaims), nil
}

/* 미들웨어: /v1/auth/* 스킵 */
func JWTAuthSkipper() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		if strings.HasPrefix(path, "/v1/auth/") {
			c.Next()
			return
		}
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "missing bearer token"})
			return
		}
		tokenStr := strings.TrimPrefix(h, "Bearer ")
		cl, err := parseAccess(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "invalid token"})
			return
		}
		c.Set("account_id", cl.AccountID)
		c.Set("uid", cl.UID)
		c.Next()
	}
}
