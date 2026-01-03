// Package handlers contains the route handlers for gin.
// They are organized by the common database table or logical route.
package handlers

import (
	"errors"
	"html/template"
	"log"
	"net/http"

	"gochat/main/internal/forms"
	"gochat/main/internal/store"
	"gochat/main/internal/utils/responses"

	"github.com/jackc/pgx/v5/pgconn"
)

func CreateLoginGetHandler(templates *template.Template) http.HandlerFunc {
	data := map[string]any{
		"errors": map[string]string{},
		"form":   forms.LogInForm{},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		responses.RenderTemplate(w, templates, "login.html", data)
	}
}

func CreateSignUpGetHandler(templates *template.Template) http.HandlerFunc {
	data := map[string]any{
		"errors": map[string]string{},
		"form":   forms.SignUpForm{},
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
			// TODO: Kill function and add explicit conditional.
			if isUniqueConstraintViolatedError(err) {

				w.WriteHeader(http.StatusBadRequest)
				responses.RenderTemplate(w, templates, "signup.html", map[string]any{
					"errors": forms.ValidationErrors{
						"Username": "A user with this username already exists.",
					},
					"form": signUpForm,
				})

			} else {
				responses.RenderInternalErrorOnTemplate(w, templates, "signup.html", map[string]any{
					"errors": map[string]string{},
					"form":   signUpForm,
				})
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
			responses.RenderTemplate(w, templates, "login.html", map[string]any{
				"form":                  loginForm,
				"areCredentialsInvalid": true,
			})
			return
		}

		isAuthenticated, err := userService.AuthenticateUser(loginForm.Username, loginForm.Password, r.Context())
		if err != nil {
			responses.RenderInternalErrorOnTemplate(w, templates, "login.html", map[string]any{})
			log.Println(err)
			return
		}

		if isAuthenticated {
			// Render home page.
			// TODO: Add session id to cookies.
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			responses.RenderTemplate(w, templates, "login.html", map[string]any{
				"form":                  loginForm,
				"areCredentialsInvalid": true,
			})
		}
	}
}
