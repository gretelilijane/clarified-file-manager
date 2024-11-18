package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

func IndexPageHandler(store *sessions.CookieStore) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("IndexPageHandler")

		session, _ := store.Get(r, "session") // Custom session name
		_, ok := session.Values["user_id"].(int32)

		if !ok {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		http.Redirect(w, r, "/files", http.StatusFound)
	}
}
