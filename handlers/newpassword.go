package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Şifre sıfırlama formunu gösterme ve güncelleme
func NewPasswordHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case http.MethodGet:
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "Geçersiz bağlantı", http.StatusBadRequest)
			return
		}

		// Token geçerli mi ve süresi dolmamış mı kontrol et
		var email string
		var expiresAt time.Time
		err := db.QueryRow("SELECT email, expires_at FROM password_resets WHERE token = ?", token).Scan(&email, &expiresAt)
		if err != nil || time.Now().After(expiresAt) {
			http.Error(w, "Geçersiz veya süresi dolmuş bağlantı", http.StatusUnauthorized)
			return
		}

		tmpl, _ := template.ParseFiles("templates/newpassword.html")
		tmpl.Execute(w, struct{ Token string }{Token: token})

	case http.MethodPost:
		r.ParseForm()
		token := r.FormValue("token")
		newPassword := r.FormValue("password")

		if len(newPassword) < 8 || !containsSpecialChar(newPassword) {
			http.Error(w, "Şifre en az 8 karakter ve özel karakter içermelidir", http.StatusBadRequest)
			return
		}

		// Token doğrulama
		var email string
		var expiresAt time.Time
		err := db.QueryRow("SELECT email, expires_at FROM password_resets WHERE token = ?", token).Scan(&email, &expiresAt)
		if err != nil || time.Now().After(expiresAt) {
			http.Error(w, "Geçersiz veya süresi dolmuş token", http.StatusUnauthorized)
			return
		}

		// Şifreyi hashle
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Şifre işlenemedi", http.StatusInternalServerError)
			return
		}

		// Kullanıcının şifresini güncelle
		_, err = db.Exec("UPDATE users SET password = ? WHERE email = ?", hashedPassword, email)
		if err != nil {
			http.Error(w, "Şifre güncellenemedi", http.StatusInternalServerError)
			return
		}

		// Token sil
		db.Exec("DELETE FROM password_resets WHERE token = ?", token)

		fmt.Fprintf(w, "Şifreniz başarıyla güncellendi.")

	}
}
