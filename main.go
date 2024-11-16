package main

import (
	"database/sql"
	"log"
	"os"
	"fmt"
	"time"
	"net/http"
	"html/template"
	"strconv"
	//"encoding/json" // TODO: think if needed

	// _ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var tmpl *template.Template
var db *sql.DB

func init() {
	log.Println("Initializing the application")
	tmpl, _ = template.ParseGlob("views/*.html")
}

func loadEnv() error {
	return godotenv.Load()
}

func initDB() (error) {
	db, err := sql.Open("sqlite3", "./files.db")
    if err != nil {
        return err
    }
    // Create table if not exists
    query := `CREATE TABLE IF NOT EXISTS files (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        filename TEXT,
		filesize INTEGER,
        created_at DATETIME,
        mime_type TEXT
    );`
    _, err = db.Exec(query)
    if err != nil {
        return err
    }

	err = db.Ping()
	if err != nil {
		// do something here
		log.Fatal("No ping: ", err)
		return  err
	} else {
		log.Println("Ping successful")
	}

    return nil
}

func main() {

	err := loadEnv()
	if err != nil {
		log.Println("No .env file found, using default configuration")
	}

	var port_str string = os.Getenv("PORT")
	var host string = os.Getenv("HOST")
	var db_user string = os.Getenv("DB_USER")
	var db_password string = os.Getenv("DB_PASSWORD")
	var db_name string = os.Getenv("DB_NAME")

	port, err := strconv.Atoi(port_str)
	if err != nil {
		log.Fatal("Invalid port number: ", port_str)
		return
	}

	log.Println("Port: ", port, "Host: ", host, "DB User: ", db_user, "DB Password: ", db_password, "DB Name: ", db_name)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable", // TODO: change sslmode
    host, port, db_user, db_password, db_name)

	db, err := sql.Open("postgres", psqlInfo) // does not open the connection, only validates the args
	if err != nil {
		panic(err)
	}
	defer db.Close() // closes db connection when main() finishes

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	gRouter := mux.NewRouter()

	// Initialize the database
    // err = initDB()
    // if err != nil {
    //     log.Fatal("Failed to initialize database:", err)
    // }
    // defer db.Close()

	gRouter.HandleFunc("/", Homepage)
	gRouter.HandleFunc("/upload", UploadHandler).Methods("POST")

	//http.HandleFunc("/login", loginHandler)
	log.Println("Server started at http://" + host + ":" + port_str)
	log.Fatal(http.ListenAndServe(host+":"+port_str, gRouter))
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type FileMetaData struct {
	Filename string `json:"filename"`
	FileSize int64 `json:"filesize"`
	MimeType string `json:"mimetype"`
	CreatedAt string `json:"createdat"`
}

func Homepage(w http.ResponseWriter, r *http.Request) {
	log.Println("Homepage")
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("UploadHandler")
	r.ParseMultipartForm(10 << 20)
	log.Println("UploadHandler2")
	file, handler, err := r.FormFile("avatar")
	log.Println("UploadHandler3")
    if err != nil {
		log.Println("Error Retrieving the File")
        // http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
        return
    }
	log.Println("UploadHandler4")
	log.Println("File loaded")

	// log.Println("File: ", fileHeader.Filename, fileHeader.Header.Get("Content-Type"))

	filename := handler.Filename
    filesize := handler.Size
    fmt.Printf("MIME Header: %+v\n", handler.Header)
	content_type := handler.Header.Get("Content-Type")
	createdAt := time.Now().Unix()

	fmt.Println("Content Type: ", content_type)
	fmt.Println("Created At: ", createdAt)

	humanReadableTime := time.Unix(createdAt, 0)

    fmt.Println("Human-readable time:", humanReadableTime)

	query := `INSERT INTO files (filename, filesize, created_at, mime_type) VALUES (?, ?, ?, ?)`
    _, err = db.Exec(query, filename, filesize, time.Now(), content_type)

	if err != nil {
        http.Error(w, "Error saving metadata", http.StatusInternalServerError)
        return
    }

    defer file.Close()


	// defer file.Close()
	// log.Printf("Uploaded File: %+v\n", handler.Filename)
	// log.Printf("File Size: %+v\n", handler.Size)
	// log.Printf("MIME Header: %+v\n", handler.Header)
	// f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// defer f.Close()
	// io.Copy(f, file)
	// fmt.Fprintf(w, "Successfully Uploaded File\n")
}
