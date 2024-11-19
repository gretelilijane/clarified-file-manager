package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/sessions"
)

func DeleteFileHandler(db *sql.DB, store *sessions.CookieStore) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		session, _ := store.Get(r, "session")
		userId, ok := session.Values["user_id"].(int32)

		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		fileID := r.URL.Path[len("/files/"):len(r.URL.Path)]

		res, err := db.Exec("DELETE FROM files WHERE id = $1 AND user_id = $2", fileID, userId)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
