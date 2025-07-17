package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"time"

	"gopkg.in/gomail.v2"
)

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/smtp.html")
		if err != nil {
			http.Error(w, "Sayfa yüklenemedi", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		email := r.FormValue("email")

		// E-posta kullanıcıya ait mi kontrol et
		var userID int
		err := db.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&userID)
		if err != nil {
			http.Error(w, "Bu e-posta ile kayıtlı kullanıcı bulunamadı.", http.StatusNotFound)
			return
		}

		// Token oluştur
		token := generateResetToken()
		expiration := time.Now().Add(1 * time.Hour)

		// Veritabanına token kaydet
		_, err = db.Exec("INSERT INTO password_resets (email, token, expires_at) VALUES (?, ?, ?)", email, token, expiration)
		if err != nil {
			http.Error(w, "Token kaydı sırasında bir hata oluştu.", http.StatusInternalServerError)
			return
		}

		// Mail gönder
		err = sendResetEmail(email, token)
		if err != nil {
			http.Error(w, "Mail gönderilirken hata oluştu.", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Şifre sıfırlama bağlantısı e-posta adresinize gönderildi.")
	}
}

// Rastgele token üretici
func generateResetToken() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func sendResetEmail(toEmail, token string) error {
	from := "senagulkara@gmail.com"
	password := "ryrv gosr usgo kugx" // Gmail uygulama şifresi

	// Şifre sıfırlama bağlantısı
	resetLink := fmt.Sprintf("http://localhost:8080/new-password?token=%s", token)

	// Mail içeriği
	subject := "Şifre Sıfırlama Talebi"
	body := fmt.Sprintf(`
Merhaba,

Şifrenizi sıfırlamak için aşağıdaki bağlantıya tıklayın:

%s

Eğer bu talebi siz yapmadıysanız, bu mesajı görmezden gelebilirsiniz.
`, resetLink)

	// Mesajı oluştur
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// Mail sunucusuna bağlan
	d := gomail.NewDialer("smtp.gmail.com", 587, from, password)

	// Gönder
	return d.DialAndSend(m)
}
