package auth

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const oauthGoogleURLAPI = "https://www.googleapis.com/oauth2/v2/userinfo"
const googleScope = "https://www.googleapis.com/auth/userinfo.email"

type GoogleAuth struct{}

func googleOAuthConfig() (oauth2.Config, error) {
	clientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_OAUTH_REDIRECT_URL")
	if clientID == "" || clientSecret == "" || redirectURL == "" {
		return oauth2.Config{}, errors.New("Google OAuth configuration is incomplete")
	}
	return oauth2.Config{
		RedirectURL:  redirectURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{googleScope},
		Endpoint:     google.Endpoint,
	}, nil
}

func (g GoogleAuth) AuthCallback(c *gin.Context) {
	config, err := googleOAuthConfig()
	if err != nil {
		log.Print(err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	oauthstate, err := c.Request.Cookie("oauthstate")
	if err != nil || oauthstate.Value == "" || c.Query("state") != oauthstate.Value {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	data, err := g.getGoogleUserInfo(c.Request, config, c.Query("code"))
	if err != nil {
		log.Print("Google OAuth callback failed")
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	_, _ = c.Writer.Write(data)
}

func (g GoogleAuth) getGoogleUserInfo(req *http.Request, config oauth2.Config, code string) ([]byte, error) {
	token, err := config.Exchange(req.Context(), code)
	if err != nil {
		return nil, fmt.Errorf("Google OAuth code exchange failed: %w", err)
	}

	resp, err := config.Client(req.Context(), token).Get(oauthGoogleURLAPI)
	if err != nil {
		return nil, fmt.Errorf("Google user info request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Google user info returned status %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func (g GoogleAuth) LoginHandler(c *gin.Context) {
	config, err := googleOAuthConfig()
	if err != nil {
		log.Print(err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}
	state := generateStateOauthCookie(c.Writer)
	c.Redirect(http.StatusTemporaryRedirect, config.AuthCodeURL(state))
}
