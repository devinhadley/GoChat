// Package handlers contains the route handlers for gin.
// They are organized by the common database table or logical route.
package handlers

import (
	"errors"
	"net/http"

	"gochat/main/internal/store"
	"gochat/main/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
)

type SignUpForm struct {
	Username        string `form:"username" binding:"required,max=100"`
	Password        string `form:"password" binding:"required"`
	ConfirmPassword string `form:"confirm-password" binding:"required,eqfield=Password"`
}

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func SignUp(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", gin.H{})
}

// TODO: Move me!
func isUniqueConstraintViolatedError(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

// CreateUser effectively allows for sign up.
func CreateUser(userService *store.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var signUpForm SignUpForm

		if err := c.ShouldBind(&signUpForm); err != nil {
			// Note to self: i.(type) is a type assertion.
			// allows us to go from interface to concrete type.
			//
			displayErrors := make(map[string]string)
			if validationError, ok := err.(validator.ValidationErrors); ok {
				for _, vError := range validationError {
					displayErrors[vError.Field()] = utils.MsgForTag(vError.Tag(), vError.Param())
				}
			}

			c.HTML(http.StatusBadRequest, "signup.html", gin.H{
				"errors": displayErrors,
				"form":   signUpForm,
			})
			return
		}

		_, err := userService.CreateUser(signUpForm.Username, signUpForm.Password, c)
		if err != nil {
			if isUniqueConstraintViolatedError(err) {
				c.HTML(http.StatusBadRequest, "signup.html", gin.H{
					"errors": gin.H{"Username": "User with that username already exists."},
					"form":   signUpForm,
				})
			} else {
				utils.HandleInternalServerError(c, err, "An internal server error occured.", "signup.html", signUpForm)
			}

			return
		}

		c.Redirect(http.StatusFound, "/login")
	}
}
