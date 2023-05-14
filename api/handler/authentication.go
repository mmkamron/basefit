package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mmkamron/basefit/pkg"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	config            = pkg.Load()
	randomState       = config.State
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  config.RedirectUri,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Scopes:       []string{"openid"},
		Endpoint:     google.Endpoint,
	}
)

func Oauth(c *gin.Context) error {
	url := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=email",
		config.ClientID,
		config.RedirectUri,
	)
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
}

func Callback(c *gin.Context) {
	if c.Query("state") != randomState {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	token, err := googleOauthConfig.Exchange(c, c.Query("code"))
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	client := googleOauthConfig.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Println(err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response: ", err)
	}
	type User struct {
		ID string `json:"id"`
	}
	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		log.Println("Error decoding response: ", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	if err != nil {
		log.Println(err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	session := sessions.Default(c)
	session.Set("SessionID", user.ID)
	session.Save()

	c.Redirect(http.StatusTemporaryRedirect, "/gym")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func Auth(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("SessionID")
	if userID == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/unauthorized")
		return
	}
	c.Set("userID", userID)
	c.Next()
}
