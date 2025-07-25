package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
)

func NewPostHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	tmpl, err := template.ParseFiles("templates/newpost.html")
	if err != nil {
		http.Error(w, "Sayfa yüklenemedi", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodGet {
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")

		userID := 1

		_, err := db.Exec("INSERT INTO posts (title, content, user_id) VALUES (?, ?, ?)", title, content, userID)
		if err != nil {
			http.Error(w, "Veritabanına eklenemedi", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/main", http.StatusSeeOther)
	}

}
