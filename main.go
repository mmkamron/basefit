package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Book struct {
	ID     int    `json:"id" schema:"-"`
	Name   string `json:"name"`
	Author string `json:"author"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", Create).Methods("POST")
	r.HandleFunc("/", Read).Methods("GET")
	r.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./html/signup.html")
	}).Methods("GET")
	r.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./html/signin.html")
	}).Methods("GET")
	r.HandleFunc("/signup", SignUp).Methods("POST")
	r.HandleFunc("/signin", SignIn).Methods("POST")
	r.HandleFunc("/github", Oauth)
	r.HandleFunc("/callback", Callback).Methods("GET")
	r.HandleFunc("/", Update).Methods("PUT")
	r.HandleFunc("/{id}", Delete).Methods("DELETE")
	r.HandleFunc("/{id}", ReadID).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
