package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

func DownloadFileHandler(db *sql.DB, store *sessions.CookieStore) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		session, _ := store.Get(r, "session")
		userId, ok := session.Values["user_id"].(int32)

		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Parse the file ID from the URL
		fileID := r.URL.Path[len("/files/"):len(r.URL.Path)]

		var fileName, mimeType string
		var content []byte

		query := "SELECT name, mime_type, content FROM files WHERE id = $1 AND user_id = $2"
		err := db.QueryRow(query, fileID, userId).Scan(&fileName, &mimeType, &content)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Set response headers
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
		w.Header().Set("Content-Type", mimeType)
		w.Write(content)
	}
}
