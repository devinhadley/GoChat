package main

import (
	"context"
	"log"

	"gochat/main/internal/application"
	"gochat/main/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Create database connection pool.
	dsn := "postgres://gochat:password@localhost:5432/gochat"
	dbConPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Failed to init database connecton pool %v", err)
	}
	defer dbConPool.Close()

	router := gin.Default()
	router.Static("/static", "./static")
	router.LoadHTMLGlob("./templates/*")

	app := &application.App{
		DB: dbConPool,
	}
	router.GET("", handlers.Home)
	addUserHandlers(router, app)

	err = router.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func addUserHandlers(router *gin.Engine, app *application.App) {
	router.GET("/login", handlers.Login)
	router.GET("/signup", handlers.SignUp)
	router.POST("/signup", handlers.CreateUser(app))
}
