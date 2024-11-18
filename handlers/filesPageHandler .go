package handlers

import (
	"clarified-file-management/types"
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

type FilesPageData struct {
	UserID       int32
	ErrorMessage string
	Success      bool
	Files        []types.File
}

func getUserFiles(db *sql.DB, userID int32) ([]types.File, error) {
	var files []types.File
	rows, err := db.Query("SELECT id, name, mime_type, size, uploaded_at FROM files WHERE user_id = $1 ORDER BY uploaded_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var file types.File
		if err := rows.Scan(&file.ID, &file.Name, &file.MimeType, &file.Size, &file.UploadedAt); err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

func FilesPageHandler(db *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
	tmpl, err := template.ParseFiles("views/base.html", "views/files.html", "views/upload-form.html")

	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {

		session, _ := store.Get(r, "session") // Custom session name
		userId, ok := session.Values["user_id"].(int32)

		if !ok {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		log.Println("User ID: ", userId)

		data := FilesPageData{
			UserID:       userId,
			Success:      false,
			Files:        nil,
			ErrorMessage: "",
		}

		if r.Method == http.MethodGet {
			files, err := getUserFiles(db, userId)
			if err != nil {
				http.Error(w, "Failed to retrieve files", http.StatusInternalServerError)
				data.ErrorMessage = "Failed to retrieve files"
				log.Fatal(err)
			} else {
				data.Files = files
			}
		}

		// Render response based on request target
		target := r.Header.Get("HX-Target")

		log.Println("Target: ", target)

		if target == "messages" {
			tmpl.ExecuteTemplate(w, target, data)
		} else if target == "file-list" {
			tmpl.ExecuteTemplate(w, target, data.Files)
		} else {
			tmpl.Execute(w, data)
		}
	}
}
