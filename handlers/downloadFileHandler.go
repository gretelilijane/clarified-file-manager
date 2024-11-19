package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

func DownloadFileHandler(db *sql.DB, store *sessions.CookieStore) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		session, _ := store.Get(r, "session") // Custom session name
		userId, ok := session.Values["user_id"].(int32)

		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Parse the file ID from the URL
		fileID := r.URL.Path[len("/files/"):len(r.URL.Path)]

		// Same steps as before to fetch file metadata (name, mime_type)
		var fileName, mimeType string
		query := "SELECT name, mime_type FROM files WHERE id = $1 AND user_id = $2"
		err := db.QueryRow(query, fileID, userId).Scan(&fileName, &mimeType)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Stream file content
		query = "SELECT content FROM files WHERE id = $1 AND user_id = $2"
		rows, err := db.Query(query, fileID, userId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		defer rows.Close()

		// Set response headers
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
		w.Header().Set("Content-Type", mimeType)

		// Write content in chunks
		for rows.Next() {
			var chunk []byte
			rows.Scan(&chunk)
			w.Write(chunk)
		}
	}
}
