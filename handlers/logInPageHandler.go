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

type LogInPageData struct {
	Username     string
	Password     string
	ErrorMessage string
	Success      bool
}

func LogInPageHandler(db *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
	tmpl, err := template.ParseFiles("views/base.html", "views/login.html")

	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data := LogInPageData{
			Username:     r.FormValue("username"),
			Password:     r.FormValue("password"),
			Success:      false,
			ErrorMessage: "",
		}

		// Basic validation
		if r.Method == http.MethodPost {
			if data.Username == "" || data.Password == "" {
				data.ErrorMessage = "All fields are required!"
			} else {
				data.Success = true
			}

			var user types.User
			// Get user, password_hash, password_salt from database
			err := db.QueryRow("SELECT id, username, password_hash, password_salt FROM users WHERE username = $1", data.Username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.PasswordSalt)

			if err != nil {
				log.Println(err)
				data.Success = false
				data.ErrorMessage = "Wrong email or password"
			}

			log.Println("User: ", user, user.ID, user.Username, user.PasswordHash, user.PasswordSalt)

			// Compare password
			argon2IDHash := auth.NewArgon2idHash(1, 32, 64*1024, 32, 256)

			err = argon2IDHash.Compare(user.PasswordHash, user.PasswordSalt, []byte(data.Password))
			if err != nil {
				data.ErrorMessage = "Wrong email or password"
				data.Success = false
			} else {
				log.Println("Password is correct")
				data.Success = true

				// Set session
				session, _ := store.Get(r, "session")
				session.Values["user_id"] = user.ID
				session.Values["username"] = user.Username

				session.Save(r, w)

				// Redirect to dashboard page if login is successful
				w.Header().Set("HX-Redirect", "/upload")
				w.WriteHeader(http.StatusOK)
			}
		}

		// Render response based on request target
		target := r.Header.Get("HX-Target")

		if target == "messages" {
			tmpl.ExecuteTemplate(w, target, data)
		} else {
			tmpl.Execute(w, data)
		}

	}
}
