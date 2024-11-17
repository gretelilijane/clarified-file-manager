package handlers

import (
	"clarified-file-management/auth"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type SignUpPageData struct {
	Username        string
	Password        string
	ConfirmPassword string
	ErrorMessage    string
	Success         bool
	UserCreated     bool
}

func SignUpPageHandler(db *sql.DB) http.HandlerFunc {
	tmpl, err := template.ParseFiles("views/base.html", "views/signup.html")

	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data := SignUpPageData{
			Username:        r.FormValue("username"),
			Password:        r.FormValue("password"),
			ConfirmPassword: r.FormValue("confirm_password"),
			Success:         false,
			UserCreated:     false,
			ErrorMessage:    "",
		}

		// Basic validation
		if r.Method == http.MethodPost {
			if data.Username == "" || data.Password == "" || data.ConfirmPassword == "" {
				data.ErrorMessage = "All fields are required!"
			} else if data.Password != data.ConfirmPassword {
				data.ErrorMessage = "Passwords do not match!"
			} else {

				// Check if username is already taken
				var userExists bool = false
				// err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", data.Username).Scan(&userExists)

				// if err != nil {
				// 	log.Println(err)
				// 	data.ErrorMessage = "Esines tõrge, palun proovi uuesti"
				// }

				if userExists {
					data.ErrorMessage = "This username is already taken."
				} else {
					data.Success = true
				}
			}
		}

		// Save to database
		if data.Success {

			argon2IDHash := auth.NewArgon2idHash(1, 32, 64*1024, 32, 256)

			hashSalt, err := argon2IDHash.GenerateHash([]byte(data.Password), nil)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			// fmt.Println(hashSalt.Hash)
			// fmt.Println(hashSalt.Salt)

			// err = argon2IDHash.Compare(hashSalt.Hash, hashSalt.Salt, []byte(data.Password))
			// if err != nil {
			// 	fmt.Fprintln(os.Stderr, err)
			// 	os.Exit(1)
			// }
			// fmt.Println("argon2IDHash Password and Hash match")

			_, err = db.Exec("INSERT INTO users (username, password_hash, password_salt) VALUES ($1, $2, $3)", data.Username, hashSalt.Hash, hashSalt.Salt)

			if err != nil {
				// log.Fatalf("Error inserting user: %v", err)
				data.ErrorMessage = "Error inserting user!"
				data.Success = false
			} else {
				log.Println("User created successfully")
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
