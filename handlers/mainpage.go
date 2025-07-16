package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
)

func slice(s string, start, end int) string {
	if len(s) < end {
		end = len(s)
	}
	return s[start:end]
}

func MainPageHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	rows, err := db.Query("SELECT id, title, content FROM posts ORDER BY id ASC")
	if err != nil {
		http.Error(w, "Veriler çekilemedi", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Post struct {
		ID      int
		Title   string
		Content string
	}

	var posts []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.Title, &p.Content)
		if err != nil {
			http.Error(w, "Veriler okunamadı", http.StatusInternalServerError)
			return
		}
		posts = append(posts, p)
	}

	funcMap := template.FuncMap{
		"slice": slice,
	}

	tmpl, err := template.New("mainpage.html").Funcs(funcMap).ParseFiles("templates/mainpage.html")
	if err != nil {
		http.Error(w, "Sayfa yüklenemedi", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, posts)
	if err != nil {
		http.Error(w, "Sayfa render edilemedi", http.StatusInternalServerError)
	}
}
