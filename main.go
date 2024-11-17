package main

import (
	"clarified-file-management/handlers"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	//"encoding/json" // TODO: think if needed

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// var tmpl *template.Template
var db *sql.DB

func init() {
	log.Println("Initializing the application")
	// tmpl, _ = template.ParseGlob("views/*.html")
}

func loadEnv() error {
	return godotenv.Load()
}

func main() {

	err := loadEnv()
	if err != nil {
		log.Println("No .env file found, using default configuration")
	}

	var server_port string = os.Getenv("SERVER_PORT")
	var server_host string = os.Getenv("SERVER_HOST")
	var db_host string = os.Getenv("DB_HOST")
	var db_port_str string = os.Getenv("DB_PORT")
	var db_user string = os.Getenv("DB_USER")
	var db_password string = os.Getenv("DB_PASSWORD")
	var db_name string = os.Getenv("DB_NAME")

	db_port, err := strconv.Atoi(db_port_str)
	if err != nil {
		log.Fatal("Invalid port number: ", db_port_str)
		return
	}

	log.Println("Port: ", db_port, "Host: ", db_host, "DB User: ", db_user, "DB Password: ", db_password, "DB Name: ", db_name)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", // TODO: change sslmode
		db_host, db_port, db_user, db_password, db_name)

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

	router := mux.NewRouter()

	// router.HandleFunc("/", Homepage)
	// router.HandleFunc("/upload", UploadHandler).Methods("POST")

	router.HandleFunc("/", handlers.IndexPageHandler()).Methods("GET", "OPTIONS")
	router.HandleFunc("/upload", handlers.UploadPageHandler()).Methods("GET")
	router.HandleFunc("/login", handlers.LogInPageHandler(db)).Methods("GET", "POST")
	router.HandleFunc("/signup", handlers.SignUpPageHandler(db)).Methods("GET", "POST")
	// router.HandleFunc("/sign-up", handlers.SignUpHandler()).Methods("POST", "OPTIONS")
	// router.HandleFunc("/login", logInPageHandler).Methods("POST", "OPTIONS")
	// // router.HandleFunc("/servelogin", handleServeLoginRequest).Methods("GET")
	// router.HandleFunc("/logout", logOutPageHandler).Methods("POST", "OPTIONS")
	// router.HandleFunc("/signup", signUpPageHandler).Methods("POST", "OPTIONS")

	//http.HandleFunc("/login", loginHandler)
	log.Println("Server started at http://" + server_host + ":" + server_port)
	log.Fatal(http.ListenAndServe(server_host+":"+server_port, router))
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type FileMetaData struct {
	Filename  string `json:"filename"`
	FileSize  int64  `json:"filesize"`
	MimeType  string `json:"mimetype"`
	CreatedAt string `json:"createdat"`
}

// func Homepage(w http.ResponseWriter, r *http.Request) {
// 	log.Println("Homepage")
// 	tmpl.ExecuteTemplate(w, "index.html", nil)
// }

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
