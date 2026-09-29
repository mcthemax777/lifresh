package auth

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"lifresh/define"
	"lifresh/internal/core/apperr"
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

func requiredJWTSecret(name string) ([]byte, error) {
	secret := os.Getenv(name)
	if len(secret) < 32 {
		return nil, fmt.Errorf("%s must contain at least 32 bytes", name)
	}
	return []byte(secret), nil
}

func SignAccess(accountID define.SnowflakeID, uid string, ttl time.Duration) (string, *apperr.AppError) {
	secret, err := requiredJWTSecret("JWT_ACCESS_SECRET")
	if err != nil {
		return "", apperr.New(201, "authentication configuration error", err)
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
	at, err := t.SignedString(secret)
	if err != nil {
		return "", apperr.New(201, "access_token invalid", err)
	}

	return at, nil
}

func SignRefresh(accountID define.SnowflakeID, uid string, ttl time.Duration) (string, string, *apperr.AppError) {
	secret, err := requiredJWTSecret("JWT_REFRESH_SECRET")
	if err != nil {
		return "", "", apperr.New(201, "authentication configuration error", err)
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
	signed, err := t.SignedString(secret)
	if err != nil {
		return "", "", apperr.New(201, "refresh_token invalid", err)
	}

	return signed, jti, nil
}

func parseAccess(tokenStr string) (*AccessClaims, error) {
	secret, err := requiredJWTSecret("JWT_ACCESS_SECRET")
	if err != nil {
		return nil, err
	}
	tok, err := jwt.ParseWithClaims(tokenStr, &AccessClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || tok == nil || !tok.Valid {
		return nil, errors.New("invalid access token")
	}
	return tok.Claims.(*AccessClaims), nil
}

func ParseRefresh(tokenStr string) (*RefreshClaims, *apperr.AppError) {
	secret, err := requiredJWTSecret("JWT_REFRESH_SECRET")
	if err != nil {
		return nil, apperr.New(201, "authentication configuration error", err)
	}
	tok, err := jwt.ParseWithClaims(tokenStr, &RefreshClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !tok.Valid {
		return nil, apperr.New(201, "invalid_refresh", err)
	}
	return tok.Claims.(*RefreshClaims), nil
}

/* 미들웨어: /v1/auth/* 스킵 */
func JWTAuthSkipper() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()

		//이미 main 에서 걸러지긴 함
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
			if _, configErr := requiredJWTSecret("JWT_ACCESS_SECRET"); configErr != nil {
				c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"code": 503, "msg": "authentication unavailable"})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "invalid token"})
			return
		}
		c.Set("account_id", cl.AccountID)
		c.Set("uid", cl.UID)
		c.Next()
	}
}
