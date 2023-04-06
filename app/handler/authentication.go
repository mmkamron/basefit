package handler

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mmkamron/library/pkg"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"log"
	"net/http"
)

var (
	config            = pkg.Load()
	ClientID          = config.ClientID
	ClientSecret      = config.ClientSecret
	randomState       = config.State
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/googlecallback",
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		Scopes:       []string{"openid"},
		Endpoint:     google.Endpoint,
	}
)

func Oauth(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL(randomState)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func Callback(c *gin.Context) {
	if c.Query("state") != randomState {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	token, err := googleOauthConfig.Exchange(c, c.Query("code"))
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	client := googleOauthConfig.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Println(err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response: ", err)
	}
	type User struct {
		ID string `json:"id"`
	}
	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		log.Println("Error decoding response: ", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	if err != nil {
		log.Println(err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	session := sessions.Default(c)
	session.Set("SessionID", user.ID)
	session.Save()

	c.Redirect(http.StatusTemporaryRedirect, "/gym")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func Auth(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("SessionID")
	if userID == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/oauth")
		return
	}
	c.Set("userID", userID)
	c.Next()
}

// import (
//
//	"database/sql"
//	"encoding/json"
//	"fmt"
//	"github.com/google/uuid"
//	"github.com/gorilla/schema"
//	"github.com/mmkamron/library/pkg"
//	"golang.org/x/crypto/bcrypt"
//	"net/http"
//	"time"
//
// )
//var (
//	config                         = pkg.Load()
//	ClientID                       = config.ClientID
//	ClientSecret                   = config.ClientSecret
//	oauth        pkg.Authenticator = pkg.New(ClientID, ClientSecret)
//)
//
//
//func Callback(c *gin.Context) {
//	//db := pkg.ConnectDB()
//	if err := c.Value("code"); err != nil {
//		log.Println(err)
//		c.Redirect(http.StatusTemporaryRedirect, "/oauth")
//		return
//	}
//	token, err := oauth.ObtainToken(fmt.Sprint(c.Value("code")))
//	if err != nil {
//		log.Printf("error obtaining token: %s\n", err)
//		c.Redirect(http.StatusTemporaryRedirect, "/oauth")
//		return
//	}
//	log.Printf("token: %s\n", token)
//}

//)
//type Credentials struct {
//	ID       int    `json:"id"`
//	Name     string `json:"name"`
//	Email    string `json:"email"`
//	Password string `json:"password"`
//}
//
//func Cookie(w http.ResponseWriter, r *http.Request) (string, error) {
//	db := pkg.ConnectDB()
//	c, err := r.Cookie("sessionID")
//	if err != nil {
//		if err == http.ErrNoCookie {
//			http.Redirect(w, r, "/signup", http.StatusFound)
//			return "Redirecting to signup", nil
//		}
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return "", err
//	}
//	sessionID := c.Value
//	var username string
//	_ = db.QueryRow("select username from session where value=$1", sessionID).Scan(&username)
//	return username, nil
//}
//
//func ClearCookie(w http.ResponseWriter) {
//	c := &http.Cookie{
//		Name:   "sessionID",
//		Value:  "",
//		MaxAge: -1,
//	}
//	http.SetCookie(w, c)
//}
//
//func SignUp(w http.ResponseWriter, r *http.Request) {
//	var credentials Credentials
//	var decoder = schema.NewDecoder()
//	if err := r.ParseForm(); err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//	if err := decoder.Decode(&credentials, r.PostForm); err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	hash, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	db := pkg.ConnectDB()
//	_, err = db.Exec("INSERT INTO users VALUES ($1, $2)", credentials.Username, string(hash))
//	if err != nil {
//		http.Redirect(w, r, "/signup", http.StatusFound)
//		return
//	}
//	http.Redirect(w, r, "/signin", http.StatusFound)
//}
//
//func SignIn(w http.ResponseWriter, r *http.Request) {
//	var creds Credentials
//	var decoder = schema.NewDecoder()
//	if err := r.ParseForm(); err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//	if err := decoder.Decode(&creds, r.PostForm); err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	db := pkg.ConnectDB()
//	dbCreds := &Credentials{}
//	if err := db.QueryRow("select password from users where username=$1",
//		creds.Username).Scan(&dbCreds.Password); err != nil {
//		http.Redirect(w, r, "/signin", http.StatusFound)
//		return
//	}
//	if err := bcrypt.CompareHashAndPassword([]byte(dbCreds.Password), []byte(creds.Password)); err != nil {
//		http.Redirect(w, r, "/signin", http.StatusFound)
//		return
//	}
//
//	sessionID := uuid.NewString()
//	expiresAt := time.Now().Add(10 * time.Minute)
//	if err := db.QueryRow("select username from session where username=$1", creds.Username).Scan(&creds.Username); err != nil {
//		if err == sql.ErrNoRows {
//			if _, err := db.Exec("INSERT INTO session VALUES ($1, $2, $3)", sessionID, expiresAt, creds.Username); err != nil {
//				fmt.Fprintf(w, "Could not insert session into database: %s", err)
//				return
//			}
//		} else {
//			fmt.Fprintf(w, "Could not query database: %s", err)
//			return
//		}
//	}
//
//	if _, err := db.Exec("UPDATE session SET value = $1, expiry = $2 WHERE username = $3", sessionID, expiresAt, creds.Username); err != nil {
//		fmt.Fprintf(w, "Could not update session into database: %s", err)
//		return
//	}
//
//	http.SetCookie(w, &http.Cookie{
//		Name:     "sessionID",
//		Value:    sessionID,
//		Expires:  expiresAt,
//		HttpOnly: true,
//	})
//
//	http.Redirect(w, r, "/", http.StatusFound)
//}
//
//func Logout(w http.ResponseWriter, r *http.Request) {
//	c := &http.Cookie{
//		Name:   "sessionID",
//		Value:  "",
//		MaxAge: -1,
//	}
//	http.SetCookie(w, c)
//	http.Redirect(w, r, "/signin", http.StatusFound)
//}
//
//
//	client := &http.Client{}
//	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
//	req.Header = http.Header{
//		"Authorization": {"Bearer " + token},
//		"Accept":        {"application/vnd.pkg+json"},
//	}
//	res, _ := client.Do(req)
//	defer res.Body.Close()
//	var user User
//	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
//		fmt.Println("could not decode response")
//	}
//
//	sessionID := uuid.NewString()
//	expiresAt := time.Now().Add(10 * time.Minute)
//	if err := db.QueryRow("select username from session where username=$1", user.Login).Scan(&user.Login); err != nil {
//		if err == sql.ErrNoRows {
//			if _, err := db.Exec("INSERT INTO session VALUES ($1, $2, $3)", sessionID, expiresAt, user.Login); err != nil {
//				fmt.Fprintf(w, "Could not insert session into database: %s", err)
//				return
//			}
//		} else {
//			fmt.Fprintf(w, "Could not query database: %s", err)
//			return
//		}
//	}
//
//	if _, err := db.Exec("UPDATE session SET value = $1, expiry = $2 WHERE username = $3", sessionID, expiresAt, user.Login); err != nil {
//		fmt.Fprintf(w, "Could not insert session into database: %s", err)
//		return
//	}
//
//	http.SetCookie(w, &http.Cookie{
//		Name:     "sessionID",
//		Value:    sessionID,
//		Expires:  expiresAt,
//		HttpOnly: true,
//	})
//
//	http.Redirect(w, r, "/", http.StatusFound)
//}
//
