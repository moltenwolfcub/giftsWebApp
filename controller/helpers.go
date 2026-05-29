package controller

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

func serveTemplate(w http.ResponseWriter, page string, data any) {
	t, err := template.ParseFiles(
		filepath.Join(templatesPath, "_base.html"),
		filepath.Join(templatesPath, page),
	)
	if err != nil {
		log.Print("Error parsing template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	t.Execute(w, data)
}
