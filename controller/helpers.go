package controller

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func serveTemplate(w http.ResponseWriter, page string, data any) {
	t, err := template.ParseFiles(
		filepath.Join(templatesPath, "_base.html"),
		filepath.Join(templatesPath, page),
	)
	if err != nil {
		log.Print("Error parsing template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w, data)
}
