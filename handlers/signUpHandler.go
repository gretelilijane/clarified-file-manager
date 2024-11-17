package handlers

import (
	"log"
	"net/http"
)

func SignUpHandler() http.HandlerFunc {
	// tmpl, err := template.ParseFiles("views/base.html", "views/signup.html")

	// if err != nil {
	// 	log.Fatalf("Error parsing template: %v", err)
	// }

	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("SignUpHandler", r.FormValue("username"))

		r.ParseForm()

		username := r.FormValue("username")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")
		//tmpl.Execute(w, nil)

		// Basic validation
		if username == "" || password == "" || confirmPassword == "" {
			http.Error(w, "<div class='alert alert-danger'>All fields are required!</div>", http.StatusBadRequest)
			return
		}

		if password != confirmPassword {
			http.Error(w, "<div class='alert alert-danger'>Passwords do not match!</div>", http.StatusBadRequest)
			return
		}

		log.Println("Request: ", username, password, confirmPassword)
	}
}
