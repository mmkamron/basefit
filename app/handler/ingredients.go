package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mmkamron/basefit/pkg"
	"io"
	"net/http"
)

func Ingredients(c *gin.Context) {
	config := pkg.Load()
	item := c.Param("item")
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.api-ninjas.com/v1/nutrition?query=%s", item), nil)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("X-Api-Key", config.ApiNinjas)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	type Data struct {
		Name          string  `json:"name"`
		Calories      float64 `json:"calories"`
		ServingSize   float64 `json:"serving_size_g"`
		FatTotal      float64 `json:"fat_total_g"`
		Protein       float64 `json:"protein_g"`
		Carbohydrates float64 `json:"carbohydrates_total_g"`
	}
	var data []Data
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, v := range data {
		c.JSON(http.StatusOK, gin.H{
			"name":          v.Name,
			"calories":      v.Calories,
			"serving_size":  v.ServingSize,
			"fat_total":     v.FatTotal,
			"protein":       v.Protein,
			"carbohydrates": v.Carbohydrates,
		})
	}
}
