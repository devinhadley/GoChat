package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleInternalServerError logs an error with a message and displays an internal server error banner
// to the user on the provided template and optional form.
func HandleInternalServerError(c *gin.Context, err error, msg string, templateName string, form any) {
	log.Printf("%s: %v", msg, err)
	c.HTML(http.StatusInternalServerError, templateName, gin.H{
		"isShowingInternalError": true,
		"form":                   form,
	})
}
