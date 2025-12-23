// Package handlers contains the route handlers for gin.
// They are organized by the common database table or logical route.
package handlers

import (
	"net/http"

	"gochat/main/internal/store"
	"gochat/main/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

// CreateUser effectively allows for sign up.
func CreateUser(userStore store.UserStore) gin.HandlerFunc {
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

		// 1. Check if user with username already exists.
		// 2. Hash the password.
		// 3. Store user with hashed password in the database.

		if err != nil {

			utils.ShowInternalServerError(c, err, "An error occurred when querying if username exists", "signup.html", signUpForm)
			return
		}

		isDupUsername, error :=  userStore.DoesUserWithUsernameExist(signUpForm.Username, c)
		if ok, e :=  {
			c.HTML(http.StatusBadRequest, "signup.html", gin.H{
				"errors": gin.H{"Username": "User with that username already exists."},
				"form":   signUpForm,
			})
			return
		}

		c.Redirect(http.StatusFound, "/login")
	}
}
