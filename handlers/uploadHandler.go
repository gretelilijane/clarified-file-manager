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

type UploadData struct {
	ErrorMessage string
	Success      bool
}

func executeUpload(db *sql.DB, store *sessions.CookieStore, r *http.Request) (UploadData, error) {
	data := UploadData{
		Success:      false,
		ErrorMessage: "",
	}

	session, _ := store.Get(r, "session") // Custom session name
	userId, ok := session.Values["user_id"].(int32)

	if !ok {
		data.ErrorMessage = "Your session has expired, please log in again"
		return data, nil
	}

	// Retrieve file from the form
	f, fileHeader, err := r.FormFile("file")

	if err != nil {
		return data, err
	}

	// Check file size before reading the content
	max_file_size := int64(10 * 1024 * 1024) // 10MiB

	if fileHeader.Size > max_file_size {
		data.ErrorMessage = "File size exceeds the 10 MiB limit"
		return data, nil
	}

	// Parse the form data (multipart/form-data)
	err = r.ParseMultipartForm(max_file_size)

	if err != nil {
		return data, err
	}

	defer f.Close()

	// Retrieve the file content (reading the file into a byte slice)
	content, err := io.ReadAll(f)

	if err != nil {
		return data, err
	}

	// Insert the file into the database
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
		return data, err
	}

	data.Success = true
	log.Printf("File inserted successfully with ID: %d", uploadedFile.ID)

	return data, nil
}

func UploadHandler(db *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
	tmpl, err := template.ParseFiles("views/base.html", "views/files.html", "views/upload-form.html")

	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data, err := executeUpload(db, store, r)

		// Handle server errors
		if err != nil {
			log.Println("Upload error: ", err)
			data.Success = false
			data.ErrorMessage = "A server error occured, please try again later"
		}

		// Trigger refresh of files list when upload is successful
		if data.Success {
			w.Header().Set("HX-Trigger", "file-uploaded")
		}

		// Render response based on request target
		target := r.Header.Get("HX-Target")
		log.Println("Target: ", target)

		if target == "upload-form" {
			tmpl.ExecuteTemplate(w, target, data)
		} else {
			tmpl.Execute(w, data)
		}
	}
}
