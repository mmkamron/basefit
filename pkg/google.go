package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Authenticator interface {
	OAuthUrl() string
	ObtainToken(code string) (string, error)
	//Authenticate(token string) (User, error)
}

type Google struct {
	clientId     string
	clientSecret string
	defaultScope string
	//userEndpoint string
	redirectUri string
}

func New(clientId, clientSecret string) Authenticator {
	return Google{
		clientId:     clientId,
		clientSecret: clientSecret,
		defaultScope: "openid",
		//userEndpoint: "https://api.github.com/user",
		redirectUri: "http://localhost:8080/googlecallback",
	}
}

func (g Google) OAuthUrl() string {
	return fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?"+
		"client_id=%s&redirect_uri=%s&response_type=code&scope=%s&access_type=offline", g.clientId, url.QueryEscape(g.redirectUri), g.defaultScope)
}

func (g Google) ObtainToken(code string) (string, error) {
	client := &http.Client{}

	payload := url.Values{}
	payload.Add("code", code)
	payload.Add("client_id", g.clientId)
	payload.Add("client_secret", g.clientSecret)
	payload.Add("redirect_uri", g.redirectUri)
	payload.Add("grant_type", "authorization_code")

	req, err := http.NewRequest(
		"POST",
		"https://oauth2.googleapis.com/token",
		strings.NewReader(payload.Encode()))
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to perform obtain token request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to obtain access token: http %d", resp.StatusCode)
	}

	var response struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode pkg response: %v", err)
	}
	return response.AccessToken, nil
}

//func (g google) Authenticate(token string) (User, error) {
//	user := User{}
//	client := &http.Client{}
//
//	req, _ := http.NewRequest("GET", g.userEndpoint, nil)
//	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
//	resp, err := client.Do(req)
//
//	if err != nil {
//		return user, fmt.Errorf("authentication request failed: %v", err)
//	}
//	defer resp.Body.Close()
//	if resp.StatusCode != http.StatusOK {
//		return user, fmt.Errorf("invalid token %v", token)
//	}
//	err = json.NewDecoder(resp.Body).Decode(&user)
//	if err != nil {
//		return user, fmt.Errorf("failed to decode user data: %v", err)
//	}
//	user.Login = strings.ToLower(user.Login)
//	return user, nil
//}
