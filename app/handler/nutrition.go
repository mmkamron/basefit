package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mmkamron/basefit/pkg"
	"net/http"
	"time"
)

type Nutrition struct {
	Date     string `json:"date_time"`
	Calories int    `json:"calories"`
	Protein  string `json:"protein"`
}

func ReadNutrition(c *gin.Context) {
	userID, _ := c.Get("userID")
	db := pkg.ConnectDB()
	var nutrition Nutrition
	rows, err := db.Query("SELECT date_time, calories, protein FROM nutrition_logs WHERE user_id = $1", userID)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var res []Nutrition
	for rows.Next() {
		if err := rows.Scan(&nutrition.Date, &nutrition.Calories, &nutrition.Protein); err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}
		res = append(res, nutrition)
	}
	c.HTML(http.StatusOK, "nutrition.html", res)
}

func CreateNutrition(c *gin.Context) {
	userID, _ := c.Get("userID")
	db := pkg.ConnectDB()
	date := time.Now().Format("01-02-2006")
	calories := c.PostForm("calories")
	protein := c.PostForm("protein")
	if _, err := db.ExecContext(c, "INSERT INTO nutrition_logs(date_time, calories, protein, user_id) VALUES ($1, $2, $3, $4)", date, calories, protein, userID); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Redirect(http.StatusFound, "/nutrition")
}

//func DeleteNutrition(c *gin.Context) {
//	db := pkg.ConnectDB()
//	ID := c.Param("id")
//	if _, err := db.Exec("DELETE FROM nutrition_logs WHERE id=$1", ID); err != nil {
//		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
//	}
//
//	// Resets id auto-increment.
//	if _, err := db.Exec("SELECT setval(pg_get_serial_sequence('nutrition_logs', 'id'), COALESCE(max(id), 0) + 1, false) FROM nutrition_logs\n"); err != nil {
//		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
//	}
//	c.Redirect(http.StatusFound, "/nutrition")
//}

//func UpdateNutrition(c *gin.Context) {
//	var id int
//	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
//		fmt.Fprintf(w, "Error: %s", err.Error())
//		return
//	}
//	db := pkg.ConnectDB()
//	if err := db.QueryRow("update books set name = $1, author = $2 where id = $3 returning id", book.Name, book.Author, book.ID).Scan(&id); err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//	}
//	c.Redirect(http.StatusFound, "/")
//}
