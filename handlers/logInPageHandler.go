package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func LogInPageHandler() http.HandlerFunc {
	tmpl, err := template.ParseFiles("views/base.html", "views/login.html")

	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("IndexPageHandler")
		tmpl.Execute(w, nil)
	}
}
