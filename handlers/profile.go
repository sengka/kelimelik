package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

var tmpl = template.Must(template.ParseFiles("templates/profile.html"))

func ProfilePageHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	session, _ := store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var name, lastname, email string
	err := db.QueryRow("SELECT name, lastname, email FROM users WHERE id = ?", userID).Scan(&name, &lastname, &email)
	if err != nil {
		http.Error(w, "Kullanıcı bulunamadı", http.StatusInternalServerError)
		return
	}

	var birthdate, phone, bio string
	err = db.QueryRow("SELECT birthdate, phone, bio FROM user_profiles WHERE user_id = ?", userID).
		Scan(&birthdate, &phone, &bio)

	if err == sql.ErrNoRows {
		birthdate, phone, bio = "", "", ""
	} else if err != nil {
		http.Error(w, "Profil verileri alınamadı", http.StatusInternalServerError)
		return
	}

	data := struct {
		Name      string
		Lastname  string
		Email     string
		Birthdate string
		Phone     string
		Bio       string
	}{
		Name:      name,
		Lastname:  lastname,
		Email:     email,
		Birthdate: birthdate,
		Phone:     phone,
		Bio:       bio,
	}

	tmpl.ExecuteTemplate(w, "profile.html", data)
}

func ensureUserProfileExists(db *sql.DB, userID int) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM user_profiles WHERE user_id = ?)", userID).Scan(&exists)
	if err != nil {
		log.Println("Profil kontrolü hatası:", err)
		return
	}

	if !exists {
		_, err := db.Exec("INSERT INTO user_profiles (user_id) VALUES (?)", userID)
		if err != nil {
			log.Println("Boş profil eklenemedi:", err)
		}
	}
}

func ProfileHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	session, _ := store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// User tablosundaki temel bilgiler
	var user struct {
		Name     string
		Lastname string
		Email    string
	}

	err := db.QueryRow("SELECT name, lastname, email FROM users WHERE id = ?", userID).Scan(&user.Name, &user.Lastname, &user.Email)
	if err != nil {
		http.Error(w, "Kullanıcı bilgileri alınamadı.", http.StatusInternalServerError)
		return
	}

	// Profilden özel bilgiler
	var profile struct {
		Birthdate string
		Phone     string
		Bio       string
	}

	err = db.QueryRow("SELECT birthdate, phone, bio FROM user_profiles WHERE user_id = ?", userID).
		Scan(&profile.Birthdate, &profile.Phone, &profile.Bio)

	if err == sql.ErrNoRows {
		// Profil kaydı yoksa boş değerlerle devam et
		profile.Birthdate = ""
		profile.Phone = ""
		profile.Bio = ""
	} else if err != nil {
		http.Error(w, "Profil bilgileri alınamadı.", http.StatusInternalServerError)
		return
	}

	// Şablona gönderilecek veri
	data := struct {
		User    interface{}
		Profile interface{}
	}{
		User:    user,
		Profile: profile,
	}

	tmpl, err := template.ParseFiles("templates/profile.html")
	if err != nil {
		http.Error(w, "Şablon yüklenemedi.", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func ProfileUpdateHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	session, _ := store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	lastname := r.FormValue("lastname")
	birthdate := r.FormValue("birthdate")
	phone := r.FormValue("phone")
	bio := r.FormValue("bio")

	_, err := db.Exec("UPDATE users SET name = ?, lastname = ? WHERE id = ?", name, lastname, userID)
	if err != nil {
		http.Error(w, "Kullanıcı bilgileri güncellenemedi.", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("UPDATE user_profiles SET birthdate = ?, phone = ?, bio = ? WHERE user_id = ?", birthdate, phone, bio, userID)
	if err != nil {
		http.Error(w, "Profil bilgileri güncellenemedi.", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
