package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type Book struct {
	ID     int    `json:"id" schema:"-"`
	Name   string `json:"name"`
	Author string `json:"author"`
}

var books = template.Must(template.ParseFiles("index.html"))

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/book", Create).Methods("POST")
	r.HandleFunc("/", Read).Methods("GET")
	r.HandleFunc("/book", Update).Methods("PUT")
	r.HandleFunc("/book/{id}", Delete).Methods("DELETE")
	r.HandleFunc("/book/{id}", ReadID).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func Create(w http.ResponseWriter, r *http.Request) {
	var decoder = schema.NewDecoder()
	if err := r.ParseForm(); err != nil {
		fmt.Fprint(w, "Try again")
	}
	var book Book
	if err := decoder.Decode(&book, r.PostForm); err != nil {
		fmt.Fprint(w, "Try again")
	}
	db := ConnectDB()
	if _, err := db.Exec("INSERT INTO books(name, author) VALUES($1, $2)", book.Name, book.Author); err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Read(w http.ResponseWriter, r *http.Request) {
	db := ConnectDB()
	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		fmt.Fprint(w, "couldn't query")
	}
	defer rows.Close()
	list := []Book{}
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Name, &book.Author); err != nil {
			fmt.Fprint(w, "couldn't scan")
		}
		list = append(list, book)
	}
	books.Execute(w, list)
}

func Update(w http.ResponseWriter, r *http.Request) {
	var book Book
	var id int
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewDecoder(r.Body).Decode(&book)
	db := ConnectDB()
	if err := db.QueryRow("update books set name = $1, author = $2 where id = $3 returning id", book.Name, book.Author, book.ID).Scan(&id); err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Book with id %d has been updated!", id)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	db := ConnectDB()
	if _, err := db.Exec("DELETE FROM books WHERE id=$1", param["id"]); err != nil {
		log.Fatal(err)
	}

	// To reset id auto-increment.
	if _, err := db.Exec("SELECT SETVAL('books_id_seq',(SELECT MAX(id) FROM books))"); err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Book with id %s has been deleted", param["id"])
}

func ReadID(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	db := ConnectDB()
	var book Book
	if err := db.QueryRow("SELECT * FROM books WHERE id=$1", param["id"]).Scan(&book.ID, &book.Name, &book.Author); err != nil {
		if err == sql.ErrNoRows {
			fmt.Fprintf(w, "Book with id %s not found\n", param["id"])
		} else {
			log.Fatal(err)
		}
	}
	fmt.Fprintf(w, "Name: %s\nauthor: %s", book.Name, book.Author)
}
