package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"library/github"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"golang.org/x/crypto/bcrypt"
)

var (
	books                             = template.Must(template.ParseFiles("./html/index.html"))
	ClientID                          = os.Getenv("GITHUB_CLIENT_ID")
	ClientSecret                      = os.Getenv("GITHUB_CLIENT_SECRET")
	oauth        github.Authenticator = github.New(ClientID, ClientSecret)
	sessions                          = map[string]session{}
)

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

func Cookie(w http.ResponseWriter, r *http.Request) (string, error) {
	db := ConnectDB()
	c, err := r.Cookie("sessionID")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Fprint(w, "You are not authorized, please sign up to view the content")
			http.Redirect(w, r, "/signup", http.StatusFound)
		}
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Sprintf("Bad request %s", err), err
	}
	sessionID := c.Value
	var username string
	_ = db.QueryRow("select username from session where value=$1", sessionID).Scan(&username)
	return username, nil
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

	sessionID := uuid.NewString()
	expiresAt := time.Now().Add(10 * time.Minute)
	if err := db.QueryRow("select username from session where username=$1", creds.Username).Scan(&creds.Username); err != nil {
		if err == sql.ErrNoRows {
			if _, err := db.Exec("INSERT INTO session VALUES ($1, $2, $3)", sessionID, expiresAt, creds.Username); err != nil {
				fmt.Fprintf(w, "Could not insert session into database: %s", err)
				return
			}
		} else {
			fmt.Fprintf(w, "Could not query database: %s", err)
			return
		}
	}

	if _, err := db.Exec("UPDATE session SET value = $1, expiry = $2 WHERE username = $3", sessionID, expiresAt, creds.Username); err != nil {
		fmt.Fprintf(w, "Could not update session into database: %s", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sessionID",
		Value:    sessionID,
		Expires:  expiresAt,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func Oauth(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, oauth.OAuthUrl(), http.StatusFound)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	db := ConnectDB()
	if err := r.ParseForm(); err != nil || r.FormValue("code") == "" {
		http.Redirect(w, r, "/github", http.StatusTemporaryRedirect)
		return
	}
	token, err := oauth.ObtainToken(r.FormValue("code"))
	if err != nil || token == "" {
		fmt.Printf("error obtaining token: %s\n", err)
		http.Redirect(w, r, "/github", http.StatusTemporaryRedirect)
		return
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header = http.Header{
		"Authorization": {"Bearer " + token},
		"Accept":        {"application/vnd.github+json"},
	}
	res, _ := client.Do(req)
	defer res.Body.Close()
	var user User
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		fmt.Println("could not decode response")
	}

	sessionID := uuid.NewString()
	expiresAt := time.Now().Add(10 * time.Minute)
	if err := db.QueryRow("select username from session where username=$1", user.Login).Scan(&user.Login); err != nil {
		if err == sql.ErrNoRows {
			if _, err := db.Exec("INSERT INTO session VALUES ($1, $2, $3)", sessionID, expiresAt, user.Login); err != nil {
				fmt.Fprintf(w, "Could not insert session into database: %s", err)
				return
			}
		} else {
			fmt.Fprintf(w, "Could not query database: %s", err)
			return
		}
	}

	if _, err := db.Exec("UPDATE session SET value = $1, expiry = $2 WHERE username = $3", sessionID, expiresAt, user.Login); err != nil {
		fmt.Fprintf(w, "Could not insert session into database: %s", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sessionID",
		Value:    sessionID,
		Expires:  expiresAt,
		HttpOnly: true,
	})

	//TODO: save user info somewhere to retrieve their books later.
	//_, err = db.Exec("INSERT INTO oauthUsers(username, sessionID) VALUES ($1, $2)", user.Login, sessionID)
	//if err != nil {
	//	fmt.Fprintf(w, "Could not insert user into database: %s", err)
	//	return
	//}

	http.Redirect(w, r, "/", http.StatusFound)
}

func Create(w http.ResponseWriter, r *http.Request) {
	username, err := Cookie(w, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Errorf("Could not get cookie: %s", err)
		return
	}
	db := ConnectDB()
	var decoder = schema.NewDecoder()
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var book Book
	if err := decoder.Decode(&book, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//if err := db.QueryRow("SELECT id FROM oauthUsers WHERE username = $1", user.Login).Scan(&user.id); err != nil {
	//	fmt.Fprintf(w, "Could not query database: %s", err)
	//}
	if _, err := db.Exec("INSERT INTO books(name, author, username) VALUES($1, $2, $3)", book.Name, book.Author, username); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func Read(w http.ResponseWriter, r *http.Request) {
	var user User
	username, err := Cookie(w, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Errorf("Could not get cookie: %s", err)
		return
	}
	db := ConnectDB()
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
	db := ConnectDB()
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
	//if !Cookie(w, r) {
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}
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
