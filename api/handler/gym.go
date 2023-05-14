package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mmkamron/basefit/pkg"
	"net/http"
	"time"
)

type Activities struct {
	ID       int    `json:"id"`
	Date     string `json:"date_time"`
	Activity string `json:"activity"`
	Weight   int    `json:"weight_lifted"`
}

func Read(c *gin.Context) {
	userID, _ := c.Get("userID")
	db := pkg.ConnectDB()
	var activities Activities
	rows, err := db.Query("SELECT id, date_time, activity, weight_lifted FROM exercise_logs WHERE user_id = $1", userID)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var res []Activities
	for rows.Next() {
		if err := rows.Scan(&activities.ID, &activities.Date, &activities.Activity, &activities.Weight); err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}
		res = append(res, activities)
	}
	c.HTML(http.StatusOK, "gym.html", res)
}

func Create(c *gin.Context) {
	userID, _ := c.Get("userID")
	db := pkg.ConnectDB()
	date := time.Now().Format("01-02-2006")
	weight := c.PostForm("weight")
	activity := c.PostForm("activity")
	if _, err := db.ExecContext(c, "INSERT INTO exercise_logs(date_time, weight_lifted, activity, user_id) VALUES ($1, $2, $3, $4)", date, weight, activity, userID); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Redirect(http.StatusFound, "/gym")
}

func Delete(c *gin.Context) {
	db := pkg.ConnectDB()
	ID := c.Param("id")
	if _, err := db.Exec("DELETE FROM exercise_logs WHERE id=$1", ID); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}

	// Resets id auto-increment.
	if _, err := db.Exec("SELECT setval(pg_get_serial_sequence('exercise_logs', 'id'), COALESCE(max(id), 0) + 1, false) FROM exercise_logs\n"); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
	c.Status(http.StatusNoContent)
}

func Update(c *gin.Context) {
	// db := pkg.ConnectDB()
	// ID := c.Param("id")
    c.Status(http.StatusOK)
}
