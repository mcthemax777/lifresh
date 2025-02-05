package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"net/http"
)

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
const redirectUrl = "http://localhost:8000/auth/google/callback"
const clientID = "1040063090576-jb2det5gcol7p62otei8k4mn51qb3ki5.apps.googleusercontent.com"
const clientSecret = "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"
const scope = "https://www.googleapis.com/auth/userinfo.email"

type GoogleAuth struct {
	googleOauthConfig oauth2.Config
	oauthUrl          string
}

func (g GoogleAuth) init() {
	g.googleOauthConfig = oauth2.Config{
		RedirectURL:  redirectUrl,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{scope},
		Endpoint:     google.Endpoint,
	}
	g.oauthUrl = oauthGoogleUrlAPI
}

func (g GoogleAuth) getRedirectUrl(c *gin.Context) string {
	state := generateStateOauthCookie(c.Writer)
	url := g.googleOauthConfig.AuthCodeURL(state)
	return url
}

func (g GoogleAuth) AuthCallback(c *gin.Context) {

	oauthstate, _ := c.Request.Cookie("oauthstate") // 12

	if c.Request.FormValue("state") != oauthstate.Value { // 13
		log.Printf("invalid google oauth state cookie:%s state:%s\n", oauthstate.Value, c.Request.FormValue("state"))
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	data, err := g.getGoogleUserInfo(c.Request.FormValue("code")) // 14
	if err != nil {                                               // 15
		log.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	//time.Sleep(time.Duration(50) * time.Second)
	_, err = fmt.Fprint(c.Writer, string(data))
	if err != nil {
		return
	} // 16
}

func (g GoogleAuth) getGoogleUserInfo(code string) ([]byte, error) { // 17

	token, err := g.googleOauthConfig.Exchange(context.Background(), code) // 18
	if err != nil {                                                        // 19
		return nil, fmt.Errorf("Failed to Exchange %s\n", err.Error())
	}

	resp, err := http.Get(g.oauthUrl + token.AccessToken) // 20
	if err != nil {                                       // 21
		return nil, fmt.Errorf("Failed to Get UserInfo %s\n", err.Error())
	}

	return ioutil.ReadAll(resp.Body) // 23
}

func (g GoogleAuth) LoginHandler(c *gin.Context) {

	state := generateStateOauthCookie(c.Writer)
	url := g.googleOauthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

//func (g GoogleAuth) GoogleForm(c *gin.Context) {
//	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(
//		"<html>"+
//			"\n<head>\n    "+
//			"<title>Go Oauth2.0 Test</title>\n"+
//			"</head>\n"+
//			"<body>\n<p>"+
//			"<a href='./auth/google/login'>Google Login</a>"+
//			"</p>\n"+
//			"</body>\n"+
//			"</html>"))
//}

//
//func main1() {
//
//	r := gin.Default()
//
//	r.GET("/", googleForm)
//	r.GET("/auth/google/login", googleLoginHandler)
//	r.GET("/auth/google/callback", googleAuthCallback)
//
//	r.Run(":8000")
//}
