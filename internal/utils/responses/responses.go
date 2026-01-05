// Package responses provides a common set of responses.
package responses

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"gochat/main/internal/store"
	"gochat/main/internal/utils/sessions"
)

// RenderTemplate renders the template with the provided data.
// It will also attempt to inject the user into the template data from the request.
// Note it buffers the template in memory in order to avoid half-rendered pages.
// It will first render and the template in its entirety, then write that to the response writer
// if it is successful.
func RenderTemplate(w http.ResponseWriter, r *http.Request, templates *template.Template, name string, data map[string]any) {
	if data == nil {
		data = make(map[string]any)
	}

	var buf bytes.Buffer

	user, ok := r.Context().Value(sessions.UserContextKey).(store.User)
	if ok {
		data["user"] = user
	}

	err := templates.ExecuteTemplate(&buf, name, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

// RenderInternalErrorOnTemplate will display a red banner stating an internal server error
// above whatever passed template.
func RenderInternalErrorOnTemplate(w http.ResponseWriter, r *http.Request, templates *template.Template, name string, data map[string]any) {
	data["isShowingInternalError"] = true
	RenderTemplate(w, r, templates, name, data)
}
