package handlers

import (
	"clarified-file-management/auth"
	"clarified-file-management/types"
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

type LogInData struct {
	Username     string
	Password     string
	ErrorMessage string
}

func executeLogInPage(db *sql.DB, store *sessions.CookieStore, w http.ResponseWriter, r *http.Request) (LogInData, error) {
	data := LogInData{
		Username:     r.FormValue("username"),
		Password:     r.FormValue("password"),
		ErrorMessage: "",
	}

	if r.Method != http.MethodPost {
		return data, nil
	}

	if data.Username == "" || data.Password == "" {
		data.ErrorMessage = "All fields are required!"
		return data, nil
	}

	// Get user, password_hash, password_salt from database
	var user types.User
	err := db.QueryRow("SELECT id, username, password_hash, password_salt FROM users WHERE username = $1", data.Username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.PasswordSalt)

	if err != nil {
		data.ErrorMessage = "Wrong username or password"
		return data, nil
	}

	// Compare password
	argon2IDHash := auth.NewArgon2idHash(1, 32, 64*1024, 32, 256)
	err = argon2IDHash.Compare(user.PasswordHash, user.PasswordSalt, []byte(data.Password))

	if err != nil {
		data.ErrorMessage = "Wrong username or password"
		return data, nil
	}

	// Set session
	session, _ := store.Get(r, "session")
	session.Values["user_id"] = user.ID
	session.Save(r, w)

	// Redirect to dashboard page if login is successful
	w.Header().Set("HX-Redirect", "/files")
	w.WriteHeader(http.StatusOK)

	return data, nil
}

func LogInPageHandler(db *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
	tmpl, err := template.ParseFiles("views/base.html", "views/login.html")

	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Render response based on request target
		target := r.Header.Get("HX-Target")

		data, err := executeLogInPage(db, store, w, r)

		// Handle server errors
		if err != nil {
			data.ErrorMessage = "A server error occured, please try again later"
		}

		if target == "messages" {
			tmpl.ExecuteTemplate(w, target, data)
		} else {
			tmpl.Execute(w, data)
		}
	}
}
