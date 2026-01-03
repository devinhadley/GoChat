package main

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"gochat/main/internal/handlers"
	"gochat/main/internal/store"

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

	// TODO: Use http.CrossOriginProtection
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	templates := template.Must(template.ParseGlob("./templates/*.html"))

	http.HandleFunc("GET /{$}", handlers.CreateHomeHandler(templates))
	addUserHandlers(store.NewUserService(dbConPool), templates)

	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func addUserHandlers(userService store.UserService, templates *template.Template) {
	http.HandleFunc("GET /login", handlers.CreateLoginGetHandler(templates))
	http.HandleFunc("POST /login", handlers.CreateLoginHandler(userService, templates))
	http.HandleFunc("GET /signup", handlers.CreateSignUpGetHandler(templates))
	http.HandleFunc("POST /signup", handlers.CreateUserHandler(userService, templates))
}
