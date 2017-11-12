package main

import (
	"SCYWWeb/controllers"
	"SCYWWeb/pool"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	pool.NewConnection()
	defer pool.Close()

	wc := controllers.NewWebController()

	//gin.Recovery()
	//gin.SetMode(gin.ReleaseMode)
	//r := gin.New()

	r := gin.Default()

	r.LoadHTMLGlob("views/html/**/*")
	r.Static("/assets", "./views/assets")
	r.Static("/image", "./data/image")

	r.GET("/hello", hello)
	r.GET("/", wc.Show)
	r.GET("/show", wc.Show)
	r.GET("/upload", wc.Upload)

	r.POST("/uploadImage", wc.UploadImage)
	r.POST("/removeImage", wc.RemoveImage)
	r.Run(":80")

	log.Println("Start SCYWWeb :8889")
}

func hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"text": "Hello SCYWWeb System.",
	})
}
