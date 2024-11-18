package handlers

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// type DeleteFileData struct {
// 	FileId       int32
// 	UserId       int32
// 	ErrorMessage string
// 	Success      bool
// }

func LogOutHandler(store *sessions.CookieStore) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		session, _ := store.Get(r, "session") // Custom session name
		session.Values["user_id"] = nil
		session.Save(r, w)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
