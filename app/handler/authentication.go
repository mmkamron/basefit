package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/schema"
	"github.com/mmkamron/library/pkg"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func Cookie(w http.ResponseWriter, r *http.Request) (string, error) {
	db := pkg.ConnectDB()
	c, err := r.Cookie("sessionID")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/signup", http.StatusFound)
			return "Redirecting to signup", err
		}
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
	db := pkg.ConnectDB()
	_, err = db.Exec("INSERT INTO users VALUES ($1, $2)", creds.Username, string(hash))
	if err != nil {
		http.Redirect(w, r, "/signup", http.StatusFound)
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
	db := pkg.ConnectDB()
	dbCreds := &Credentials{}
	if err := db.QueryRow("select password from users where username=$1",
		creds.Username).Scan(&dbCreds.Password); err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(dbCreds.Password), []byte(creds.Password)); err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
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
	db := pkg.ConnectDB()
	if err := r.ParseForm(); err != nil || r.FormValue("code") == "" {
		http.Redirect(w, r, "/pkg", http.StatusTemporaryRedirect)
		return
	}
	token, err := oauth.ObtainToken(r.FormValue("code"))
	if err != nil || token == "" {
		fmt.Printf("error obtaining token: %s\n", err)
		http.Redirect(w, r, "/pkg", http.StatusTemporaryRedirect)
		return
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header = http.Header{
		"Authorization": {"Bearer " + token},
		"Accept":        {"application/vnd.pkg+json"},
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

	http.Redirect(w, r, "/", http.StatusFound)
}

func Create(w http.ResponseWriter, r *http.Request) {
	username, err := Cookie(w, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Errorf("could not get cookie: %s", err)
		return
	}
	db := pkg.ConnectDB()
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
