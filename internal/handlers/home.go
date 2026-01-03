package handlers

import (
	"html/template"
	"net/http"

	"gochat/main/internal/utils/responses"
)

func CreateHomeHandler(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responses.RenderTemplate(w, templates, "home.html", map[string]any{
			"user": nil,
		})
	}
}
