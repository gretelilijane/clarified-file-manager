package handlers

import (
	"database/sql"
	"fmt"
	"log"
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
		fileId := r.URL.Path[len("/files/"):len(r.URL.Path)]

		var fileName, mimeType string
		var content []byte

		// query := "SELECT name, mime_type, content FROM files WHERE id = $1 AND user_id = $2"
		// err := db.QueryRow(query, fileId, userId).Scan(&fileName, &mimeType, &content)
		query := fmt.Sprintf("SELECT name, mime_type, content FROM files WHERE id = %s AND user_id = %d", fileId, userId)
		log.Println(query)
		err := db.QueryRow(query).Scan(&fileName, &mimeType, &content)
		log.Println(fileName, mimeType, err)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Set response headers
		//w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
		w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", fileName))
		w.Header().Set("Content-Type", mimeType)
		w.Write(content)
	}
}
