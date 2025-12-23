package utils

func MsgForTag(tag string, param string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email address"
	case "min":
		return "Must be at least " + param + " characters"
	case "max":
		return "Must be at most " + param + " characters"
	case "eqfield":
		return "Passwords do not match"
	}
	return "Invalid input"
}
