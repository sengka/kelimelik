package handlers

import (
	"html/template"
	"net/http"
)

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/mainpage.html")
	if err != nil {
		http.Error(w, "Şablon yüklenemedi", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
