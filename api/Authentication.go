package api

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"github.com/mmkamron/basefit/db"
	"golang.org/x/crypto/bcrypt"
)

var activeSessions = make(map[string]string)

func Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		return
	}
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	db := db.ConnectDB()
	row := db.QueryRow("SELECT username FROM users WHERE username=(?)", username)
	var user string
	if err := row.Scan(&user); err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return
	}
	if user == username {
		log.Println("username is already taken")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return
	}
	stmt, err := db.Prepare("INSERT INTO users(username, password) values(?,?)")
	if err != nil {
		log.Println(err)
		return
	}
	_, err = stmt.Exec(username, hash)
	if err != nil {
		log.Println(err)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		return
	}
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	db := db.ConnectDB()
	var hash string
	row := db.QueryRow("SELECT password FROM users WHERE username=(?)", username)
	if err := row.Scan(&hash); err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return
	} else if err == sql.ErrNoRows {
		log.Println("username or password incorrect")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		log.Println("username or password incorrect")
		log.Println(err)
		return
	}
	sessionID := generateSessionID()
	activeSessions[sessionID] = username
	cookie := http.Cookie{
		Name:     "sid",
		Value:    sessionID,
		Expires:  time.Now().Add(72 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	w.Write([]byte("logged in"))
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("sid")
	if err != nil {
		log.Println(err)
		return
	}

	sessionData, ok := activeSessions[session.Value]
	if !ok {
		log.Println("Unknown cookie")
		return
	}
	log.Printf("Welcome to the dashboard, %s", sessionData)
}

func generateSessionID() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		log.Println(err)
		return ""
	}
	sessionID := base64.URLEncoding.EncodeToString(bytes)
	return sessionID
}
