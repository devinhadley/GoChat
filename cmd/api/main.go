package main

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"gochat/main/internal/handlers"
	"gochat/main/internal/middleware"
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

	// Templates, and static serve setup.
	fs := http.FileServer(http.Dir("./static"))
	templates := template.Must(template.ParseGlob("./templates/*.html"))

	// Initialize services.
	userService := store.NewUserService(dbConPool)
	sessionService := store.NewSessionService(dbConPool)

	// Add routes and handlers to multiplexer.
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("GET /{$}", handlers.CreateHomeHandler(templates))
	addUserHandlers(
		mux,
		userService,
		sessionService,
		templates,
	)

	// Add middleware.
	handler := middleware.AuthMiddleware(mux, userService)
	crossOriginProtection := http.NewCrossOriginProtection()
	handler = crossOriginProtection.Handler(handler)

	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}

func addUserHandlers(mux *http.ServeMux, userService store.UserService, sessionService store.SessionService, templates *template.Template) {
	mux.HandleFunc("GET /login", handlers.CreateLoginGetHandler(templates))
	mux.HandleFunc("POST /login", handlers.CreateLoginHandler(userService, sessionService, templates))
	mux.HandleFunc("GET /logout", handlers.CreateLogoutHandler(userService, sessionService, templates))
	mux.HandleFunc("GET /signup", handlers.CreateSignUpGetHandler(templates))
	mux.HandleFunc("POST /signup", handlers.CreateUserHandler(userService, templates))
}
