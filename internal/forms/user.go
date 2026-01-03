// Package forms handels form struct mapping and form validation.
package forms

import "net/http"

type SignUpForm struct {
	Username        string
	Password        string
	ConfirmPassword string
}

func NewSignUpFormFromRequest(r *http.Request) SignUpForm {
	return SignUpForm{
		Username:        r.FormValue("username"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirm-password"),
	}
}

func (form *SignUpForm) Validate() ValidationErrors {
	validationErrors := make(ValidationErrors)

	if len(form.Username) == 0 {
		validationErrors["Username"] = "Username can not be empty."
	} else if len(form.Username) > 30 {
		validationErrors["Username"] = "Username can not be greater than 30 characters."
	}

	if len(form.Password) == 0 {
		validationErrors["Password"] = "Password can not be empty."
	} else if len(form.Password) > 64 {
		validationErrors["Password"] = "Password can not be greater than 64 characters."
	}

	if len(form.ConfirmPassword) == 0 {
		validationErrors["ConfirmPassword"] = "Confirm password can not be empty."
	} else if form.ConfirmPassword != form.Password {
		validationErrors["ConfirmPassword"] = "Passwords do not match."
	}

	return validationErrors
}

type LogInForm struct {
	Username string
	Password string
}

func NewLogInFormFromRequest(r *http.Request) LogInForm {
	return LogInForm{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}
}

func (form *LogInForm) Validate() ValidationErrors {
	validationErrors := make(ValidationErrors)

	if len(form.Username) == 0 {
		validationErrors["Username"] = "Username can not be empty."
	} else if len(form.Username) > 30 {
		validationErrors["Username"] = "Username can not be greater than 30 characters."
	}

	if len(form.Password) < 8 {
		validationErrors["Password"] = "Password must be 8 or more characters."
	} else if len(form.Password) > 64 {
		validationErrors["Password"] = "Password can not be greater than 64 characters."
	}

	return validationErrors
}
