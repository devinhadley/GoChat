package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{})
}

func login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func signup(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", gin.H{})
}

func main() {
	router := gin.Default()

	router.Static("/static", "./static")

	router.LoadHTMLGlob("./templates/*")

	// Routes...
	router.GET("", home)
	router.GET("/login", login)
	router.GET("/signup", signup)

	err := router.Run()
	if err != nil {
		log.Fatal(err)
	}
}
