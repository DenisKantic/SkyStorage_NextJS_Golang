package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Response struct {
	Message string `json:"message"`
	Path    string `json:"path"`
}

type File struct {
	ID         int    `json:"id"`
	Filename   string `json:"filename"`
	Filetype   string `json:"filetype"`
	Filesize   int64  `json:"filesize"`
	UploadedAt string `json:"uploaded_at"`
}

var (
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	DB_HOST     string
	DB_PORT     string
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// Allow CORS
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	if r.Method == http.MethodOptions {
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		fmt.Println("Error happened in processing form", err)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error happened in getting form file", err)
		return
	}
	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file", err)
		return
	}

	savedFilePath := fmt.Sprintf("db_files/%s", handler.Filename)
	err = os.WriteFile(savedFilePath, fileBytes, 0644)
	if err != nil {
		fmt.Println("Error saving file to disk", err)
		return
	}

	fmt.Printf("Successfully uploaded file. Saved path: %s\n", savedFilePath)

	err = saveFilePathToDB(handler.Filename, handler.Header.Get("Content-Type"), handler.Size)
	if err != nil {
		fmt.Println("Error saving file to database", err)
		return
	}

	response := Response{
		Message: "Successfully uploaded file",
		Path:    savedFilePath,
	}

	json.NewEncoder(w).Encode(response)
}

func saveFilePathToDB(filename, filetype string, filesize int64) error {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME, DB_HOST, DB_PORT)

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return fmt.Errorf("Error connecting to database: %v", err)
	}
	defer db.Close()

	query := `INSERT INTO dbFiles (filename, filetype, filesize) VALUES ($1, $2, $3)`
	_, err = db.Exec(query, filename, filetype, filesize)
	if err != nil {
		return fmt.Errorf("Error saving file to database: %v", err)
	}

	return nil
}

func getFiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME, DB_HOST, DB_PORT)

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, filename, filetype, filesize, uploaded_at FROM dbFiles")
	if err != nil {
		http.Error(w, "Error retrieving files from database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var files []File

	for rows.Next() {
		var file File
		err := rows.Scan(&file.ID, &file.Filename, &file.Filetype, &file.Filesize, &file.UploadedAt)
		if err != nil {
			http.Error(w, "Error scanning file data", http.StatusInternalServerError)
			return
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error with rows", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(files)
}

func deleteFiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")

	if r.Method == http.MethodOptions {
		return
	}

	idString := r.URL.Query().Get("id")
	if idString == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idString) // need to google what is this
	if err != nil {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME, DB_HOST, DB_PORT)

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	query := `DELETE FROM dbFiles WHERE id = $1`
	_, err = db.Exec(query, id)
	if err != nil {
		http.Error(w, "Error deleting file", http.StatusInternalServerError)
		return
	}

	response := Response{
		Message: "Successfully deleted file",
	}

	json.NewEncoder(w).Encode(response)
}
func setupRoutes() {

	mux := http.NewServeMux()

	mux.HandleFunc("/upload", uploadFile)
	mux.HandleFunc("/files", getFiles)
	mux.HandleFunc("/delete", deleteFiles)
	corsHandler := cors.Default().Handler(mux)

	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	DB_USER = os.Getenv("DB_USER")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_NAME = os.Getenv("DB_NAME")
	DB_HOST = os.Getenv("DB_HOST")
	DB_PORT = os.Getenv("DB_PORT")

	fmt.Println("Server is running on port 8080")
	setupRoutes()
}
