package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mmkamron/basefit/google"
)

// endpoint
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, fmt.Sprintf(
			"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=https://www.googleapis.com/auth/userinfo.profile",
			os.Getenv("CLIENT_ID"),
			os.Getenv("REDIRECT_URI"),
		), http.StatusTemporaryRedirect)
	})
	http.HandleFunc("/callback", google.Authenticate)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
