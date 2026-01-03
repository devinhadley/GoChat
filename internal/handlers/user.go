// Package handlers contains the route handlers for gin.
// They are organized by the common database table or logical route.
package handlers

import (
	"errors"
	"html/template"
	"net/http"

	"gochat/main/internal/forms"
	"gochat/main/internal/store"
	"gochat/main/internal/utils/responses"

	"github.com/jackc/pgx/v5/pgconn"
)

func CreateLoginGetHandler(templates *template.Template) http.HandlerFunc {
	data := map[string]any{
		"errors": map[string]string{},
		"form":   map[string]string{},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		responses.RenderTemplate(w, templates, "login.html", data)
	}
}

func CreateSignUpGetHandler(templates *template.Template) http.HandlerFunc {
	data := map[string]any{
		"errors": map[string]string{},
		"form":   map[string]string{},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		responses.RenderTemplate(w, templates, "signup.html", data)
	}
}

func isUniqueConstraintViolatedError(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func CreateUserHandler(userService store.UserService, templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		signUpForm := forms.NewSignUpFormFromRequest(r)

		validationErrors := signUpForm.Validate()
		if len(validationErrors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			responses.RenderTemplate(w, templates, "signup.html", map[string]any{
				"errors": validationErrors,
				"form":   signUpForm,
			})
			return
		}
		_, err := userService.CreateUser(signUpForm.Username, signUpForm.Password, r.Context())
		if err != nil {
			// Username is the only user populated field with a unique constraint.
			if isUniqueConstraintViolatedError(err) {

				w.WriteHeader(http.StatusBadRequest)
				responses.RenderTemplate(w, templates, "signup.html", map[string]any{
					"errors": forms.ValidationErrors{
						"Username": "A user with this username already exists.",
					},
					"form": signUpForm,
				})

			} else {
				// TODO: Show error banner instead.
				http.Error(w, "An internal error occurred.", http.StatusInternalServerError)
			}

			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func CreateLoginHandler(userService store.UserService, templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loginForm := forms.NewLogInFormFromRequest(r)

		validationErrors := loginForm.Validate()
		if len(validationErrors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			responses.RenderTemplate(w, templates, "signup.html", map[string]any{
				"errors": validationErrors,
				"form":   loginForm,
			})
			return
		}

		isAuthenticated, err := userService.AuthenticateUser(loginForm.Username, loginForm.Password, r.Context())
		if err != nil {
			// TODO: Show error banner instead.
			http.Error(w, "An internal error occurred.", http.StatusInternalServerError)
			return
		}
		if isAuthenticated {
			// Render home page.
		} else {
			// Render login page with invalid credentias error.
		}
	}
}
