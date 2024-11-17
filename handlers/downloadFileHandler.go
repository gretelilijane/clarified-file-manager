package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

func DownloadFileHandler(db *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
	// tmpl, err := template.ParseFiles("views/base.html", "views/login.html")

	// if err != nil {
	// 	log.Fatalf("Error parsing template: %v", err)
	// }

	return func(w http.ResponseWriter, r *http.Request) {

		session, _ := store.Get(r, "session") // Custom session name
		userId, ok := session.Values["user_id"].(int32)

		if !ok {
			w.Header().Set("HX-Redirect", "/login")
			w.WriteHeader(http.StatusFound) //TODO: check if this is the correct status code
		}

		log.Println("User ID: ", userId)

		// Parse the file ID from the URL
		fileID := r.URL.Path[len("/files/"):len(r.URL.Path)]
		userId = 1

		// Same steps as before to fetch file metadata (name, mime_type)
		var fileName, mimeType string
		query := "SELECT name, mime_type FROM files WHERE id = $1 AND user_id = $2"
		err := db.QueryRow(query, fileID, userId).Scan(&fileName, &mimeType)
		if err != nil {
			// Handle errors as before
			log.Println("Not possible", err)
		}

		// Stream file content
		query = "SELECT content FROM files WHERE id = $1 AND user_id = $2"
		rows, err := db.Query(query, fileID, userId)
		if err != nil {
			// Handle errors as before
			log.Println("Not possible", err)
		}
		defer rows.Close()

		// Set response headers
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
		w.Header().Set("Content-Type", mimeType)

		// Write content in chunks
		for rows.Next() {
			var chunk []byte
			err := rows.Scan(&chunk)
			if err != nil {
				// Handle error
				log.Println("Not possible", err)
			}
			_, err = w.Write(chunk)
			if err != nil {
				// Handle error
				log.Println("Not possible", err)
			}
		}
	}
}
