package handlers

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func LogOutHandler(store *sessions.CookieStore) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		session, _ := store.Get(r, "session")
		session.Values["user_id"] = nil
		session.Save(r, w)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
