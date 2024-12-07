package handlers

import (
	"clarified-file-management/auth"
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

type SignUpData struct {
	Username        string
	Password        string
	ConfirmPassword string
	ErrorMessage    string
	Success         bool
}

func executeSignUp(db *sql.DB, r *http.Request) (SignUpData, error) {
	data := SignUpData{
		Username:        r.FormValue("username"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirm_password"),
		Success:         false,
		ErrorMessage:    "",
	}

	if r.Method != http.MethodPost {
		return data, nil
	}

	if data.Username == "" || data.Password == "" || data.ConfirmPassword == "" {
		data.ErrorMessage = "All fields are required!"
		return data, nil
	}

	if data.Password != data.ConfirmPassword {
		data.ErrorMessage = "Passwords do not match!"
		return data, nil
	}

	// Generate hash and salt
	argon2IDHash := auth.NewArgon2idHash(1, 32, 64*1024, 32, 256)
	hashSalt, err := argon2IDHash.GenerateHash([]byte(data.Password), nil)

	if err != nil {
		return data, err
	}

	// Try to insert user into database
	_, err = db.Exec("INSERT INTO users (username, password_hash, password_salt) VALUES ($1, $2, $3)", data.Username, hashSalt.Hash, hashSalt.Salt)

	if err != nil {
		data.ErrorMessage = "Username is already taken"
		return data, nil
	}

	data.Success = true

	return data, nil
}

func SignUpPageHandler(db *sql.DB) http.HandlerFunc {
	tmpl, err := template.ParseFiles("views/base.html", "views/signup.html")

	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {

		data, err := executeSignUp(db, r)

		// Handle server errors
		if err != nil {
			data.ErrorMessage = "A server error occured, please try again later"
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
