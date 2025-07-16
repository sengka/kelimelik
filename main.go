package main

import (
	"database/sql"
	"kelimelik/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := initDB()
	defer db.Close()

	r := mux.NewRouter()

	// Anasayfa
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/homepage.html")
	}).Methods("GET")

	// Register sayfası
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterHandler(w, r, db)
	}).Methods("GET", "POST")

	// Login sayfası
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginHandler(w, r, db)
	}).Methods("GET", "POST")

	// Anasayfa / Main
	r.HandleFunc("/main", func(w http.ResponseWriter, r *http.Request) {
		handlers.MainPageHandler(w, r, db)
	}).Methods("GET")

	// Post detay sayfası
	r.HandleFunc("/post/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.PostDetailHandler(w, r, db)
	}).Methods("GET")

	// Yeni post ekleme
	r.HandleFunc("/newpost", func(w http.ResponseWriter, r *http.Request) {
		handlers.NewPostHandler(w, r, db)
	}).Methods("GET", "POST")

	// Çıkış
	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		handlers.LogoutHandler(w, r, db)
	}).Methods("GET")

	log.Println("Server çalışıyor : http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal(err)
	}

	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`

	sessionTable := `
	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		session_token TEXT NOT NULL UNIQUE,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	postTable := `
	CREATE TABLE IF NOT EXISTS posts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	user_id INTEGER,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(user_id) REFERENCES users(id)
);`
	_, err = db.Exec(userTable)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(sessionTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(postTable)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
