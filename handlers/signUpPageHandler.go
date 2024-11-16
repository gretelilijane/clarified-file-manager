package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func SignUpPageHandler() http.HandlerFunc {
	tmpl, err := template.ParseFiles("views/base.html", "views/signup.html")

	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("IndexPageHandler")
		tmpl.Execute(w, nil)
	}
}
