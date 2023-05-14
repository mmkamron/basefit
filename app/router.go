package app

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"
	"github.com/mmkamron/basefit/app/handler"
	"github.com/mmkamron/basefit/pkg"
	"log"
	"net/http"
)

func Init() {
	r := gin.Default()
	config := pkg.Load()

	r.Static("/static", "./web")
	r.LoadHTMLGlob("web/*")

	db := pkg.ConnectDB()
	store, err := postgres.NewStore(db, []byte(config.CookieSecret))
	if err != nil {
		log.Println(err)
	}

	r.Use(sessions.Sessions("session", store))

	gym := r.Group("/gym")
	gym.Use(handler.Auth)
	gym.GET("/", handler.Read)
	gym.POST("/", handler.Create)
	//gym.DELETE("/:id", handler.Delete)
	//gym.PUT("/:id", handler.Update)

	nutrition := r.Group("/nutrition")
	nutrition.Use(handler.Auth)
	nutrition.GET("/", handler.ReadNutrition)
	nutrition.POST("/", handler.CreateNutrition)
	//nutrition.DELETE("/", handler.DeleteNutrition)

	// TODO: group these routes together
	r.GET("/oauth", handler.Oauth)
	r.GET("/googlecallback", handler.Callback)
	r.GET("/logout", handler.Logout)
	r.GET("/unauthorized", func(c *gin.Context) {
		c.HTML(http.StatusOK, "401.html", nil)
	})

	ingredients := r.Group("/ingredients")
	ingredients.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "ingredients.html", nil)
	})
	ingredients.POST("/", handler.Ingredients)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.Run()
}
