package main

import (
	"github.com/mmkamron/library/app/handler"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handler.Create).Methods("POST")
	r.HandleFunc("/", handler.Read).Methods("GET")
	r.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/signup.html")
	}).Methods("GET")
	r.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/signin.html")
	}).Methods("GET")
	r.HandleFunc("/signup", handler.SignUp).Methods("POST")
	r.HandleFunc("/signin", handler.SignIn).Methods("POST")
	r.HandleFunc("/logout", handler.Logout)
	r.HandleFunc("/pkg", handler.Oauth)
	r.HandleFunc("/callback", handler.Callback).Methods("GET")
	r.HandleFunc("/", handler.Update).Methods("PUT")
	r.HandleFunc("/{id}", handler.Delete).Methods("DELETE")
	r.HandleFunc("/{id}", handler.ReadID).Methods("GET")
	http.Handle("/", r)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
