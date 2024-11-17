package handlers

import (
	"clarified-file-management/types"
	"database/sql"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

type UploadPageData struct {
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

func UploadPageHandler(db *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
	tmpl, err := template.ParseFiles("views/base.html", "views/upload.html")

	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {

		session, _ := store.Get(r, "session") // Custom session name
		userId, ok := session.Values["user_id"].(int32)

		if !ok {
			w.Header().Set("HX-Redirect", "/login")
			w.WriteHeader(http.StatusFound) //TODO: check if this is the correct status code
		}

		log.Println("User ID: ", userId)

		data := UploadPageData{
			UserID:       userId,
			Success:      false,
			Files:        nil,
			ErrorMessage: "",
		}

		// POST
		if r.Method == http.MethodPost {

			// Parse the form data (multipart/form-data)
			err := r.ParseMultipartForm(10 << 20) // 10 MB limit for file size
			if err != nil {
				http.Error(w, "Unable to parse form", http.StatusBadRequest)
				return
			}

			// Retrieve the f from the form
			f, fileHeader, err := r.FormFile("avatar")
			if err != nil {
				http.Error(w, "Unable to retrieve the file", http.StatusBadRequest)
				return
			}
			defer f.Close()

			// Retrieve the file content (reading the file into a byte slice)
			content, err := io.ReadAll(f)
			if err != nil {
				http.Error(w, "Failed to read file content", http.StatusInternalServerError)
				return
			}

			uploadedFile := types.File{
				UserID:   userId, // For example, assuming user ID is 1
				Name:     fileHeader.Filename,
				MimeType: fileHeader.Header.Get("Content-Type"),
				Content:  content, // Store file content as []byte
				Size:     fileHeader.Size,
			}

			query := `
			INSERT INTO files (user_id, name, mime_type, content, size)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
			`
			// Execute the query and retrieve the ID of the inserted file
			err = db.QueryRow(query,
				uploadedFile.UserID,
				uploadedFile.Name,
				uploadedFile.MimeType,
				uploadedFile.Content,
				uploadedFile.Size,
			).Scan(&uploadedFile.ID)

			if err != nil {
				log.Fatal("could not insert file", err)
			} else {
				data.Success = true
				log.Printf("File inserted successfully with ID: %d", uploadedFile.ID)
				w.Header().Set("HX-Trigger", "file-uploaded")
			}
			// GET
		} else if r.Method == http.MethodGet {
			// Get all user created files from DB
			// Get the files uploaded by the user
			files, err := getUserFiles(db, userId)
			if err != nil {
				http.Error(w, "Failed to retrieve files", http.StatusInternalServerError)
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
