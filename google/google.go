package google

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func Redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=https://www.googleapis.com/auth/userinfo.profile",
		os.Getenv("CLIENT_ID"),
		os.Getenv("REDIRECT_URI"),
	), http.StatusTemporaryRedirect)
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	if code == "" {
		http.Redirect(w, r, "/oauth", http.StatusTemporaryRedirect)
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST",
		fmt.Sprintf(
			"https://oauth2.googleapis.com/token?"+
				"client_id=%s&client_secret=%s&code=%s&grant_type=authorization_code&redirect_uri=%s",
			os.Getenv("CLIENT_ID"),
			os.Getenv("CLIENT_SECRET"),
			code,
			os.Getenv("REDIRECT_URI"),
		), nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		http.Redirect(w, r, "/oauth", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	var response struct {
		Token string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		http.Redirect(w, r, "/oauth", http.StatusTemporaryRedirect)
		return
	}
	userdata, err := http.Get(
		fmt.Sprintf(
			"https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s",
			response.Token,
		),
	)
	if err != nil {
		http.Redirect(w, r, "/oauth", http.StatusTemporaryRedirect)
		return
	}
	defer userdata.Body.Close()

	var user struct {
		Name string `json:"given_name"`
	}

	if err := json.NewDecoder(userdata.Body).Decode(&user); err != nil {
		http.Redirect(w, r, "/oauth", http.StatusTemporaryRedirect)
		return
	}
	fmt.Println(user.Name)
}
