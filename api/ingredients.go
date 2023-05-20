package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mmkamron/basefit/pkg"
)

func Ingredients(c *gin.Context) {
	config := pkg.Load()
	food := c.PostForm("food")
	size := c.PostForm("size")
	if size != "" {
		size += string('g')
	}
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"https://api.api-ninjas.com/v1/nutrition?query=%s+%s",
			size,
			strings.ReplaceAll(food, " ", "+"),
		),
		nil,
	)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("X-Api-Key", config.ApiNinjas)
	req.Header.Set("application", "x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	type Data struct {
		Name     string  `json:"name"`
		Size     float64 `json:"serving_size_g"`
		Calories float64 `json:"calories"`
		Protein  float64 `json:"protein_g"`
		Fat      float64 `json:"fat_total_g"`
		Carbs    float64 `json:"carbohydrates_total_g"`
	}
	var data []Data
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	c.HTML(http.StatusOK, "ingredients.html", data)
	//for _, v := range data {
	//	c.JSON(http.StatusOK, gin.H{
	//		"name":          v.Name,
	//		"calories":      v.Calories,
	//		"serving_size":  v.ServingSize,
	//		"fat_total":     v.FatTotal,
	//		"protein":       v.Protein,
	//		"carbohydrates": v.Carbohydrates,
	//	})
	//}
}
