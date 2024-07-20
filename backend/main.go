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
	"time"
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
	Filepath   string `json:"filepath"`
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
		http.Error(w, "Error processing form", http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error getting form file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Create a file on the server
	savedFilePath := fmt.Sprintf("db_files/%s", handler.Filename)
	dst, err := os.Create(savedFilePath)
	if err != nil {
		http.Error(w, "Error saving file to disk", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the server file
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error saving file to disk", http.StatusInternalServerError)
		return
	}

	// Save the file path in the database
	err = saveFilePathToDB(handler.Filename, handler.Header.Get("Content-Type"), handler.Size, savedFilePath)
	if err != nil {
		http.Error(w, "Error saving file to database", http.StatusInternalServerError)
		return
	}

	response := Response{
		Message: "Successfully uploaded file",
		Path:    savedFilePath,
	}

	json.NewEncoder(w).Encode(response)
}

func saveFilePathToDB(filename, filetype string, filesize int64, filepath string) error {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME, DB_HOST, DB_PORT)

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return fmt.Errorf("Error connecting to database: %v", err)
	}
	defer db.Close()

	query := `INSERT INTO dbFiles (filename, filetype, filesize, filepath, uploaded_at) VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(query, filename, filetype, filesize, filepath, time.Now())
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

	rows, err := db.Query("SELECT id, filename, filetype, filesize, filepath, uploaded_at FROM dbFiles")
	if err != nil {
		http.Error(w, "Error retrieving files from database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var files []File

	for rows.Next() {
		var file File
		err := rows.Scan(&file.ID, &file.Filename, &file.Filetype, &file.Filesize, &file.Filepath, &file.UploadedAt)
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

	id, err := strconv.Atoi(idString)
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

// download file Function
func downloadFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		return
	}

	idString := r.URL.Query().Get("id")
	if idString == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "invalid file ID", http.StatusBadRequest)
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

	var filename, filepath string
	query := `SELECT filename, filepath FROM dbFiles WHERE id = $1`
	err = db.QueryRow(query, id).Scan(&filename, &filepath)
	if err != nil {
		http.Error(w, "Error retrieving file from database", http.StatusInternalServerError)
		log.Printf("Error retrieving file from database: %v", err)
		return
	}

	// extract file extension from filename for MIME types
	//ext := pf.Ext(filename)
	//if ext == "" {
	//	http.Error(w, "File extension not found", http.StatusInternalServerError)
	//	fmt.Println("File extension not found")
	//	log.Printf("File extension not found for filename: %s", filename)
	//	return
	//}
	//
	//fmt.Println("file is found", ext)
	//fmt.Fprintf(w, "File extension is %s", ext)
	// Get MIME type by extension
	//mimeType := mime.TypeByExtension(ext)
	//if mimeType == "" {
	//	http.Error(w, "Unknown MIME type", http.StatusInternalServerError)
	//	log.Printf("Uknown MIME type: %s", ext)
	//	return
	//}
	//fmt.Println("mime type is:", mimeType)

	http.ServeFile(w, r, filepath)
}

func setupRoutes() {
	mux := http.NewServeMux()

	mux.HandleFunc("/upload", uploadFile)
	mux.HandleFunc("/files", getFiles)
	mux.HandleFunc("/delete", deleteFiles)
	mux.HandleFunc("/download", downloadFile)
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
