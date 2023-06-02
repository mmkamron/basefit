package main

import (
	"log"
	"net/http"

	"github.com/mmkamron/basefit/api"
	"github.com/mmkamron/basefit/google"
)

func main() {
	http.HandleFunc("/oauth", google.Redirect)
	http.HandleFunc("/callback", google.Authenticate)

	http.HandleFunc("/register", api.Register)
	http.HandleFunc("/login", api.Login)
	http.HandleFunc("/dashboard", api.Dashboard)
	// Start the server
	log.Fatal(http.ListenAndServe(":8000", nil))
}
