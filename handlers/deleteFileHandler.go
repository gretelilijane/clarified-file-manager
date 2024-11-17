package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

// type DeleteFileData struct {
// 	FileId       int32
// 	UserId       int32
// 	ErrorMessage string
// 	Success      bool
// }

func DeleteFileHandler(db *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
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

		log.Println("File ID: ", fileID)

		userId = 1

		// delete file with userId and fileID
		res, err := db.Exec("DELETE FROM files WHERE id = $1 AND user_id = $2", fileID, userId)

		if err != nil {
			log.Fatal("Not possible", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			log.Fatal("Not possible", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		if rowsAffected == 0 {
			log.Println("No rows affected")
			w.WriteHeader(http.StatusNotFound)
		}
	}
}
