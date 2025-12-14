package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{})
}

func main() {
	router := gin.Default()

	router.Static("/static", "./css")

	router.LoadHTMLGlob("./templates/*")

	// Routes...
	router.GET("", home)

	err := router.Run()
	if err != nil {
		log.Fatal(err)
	}
}
