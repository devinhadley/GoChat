// Package responses provides a common set of responses.
package responses

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
)

// RenderTemplate renders the template with the provided data.
// Note it buffers the template in memory in order to avoid half-rendered pages.
// It will first render and the template in its entirety, then write that to the response writer
// if it is successful.
func RenderTemplate(w http.ResponseWriter, templates *template.Template, name string, data map[string]any) {
	var buf bytes.Buffer

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
func RenderInternalErrorOnTemplate(w http.ResponseWriter, templates *template.Template, name string, data map[string]any) {
	data["isShowingInternalError"] = true
	RenderTemplate(w, templates, name, data)
}
