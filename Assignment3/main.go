package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Book struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Cost        float32 `json:"cost"`
}

func main() {
	dsn := "host=localhost user=postgres password=12345 dbname=bookstore port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&Book{})

	r := mux.NewRouter()

	// Define the REST endpoints
	r.HandleFunc("/", homepage)
	r.HandleFunc("/books/{id}", getBookByIDHandler(db)).Methods("GET")
	r.HandleFunc("/books", getAllBooksHandler(db)).Methods("GET")
	r.HandleFunc("/books/{id}", updateBookByIDHandler(db)).Methods("PUT")
	r.HandleFunc("/books/{id}", deleteBookByIDHandler(db)).Methods("DELETE")
	r.HandleFunc("/search/{title}", searchBookByTitleHandler(db)).Methods("GET")
	r.HandleFunc("/books", addBookHandler(db)).Methods("POST")
	r.HandleFunc("/books/sort/{order}", sortBooksByCostHandler(db)).Methods("GET")

	fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", r)
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!")
}

// Handler for GET /books/{id}
func getBookByIDHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		var book Book
		if err := db.First(&book, params["id"]).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(book)
	}
}

// Handler for GET /books
func getAllBooksHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var books []Book
		if err := db.Find(&books).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(books)
	}
}

// Handler for PUT /books/{id}
func updateBookByIDHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		var book Book
		if err := db.First(&book, params["id"]).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var updatedBook Book
		if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		book.Title = updatedBook.Title
		book.Description = updatedBook.Description
		book.Cost = updatedBook.Cost
		if err := db.Save(&book).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(book)
	}
}

// Handler for DELETE /books/{id}
func deleteBookByIDHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		var book Book
		if err := db.First(&book, params["id"]).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := db.Delete(&book).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// Handler for GET /search?title={title}
func searchBookByTitleHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// title := r.URL.Query().Get("title")
		// if query == "Yera" {
		// 	fmt.Fprint(w, "Yera!")
		// }
		params := mux.Vars(r)
		title := params["title"]
		if title == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var books []Book
		if err := db.Where("title LIKE ?", "%"+title+"%").Find(&books).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(books)
	}

}

// Handler for POST /books
func addBookHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := db.Create(&book).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(book)
	}
}

// Handler for GET /books/sort/{order}
func sortBooksByCostHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		order := strings.ToLower(params["order"])
		if order != "asc" && order != "desc" {
			fmt.Fprint(w, "ERROR!")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var books []Book
		var err error
		if order == "asc" {
			err = db.Order("id asc").Find(&books).Error
		} else {
			err = db.Order("id desc").Find(&books).Error
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(books)
	}
}
