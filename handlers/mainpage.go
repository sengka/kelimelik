package handlers

import (
	"database/sql"
	"html/template"
	"kelimelik/models"
	"net/http"
)

func MainPageHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	tmpl, err := template.ParseFiles("templates/main.html")
	if err != nil {
		http.Error(w, "Sayfa yüklenemedi", http.StatusInternalServerError)
		return
	}

	rows, err := db.Query("SELECT id, title, content FROM posts ORDER BY id DESC")
	if err != nil {
		http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content)
		if err != nil {
			http.Error(w, "Veri okunamadı", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	tmpl.Execute(w, posts)
}
