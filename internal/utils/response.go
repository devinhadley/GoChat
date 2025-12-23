package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowInternalServerError(c *gin.Context, err error, msg string, templateName string, form any) {
	log.Printf("%s: %v", msg, err)
	c.HTML(http.StatusInternalServerError, templateName, gin.H{
		"isShowingInternalError": true,
		"form":                   form,
	})
}
