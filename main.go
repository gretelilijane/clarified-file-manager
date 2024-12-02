package main

import (
	"clarified-file-management/handlers"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default configuration")
	}

	var server_port string = os.Getenv("SERVER_PORT")
	var server_host string = os.Getenv("SERVER_HOST")
	var session_key string = os.Getenv("SESSION_KEY")
	var db_host string = os.Getenv("DATABASE_HOST")
	var db_port_str string = os.Getenv("DATABASE_PORT")
	var db_user string = os.Getenv("POSTGRES_USER")
	var db_password string = os.Getenv("POSTGRES_PASSWORD")
	var db_name string = os.Getenv("POSTGRES_DB")

	// Set up cookie session store
	store := sessions.NewCookieStore([]byte(session_key))

	db_port, err := strconv.Atoi(db_port_str)
	if err != nil {
		log.Fatal("Invalid port number: ", db_port_str)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", // TODO: change sslmode
		db_host, db_port, db_user, db_password, db_name)

	db, err := sql.Open("postgres", psqlInfo) // does not open the connection, only validates the args
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // closes db connection when main() finishes

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/", handlers.IndexPageHandler(store)).Methods("GET", "OPTIONS")
	router.HandleFunc("/login", handlers.LogInPageHandler(db, store)).Methods("GET", "POST")
	router.HandleFunc("/signup", handlers.SignUpPageHandler(db)).Methods("GET", "POST")
	router.HandleFunc("/logout", handlers.LogOutHandler(store))
	router.HandleFunc("/files", handlers.FilesPageHandler(db, store)).Methods("GET")
	router.HandleFunc("/files", handlers.UploadHandler(db, store)).Methods("POST")
	router.HandleFunc("/files/{id}", handlers.DeleteFileHandler(db, store)).Methods("DELETE")
	router.HandleFunc("/files/{id}", handlers.DownloadFileHandler(db, store)).Methods("GET")

	// static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server started at http://" + server_host + ":" + server_port)
	log.Fatal(http.ListenAndServe(server_host+":"+server_port, router))
}
