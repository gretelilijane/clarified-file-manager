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
	Files        []types.File
	Sort         types.FileSortableColumn
	Direction    types.SortDirection
	UploadData   UploadData
	TableHeaders []types.TableHeader
}

func getUserFiles(db *sql.DB, userID int32, sortColumn types.FileSortableColumn, sortDirection types.SortDirection) ([]types.File, error) {
	var files []types.File

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

func executeFilesPage(db *sql.DB, store *sessions.CookieStore, w http.ResponseWriter, r *http.Request) (FilesPageData, error) {

	data := FilesPageData{
		Files:     nil,
		Sort:      "uploaded_at",
		Direction: "desc",
	}

	session, _ := store.Get(r, "session")
	userId, ok := session.Values["user_id"].(int32)

	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return data, nil
	}

	// Get sort and direction from query params
	query := r.URL.Query()

	if query.Has("sort") {
		data.Sort = types.FileSortableColumnFromString(query.Get("sort"))
	}

	if query.Has("dir") {
		data.Direction = types.SortDirectionFromString(query.Get("dir"))
	}

	// Set table headers
	headers := []types.TableHeader{
		{
			Title:   "Name",
			SortKey: "name",
		},
		{
			Title:   "Mime Type",
			SortKey: "mime_type",
		},
		{
			Title:   "Size",
			SortKey: "size",
		},
		{
			Title:   "Uploaded At",
			SortKey: "uploaded_at",
		},
		{
			Title: "Actions",
		},
	}

	for _, header := range headers {
		if header.SortKey == "" {
			continue
		}

		// Is this header currently sorted by
		if header.SortKey == data.Sort {
			if data.Direction == "asc" {
				header.Icon = "fa-sort-up"
				header.Link = fmt.Sprintf("/files?sort=%s&dir=desc", header.SortKey)
			} else {
				header.Icon = "fa-sort-down"
				header.Link = fmt.Sprintf("/files?sort=%s&dir=asc", header.SortKey)
			}
		} else {
			header.Link = fmt.Sprintf("/files?sort=%s&dir=asc", header.SortKey)
		}

		data.TableHeaders = append(data.TableHeaders, header)
	}

	// Fetch user files
	files, err := getUserFiles(db, userId, data.Sort, data.Direction)
	if err != nil {
		return data, err
	}

	data.Files = files

	return data, nil
}

func FilesPageHandler(db *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
	tmpl, err := template.ParseFiles("views/base.html", "views/files.html", "views/upload-form.html")

	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {

		data, err := executeFilesPage(db, store, w, r)

		// Handle server errors
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Render response based on request target
		target := r.Header.Get("HX-Target")

		if target == "messages" || target == "files-table" {
			tmpl.ExecuteTemplate(w, target, data)
		} else if target == "file-list" {
			tmpl.ExecuteTemplate(w, target, data.Files)
		} else {
			tmpl.Execute(w, data)
		}
	}
}
