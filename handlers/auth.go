package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"html/template"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

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

// Rastgele session token üret
func generateSessionToken() (string, error) {
	b := make([]byte, 16) // 128 bit
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Kullanıcıdan gelen session_token çerezine göre kullanıcı id çek
func getUserIDFromSession(r *http.Request, db *sql.DB) (int, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return 0, errors.New("session yok")
	}

	var userID int
	err = db.QueryRow("SELECT user_id FROM sessions WHERE session_token = ?", cookie.Value).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("geçersiz session token")
		}
		return 0, err
	}

	return userID, nil
}

// Register işlemi
func RegisterHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		if len(password) < 8 || !containsSpecialChar(password) {
			http.Error(w, "Şifre en az 8 karakter olmalı ve en az bir özel karakter içermelidir.", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Şifre işlenemedi", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO users (email, password) VALUES (?, ?)", email, hashedPassword)
		if err != nil {
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

// Login işlemi
func LoginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, "Şablon yüklenemedi.", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		var userID int
		var hashedPassword string
		err := db.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&userID, &hashedPassword)
		if err != nil {
			http.Error(w, "Kullanıcı bulunamadı.", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			http.Error(w, "Şifre yanlış.", http.StatusUnauthorized)
			return
		}

		token, err := generateSessionToken()
		if err != nil {
			http.Error(w, "Session token oluşturulamadı", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO sessions (user_id, session_token) VALUES (?, ?)", userID, token)
		if err != nil {
			http.Error(w, "Session kaydedilemedi", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "session_token",
			Value: token,
			Path:  "/",
		})

		http.Redirect(w, r, "/main", http.StatusSeeOther)
		return
	}

	tmpl.Execute(w, nil)
}
func LogoutHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Cookie'den session token'ı al
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Session token veritabanından sil
	_, err = db.Exec("DELETE FROM sessions WHERE session_token = ?", cookie.Value)
	if err != nil {
		http.Error(w, "Oturum silinemedi", http.StatusInternalServerError)
		return
	}

	// Cookie'yi tarayıcıdan sil
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // tarayıcıdan silmek için
	})

	// Ana sayfaya veya giriş sayfasına yönlendir
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
