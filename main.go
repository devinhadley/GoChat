package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func main() {
	router := gin.Default()

	router.Static("/static", "./css")

	router.GET("/ping", pong)

	err := router.Run()
	if err != nil {
		log.Fatal(err)
	}
}
