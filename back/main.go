package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Article represents the data structure for articles
type Article struct {
	ArticleID int    `json:"article_id"`
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	Link      string `json:"link"`
}

// Comment represents the data structure for comments
type Comment struct {
	CommentID int    `json:"comment_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	ArticleID int    `json:"article_id"`
}

// Links represents the data structure for hypermedia links
type Links struct {
	Self   string `json:"self"`
	Parent string `json:"parent"`
}

var (
	db *sql.DB
)

func initDB() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/tpapigolang")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initDB()
	r := mux.NewRouter()

	// API routes
	apiRouter := r.PathPrefix("/v1").Subrouter()
	apiRouter.HandleFunc("/", RootHandler).Methods("GET")
	apiRouter.HandleFunc("/articles/search", SearchArticlesHandler).Methods("GET")
	apiRouter.HandleFunc("/articles", GetArticlesHandler).Methods("GET")
	apiRouter.HandleFunc("/articles/{article_id}", GetArticleHandler).Methods("GET")
	apiRouter.HandleFunc("/articles", PostArticleHandler).Methods("POST")
	apiRouter.HandleFunc("/articles/{article_id}", PutArticleHandler).Methods("PUT")
	apiRouter.HandleFunc("/articles/{article_id}", DeleteArticleHandler).Methods("DELETE")
	apiRouter.HandleFunc("/articles/{article_id}/comments", GetCommentsHandler).Methods("GET")
	apiRouter.HandleFunc("/articles/{article_id}/comments", PostCommentHandler).Methods("POST")

	// Start the server
	http.Handle("/", r)
	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}

// RootHandler handles the root endpoint
func RootHandler(w http.ResponseWriter, r *http.Request) {
	paths := []string{"/v1/articles"}
	jsonResponse(w, map[string]interface{}{"paths": paths})
}

// SearchArticlesHandler handles the search articles endpoint
func SearchArticlesHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for search articles
	// ...
}

func GetArticlesHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve query parameters (e.g., page, filter) from the request
	params := r.URL.Query()
	pageStr := params.Get("page")

	// Convert page parameter to integer
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		http.Error(w, "Invalid page parameter", http.StatusBadRequest)
		return
	}

	// Calculate the offset for pagination
	offset := (page - 1) * 5

	// Query the database to get a list of articles with pagination
	rows, err := db.Query("SELECT * FROM Article LIMIT 5 OFFSET ?", offset)
	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		return
	}

	// Iterate over the rows and build the list of articles
	var articles []Article
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ArticleID, &article.Title, &article.Slug, &article.Content, &article.Author)
		if err != nil {
			http.Error(w, "Error scanning database rows", http.StatusInternalServerError)
			return
		}
		article.Link = fmt.Sprintf("/v1/articles/%d", article.ArticleID)
		articles = append(articles, article)
	}

	// Build and send the JSON response
	response := map[string]interface{}{
		"articles": articles,
		"_links": map[string]string{
			"self":   fmt.Sprintf("/v1/articles?page=%d", page),
			"parent": "/v1/",
			"prev":   fmt.Sprintf("/v1/articles?page=%d", page-1),
			"next":   fmt.Sprintf("/v1/articles?page=%d", page+1),
			"first":  "/v1/articles?page=1",
			"last":   fmt.Sprintf("/v1/articles?page=%d", page),
		},
	}
	jsonResponse(w, response)
}

// GetArticleHandler handles the get article endpoint
func GetArticleHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the article ID from the URL parameters
	vars := mux.Vars(r)
	articleIDStr := vars["article_id"]

	// Convert article ID parameter to integer
	articleID, err := strconv.Atoi(articleIDStr)
	if err != nil {
		http.Error(w, "Invalid article_id parameter", http.StatusBadRequest)
		return
	}

	// Open a connection to the SQLite database
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}

	// Query the database to get the specific article
	row := db.QueryRow("SELECT * FROM Article WHERE article_id = ?", articleID)

	// Build the article from the row
	var article Article
	err = row.Scan(&article.ArticleID, &article.Title, &article.Slug, &article.Content, &article.Author)
	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	// Build and send the JSON response
	jsonResponse(w, article)
}

// PostArticleHandler handles the post article endpoint
func PostArticleHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the Article data
	var article Article
	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Insert the new article into the database
	result, err := db.Exec("INSERT INTO Article(title, slug, content, author) VALUES(?, ?, ?, ?)",
		article.Title, article.Slug, article.Content, article.Author)
	if err != nil {
		http.Error(w, "Error inserting into the database", http.StatusInternalServerError)
		return
	}

	// Get the last inserted ID
	articleID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Error getting the last inserted ID", http.StatusInternalServerError)
		return
	}

	// Set the article ID and construct the response
	article.ArticleID = int(articleID)
	article.Link = fmt.Sprintf("/v1/articles/%d", article.ArticleID)

	// Build and send the JSON response
	jsonResponse(w, article)
}

// PutArticleHandler handles the put article endpoint
func PutArticleHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the article ID from the URL parameters
	vars := mux.Vars(r)
	articleIDStr := vars["article_id"]

	// Convert article ID parameter to integer
	articleID, err := strconv.Atoi(articleIDStr)
	if err != nil {
		http.Error(w, "Invalid article_id parameter", http.StatusBadRequest)
		return
	}

	// Parse the request body to get the updated Article data
	var updatedArticle Article
	err = json.NewDecoder(r.Body).Decode(&updatedArticle)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update the article in the database
	_, err = db.Exec("UPDATE Article SET title=?, slug=?, content=?, author=? WHERE article_id=?",
		updatedArticle.Title, updatedArticle.Slug, updatedArticle.Content, updatedArticle.Author, articleID)
	if err != nil {
		http.Error(w, "Error updating the database", http.StatusInternalServerError)
		return
	}

	// Build and send the JSON response
	updatedArticle.ArticleID = articleID
	updatedArticle.Link = fmt.Sprintf("/v1/articles/%d", articleID)
	jsonResponse(w, updatedArticle)
}

// DeleteArticleHandler handles the delete article endpoint
func DeleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the article ID from the URL parameters
	vars := mux.Vars(r)
	articleIDStr := vars["article_id"]

	// Convert article ID parameter to integer
	articleID, err := strconv.Atoi(articleIDStr)
	if err != nil {
		http.Error(w, "Invalid article_id parameter", http.StatusBadRequest)
		return
	}

	// Delete the article from the database
	_, err = db.Exec("DELETE FROM Article WHERE article_id=?", articleID)
	if err != nil {
		http.Error(w, "Error deleting from the database", http.StatusInternalServerError)
		return
	}

	// Send a 204 No Content response
	w.WriteHeader(http.StatusNoContent)
}

// GetCommentsHandler handles the get comments endpoint
func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the article ID from the URL parameters
	vars := mux.Vars(r)
	articleIDStr := vars["article_id"]

	// Convert article ID parameter to integer
	articleID, err := strconv.Atoi(articleIDStr)
	if err != nil {
		http.Error(w, "Invalid article_id parameter", http.StatusBadRequest)
		return
	}

	// Query the database to get the comments for the specified article
	rows, err := db.Query("SELECT * FROM Comment WHERE article_id=?", articleID)
	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		return
	}

	// Build the list of comments
	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.CommentID, &comment.Title, &comment.Content, &comment.ArticleID)
		if err != nil {
			http.Error(w, "Error scanning database rows", http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}

	// Build and send the JSON response
	jsonResponse(w, comments)
}

// PostCommentHandler handles the post comment endpoint
func PostCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the article ID from the URL parameters
	vars := mux.Vars(r)
	articleIDStr := vars["article_id"]

	// Convert article ID parameter to integer
	articleID, err := strconv.Atoi(articleIDStr)
	if err != nil {
		http.Error(w, "Invalid article_id parameter", http.StatusBadRequest)
		return
	}

	// Parse the request body to get the Comment data
	var comment Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Insert the new comment into the database
	result, err := db.Exec("INSERT INTO Comment(title, content, article_id) VALUES(?, ?, ?)",
		comment.Title, comment.Content, articleID)
	if err != nil {
		http.Error(w, "Error inserting into the database", http.StatusInternalServerError)
		return
	}

	// Get the last inserted ID
	commentID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Error getting the last inserted ID", http.StatusInternalServerError)
		return
	}

	// Set the comment ID and article ID, and construct the response
	comment.CommentID = int(commentID)
	comment.ArticleID = articleID

	// Build and send the JSON response
	jsonResponse(w, comment)
}

// jsonResponse writes a JSON response to the given http.ResponseWriter
func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
