package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/mmkamron/library/pkg"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

var (
	config                         = pkg.Load()
	books                          = template.Must(template.ParseFiles("./web/index.html"))
	ClientID                       = config.ClientID
	ClientSecret                   = config.ClientSecret
	oauth        pkg.Authenticator = pkg.New(ClientID, ClientSecret)
	sessions                       = map[string]session{}
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

type session struct {
	username string
	expiry   time.Time
}

type User struct {
	id    int
	Name  string `json:"name"`
	Login string `json:"login"`
}

func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

func Read(w http.ResponseWriter, r *http.Request) {
	var user User
	username, err := Cookie(w, r)
	if err != nil {
		fmt.Printf("Could not get cookie: %s", err)
		return
	}
	db := pkg.ConnectDB()
	rows, err := db.Query("SELECT * FROM books where username = $1", username)
	if err != nil {
		fmt.Fprint(w, "Could not query database")
		return
	}
	defer rows.Close()
	list := []Book{}
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Name, &book.Author, &user.Login); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		list = append(list, book)
	}
	books.Execute(w, list)
}

func Update(w http.ResponseWriter, r *http.Request) {
	//if !Cookie(w, r) {
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}
	var book Book
	var id int
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}
	db := pkg.ConnectDB()
	if err := db.QueryRow("update books set name = $1, author = $2 where id = $3 returning id", book.Name, book.Author, book.ID).Scan(&id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "Book with id %d has been updated!", id)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	//if !Cookie(w, r) {
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}
	param := mux.Vars(r)
	db := pkg.ConnectDB()
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
	//if !Cookie(w, r) {
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}
	param := mux.Vars(r)
	db := pkg.ConnectDB()
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
