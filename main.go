package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Book struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/book", Create).Methods("POST")
	r.HandleFunc("/book", Read).Methods("GET")
	r.HandleFunc("/book", Update).Methods("PUT")
	r.HandleFunc("/book/{id}", Delete).Methods("DELETE")
	r.HandleFunc("/book/{id}", ReadID).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func Create(w http.ResponseWriter, r *http.Request) {
	var book Book
	var id int
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewDecoder(r.Body).Decode(&book)
	db := ConnectDB()
	if err := db.QueryRow("INSERT INTO books(name, author) VALUES($1, $2) RETURNING id", book.Name, book.Author).Scan(&id); err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Book with id %d has been added!", id)
}

func Read(w http.ResponseWriter, r *http.Request) {
	var books []string
	db := ConnectDB()
	rows, err := db.Query("SELECT name, author FROM books")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.Name, &book.Author); err != nil {
			log.Fatal(err)
		}
		books = append(books, book.Name)
	}
	for _, book := range books {
		fmt.Fprintf(w, "%s\n", book)
	}
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
	_, err := db.Exec("DELETE FROM books WHERE id=$1", param["id"])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Book with id %s has been deleted", param["id"])
}

func ReadID(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	db := ConnectDB()
	var id int
	var book, author string
	if err := db.QueryRow("SELECT * FROM books WHERE id=$1", param["id"]).Scan(&id, &book, &author); err != nil {
		if err == sql.ErrNoRows {
			fmt.Fprintf(w, "Book with id %s not found\n", param["id"])
		} else {
			log.Fatal(err)
		}
	}
	fmt.Fprintf(w, "Name: %s\nauthor: %s", book, author)
}
