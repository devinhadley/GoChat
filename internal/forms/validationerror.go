package forms

type ValidationErrors map[string]string

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return "validation failed"
	}
	for field, msg := range v {
		return field + ": " + msg
	}
	return ""
}
