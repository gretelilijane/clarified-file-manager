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

func allowAllOriginsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		origin := r.Header.Get("Origin")

		if origin != "" {
			// Set CORS headers to allow credentials
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			//w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		}

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default configuration")
	}

	var server_port string = os.Getenv("SERVER_PORT")
	var server_host string = os.Getenv("SERVER_HOST")
	var session_key string = os.Getenv("SESSION_KEY")
	var session_secure_str string = os.Getenv("SESSION_SECURE")
	var session_http_only_str string = os.Getenv("SESSION_HTTP_ONLY")
	var session_max_age_str string = os.Getenv("SESSION_MAX_AGE")
	// var session_same_site string = os.Getenv("SESSION_SAME_SITE")
	var db_host string = os.Getenv("DATABASE_HOST")
	var db_port_str string = os.Getenv("DATABASE_PORT")
	var db_user string = os.Getenv("POSTGRES_USER")
	var db_password string = os.Getenv("POSTGRES_PASSWORD")
	var db_name string = os.Getenv("POSTGRES_DB")

	// Set up cookie session store
	session_secure, err := strconv.ParseBool(session_secure_str)
	if err != nil {
		log.Fatal("Invalid boolean value for SESSION_SECURE: ", session_secure_str)
	}

	session_http_only, err := strconv.ParseBool(session_http_only_str)
	if err != nil {
		log.Fatal("Invalid boolean value for SESSION_HTTP_ONLY: ", session_http_only_str)
	}

	session_max_age, err := strconv.Atoi(session_max_age_str)
	if err != nil {
		log.Fatal("Invalid integer value for SESSION_MAX_AGE: ", session_max_age_str)
	}

	store := sessions.NewCookieStore([]byte(session_key))
	store.Options = &sessions.Options{Secure: session_secure, HttpOnly: session_http_only, MaxAge: session_max_age, SameSite: http.SameSiteLaxMode}
	//store.Options = &sessions.Options{Secure: session_secure, HttpOnly: session_http_only, MaxAge: session_max_age, SameSite: http.SameSiteNoneMode}

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
	//log.Fatal(http.ListenAndServe(server_host+":"+server_port, router))
	log.Fatal(http.ListenAndServe(server_host+":"+server_port, allowAllOriginsMiddleware(router)))

}
