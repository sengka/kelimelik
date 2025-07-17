package handlers

import (
	// "crypto/rand"
	"database/sql"
	// "encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("your-very-secret-key"))

// Şifre içerisinde özel karakter var mı kontrolü
func containsSpecialChar(s string) bool {
	specialChars := "!@#$%^&*()-_=+[]{}|;:'\",.<>/?`~"
	for _, ch := range s {
		if strings.ContainsRune(specialChars, ch) {
			return true
		}
	}
	return false
}

// // Rastgele session token üret (İstersen, artık kullanmana gerek kalmaz)
// func generateSessionToken() (string, error) {
// 	b := make([]byte, 16) // 128 bit
// 	_, err := rand.Read(b)
// 	if err != nil {
// 		return "", err
// 	}
// 	return hex.EncodeToString(b), nil
// }

// Kayıt işlemi
func RegisterHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		lastname := r.FormValue("lastname")
		email := strings.TrimSpace(strings.ToLower(r.FormValue("email")))
		password := r.FormValue("password")

		if len(password) < 8 || !containsSpecialChar(password) {
			http.Error(w, "Şifre en az 8 karakter olmalı ve en az bir özel karakter içermelidir.", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Kayıt hatası:", err)
			http.Error(w, fmt.Sprintf("Kayıt başarısız: %v", err), http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO users (name, lastname, email, password) VALUES (?, ?, ?, ?)",
			name, lastname, email, hashedPassword)
		if err != nil {
			log.Println("Kayıt hatası:", err)
			http.Error(w, "Kayıt başarısız, email zaten kullanılıyor olabilir.", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("templates/register.html")
	if err != nil {
		http.Error(w, "Şablon yüklenemedi.", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Giriş işlemi (gorilla/sessions ile)
func LoginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, "Şablon yüklenemedi.", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		email := strings.TrimSpace(strings.ToLower(r.FormValue("email")))
		password := r.FormValue("password")

		var userID int
		var hashedPassword string
		err := db.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&userID, &hashedPassword)
		if err != nil {
			log.Println("Kullanıcı bulunamadı:", err)
			http.Error(w, "Kullanıcı bulunamadı.", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			http.Error(w, "Şifre yanlış.", http.StatusUnauthorized)
			return
		}

		// Oturum aç
		session, err := store.Get(r, "session-name")
		if err != nil {
			http.Error(w, "Oturum alınamadı.", http.StatusInternalServerError)
			return
		}
		session.Values["user_id"] = userID

		err = session.Save(r, w)
		if err != nil {
			http.Error(w, "Oturum kaydedilemedi.", http.StatusInternalServerError)
			return
		}

		ensureUserProfileExists(db, userID)

		http.Redirect(w, r, "/main", http.StatusSeeOther)
		return
	}

	tmpl.Execute(w, nil)
}

// Logout işlemi (gorilla/sessions ile)
func LogoutHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Oturumu temizle
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Oturum kapatılamadı.", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
