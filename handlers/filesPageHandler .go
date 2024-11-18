package handlers

import (
	"clarified-file-management/types"
	"database/sql"
	"fmt"
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
	Sort         types.FileSortableColumn
	Direction    types.SortDirection
}

func getUserFiles(db *sql.DB, userID int32, sortColumn types.FileSortableColumn, sortDirection types.SortDirection) ([]types.File, error) {
	var files []types.File
	//         query := fmt.Sprintf("SELECT id, name, age, created_at FROM users ORDER BY %s %s", sortColumn, sortOrder)
	// query := fmt.Sprintf("SELECT id, name, mime_type, size, uploaded_at FROM files WHERE user_id = $1 ORDER BY $2 $3", userID, sort, desc)

	// Construct the query safely
	query := fmt.Sprintf("SELECT id, name, mime_type, size, uploaded_at FROM files WHERE user_id = $1 ORDER BY %s %s", sortColumn, sortDirection)
	rows, err := db.Query(query, userID)

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
			Sort:         "uploaded_at",
			Direction:    "desc",
		}

		if r.Method == http.MethodGet {
			query := r.URL.Query()

			// Get sort and direction from query params
			if query.Has("sort") {
				data.Sort = types.FileSortableColumnFromString(query.Get("sort"))
			}

			if query.Has("dir") {
				data.Direction = types.SortDirectionFromString(query.Get("dir"))
			}

			log.Println("Sort: ", data.Sort, "Direction: ", data.Direction)

			files, err := getUserFiles(db, userId, data.Sort, data.Direction)
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
