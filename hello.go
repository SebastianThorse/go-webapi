package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go-webapi/models"
	"log"
	"net/http"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "go-webapi"
)

func insertPost(db *sql.DB, post models.Post) error {
	stmt, err := db.Prepare("INSERT INTO posts (title, description, created_at, updated_at, author_name) VALUES ($1, $2, $3, $4, $5) RETURNING id")
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := time.Now()

	var updatedTime sql.NullTime
	if post.UpdatedAt == "" {
		updatedTime.Valid = false
	} else {
		updatedAt, err := time.Parse(time.RFC3339, post.UpdatedAt)
		if err != nil {
			return err
		}
		updatedTime.Time = updatedAt
		updatedTime.Valid = true
	}

	_, err = stmt.Exec(post.Title, post.Description, now, updatedTime, post.AuthorName)
	if err != nil {
		return err
	}

	return nil
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	var postRequest models.PostRequest
	err := json.NewDecoder(r.Body).Decode(&postRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if postRequest.Title == "" || postRequest.Description == "" || postRequest.AuthorName == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	post := models.Post{
		Title:       postRequest.Title,
		Description: postRequest.Description,
		AuthorName:  postRequest.AuthorName,
		UpdatedAt:   postRequest.UpdatedAt,
	}

	sqlcode := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", sqlcode)
	err = insertPost(db, post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Post added")
}

func main() {
	sqlcode := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", sqlcode)
	CheckError(err)

	defer db.Close()

	err = db.Ping()
	CheckError(err)

	r := mux.NewRouter()

	r.HandleFunc("/posts", createPostHandler).Methods("POST")

	fmt.Println("Connected!")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
