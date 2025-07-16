package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
)

func PostDetailHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	idStr := r.URL.Path[len("/post/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Geçersiz ID", http.StatusBadRequest)
		return
	}

	var title, content string
	err = db.QueryRow("SELECT title, content FROM posts WHERE id = ?", id).Scan(&title, &content)
	if err != nil {
		http.Error(w, "Yazı bulunamadı", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/post_detail.html")
	if err != nil {
		http.Error(w, "Şablon yüklenemedi", http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"Title":   title,
		"Content": content,
	}

	tmpl.Execute(w, data)
}
