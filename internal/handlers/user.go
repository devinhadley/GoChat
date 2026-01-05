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
	"gochat/main/internal/utils/sessions"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func CreateLoginGetHandler(templates *template.Template) http.HandlerFunc {
	data := map[string]any{
		"errors": map[string]string{},
		"form":   forms.LogInForm{},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		responses.RenderTemplate(w, r, templates, "login.html", data)
	}
}

func CreateSignUpGetHandler(templates *template.Template) http.HandlerFunc {
	data := map[string]any{
		"errors": map[string]string{},
		"form":   forms.SignUpForm{},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		responses.RenderTemplate(w, r, templates, "signup.html", data)
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
			responses.RenderTemplate(w, r, templates, "signup.html", map[string]any{
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
				responses.RenderTemplate(w, r, templates, "signup.html", map[string]any{
					"errors": forms.ValidationErrors{
						"Username": "A user with this username already exists.",
					},
					"form": signUpForm,
				})

			} else {
				responses.RenderInternalErrorOnTemplate(w, r, templates, "signup.html", map[string]any{
					"errors": map[string]string{},
					"form":   signUpForm,
				})
			}

			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func CreateLoginHandler(userService store.UserService, sessionService store.SessionService, templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loginForm := forms.NewLogInFormFromRequest(r)

		validationErrors := loginForm.Validate()
		if len(validationErrors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			responses.RenderTemplate(w, r, templates, "login.html", map[string]any{
				"form":                  loginForm,
				"areCredentialsInvalid": true,
			})
			return
		}

		user, err := userService.AuthenticateUser(r.Context(), loginForm.Username, loginForm.Password)
		if err != nil {
			if errors.Is(err, store.ErrInvalidCredentials) {
				responses.RenderTemplate(w, r, templates, "login.html", map[string]any{
					"form":                  loginForm,
					"areCredentialsInvalid": true,
				})
			} else {
				responses.RenderInternalErrorOnTemplate(w, r, templates, "login.html", map[string]any{})
				log.Println(err)
			}
			return
		}

		sessionCookie, err := sessions.CreateSessionCookie()
		if err != nil {
			responses.RenderInternalErrorOnTemplate(w, r, templates, "login.html", map[string]any{})
			log.Println(err)
			return
		}

		_, err = sessionService.CreateSession(r.Context(), sessionCookie.Value, user.ID, sessionCookie.Expires)
		if err != nil {
			responses.RenderInternalErrorOnTemplate(w, r, templates, "login.html", map[string]any{})
			log.Println(err)
			return
		}

		http.SetCookie(w, &sessionCookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func CreateLogoutHandler(userService store.UserService, sessionService store.SessionService, templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie(sessions.SessionCookieName)
		if err != nil {
			if !errors.Is(err, http.ErrNoCookie) {
				log.Printf("Error getting session cookie when deleting: %v", err)
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		err = sessionService.DeleteSession(r.Context(), sessionCookie.Value)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				log.Printf("Error deleting session cookie from db: %v", err)
			}
		}

		clearSessionCookie := sessions.CreateClearSessionCookie()
		http.SetCookie(w, &clearSessionCookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
