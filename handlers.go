package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"golang.org/x/crypto/bcrypt"
)

var books = template.Must(template.ParseFiles("./html/index.html"))

var sessions = map[string]session{}

type session struct {
	username string
	expiry   time.Time
}

func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	var decoder = schema.NewDecoder()
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := decoder.Decode(&creds, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	db := ConnectDB()
	_, err = db.Exec("INSERT INTO users VALUES ($1, $2)", creds.Username, string(hash))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/signin", http.StatusFound)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	var decoder = schema.NewDecoder()
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := decoder.Decode(&creds, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	db := ConnectDB()
	dbCreds := &Credentials{}
	if err := db.QueryRow("select password from users where username=$1",
		creds.Username).Scan(&dbCreds.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(dbCreds.Password), []byte(creds.Password)); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(10 * time.Minute)

	sessions[sessionToken] = session{
		username: creds.Username,
		expiry:   expiresAt,
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiresAt,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/book", http.StatusFound)
}

func Create(w http.ResponseWriter, r *http.Request) {
	if !hasCookie(w, r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var decoder = schema.NewDecoder()
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var book Book
	if err := decoder.Decode(&book, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	db := ConnectDB()
	if _, err := db.Exec("INSERT INTO books(name, author) VALUES($1, $2)", book.Name, book.Author); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/book", http.StatusFound)
}

func Read(w http.ResponseWriter, r *http.Request) {
	if !hasCookie(w, r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	db := ConnectDB()
	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()
	list := []Book{}
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Name, &book.Author); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		list = append(list, book)
	}
	books.Execute(w, list)
}

func Update(w http.ResponseWriter, r *http.Request) {
	if !hasCookie(w, r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var book Book
	var id int
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewDecoder(r.Body).Decode(&book)
	db := ConnectDB()
	if err := db.QueryRow("update books set name = $1, author = $2 where id = $3 returning id", book.Name, book.Author, book.ID).Scan(&id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "Book with id %d has been updated!", id)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	if !hasCookie(w, r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	param := mux.Vars(r)
	db := ConnectDB()
	if _, err := db.Exec("DELETE FROM books WHERE id=$1", param["id"]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Resets id auto-increment.
	if _, err := db.Exec("SELECT SETVAL('books_id_seq',(SELECT MAX(id) FROM books))"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "Book with id %s has been deleted", param["id"])
}

func ReadID(w http.ResponseWriter, r *http.Request) {
	if !hasCookie(w, r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	param := mux.Vars(r)
	db := ConnectDB()
	var book Book
	if err := db.QueryRow("SELECT * FROM books WHERE id=$1", param["id"]).Scan(&book.ID, &book.Name, &book.Author); err != nil {
		if err == sql.ErrNoRows {
			fmt.Fprintf(w, "Book with id %s not found\n", param["id"])
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	fmt.Fprintf(w, "Name: %s\nauthor: %s", book.Name, book.Author)
}

func hasCookie(w http.ResponseWriter, r *http.Request) bool {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "You are not authorized, please sign in to view the content")
			return false
		}
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	sessionToken := c.Value
	userSession, ok := sessions[sessionToken]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "You are not authorized, please sign in to view the content")
		return false
	}
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "You are not authorized, please sign in to view the content")
		return false
	}
	return true
}
